package usecase

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/domain/user"
	"context"
)

type RefreshTokenUserUserUseCase struct {
	userService user.Service
}

func NewRefreshTokenUserUserUseCase(userService user.Service) *RefreshTokenUserUserUseCase {
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
