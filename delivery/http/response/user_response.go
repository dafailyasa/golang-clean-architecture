package response

import (
	"auth-service/domain/enum"
	"auth-service/domain/aggregate"
	"time"
)

type UserResponse struct {
	ID           uint        `json:"id"`
	KeycloakUUID string      `json:"keycloakUUID"`
	IsAdmin      bool        `json:"isAdmin"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Email        string      `json:"email"`
	Status       enum.Status `json:"status"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

func NewUserResponse(u *aggregate.User) *UserResponse {
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

func NewUserResponseList(users []*aggregate.User) []*UserResponse {
	var response []*UserResponse
	for _, u := range users {
		response = append(response, NewUserResponse(u))
	}
	return response
}
