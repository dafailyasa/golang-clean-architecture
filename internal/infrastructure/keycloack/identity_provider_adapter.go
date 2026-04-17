package keycloack

import (
	"auth-service/internal/application/port"
	"context"
)

type KeycloakAdapter struct {
	svc *KeycloakService
}

func NewKeycloakAdapter(svc *KeycloakService) port.IdentityProvider {
	return &KeycloakAdapter{svc: svc}
}

func (a *KeycloakAdapter) CreateUser(ctx context.Context, param port.RegisterUserParam) (string, error) {
	cmd := RegisterUserDTO{
		Username:      param.Username,
		Enabled:       true,
		EmailVerified: true,
		FirstName:     param.FirstName,
		LastName:      param.LastName,
		Email:         param.Email,
		Credentials: []CredentialsDTO{
			{
				Type:      "password",
				Value:     param.Password,
				Temporary: false,
			},
		},
	}
	return a.svc.CreateUser(ctx, cmd)
}

func (a *KeycloakAdapter) GetAccessToken(ctx context.Context, email, password, clientID, grantType, scope string) (*port.TokenResponse, error) {
	res, err := a.svc.GetAccessToken(ctx, email, password, clientID, grantType, scope)
	if err != nil {
		return nil, err
	}
	return &port.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *KeycloakAdapter) RefreshToken(ctx context.Context, refreshToken string) (*port.TokenResponse, error) {
	res, err := a.svc.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return &port.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *KeycloakAdapter) GetUserInfo(ctx context.Context, accessToken string) (*port.UserInfoResponse, error) {
	res, err := a.svc.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return &port.UserInfoResponse{
		Sub:   res.Sub,
		Email: res.Email,
	}, nil
}

func (a *KeycloakAdapter) UpdateUser(ctx context.Context, keycloakUUID string, param port.UpdateUserParam) error {
	cmd := UpdateUserDTO{
		FirstName:     param.FirstName,
		LastName:      param.LastName,
		Email:         param.Email,
		Enabled:       param.Enabled,
		EmailVerified: param.Enabled,
	}

	if param.Password != "" {
		cmd.Credentials = []CredentialsDTO{
			{
				Type:      "password",
				Value:     param.Password,
				Temporary: false,
			},
		}
	}

	return a.svc.UpdateUser(ctx, keycloakUUID, cmd)
}

func (a *KeycloakAdapter) DeleteUser(ctx context.Context, keycloakUUID string) error {
	return a.svc.DeleteUser(ctx, keycloakUUID)
}
