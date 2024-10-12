package handlers

import (
	"encoding/json"
	e "errors"
	"net/http"

	"github.com/mohdjishin/SplitWise/helper"
	"github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/errors"
	"github.com/mohdjishin/SplitWise/internal/middleware"
	"github.com/mohdjishin/SplitWise/internal/models"
	"github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const done = "DONE"

func MarkPayment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GroupID uint   `json:"groupId"`
		Remarks string `json:"remarks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.LoggerInstance.Error("Invalid request payload", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.ErrInvalidInput)
		return
	}

	userID := middleware.GetCurrentUserId(r)

	var groupMember models.GroupMember
	if err := db.GetDb().Where("group_id = ? AND user_id = ?", input.GroupID, userID).First(&groupMember).Error; err != nil {
		logger.LoggerInstance.Error("Group member not found", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupNotFound)
		return
	}

	if groupMember.HasPaid {
		logger.LoggerInstance.Warn("Payment already made by user", zap.Float64("user_id", userID))
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(errors.ErrPaymentAlreadyMade)
		return
	}

	groupMember.HasPaid = true
	groupMember.Remarks = input.Remarks
	if err := db.GetDb().Save(&groupMember).Error; err != nil {
		logger.LoggerInstance.Error("Failed to update payment status", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrPaymentFailed)
		return
	}

	var group models.Group
	if err := db.GetDb().Where("id = ?", input.GroupID).First(&group).Error; err != nil {
		logger.LoggerInstance.Error("Group not found", zap.Error(err))
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
				logger.LoggerInstance.Error("Failed to mark bill as completed", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(errors.ErrBillCompletionFailed)
				return
			}
		}
		group.Status = done
	}

	if err := db.GetDb().Save(&group).Error; err != nil {
		logger.LoggerInstance.Error("Failed to update group paid amount", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.ErrGroupUpdateFailed)
		return
	}
	var user models.User
	if err := db.GetDb().Select("name").Where("id = ?", userID).First(&user).Error; err != nil {
		if e.Is(err, gorm.ErrRecordNotFound) {
			logger.LoggerInstance.Warn("User not found", zap.Float64("userID", userID))
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			logger.LoggerInstance.Error("Error retrieving user", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	_ = helper.LogBillHistory(group.BillID, groupMember.SplitAmount, user.Name)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment marked successfully",
	})
}
