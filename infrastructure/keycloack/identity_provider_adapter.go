package keycloack

import (
	"auth-service/domain/service"
	"auth-service/domain/valueobject"

	"context"
)

type KeycloakAdapter struct {
	svc *KeycloakService
}

func NewKeycloakAdapter(svc *KeycloakService) service.AuthProvider {
	return &KeycloakAdapter{svc: svc}
}

func (a *KeycloakAdapter) CreateUser(ctx context.Context, param valueobject.RegisterUserParam) (string, error) {
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

func (a *KeycloakAdapter) GetAccessToken(ctx context.Context, email, password, clientID, grantType, scope string) (*valueobject.TokenResponse, error) {
	res, err := a.svc.GetAccessToken(ctx, email, password, clientID, grantType, scope)
	if err != nil {
		return nil, err
	}
	return &valueobject.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *KeycloakAdapter) RefreshToken(ctx context.Context, refreshToken string) (*valueobject.TokenResponse, error) {
	res, err := a.svc.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return &valueobject.TokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *KeycloakAdapter) GetUserInfo(ctx context.Context, accessToken string) (*valueobject.UserInfoResponse, error) {
	res, err := a.svc.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return &valueobject.UserInfoResponse{
		Sub:   res.Sub,
		Email: res.Email,
	}, nil
}

func (a *KeycloakAdapter) UpdateUser(ctx context.Context, keycloakUUID string, param valueobject.UpdateUserParam) error {
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
