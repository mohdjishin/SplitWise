package models

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID                 uint           `gorm:"primarykey"`
	CreatedAt          time.Time      `json:"createdAt,omitempty"`
	UpdatedAt          time.Time      `json:"updatedAt,omitempty"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	Name               string         `json:"name,omitempty"`
	CreatedBy          uint           `json:"createdBy,omitempty"`
	BillID             uint           `json:"billId,omitempty"`
	Bill               *Bill          `json:"bill,omitempty" gorm:"constraint:OnDelete:SET NULL;"`
	TotalAmount        float64        `json:"totalAmount,omitempty"`
	PerUserSplitAmount float64        `json:"perUserSplitAmount,omitempty"`
	PaidAmount         float64        `json:"paidAmount,omitempty"`
	Status             string         `json:"status,omitempty" gorm:"default:PENDING"`
}

type GroupMember struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      `json:"createdAt,omitempty"`
	UpdatedAt   time.Time      `json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	GroupID     uint           `json:"groupId"`
	UserID      uint           `json:"userId"`
	HasPaid     bool           `json:"hasPaid"` // Tracks if the member has paid
	SplitAmount float64        `json:"splitAmount"`
	Remarks     string         `json:"remarks"`
}
