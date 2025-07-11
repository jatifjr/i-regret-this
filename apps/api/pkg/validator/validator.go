package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate

func init() {
	v = validator.New()
	RegisterCustomValidators(v)
}

func New() *validator.Validate {
	return v
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("notpastdate", validateNotPastDate)
}

func validateNotPastDate(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	now := time.Now().Truncate(24 * time.Hour)
	input := date.Truncate(24 * time.Hour)

	return !input.Before(now)
}
