package repository

import (
	"auth-service/domain/entity"

	"auth-service/pkg/pagination"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	ExistsByEmailAndUsername(ctx context.Context, email, username string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByKeycloakUUID(ctx context.Context, keycloakUUID string) (*entity.User, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	ExistEmailByUserID(ctx context.Context, id uint, email string) (bool, error)
	DeleteByID(ctx context.Context, id uint) error
	FindAll(ctx context.Context, req *pagination.PaginationRequest) ([]*entity.User, error)
}
