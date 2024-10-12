package dto

type LoginRequest struct {
	Email    string `json:"email" example:"jis@jish.com" validate:"required"`
	Password string `json:"password" example:"Passw0rd@123" validate:"required"`
}
