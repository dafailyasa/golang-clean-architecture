package usecase

import (
	"auth-service/domain/service"
	"auth-service/delivery/http/request"
	"auth-service/delivery/http/response"

	"context"
)

type UpdateUserUseCase struct {
	userService service.Service
}

func NewUpdateUserUseCase(userService service.Service) *UpdateUserUseCase {
	return &UpdateUserUseCase{userService: userService}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, id uint, req *request.UpdateUserRequest) (*response.UserResponse, error) {
	u, err := uc.userService.UpdateUser(
		ctx,
		id,
		req.Email,
		req.FirstName,
		req.LastName,
		req.Status,
		req.Password,
		req.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	return response.NewUserResponse(u), nil
}
