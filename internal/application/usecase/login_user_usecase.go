package usecase

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/domain/user"
	"auth-service/pkg/constant"
	"context"
)

type LoginUserUseCase struct {
	userService    user.Service
	userRepository user.Repository
}

func NewLoginUserUseCase(userService user.Service, userRepository user.Repository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userService:    userService,
		userRepository: userRepository,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, body *dto.LoginUserRequest) (*dto.LoginUserResponse, error) {
	u, err := uc.userService.ValidateAuthenticateUser(ctx, body.Email)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := uc.userService.GenerateToken(ctx, body.Email, body.Password, constant.KeycloakGrantTypePasswordConst)
	if err != nil {
		return nil, err
	}

	return dto.NewLoginUserResponse(*accessToken, *refreshToken, u.GetFullName()), nil
}
