package models

import "time"

// User represents a user in the system.
// @Description User model for registration and login.
// @Name User
// @Property email string true "Email" format(email)
// @Property password string true "Password"
// @Property name string true "Name"
type User struct {
	ID        uint      `json:"id" example:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email" gorm:"unique" example:"user@example.com"`
	Password  string    `json:"password" example:"password123"`
	Name      string    `json:"name" example:"John Doe"`
}
