package mysql

import (
	"auth-service/domain/aggregate"
	"auth-service/domain/repository"

	"auth-service/infrastructure/persistence/mysql/model"
	pkgErrors "auth-service/pkg/errors"
	"auth-service/pkg/pagination"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.Repository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, u *aggregate.User) error {
	entity := model.ToEntity(u)

	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).Model(&model.UserEntity{}).Create(&entity).Error
	if err != nil {
		return pkgErrors.NewInfrastructureError("DB_", "CREATE_FAILED", err.Error())
	}

	u.SetID(entity.ID)
	return nil
}

func (r *UserRepository) ExistsByEmailAndUsername(ctx context.Context, email, username string) (bool, error) {
	var count int64
	db := GetDB(ctx, r.db)

	if err := db.WithContext(ctx).
		Model(&model.UserEntity{}).
		Where("email = ? OR username = ?", email, username).
		Count(&count).Error; err != nil {
		return false, pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_EBE", err.Error())
	}
	return count > 0, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*aggregate.User, error) {
	var mdl model.UserEntity
	db := GetDB(ctx, r.db)

	if err := db.WithContext(ctx).Where("id = ?", id).First(&mdl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewBusinessError("NOT_FOUND_", "USER_BYID", "User not found")
		}
		return nil, pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_FBYID", err.Error())
	}

	return mdl.ToDomain()
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*aggregate.User, error) {
	var mdl model.UserEntity
	db := GetDB(ctx, r.db)

	if err := db.WithContext(ctx).Where("email = ?", email).First(&mdl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewBusinessError("NOT_FOUND_", "USER_FBE", "User not found")
		}
		return nil, pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_FBE", err.Error())
	}

	return mdl.ToDomain()
}

func (r *UserRepository) FindByKeycloakUUID(ctx context.Context, keycloakUUID string) (*aggregate.User, error) {
	var mdl model.UserEntity
	db := GetDB(ctx, r.db)

	if err := db.WithContext(ctx).Where("keycloak_uuid = ?", keycloakUUID).First(&mdl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkgErrors.NewBusinessError("NOT_FOUND_", "USER_FBKKUUID", "User not found")
		}
		return nil, pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_FBKKUUID", err.Error())
	}

	return mdl.ToDomain()
}

func (r *UserRepository) Update(ctx context.Context, u *aggregate.User) error {
	entity := model.ToEntity(u)
	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).Model(&model.UserEntity{}).
		Where("id = ?", u.ID).
		Updates(entity).Error
	if err != nil {
		return pkgErrors.NewInfrastructureError("DB_", "UPDATE_FAILED", "Failed to update entity")
	}

	return nil
}

func (r *UserRepository) ExistEmailByUserID(ctx context.Context, id uint, email string) (bool, error) {
	var count int64

	db := GetDB(ctx, r.db)

	err := db.WithContext(ctx).Model(&model.UserEntity{}).Where("id != ?", id).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_EBE", err.Error())
	}

	return count > 0, nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id uint) error {
	db := GetDB(ctx, r.db)
	err := db.WithContext(ctx).Model(model.UserEntity{}).Where("id = ?", id).Delete(&model.UserEntity{}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.NewBusinessError("NOT_FOUND_", "USER_DBYID", "User not found")
		}
		return pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_USER_DBYID", err.Error())
	}

	return nil
}

func (r *UserRepository) FindAll(ctx context.Context, req *pagination.PaginationRequest) ([]*aggregate.User, error) {
	var models []model.UserEntity

	db := GetDB(ctx, r.db)

	q := db.WithContext(ctx).Model(&model.UserEntity{})

	if req.GetSearch() != "" {
		searchTerm := "%" + req.GetSearch() + "%"
		q = q.Where("email LIKE ? ", searchTerm)
	}

	if err := q.Count(&req.Total).Error; err != nil {
		return nil, pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_COUNT", err.Error())
	}

	if err := q.Offset(req.GetOffset()).Limit(req.GetLimit()).Order(req.GetOrderClause()).Find(&models).Error; err != nil {
		return nil, pkgErrors.NewInfrastructureError("DB_", "QUERY_FAILED_FINDALL", err.Error())
	}

	var users []*aggregate.User
	for _, m := range models {
		domainUser, err := m.ToDomain()
		if err != nil {
			return nil, pkgErrors.NewInfrastructureError("DB_", "MAPPING_FAILED_FINDALL", err.Error())
		}
		users = append(users, domainUser)
	}

	return users, nil
}
