package usecase

import (
	"auth-service/internal/application/dto"
	"auth-service/internal/domain/user"
	"context"
)

type UpdateUserUseCase struct {
	userService user.Service
}

func NewUpdateUserUseCase(userService user.Service) *UpdateUserUseCase {
	return &UpdateUserUseCase{userService: userService}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
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

	return dto.NewUserResponse(u), nil
}
