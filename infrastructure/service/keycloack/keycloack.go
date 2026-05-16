package keycloack

import "context"

type KeycloackService interface {
	GetAccessToken(ctx context.Context, email, password, clientID, grantType, scope string) (*AccessTokenResponseDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AccessTokenResponseDTO, error)
	CreateUser(ctx context.Context, body RegisterUserDTO) (keycloakUUID string, err error)
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfoResponseDTO, error)
	UpdateUser(ctx context.Context, keycloakUUID string, body UpdateUserDTO) error
	DeleteUser(ctx context.Context, keycloakUUID string) error
}
