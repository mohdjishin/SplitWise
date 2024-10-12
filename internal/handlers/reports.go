package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	e "errors"

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

	pdfResponse := pdf.GeneratePDFReport(grpInfo, fromDate, toDate, user.Name)

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
}
