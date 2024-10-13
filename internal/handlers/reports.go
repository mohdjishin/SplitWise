package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	e "errors"

	"github.com/go-chi/chi"
	"github.com/mohdjishin/SplitWise/helper/pdf"
	"github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/errors"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	"github.com/mohdjishin/SplitWise/internal/models"
	"github.com/mohdjishin/SplitWise/internal/models/dto"
	log "github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

// GetGroupReport handles downloading a PDF report for a user's groups
// @Summary Download PDF report of user's groups
// @Description Generates and downloads a PDF report for the groups created by the user within a specified date range.
// @Tags reports
// @Accept json
// @Produce application/pdf
// @Param Authorization header string true "Bearer token"
// @Param from query string false "from date in the format YYYY-MM-DD"
// @Param to query string false "to date in the format YYYY-MM-DD"
// @Param request body dto.GetGroupReportRequest true "GetGroupReportRequest details"
// @Success 200 {file} file "PDF report generated and downloaded"
// @Failure 400 {object} errors.Error "Bad Request"
// @Failure 404 {object} errors.Error "Not Found"
// @Failure 500 {object} errors.Error "Internal Server Error"
// @Router /v1/groups/report [post]
func GetGroupReport(w http.ResponseWriter, r *http.Request) {
	log.Info("GetGroupReport handler called")

	var req dto.GetGroupReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Invalid request payload", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrInvalidInput)
		return
	}
	log.Debug("GetGroupReport request", zap.Any("request", req))
	userId := middleware.GetCurrentUserId(r)
	var user models.User
	if err := db.GetDb().Select("name").Where("id = ?", userId).First(&user).Error; err != nil {
		log.Error("Database query failed", zap.Error(err))
		if e.Is(err, gorm.ErrRecordNotFound) {
			log.Info("No user found", zap.Float64("user_id", userId))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errors.ErrUserNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	var (
		toDate   time.Time
		fromDate time.Time
		err      error
	)

	if req.From != nil {
		fromDate, err = time.Parse("2006-01-02", *req.From)
		if err != nil {
			log.Error("Invalid from date format", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errors.ErrInvalid("invalid from date format"))
			return
		}
	} else {
		fromDate = time.Now().AddDate(0, 0, -7)
	}

	if req.To != nil {
		toDate, err = time.Parse("2006-01-02", *req.To)
		if err != nil {
			log.Error("Invalid to date format", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errors.ErrInvalid("invalid to date format"))
			return
		}
	} else {
		toDate = time.Now()
	}

	rows, err := db.GetDb().Table("groups").
		Select(`
		groups.id, groups.name, groups.status, groups.total_amount, groups.per_user_split_amount, groups.paid_amount,
		bills.amount AS bill_amount, bills.completed AS bill_paid, bills.created_at AS bill_date,
		COUNT(group_members.id) AS member_count, groups.total_amount AS total_split, groups.status 
	`).
		Joins("LEFT JOIN bills ON bills.group_id = groups.id").
		Joins("LEFT JOIN group_members ON group_members.group_id = groups.id").
		Where("groups.created_by = ? AND groups.created_at BETWEEN ? AND ?", userId, fromDate, toDate).
		Group("groups.id, bills.id").
		Rows()
	if err != nil {
		log.Error("Database query failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}
	defer rows.Close()

	var grpInfo []dto.Group
	for rows.Next() {
		var group dto.Group
		var bill dto.Bill
		var memberCount int
		var totalSplit float64
		var status string

		err = rows.Scan(
			&group.ID,
			&group.Name,
			&group.Status,
			&group.TotalAmount,
			&group.PerUserSplitAmount,
			&group.PaidAmount,
			&bill.Amount,
			&bill.Paid,
			&bill.Date,
			&memberCount,
			&totalSplit,
			&status,
		)
		if err != nil {
			log.Error("Failed to scan row data", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}

		group.Bills = bill
		group.Members = memberCount
		group.Owner = user.Name
		grpInfo = append(grpInfo, group)
	}
	if len(grpInfo) == 0 {
		log.Info("No groups found for user", zap.Float64("user_id", userId))
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	pdfResponse := pdf.GenerateOwnerGroupsReportPDF(grpInfo, fromDate, toDate, user.Name)

	var buf bytes.Buffer
	if err := pdfResponse.Output(&buf); err != nil {
		log.Error("Failed to generate PDF", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	fileName := fmt.Sprintf("%s_%s_%s_report.pdf", user.Name, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.WriteHeader(http.StatusOK)
	log.Info("PDF generated successfully")
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Error("Failed to write PDF to response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)

	}
	log.Info("PDF generated successfully")
}

func GenerateSingleGroupReport(w http.ResponseWriter, r *http.Request) {
	log.Debug("GenerateSingleGroupReport handler called")

	grpId := chi.URLParam(r, "id")
	userId := middleware.GetCurrentUserId(r)

	var group models.Group
	err := db.GetDb().Where("id = ? AND created_by = ?", grpId, userId).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("Group not found for user:", zap.Any("group_id", grpId), zap.Any("user_id", userId))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
			return
		}
		log.Error("Database error while fetching group:", zap.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	var bill models.Bill
	err = db.GetDb().Where("group_id = ?", group.ID).First(&bill).Error
	if err != nil {
		log.Error("Database error while fetching bill:", zap.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	var grpMembers []models.GroupMember
	err = db.GetDb().Where("group_id = ?", group.ID).Find(&grpMembers).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Database error while fetching group members:", zap.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}
	if len(grpMembers) == 0 {
		log.Warn("No group members found for group:", zap.Any("group_id", group.ID))
	}

	var billHistory []models.BillHistory
	err = db.GetDb().Where("bill_id = ?", bill.ID).Find(&billHistory).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Database error while fetching bill history:", zap.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}
	userMap := make(map[uint]string, len(grpMembers))
	if len(billHistory) == 0 {
		log.Warn("No bill history found for bill:", zap.Any("bill_id", bill.ID))
	} else {
		var memberIDs []int
		for _, member := range grpMembers {
			memberIDs = append(memberIDs, int(member.UserID))
		}
		var users []models.User
		err = db.GetDb().Select("id,name").Where("id IN ?", memberIDs).Find(&users).Error
		if err != nil {
			log.Error("Database error while fetching user details:", zap.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}

		for _, user := range users {
			userMap[user.ID] = user.Name
		}

	}

	log.Debug("[+]--->Group and associated data fetched successfully", zap.Any("group", group), zap.Any("bill", bill), zap.Any("members", grpMembers), zap.Any("history", billHistory))
	req := dto.GroupReportRequest{Group: group,
		Bill:     bill,
		Members:  grpMembers,
		History:  billHistory,
		UserInfo: userMap}
	pdfResponse := pdf.GenerateGroupDetailedReportPDF(req)

	var buf bytes.Buffer
	if err := pdfResponse.Output(&buf); err != nil {
		log.Error("Failed to generate PDF", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		return
	}

	fileName := fmt.Sprintf("%s_%s_report.pdf", userMap[uint(userId)], group.Name)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Error("Failed to write PDF to response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
	}
	log.Info("PDF generated successfully")
}
