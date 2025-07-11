package model

import (
	"time"
)

// Base model - Simplified for MVP without SSO
type Student struct {
	ID            int64     `json:"id"`
	StudentNumber string    `json:"student_number"` // Will be linked to SSO later
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Major         string    `json:"major"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Create model - Only essential fields for MVP
type CreateStudent struct {
	StudentNumber string `json:"student_number" validate:"required,min=8,max=20"`
	FullName      string `json:"full_name" validate:"required,min=3,max=100"`
	Phone         string `json:"phone" validate:"required,min=10,max=15"`
	Email         string `json:"email" validate:"required,email"`
	Major         string `json:"major" validate:"required,max=100"`
}

// Update model - Only fields that might need updating
type UpdateStudent struct {
	Phone string `json:"phone" validate:"required,min=10,max=15"`
	Email string `json:"email" validate:"required,email"`
}
