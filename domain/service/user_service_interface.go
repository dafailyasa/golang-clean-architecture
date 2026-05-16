package service

import (
	"auth-service/domain/aggregate"
	"context"
)

type Service interface {
	RegisterUser(ctx context.Context, email, firstName, lastName, password string, isAdmin *bool) (*aggregate.User, error)
	ValidateAuthenticateUser(ctx context.Context, email string) (*aggregate.User, error)
	GenerateToken(ctx context.Context, email, password, grantType string) (accessToken *string, refreshToken *string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (user *aggregate.User, accessToken *string, refreshTokenRes *string, err error)
	UpdateUser(ctx context.Context, id uint, email, firstsName, lastName, status, password string, isAdmin *bool) (*aggregate.User, error)
	DeleteUser(ctx context.Context, id uint) error
}
