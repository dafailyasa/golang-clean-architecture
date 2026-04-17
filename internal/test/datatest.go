package test

import (
	"auth-service/internal/domain/user"
	"time"

	"github.com/google/uuid"
)

var (
	newEmail, _  = user.NewEmail("testing@mail.com")
	UserTestData = user.User{
		KeycloakUUID: uuid.NewString(),
		Email:        newEmail,
		IsAdmin:      nil,
		FirstName:    "Testing First Name",
		LastName:     "Testing Last Name",
		Username:     "Testing_Username",
		Status:       user.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
)
