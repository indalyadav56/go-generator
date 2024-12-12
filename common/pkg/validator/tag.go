package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// validateUUID validates UUID strings
func validateUUID(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	_, err := uuid.Parse(field.String())
	return err == nil
}

// trimString trims whitespace from the beginning and end of a string.
func trimString(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	trimmed := strings.TrimSpace(field.String())
	field.SetString(trimmed)

	return true
}

// validateEmail checks if a string is a valid email address format.
func validateEmail(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(field.String())
}

// minLength checks if a string has at least the specified minimum length.
func minLength(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param() // Get the length parameter
	if field.Kind() != reflect.String {
		return false
	}

	min, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	return len(field.String()) >= min
}

// maxLength checks if a string does not exceed the specified maximum length.
func maxLength(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()
	if field.Kind() != reflect.String {
		return false
	}

	max, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	return len(field.String()) <= max
}

// isNumeric checks if a string contains only numeric characters.
func isNumeric(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	numericRegex := `^[0-9]+$`
	return regexp.MustCompile(numericRegex).MatchString(field.String())
}

// isAlpha checks if a string contains only alphabetic characters (a-z, A-Z).
func isAlpha(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	alphaRegex := `^[a-zA-Z]+$`
	return regexp.MustCompile(alphaRegex).MatchString(field.String())
}

// isAlphaNumeric checks if a string contains only alphanumeric characters.
func isAlphaNumeric(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	alphaNumericRegex := `^[a-zA-Z0-9]+$`
	return regexp.MustCompile(alphaNumericRegex).MatchString(field.String())
}
