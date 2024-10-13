package dto

import "github.com/mohdjishin/SplitWise/internal/models"

// ONLY USER FOR INTERNAL USE
type GroupReportRequest struct {
	Bill     models.Bill          `json:"bill"`
	Group    models.Group         `json:"group"`
	History  []models.BillHistory `json:"history"`
	Members  []models.GroupMember `json:"members"`
	UserInfo map[uint]string      `json:"-"`
}
