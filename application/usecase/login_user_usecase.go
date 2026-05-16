package usecase

import (
	"auth-service/domain/repository"
	"auth-service/domain/service"
	"auth-service/delivery/http/request"
	"auth-service/delivery/http/response"

	"auth-service/pkg/constant"
	"context"
)

type LoginUserUseCase struct {
	userService    service.Service
	userRepository repository.Repository
}

func NewLoginUserUseCase(userService service.Service, userRepository repository.Repository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userService:    userService,
		userRepository: userRepository,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, body *request.LoginUserRequest) (*response.LoginUserResponse, error) {
	u, err := uc.userService.ValidateAuthenticateUser(ctx, body.Email)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := uc.userService.GenerateToken(ctx, body.Email, body.Password, constant.KeycloakGrantTypePasswordConst)
	if err != nil {
		return nil, err
	}

	return response.NewLoginUserResponse(*accessToken, *refreshToken, u.GetFullName()), nil
}
