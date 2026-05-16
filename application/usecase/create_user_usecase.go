package usecase

import (
	"auth-service/domain/service"

	"auth-service/application/dto"

	"context"
)

type CreateUserUseCase struct {
	userService service.Service
}

func NewCreateUserUseCase(userService service.Service) *CreateUserUseCase {
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
