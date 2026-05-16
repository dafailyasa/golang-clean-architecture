package valueobject

import (
	"auth-service/pkg/errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(
	`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
)

func NewEmail(email string) (*Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || !emailRegex.MatchString(email) {
		return nil, errors.NewBusinessError(
			"VALIDATION_",
			"INVALID_EMAIL",
			"email must be a valid email address",
		)
	}

	return &Email{value: email}, nil

}

func (e *Email) String() string {
	return e.value
}

type Status string

const (
	UserStatusActive    Status = "active"
	UserStatusInactive  Status = "inactive"
	UserStatusSuspended Status = "suspended"
)

func (s Status) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		return true
	}
	return false
}
