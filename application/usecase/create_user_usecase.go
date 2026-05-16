package usecase

import (
	"auth-service/domain/service"
	"auth-service/delivery/http/request"
	"auth-service/delivery/http/response"

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

func (uc *CreateUserUseCase) Execute(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error) {
	u, err := uc.userService.RegisterUser(ctx, req.Email, req.FirstName, req.LastName, req.Password, req.IsAdmin)
	if err != nil {
		return nil, err
	}

	return response.NewUserResponse(u), nil
}
