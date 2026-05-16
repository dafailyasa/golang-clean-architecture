package service

import (
	"auth-service/domain/valueobject"
	"context"
)

type AuthProvider interface {
	CreateUser(ctx context.Context, param valueobject.RegisterUserParam) (string, error)
	GetAccessToken(ctx context.Context, email, password, clientID, grantType, scope string) (*valueobject.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*valueobject.TokenResponse, error)
	GetUserInfo(ctx context.Context, accessToken string) (*valueobject.UserInfoResponse, error)
	UpdateUser(ctx context.Context, keycloakUUID string, param valueobject.UpdateUserParam) error
	DeleteUser(ctx context.Context, keycloakUUID string) error
}
