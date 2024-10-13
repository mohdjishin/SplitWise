package dto

import "github.com/mohdjishin/SplitWise/internal/models"

// ListOwnedGroupsResponse represents the response body for listing owned groups.
// @Description ListOwnedGroupsResponse is the response model for listing owned groups.
// @Name ListOwnedGroupsResponse"
// @Property group Group "Group details"
// @Property members []GroupMember "Group members"
type ListOwnedGroupsResponse struct {
	Group   models.Group         `json:"group"`
	Members []models.GroupMember `json:"members"`
}
