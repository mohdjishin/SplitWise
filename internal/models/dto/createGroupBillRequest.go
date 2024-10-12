package dto

type CreateGroupWithBillRequest struct {
	GroupName string `json:"groupName" validate:"required"`
	Bill      struct {
		Name   string  `json:"name" validate:"required"`
		Amount float64 `json:"amount" validate:"required"`
	} `json:"bill" validate:"required"`
}
