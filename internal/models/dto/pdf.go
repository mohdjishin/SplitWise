package dto

import "time"

type Bill struct {
	Amount float64
	Paid   bool
	Date   time.Time
}

type Group struct {
	ID                 string
	Name               string
	Bills              Bill
	Owner              string
	Status             string
	Members            int
	TotalAmount        float64 `json:"totalAmount,omitempty"`
	PerUserSplitAmount float64 `json:"perUserSplitAmount,omitempty"`
	PaidAmount         float64 `json:"paidAmount,omitempty"`
}
