package models

import (
	"github.com/go-playground/validator/v10"
)

// newValidator creates a new validator instance with default settings
func newValidator() *validator.Validate {
	return validator.New()
}

// validateStruct validates a struct using the validator
func validateStruct(v *validator.Validate, s interface{}) error {
	return v.Struct(s)
}
