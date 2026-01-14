package entity

import "gorm.io/gorm"

// Admin User
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

// Login Payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
