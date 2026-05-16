package usecase

import (
	"auth-service/domain/repository"

	"auth-service/application/dto"

	"context"
)

type GetUserUseCase struct {
	userRepository repository.Repository
}

func NewDetailUserUseCase(userRepository repository.Repository) *GetUserUseCase {
	return &GetUserUseCase{userRepository: userRepository}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, id uint) (*dto.UserResponse, error) {
	u, err := uc.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(u), nil
}
