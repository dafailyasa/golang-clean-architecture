package user

import (
	"auth-service/pkg/pagination"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	ExistsByEmailAndUsername(ctx context.Context, email, username string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByKeycloakUUID(ctx context.Context, keycloakUUID string) (*User, error)
	FindByID(ctx context.Context, id uint) (*User, error)
	Update(ctx context.Context, user *User) error
	ExistEmailByUserID(ctx context.Context, id uint, email string) (bool, error)
	DeleteByID(ctx context.Context, id uint) error
	FindAll(ctx context.Context, req *pagination.PaginationRequest) ([]*User, error)
}
