package test

import (
	"auth-service/domain/enum"

	"auth-service/domain/aggregate"
	"auth-service/domain/valueobject"

	"time"

	"github.com/google/uuid"
)

var (
	newEmail, _  = valueobject.NewEmail("testing@mail.com")
	UserTestData = aggregate.User{
		KeycloakUUID: uuid.NewString(),
		Email:        newEmail,
		IsAdmin:      nil,
		FirstName:    "Testing First Name",
		LastName:     "Testing Last Name",
		Username:     "Testing_Username",
		Status:       enum.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
)
