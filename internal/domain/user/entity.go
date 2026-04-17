package user

import (
	pkgErrors "auth-service/pkg/errors"
	"auth-service/pkg/helpers"
	"fmt"
	"time"
)

type User struct {
	ID           uint
	KeycloakUUID string
	IsAdmin      *bool
	Email        *Email
	FirstName    string
	LastName     string
	Username     string
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, firstName, lastName string, isAdmin *bool) (*User, error) {
	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	if err := validateName("first name", firstName); err != nil {
		return nil, err
	}

	if err := validateName("last name", lastName); err != nil {
		return nil, err
	}

	return &User{
		IsAdmin:   isAdmin,
		Email:     emailVO,
		FirstName: firstName,
		LastName:  lastName,
		Status:    UserStatusActive,
		Username:  helpers.GenerateUsername(firstName, lastName),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) Rebuild(firstName string, lastName string, email string, status Status, isAdmin *bool) error {

	if !status.IsValid() {
		return pkgErrors.NewBusinessError(
			"VALIDATION_",
			"INVALID_STATUS",
			"Invalid user status",
		)
	}

	err := u.ChangeEmail(email)
	if err != nil {
		return err
	}

	if err := validateName("first name", firstName); err != nil {
		return err
	}

	if err := validateName("last name", lastName); err != nil {
		return err
	}

	u.SetFirstName(firstName)
	u.SetLastName(lastName)
	u.SetStatus(status)
	u.SetIsAdmin(isAdmin)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) SetFirstName(firstName string) {
	u.FirstName = firstName
}

func (u *User) SetKeycloakUUID(keycloakUUID string) {
	u.KeycloakUUID = keycloakUUID
}

func (u *User) SetLastName(lastName string) {
	u.LastName = lastName
}

func (u *User) SetStatus(status Status) {
	u.Status = status
}

func (u *User) ChangeEmail(newEmail string) error {
	emailVO, err := NewEmail(newEmail)
	if err != nil {
		return err
	}

	u.Email = emailVO
	return nil
}

func (u *User) SetIsAdmin(isAdmin *bool) {
	u.IsAdmin = isAdmin
}

func (u *User) IsInActiveAndSuspend() bool {
	return u.Status != UserStatusActive
}

func (u *User) GetStatus() Status {
	return u.Status
}

func (u *User) SetID(id uint) {
	u.ID = id
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) GetKeycloakUUID() string {
	return u.KeycloakUUID
}

func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

func validateName(fieldName string, value string) error {
	if len(value) < 2 {
		return pkgErrors.NewBusinessError(
			"VALIDATION_",
			"INVALID_NAME",
			fmt.Sprintf("%s must be at least 2 characters", fieldName),
		)
	}

	return nil
}
