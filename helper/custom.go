package helper

import (
	"time"

	"github.com/mohdjishin/SplitWise/internal/db"
	"github.com/mohdjishin/SplitWise/internal/models"
	"github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap"
)

func LogBillHistory(billID uint, amount float64, paidBy string) error {
	history := models.BillHistory{
		BillID:    billID,
		Amount:    amount,
		PaidBy:    paidBy,
		PaidAt:    time.Now(),
		CreatedAt: time.Now(),
	}

	if err := db.GetDb().Create(&history).Error; err != nil {
		logger.LoggerInstance.Error("Failed to log bill history", zap.Error(err))
		return err
	}
	return nil
}
