package repository

import (
	"auth-service/domain/aggregate"

	"auth-service/pkg/pagination"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *aggregate.User) error
	ExistsByEmailAndUsername(ctx context.Context, email, username string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
	FindByKeycloakUUID(ctx context.Context, keycloakUUID string) (*aggregate.User, error)
	FindByID(ctx context.Context, id uint) (*aggregate.User, error)
	Update(ctx context.Context, user *aggregate.User) error
	ExistEmailByUserID(ctx context.Context, id uint, email string) (bool, error)
	DeleteByID(ctx context.Context, id uint) error
	FindAll(ctx context.Context, req *pagination.PaginationRequest) ([]*aggregate.User, error)
}
