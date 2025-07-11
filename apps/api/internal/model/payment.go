package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// PaymentMethod represents the type of payment
type PaymentMethod string

const (
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	// PaymentMethodVirtualAccount PaymentMethod = "virtual_account"  // For future use
	// PaymentMethodPaymentGateway PaymentMethod = "payment_gateway"  // For future use
)

// PaymentStatus represents the status of payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"    // Initial state when payment method is selected
	PaymentStatusProcessing PaymentStatus = "processing" // Payment is being processed
	PaymentStatusPaid       PaymentStatus = "paid"       // Payment has been made
	PaymentStatusVerified   PaymentStatus = "verified"   // Admin has verified the payment
	PaymentStatusRejected   PaymentStatus = "rejected"   // Payment was rejected
	PaymentStatusExpired    PaymentStatus = "expired"    // Payment expired
	PaymentStatusFailed     PaymentStatus = "failed"     // Payment failed
	PaymentStatusCancelled  PaymentStatus = "cancelled"  // Payment was cancelled by user
)

// Implement sql.Scanner and driver.Valuer for PaymentMethod
func (pm *PaymentMethod) Scan(value interface{}) error {
	if value == nil {
		*pm = ""
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid payment method type: %T", value)
	}
	*pm = PaymentMethod(str)
	return nil
}

func (pm PaymentMethod) Value() (driver.Value, error) {
	return string(pm), nil
}

// Implement sql.Scanner and driver.Valuer for PaymentStatus
func (ps *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		*ps = ""
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid payment status type: %T", value)
	}
	*ps = PaymentStatus(str)
	return nil
}

func (ps PaymentStatus) Value() (driver.Value, error) {
	return string(ps), nil
}

// Base model
type Payment struct {
	ID             int64         `json:"id"`
	RegistrationID int64         `json:"registration_id"`
	Amount         float64       `json:"amount"`
	PaymentMethod  PaymentMethod `json:"payment_method"`
	PaymentStatus  PaymentStatus `json:"payment_status"`

	// Common fields
	ReceiptImage string    `json:"receipt_image,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	ExpiredAt    time.Time `json:"expired_at,omitempty"`
	PaidAt       time.Time `json:"paid_at,omitempty"`
	VerifiedBy   string    `json:"verified_by,omitempty"`
	VerifiedAt   time.Time `json:"verified_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Bank Transfer specific fields
	BankName      string    `json:"bank_name,omitempty"`
	AccountNumber string    `json:"account_number,omitempty"`
	AccountName   string    `json:"account_name,omitempty"`
	TransferDate  time.Time `json:"transfer_date,omitempty"`

	// Virtual Account specific fields (kept for future use)
	VirtualAccountNumber string    `json:"virtual_account_number,omitempty"`
	BankCode             string    `json:"bank_code,omitempty"`
	VAExpiredAt          time.Time `json:"va_expired_at,omitempty"`

	// Payment Gateway specific fields (kept for future use)
	GatewayTransactionID string          `json:"gateway_transaction_id,omitempty"`
	GatewayName          string          `json:"gateway_name,omitempty"`
	GatewayResponse      json.RawMessage `json:"gateway_response,omitempty"`
	GatewayRedirectURL   string          `json:"gateway_redirect_url,omitempty"`
	GatewayCallbackURL   string          `json:"gateway_callback_url,omitempty"`
}

// Create model for initial payment method selection
type CreatePayment struct {
	RegistrationID int64         `json:"registration_id" validate:"required"`
	Amount         float64       `json:"amount" validate:"required,min=0"`
	PaymentMethod  PaymentMethod `json:"payment_method" validate:"required,eq=bank_transfer"` // Only allow bank transfer
}

// Update model for payment processing
type UpdatePayment struct {
	PaymentStatus PaymentStatus `json:"payment_status,omitempty"`
	ReceiptImage  string        `json:"receipt_image,omitempty"`
	Notes         string        `json:"notes,omitempty"`
	VerifiedBy    string        `json:"verified_by,omitempty"`

	// Bank Transfer updates
	BankName      string    `json:"bank_name,omitempty" validate:"required"`
	AccountNumber string    `json:"account_number,omitempty" validate:"required"`
	AccountName   string    `json:"account_name,omitempty" validate:"required"`
	TransferDate  time.Time `json:"transfer_date,omitempty" validate:"required"`
}
