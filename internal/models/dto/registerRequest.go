package dto

type RegisterRequest struct {
	Email    string `json:"email" gorm:"unique" example:"user@example.com" validate:"required,email"`
	Password string `json:"password" example:"password123" validate:"required,password_complexity"`
	Name     string `json:"name" example:"John Doe" validate:"required"`
}
