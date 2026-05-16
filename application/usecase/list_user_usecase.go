package usecase

import (
	"auth-service/domain/repository"

	"auth-service/application/dto"

	"context"
)

type ListUserUseCase struct {
	userRepository repository.Repository
}

func NewListUserUseCase(userRepository repository.Repository) *ListUserUseCase {
	return &ListUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *ListUserUseCase) Execute(ctx context.Context, req *dto.ListUserFilterRequest) ([]*dto.UserResponse, error) {
	users, err := uc.userRepository.FindAll(ctx, req.PaginationRequest)
	if err != nil {
		return nil, err
	}

	response := dto.NewUserResponseList(users)
	return response, nil
}
