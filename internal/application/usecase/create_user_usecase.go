package usecase

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/domain/user"
	"context"
)

type CreateUserUseCase struct {
	userService user.Service
}

func NewCreateUserUseCase(userService user.Service) *CreateUserUseCase {
	return &CreateUserUseCase{
		userService: userService,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	u, err := uc.userService.RegisterUser(ctx, req.Email, req.FirstName, req.LastName, req.Password, req.IsAdmin)
	if err != nil {
		return nil, err
	}

	return dto.NewUserResponse(u), nil
}
