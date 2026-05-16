package usecase

import (
	"auth-service/domain/service"

	"auth-service/application/dto"

	"context"
)

type RefreshTokenUserUserUseCase struct {
	userService service.Service
}

func NewRefreshTokenUserUserUseCase(userService service.Service) *RefreshTokenUserUserUseCase {
	return &RefreshTokenUserUserUseCase{
		userService: userService,
	}
}

func (uc *RefreshTokenUserUserUseCase) Execute(ctx context.Context, body *dto.RefreshTokenUserDTO) (*dto.LoginUserResponse, error) {
	u, accessToken, refreshToken, err := uc.userService.RefreshToken(ctx, body.RefreshToken)
	if err != nil {
		return nil, err
	}

	return dto.NewLoginUserResponse(*accessToken, *refreshToken, u.GetFullName()), nil
}
