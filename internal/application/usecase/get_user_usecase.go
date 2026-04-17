package usecase

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/domain/user"
	"context"
)

type GetUserUseCase struct {
	userRepository user.Repository
}

func NewDetailUserUseCase(userRepository user.Repository) *GetUserUseCase {
	return &GetUserUseCase{userRepository: userRepository}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id uint) (*dto.UserResponse, error) {
	u, err := uc.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(u), nil
}
