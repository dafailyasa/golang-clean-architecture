package exception

import (
	"errors"
	"net/http"
)

var (
	ErrExistEmailOrUsername = errors.New("Email or Username already registered")
	ErrEmailRegistered      = errors.New("Email is already registered")
)

var HTTPStatusCodes = map[error]int{
	ErrExistEmailOrUsername: http.StatusConflict,
	ErrEmailRegistered:      http.StatusConflict,
}

func GetHTTPStatus(err error) *int {
	if status, exists := HTTPStatusCodes[err]; exists {
		return &status
	}

	return nil
}
