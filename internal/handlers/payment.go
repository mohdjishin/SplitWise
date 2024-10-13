package handlers

import (
	"encoding/json"
	e "errors"
	"net/http"

	"github.com/mohdjishin/SplitWise/helper"
	"github.com/mohdjishin/SplitWise/helper/validate"
	"github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/errors"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	"github.com/mohdjishin/SplitWise/internal/models"
	"github.com/mohdjishin/SplitWise/internal/models/dto"
	log "github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const done = "DONE"

// MarkPayment marks a payment for a specific group.
// @Summary Marks a payment for a group.
// @Description Marks a payment for a specific group and updates the group's payment status.
// @Tags payments
// @Accept json
// @Produce json
// @Param request body dto.MarkPaymentRequest true "Mark Payment Request"
// @Param groupId body uint true "groupId of the group for which the payment is marked"  // required
// @Success 200 {object} dto.MarkPaymentResponse
// @Failure 400 {object} errors.Error "Invalid input"
// @Failure 404 {object} errors.Error "Group not found or User not found"
// @Failure 409 {object} errors.Error "Payment already made"
// @Failure 500 {object} errors.Error "Internal server error"
// @Router /payments [post]
func MarkPayment(w http.ResponseWriter, r *http.Request) {
	var input dto.MarkPaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Error("Invalid request payload", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrInvalidInput)
		return
	}
	log.Debug("MarkPayment request", zap.Any("request", input))

	userID := middleware.GetCurrentUserId(r)
	if err := validate.ValidateStruct(input); err != nil {

		log.Error("Error validating request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrValidationFailed(err.Error()))
		return
	}

	var groupMember models.GroupMember
	if err := db.GetDb().Where("group_id = ? AND user_id = ?", input.GroupID, userID).First(&groupMember).Error; err != nil {
		log.Error("Group member not found", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	if groupMember.HasPaid {
		log.Warn("Payment already made by user", zap.Float64("user_id", userID))
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(errors.ErrPaymentAlreadyMade)
		return
	}

	groupMember.HasPaid = true
	groupMember.Remarks = input.Remarks
	if err := db.GetDb().Save(&groupMember).Error; err != nil {
		log.Error("Failed to update payment status", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrPaymentFailed)
		return
	}

	var group models.Group
	if err := db.GetDb().Where("id = ?", input.GroupID).First(&group).Error; err != nil {
		log.Error("Group not found", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	group.PaidAmount += groupMember.SplitAmount

	if group.PaidAmount >= group.TotalAmount {
		var bill models.Bill
		if err := db.GetDb().Where("id = ?", group.BillID).First(&bill).Error; err == nil {
			bill.Completed = true

			if err := db.GetDb().Save(&bill).Error; err != nil {
				log.Error("Failed to mark bill as completed", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(errors.ErrBillCompletionFailed)
				return
			}
		}
		group.Status = done
	}

	if err := db.GetDb().Save(&group).Error; err != nil {
		log.Error("Failed to update group paid amount", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupUpdateFailed)
		return
	}
	var user models.User
	if err := db.GetDb().Select("name").Where("id = ?", userID).First(&user).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found", zap.Float64("userID", userID))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errors.ErrUserNotFound)
		} else {
			log.Error("Error retrieving user", zap.Error(err))
			json.NewEncoder(w).Encode(errors.ErrInternalError)
		}
		return
	}
	_ = helper.LogBillHistory(group.BillID, groupMember.SplitAmount, user.Name)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.MarkPaymentResponse{Message: "Payment marked successfully"})
}

// GetPendingPayments retrieves all pending payments associated with the authenticated user.
// @Summary Retrieve Pending Payments
// @Description Fetches all pending payments for the current user that have not been paid yet, including group ID, group name, bill ID, and amount owed.
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} dto.PendingPaymentsWithTotalResponse "Successful response containing the list of pending payments and total amount."
// @Failure 404 {object} errors.Error "No pending payments found for the user."
// @Failure 500 {object} errors.Error "Internal server error occurred while fetching pending payments."
// @Router /v1/payments/pending-payments [get]
func GetPendingPayments(w http.ResponseWriter, r *http.Request) {
	userId := middleware.GetCurrentUserId(r)
	var groupMembers []models.GroupMember
	log.Debug("GetPendingPayments request", zap.Any("userId", userId))
	if err := db.GetDb().Where("user_id = ? AND has_paid = ?", userId, false).Find(&groupMembers).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("No pending payments found", zap.Float64("userId", userId))
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(errors.ErrNoPendingPayments)
		} else {
			log.Error("Failed to fetch pending payments", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
		}
		return
	}

	log.Debug("GetPendingPayments request", zap.Any("groupMembers", groupMembers))
	var pendingPayments []dto.PendingPayments
	totalAmount := 0.0

	for _, member := range groupMembers {
		var group models.Group
		if err := db.GetDb().Where("id = ?", member.GroupID).First(&group).Error; err != nil {
			log.Error("Failed to fetch group", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}

		var bill models.Bill
		if err := db.GetDb().Where("id = ?", group.BillID).First(&bill).Error; err != nil {
			log.Error("Failed to fetch bill", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errors.ErrInternalError)
			return
		}

		pendingPayments = append(pendingPayments, dto.PendingPayments{
			GroupID:   group.ID,
			GroupName: group.Name,
			BillID:    bill.ID,
			Amount:    bill.Amount,
		})
		totalAmount += bill.Amount
	}
	response := dto.PendingPaymentsWithTotalResponse{
		PendingPayments: pendingPayments,
		TotalAmount:     totalAmount,
		Message:         "",
	}
	if len(pendingPayments) == 0 {
		response.Message = "No pending payments found"
		// response.TotalAmount = 0
	}

	log.Debug("GetPendingPayments request", zap.Any("response", response))

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
