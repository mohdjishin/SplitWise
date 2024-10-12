package models

import (
	"time"

	"gorm.io/gorm"
)

// Bill represents an expense associated with a group
type Bill struct {
	gorm.Model `json:"-"`
	Name       string        `json:"name"`
	Amount     float64       `json:"amount"`    // Total amount
	GroupID    uint          `json:"groupId"`   // Reference to the associated group
	Completed  bool          `json:"completed"` // Overall bill payment status
	History    []BillHistory `json:"history"`   // Bill payment history
}
type BillHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BillID    uint      `json:"billId"`    // Automatically inferred foreign key
	Amount    float64   `json:"amount"`    // Amount related to this history entry
	PaidBy    string    `json:"paidBy"`    // User who made the payment
	PaidAt    time.Time `json:"paidAt"`    // Time of payment
	CreatedAt time.Time `json:"createdAt"` // Auto-create timestamp
}
