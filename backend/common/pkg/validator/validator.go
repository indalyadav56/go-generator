package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateStruct(s interface{}) error
	RegisterCustomValidation(tag string, fn validator.Func) error
}

type customValidator struct {
	validator *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func New() *customValidator {
	v := validator.New()
	customValidator := &customValidator{
		validator: v,
	}
	customValidator.RegisterCustomValidation("uuid", validateUUID)
	customValidator.RegisterCustomValidation("trim", trimString)
	customValidator.RegisterCustomValidation("email", validateEmail)
	customValidator.RegisterCustomValidation("min", minLength)
	customValidator.RegisterCustomValidation("max", maxLength)
	customValidator.RegisterCustomValidation("numeric", isNumeric)
	customValidator.RegisterCustomValidation("alpha", isAlpha)
	customValidator.RegisterCustomValidation("alphanumeric", isAlphaNumeric)

	return customValidator
}

func (v *customValidator) ValidateStruct(s interface{}) error {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := getJSONFieldName(s, err.StructField())
		validationErrors = append(validationErrors, ValidationError{
			Field:   fieldName,
			Message: getErrorMessage(err),
		})
	}
	return validationErrors
}

func (v *customValidator) RegisterCustomValidation(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}

func (ve ValidationErrors) Error() string {
	var errMsgs []string
	for _, err := range ve {
		errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(errMsgs, "; ")
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "phone":
		return "Invalid phone number format"
	case "password":
		return "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number and one special character"
	case "uuid":
		return "Invalid UUID format"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Must not be longer than %s characters", err.Param())
	default:
		return fmt.Sprintf("Invalid value for %s", err.Field())
	}
}

// getJSONFieldName retrieves the JSON tag for a struct field, falling back to the field name if no JSON tag is found
func getJSONFieldName(s interface{}, fieldName string) string {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field, found := val.Type().FieldByName(fieldName)
	if !found {
		return fieldName
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return fieldName
	}

	return strings.Split(jsonTag, ",")[0]
}
