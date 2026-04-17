package model

import (
	"auth-service/internal/domain/user"
	"time"
)

type UserEntity struct {
	ID           uint        `gorm:"column:id;primaryKey;autoIncrement"`
	KeycloakUUID string      `gorm:"column:keycloak_uuid;type:char(36);uniqueIndex;not null"`
	Email        string      `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	FirstName    string      `gorm:"column:first_name;type:varchar(255);not null"`
	LastName     string      `gorm:"column:last_name;type:varchar(255);not null"`
	Username     string      `gorm:"column:username;type:varchar(255);uniqueIndex;not null"`
	Status       user.Status `gorm:"column:status;type:varchar(20);default:'active'"`
	IsAdmin      *bool       `gorm:"column:is_admin;type:bool;not null;default:0"`
	CreatedAt    time.Time   `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time   `gorm:"column:updated_at;not null"`
}

func (UserEntity) TableName() string {
	return "users"
}

func ToEntity(u *user.User) *UserEntity {
	return &UserEntity{
		ID:           u.ID,
		KeycloakUUID: u.KeycloakUUID,
		Email:        u.Email.String(),
		IsAdmin:      u.IsAdmin,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Status:       u.Status,
		Username:     u.Username,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (r UserEntity) ToDomain() (*user.User, error) {
	email, err := user.NewEmail(r.Email)
	if err != nil {
		return nil, err
	}

	return &user.User{
		ID:           r.ID,
		KeycloakUUID: r.KeycloakUUID,
		IsAdmin:      r.IsAdmin,
		FirstName:    r.FirstName,
		LastName:     r.LastName,
		Email:        email,
		Username:     r.Username,
		Status:       user.Status(r.Status),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}, nil
}
