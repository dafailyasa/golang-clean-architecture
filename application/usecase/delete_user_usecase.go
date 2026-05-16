package usecase

import (
	"auth-service/domain/service"

	"context"
)

type DeleteUserUseCase struct {
	userService service.Service
}

func NewDeleteUserUseCase(userService service.Service) *DeleteUserUseCase {
	return &DeleteUserUseCase{userService: userService}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id uint) error {
	return uc.userService.DeleteUser(ctx, id)
}
