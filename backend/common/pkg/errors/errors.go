package errors

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrDuplicateRecord = errors.New("record already exists")
	ErrNotFound        = errors.New("record not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrUnauthorized    = errors.New("unauthorized")
)

type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(err error, message string, statusCode int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: statusCode,
	}
}

func IsDuplicateError(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}

func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
