package errors

import (
	"errors"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:-`
	Code       string `json:code`
	Message    string `json:message`
}

func (e AppError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) error {
	return AppError{
		StatusCode: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

func NewNotFoundError(message string) error {
	return AppError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}

func NewInternalServerError(message string) error {
	return AppError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    message,
	}
}

func GetStatusCode(err error) int {
	var appError AppError
	if errors.As(err, &appError) {
		return appError.StatusCode
	}

	return http.StatusInternalServerError
}
