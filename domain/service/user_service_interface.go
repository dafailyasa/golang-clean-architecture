package service

import (
	"auth-service/domain/entity"
	"context"
)

type Service interface {
	RegisterUser(ctx context.Context, email, firstName, lastName, password string, isAdmin *bool) (*entity.User, error)
	ValidateAuthenticateUser(ctx context.Context, email string) (*entity.User, error)
	GenerateToken(ctx context.Context, email, password, grantType string) (accessToken *string, refreshToken *string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (user *entity.User, accessToken *string, refreshTokenRes *string, err error)
	UpdateUser(ctx context.Context, id uint, email, firstsName, lastName, status, password string, isAdmin *bool) (*entity.User, error)
	DeleteUser(ctx context.Context, id uint) error
}
