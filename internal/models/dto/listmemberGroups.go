package dto

import "github.com/mohdjishin/SplitWise/internal/models"

// ListMemberGroupsResponse represents the response for listing member groups.
// @Description Response model for listing groups the user belongs to, including group details and member information.
// @Name ListMemberGroupsResponse"
type ListMemberGroupsResponse struct {
	Group   models.Group         `json:"group"`
	Members []models.GroupMember `json:"members"`
}
