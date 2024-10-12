package dto

type AddUsersToGroupRequest struct {
	// UserIds      []uint   `json:"userIds"`
	UserEmailIds []string `json:"userEmailIds" validate:"required,dive,email"`
}
