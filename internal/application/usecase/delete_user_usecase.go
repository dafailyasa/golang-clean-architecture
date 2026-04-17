package usecase

import (
	"auth-service/internal/domain/user"
	"context"
)

type DeleteUserUseCase struct {
	userService user.Service
}

func NewDeleteUserUseCase(userService user.Service) *DeleteUserUseCase {
	return &DeleteUserUseCase{userService: userService}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id uint) error {
	return uc.userService.DeleteUser(ctx, id)
}
