package keycloack

import (
	"auth-service/config"
	pkgHttp "auth-service/pkg/app_http"
	"auth-service/pkg/constant"
	pkgErrors "auth-service/pkg/errors"
	"auth-service/pkg/helpers"
	"context"
	"fmt"
	"net/http"
)

type keycloakService struct {
	HttpClient  *pkgHttp.AppHttp
	KeycloakCfg config.KeycloakConfig
}

func NewKeycloakService(httpClient *pkgHttp.AppHttp, keycloakCfg config.KeycloakConfig) KeycloackService {
	return keycloakService{
		HttpClient:  httpClient,
		KeycloakCfg: keycloakCfg,
	}
}

func (s keycloakService) GetAccessToken(ctx context.Context, email, password, clientID, grantType, scope string) (*AccessTokenResponseDTO, error) {
	body := AccessTokenDTO{
		Username:     email,
		Password:     password,
		GrantType:    grantType,
		ClientID:     clientID,
		ClientSecret: s.KeycloakCfg.ClientSecret,
	}

	if scope != "" {
		body.Scope = scope
	}

	formData, err := helpers.StructToFormData(body)
	if err != nil {
		return nil, pkgErrors.NewTechnicalError("KCGT", "001", "invalid payload request")
	}

	var res AccessTokenResponseDTO
	httpResp, err := s.HttpClient.DoHttpRequest(ctx, pkgHttp.Request{
		Method: http.MethodPost,
		Endpoint: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
			s.KeycloakCfg.KeycloakBaseURL,
			s.KeycloakCfg.Realm,
		),
		Headers: map[string]string{
			"Content-Type": constant.FormUrlEncodedConst,
		},
		ToURLEncoded: true,
		FormData:     formData,
	}, &res)

	if err != nil {
		return nil, pkgErrors.NewInfrastructureError("KCGT", "002", err.Error())
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, pkgErrors.NewBusinessError("KCGT", "003", res.ErrorDescription)
	}

	return &res, nil
}

func (s keycloakService) RefreshToken(ctx context.Context, refreshToken string) (*AccessTokenResponseDTO, error) {
	body := RefreshTokenDTO{
		ClientID:     s.KeycloakCfg.ClientID,
		ClientSecret: s.KeycloakCfg.ClientSecret,
		RefreshToken: refreshToken,
		GrantType:    constant.KeycloakGrantTypeRefreshTokenConst,
	}

	formData, err := helpers.StructToFormData(body)
	if err != nil {
		return nil, pkgErrors.NewTechnicalError("KCRT", "001", "invalid payload request")
	}

	var res AccessTokenResponseDTO
	httpResp, err := s.HttpClient.DoHttpRequest(ctx, pkgHttp.Request{
		Method: http.MethodPost,
		Endpoint: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
			s.KeycloakCfg.KeycloakBaseURL,
			s.KeycloakCfg.Realm,
		),
		Headers: map[string]string{
			"Content-Type": constant.FormUrlEncodedConst,
		},
		ToURLEncoded: true,
		FormData:     formData,
	}, &res)

	if err != nil {
		return nil, pkgErrors.NewInfrastructureError("KCRT", "002", err.Error())
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, pkgErrors.NewBusinessError("KCRT", "003", res.ErrorDescription)
	}

	return &res, nil
}

