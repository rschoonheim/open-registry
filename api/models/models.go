package models

import (
	"time"

	"gorm.io/gorm"
)

// Base contains common columns for all tables
type Base struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User represents a user in the system
type User struct {
	Base
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Password is not exposed in JSON responses
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	IsAdmin  bool   `gorm:"default:false" json:"is_admin"`
}

// Authentication request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Authentication response
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
