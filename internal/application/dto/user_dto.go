package dto

import (
	"auth-service/internal/domain/user"
	"time"
)

type UserResponse struct {
	ID           uint        `json:"id"`
	KeycloakUUID string      `json:"keycloakUUID"`
	IsAdmin      bool        `json:"isAdmin"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Email        string      `json:"email"`
	Status       user.Status `json:"status"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

func NewUserResponse(u *user.User) *UserResponse {
	return &UserResponse{
		ID:           u.ID,
		KeycloakUUID: u.KeycloakUUID,
		IsAdmin:      *u.IsAdmin,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email.String(),
		Status:       u.Status,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func NewUserResponseList(users []*user.User) []*UserResponse {
	var response []*UserResponse
	for _, u := range users {
		response = append(response, NewUserResponse(u))
	}
	return response
}
