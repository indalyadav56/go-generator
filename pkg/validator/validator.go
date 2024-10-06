package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

// Validator struct that wraps the validator.Validate object
type Validator struct {
	validate *validator.Validate
}

var (
	once     sync.Once
	instance *Validator
)

// GetValidator returns a singleton instance of the Validator
func GetValidator() *Validator {
	once.Do(func() {
		instance = &Validator{validate: validator.New()}
	})
	return instance
}

// ValidateStruct validates a given struct and returns validation errors if any
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// RegisterCustomValidation allows registration of custom validators
func (v *Validator) RegisterCustomValidation(tag string, fn validator.Func) error {
	return v.validate.RegisterValidation(tag, fn)
}
