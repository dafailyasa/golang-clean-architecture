package port

import "context"

type RegisterUserParam struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UpdateUserParam struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Enabled   bool
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}

type UserInfoResponse struct {
	Sub   string
	Email string
}

type IdentityProvider interface {
	CreateUser(ctx context.Context, param RegisterUserParam) (string, error)
	GetAccessToken(ctx context.Context, email, password, clientID, grantType, scope string) (*TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	GetUserInfo(ctx context.Context, accessToken string) (*UserInfoResponse, error)
	UpdateUser(ctx context.Context, keycloakUUID string, param UpdateUserParam) error
	DeleteUser(ctx context.Context, keycloakUUID string) error
}
