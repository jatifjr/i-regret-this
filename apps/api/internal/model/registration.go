package model

import (
	"time"
)

// Base model
type Registration struct {
	ID           int64     `json:"id"`
	RegNumber    string    `json:"reg_number"` // Unique: format: reg_order_of_the_month/month_in_roman/year = 001/V/2025
	StudentID    int64     `json:"student_id"`
	TestPlotID   int64     `json:"test_plot_id"`
	PaymentID    int64     `json:"payment_id,omitempty"`
	Status       string    `json:"status"` // pending, payment_verified, approved, rejected
	TestDate     time.Time `json:"test_date,omitempty"`
	TestLocation string    `json:"test_location,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	ApprovedAt   time.Time `json:"approved_at,omitempty"`
	ApprovedBy   string    `json:"approved_by,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// History model for registration status changes
type RegistrationHistory struct {
	ID             int64     `json:"id"`
	RegistrationID int64     `json:"registration_id"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes,omitempty"`
	ChangedBy      string    `json:"changed_by"`
	CreatedAt      time.Time `json:"created_at"`
}

// Create model
type CreateRegistration struct {
	StudentID  int64 `json:"student_id" validate:"required"`
	TestPlotID int64 `json:"test_plot_id" validate:"required"`
}

// Update model
type UpdateRegistration struct {
	Status       string    `json:"status,omitempty"`
	TestDate     time.Time `json:"test_date,omitempty"`
	TestLocation string    `json:"test_location,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	ApprovedBy   string    `json:"approved_by,omitempty"`
}

// History create model
type CreateRegistrationHistory struct {
	RegistrationID int64  `json:"registration_id" validate:"required"`
	Status         string `json:"status" validate:"required"`
	Notes          string `json:"notes,omitempty"`
	ChangedBy      string `json:"changed_by" validate:"required"`
}
