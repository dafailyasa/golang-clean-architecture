package user

import (
	"context"
)

type Service interface {
	RegisterUser(ctx context.Context, email, firstName, lastName, password string, isAdmin *bool) (*User, error)
	ValidateAuthenticateUser(ctx context.Context, email string) (*User, error)
	GenerateToken(ctx context.Context, email, password, grantType string) (accessToken *string, refreshToken *string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (user *User, accessToken *string, refreshTokenRes *string, err error)
	UpdateUser(ctx context.Context, id uint, email, firstsName, lastName, status, password string, isAdmin *bool) (*User, error)
	DeleteUser(ctx context.Context, id uint) error
}
