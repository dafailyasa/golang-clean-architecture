package usecase

import (
	"auth-service/domain/repository"
	"auth-service/delivery/http/request"
	"auth-service/delivery/http/response"

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

func (uc *ListUserUseCase) Execute(ctx context.Context, req *request.ListUserFilterRequest) ([]*response.UserResponse, error) {
	users, err := uc.userRepository.FindAll(ctx, req.PaginationRequest)
	if err != nil {
		return nil, err
	}

	res := response.NewUserResponseList(users)
	return res, nil
}
