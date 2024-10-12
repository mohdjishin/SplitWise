package models

import "gorm.io/gorm"

type Spending struct {
	gorm.Model
	ID          uint
	Amount      float64
	GroupID     uint
	CreatedBy   uint
	SplitMethod string // "equal", "exact", or "percentage"
}
