package usecase

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/domain/user"
	"context"
)

type ListUserUseCase struct {
	userRepository user.Repository
}

func NewListUserUseCase(userRepository user.Repository) *ListUserUseCase {
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