func (s keycloakService) CreateUser(ctx context.Context, body RegisterUserDTO) (keycloakUUID string, err error) {
	tokenRes, err := s.GetAccessToken(ctx,
		s.KeycloakCfg.AdminKeycloak.Email,
		s.KeycloakCfg.AdminKeycloak.Password,
		constant.KeycloakAdminCLIConst,
		constant.KeycloakGrantTypePasswordConst,
		"",
	)
	if err != nil {
		return "", err
	}

	var res RegisterUserResponseDTO

	httpResp, err := s.HttpClient.DoHttpRequest(ctx, pkgHttp.Request{
		Method: http.MethodPost,
		Endpoint: fmt.Sprintf("%s/admin/realms/%s/users",
			s.KeycloakCfg.KeycloakBaseURL,
			s.KeycloakCfg.Realm,
		),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", tokenRes.AccessToken),
		},
		Body: body,
	}, &res)

	if err != nil || httpResp.StatusCode == http.StatusInternalServerError {
		return "", pkgErrors.NewInfrastructureError("KCCU", "001", err.Error())
	}

	if httpResp.StatusCode == http.StatusConflict {
		return "", pkgErrors.NewBusinessError("KCCU", "002", res.ErrorMessage)
	}

	keycloakUUID = helpers.ExtractUUIDLocation(httpResp.Headers["Location"])
	if keycloakUUID == "" {
		return "", pkgErrors.NewInfrastructureError("KCCU", "003", "failed to extract keycloak uuid")
	}

	return keycloakUUID, nil
}

func (s keycloakService) GetUserInfo(ctx context.Context, accessToken string) (*UserInfoResponseDTO, error) {
	var res UserInfoResponseDTO
	httpResp, err := s.HttpClient.DoHttpRequest(ctx, pkgHttp.Request{
		Method: http.MethodGet,
		Endpoint: fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo",
			s.KeycloakCfg.KeycloakBaseURL,
			s.KeycloakCfg.Realm,
		),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
			"Content-Type":  constant.ApplicationJsonConst,
		},
	}, &res)
	if err != nil {
		return nil, pkgErrors.NewInfrastructureError("KCGU", "001", err.Error())
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, pkgErrors.NewBusinessError("KCGU", "002", "failed to fetch user info from keycloak")
	}

	return &res, nil
}

func (s keycloakService) UpdateUser(ctx context.Context, keycloakUUID string, body UpdateUserDTO) error {
	tokenRes, err := s.GetAccessToken(ctx,
		s.KeycloakCfg.AdminKeycloak.Email,
		s.KeycloakCfg.AdminKeycloak.Password,
		constant.KeycloakAdminCLIConst,
		constant.KeycloakGrantTypePasswordConst,
		"",
	)
	if err != nil {
		return err
	}

	httpResp, err := s.HttpClient.DoHttpRequest(ctx, pkgHttp.Request{
		Method: http.MethodPut,
		Endpoint: fmt.Sprintf("%s/admin/realms/%s/users/%s",
			s.KeycloakCfg.KeycloakBaseURL,
			s.KeycloakCfg.Realm,
			keycloakUUID,
		),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", tokenRes.AccessToken),
			"Content-Type":  constant.ApplicationJsonConst,
		},
		Body: body,
	}, &struct{}{})

	if err != nil {
		return pkgErrors.NewInfrastructureError("KCUU", "001", err.Error())
	}

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		return pkgErrors.NewBusinessError("KCUU", "002", "failed to update user in keycloak")
	}

	return nil
}

func (s keycloakService) DeleteUser(ctx context.Context, keycloakUUID string) error {
	tokenRes, err := s.GetAccessToken(ctx,
		s.KeycloakCfg.AdminKeycloak.Email,
		s.KeycloakCfg.AdminKeycloak.Password,
		constant.KeycloakAdminCLIConst,
		constant.KeycloakGrantTypePasswordConst,
		"",
	)
	if err != nil {
		return err
	}

	httpResp, err := s.HttpClient.DoHttpRequest(ctx, pkgHttp.Request{
		Method: http.MethodDelete,
		Endpoint: fmt.Sprintf("%s/admin/realms/%s/users/%s",
			s.KeycloakCfg.KeycloakBaseURL,
			s.KeycloakCfg.Realm,
			keycloakUUID,
		),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", tokenRes.AccessToken),
		},
	}, &struct{}{})

	if err != nil {
		return pkgErrors.NewInfrastructureError("KCDU", "001", err.Error())
	}

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		return pkgErrors.NewBusinessError("KCDU", "002", "failed to delete user in keycloak")
	}

	return nil
}
