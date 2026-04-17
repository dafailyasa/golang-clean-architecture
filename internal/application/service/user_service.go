package service

import (
	"auth-service/config"
	"auth-service/internal/application/port"
	"auth-service/internal/domain/user"
	"auth-service/pkg/constant"
	pkgErrors "auth-service/pkg/errors"
	"context"
	"fmt"
)

type userService struct {
	userRepo         user.Repository
	txManager        port.TransactionManager
	cfg              config.Config
	identityProvider port.IdentityProvider
}

func NewUserService(userRepo user.Repository, txManager port.TransactionManager, jwtConfig config.Config, identityProvider port.IdentityProvider) user.Service {
	return &userService{
		userRepo:         userRepo,
		txManager:        txManager,
		cfg:              jwtConfig,
		identityProvider: identityProvider,
	}
}

func (s *userService) createUserKeycloak(ctx context.Context, user *user.User, password string) (string, error) {
	data := port.RegisterUserParam{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email.String(),
		Password:  password,
	}

	keycloakUUID, err := s.identityProvider.CreateUser(ctx, data)
	if err != nil {
		return "", err
	}

	return keycloakUUID, nil
}

func (s *userService) RegisterUser(ctx context.Context, email, firstName, lastName, password string, isAdmin *bool) (*user.User, error) {
	var createdUser *user.User
	newUser, err := user.NewUser(email, firstName, lastName, isAdmin)
	if err != nil {
		return nil, err
	}

	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.userRepo.ExistsByEmailAndUsername(txCtx, newUser.Email.String(), newUser.Username)
		if err != nil {
			return err
		}

		if exists {
			return pkgErrors.NewBusinessError("SRUS", "001", "Email or Username already registered")
		}

		keycloakUUID, err := s.createUserKeycloak(ctx, newUser, password)
		if err != nil {
			return err
		}
		newUser.SetKeycloakUUID(keycloakUUID)

		err = s.userRepo.Create(txCtx, newUser)
		if err != nil {
			return err
		}

		createdUser = newUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *userService) ValidateAuthenticateUser(ctx context.Context, email string) (*user.User, error) {
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err.(*pkgErrors.ErrorCustomize).IsBusinessError() {
			return nil, pkgErrors.NewBusinessError("SAUTU", "001", "Email was not registered")
		}
		return nil, err
	}

	if u.IsInActiveAndSuspend() {
		return nil, pkgErrors.NewBusinessError("SAUTU", "002", fmt.Sprintf("Your account is %s. Please contact admin", u.GetStatus()))
	}

	return u, nil
}

func (s *userService) GenerateToken(ctx context.Context, email, password, grantType string) (accessToken *string, refreshToken *string, err error) {
	token, err := s.identityProvider.GetAccessToken(ctx, email, password, s.cfg.Keycloak.ClientID, grantType, constant.KeycloakScope)
	if err != nil {
		return nil, nil, err
	}

	return &token.AccessToken, &token.RefreshToken, nil
}

func (s *userService) RefreshToken(ctx context.Context, refreshToken string) (*user.User, *string, *string, error) {
	tokenRes, err := s.identityProvider.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, nil, err
	}

	userInfo, err := s.identityProvider.GetUserInfo(ctx, tokenRes.AccessToken)
	if err != nil {
		return nil, nil, nil, err
	}

	u, err := s.userRepo.FindByKeycloakUUID(ctx, userInfo.Sub)
	if err != nil {
		return nil, nil, nil, err
	}

	return u, &tokenRes.AccessToken, &tokenRes.RefreshToken, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, email, firstsName, lastName, status, password string, isAdmin *bool) (*user.User, error) {
	var updatedUser *user.User

	err := s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		u, err := s.userRepo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		existEmail, err := s.userRepo.ExistEmailByUserID(txCtx, id, email)
		if err != nil {
			return pkgErrors.NewInfrastructureError("SUPU", "001", err.Error())
		}

		if existEmail {
			return pkgErrors.NewBusinessError("SUPU", "002", "Email is already registered")
		}

		if err := u.Rebuild(firstsName, lastName, email, user.Status(status), isAdmin); err != nil {
			return err
		}

		if err := s.userRepo.Update(txCtx, u); err != nil {
			return err
		}

		updatePayload := port.UpdateUserParam{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email.String(),
			Password:  password,
			Enabled:   u.Status == user.UserStatusActive,
		}

		if err := s.identityProvider.UpdateUser(ctx, u.KeycloakUUID, updatePayload); err != nil {
			return err
		}

		updatedUser = u
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		u, err := s.userRepo.FindByID(txCtx, id)
		if err != nil {
			return err
		}

		if err := s.identityProvider.DeleteUser(ctx, u.KeycloakUUID); err != nil {
			return err
		}

		if err := s.userRepo.DeleteByID(txCtx, id); err != nil {
			return err
		}

		return nil
	})
}
