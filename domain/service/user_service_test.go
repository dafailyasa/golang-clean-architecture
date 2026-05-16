package service

/*
import (
	"auth-service/domain/entity"
	"auth-service/domain/valueobject"
	"auth-service/domain/repository"
	"auth-service/domain/service"

	"auth-service/config"
	"auth-service/application/port"

	"auth-service/pkg/constant"
	pkgErrors "auth-service/pkg/errors"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx  context.Context
	repo *MockUserRepository
	tx   *MockTransactionManager
	ip   *MockIdentityProvider
	svc  service.Service
}

func (t *UserServiceTestSuite) SetupTest() {
	t.ctx = context.Background()
	t.repo = new(MockUserRepository)
	t.tx = new(MockTransactionManager)
	t.ip = new(MockIdentityProvider)
	t.svc = NewUserService(t.repo, t.tx, config.Config{}, t.ip)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// ─── RegisterUser ────────────────────────────────────────────────────────────

func (t *UserServiceTestSuite) TestRegisterUser() {
	t.Run("it should register user successfully", func() {
		isAdmin := false
		t.tx.On("WithTransaction", t.ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		t.repo.On("ExistsByEmailAndUsername", t.ctx, "test@example.com", mock.AnythingOfType("string")).Return(false, nil)
		t.ip.On("CreateUser", t.ctx, mock.MatchedBy(func(p port.RegisterUserParam) bool {
			return p.Email == "test@example.com" && p.FirstName == "John"
		})).Return("uuid-1234", nil)
		t.repo.On("Create", t.ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		res, err := t.svc.RegisterUser(t.ctx, "test@example.com", "John", "Doe", "pass@1234A", &isAdmin)

		t.NoError(err)
		t.NotNil(res)
		t.Equal("uuid-1234", res.KeycloakUUID)
		t.Equal("test@example.com", res.Email.String())
		t.repo.AssertExpectations(t.T())
		t.tx.AssertExpectations(t.T())
		t.ip.AssertExpectations(t.T())
	})

	t.Run("it should return error when email or username already registered", func() {
		isAdmin := false
		t.tx.On("WithTransaction", t.ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(pkgErrors.NewBusinessError("SRUS", "001", "Email or Username already registered"))
		t.repo.On("ExistsByEmailAndUsername", t.ctx, "taken@example.com", mock.AnythingOfType("string")).Return(true, nil)

		res, err := t.svc.RegisterUser(t.ctx, "taken@example.com", "John", "Doe", "pass@1234A", &isAdmin)

		t.Error(err)
		t.Nil(res)
		t.repo.AssertExpectations(t.T())
		t.tx.AssertExpectations(t.T())
	})

	t.Run("it should return error when identity provider (keycloak) fails", func() {
		isAdmin := false
		t.tx.On("WithTransaction", t.ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(errors.New("keycloak unavailable"))
		t.repo.On("ExistsByEmailAndUsername", t.ctx, "test@example.com", mock.AnythingOfType("string")).Return(false, nil)
		t.ip.On("CreateUser", t.ctx, mock.AnythingOfType("port.RegisterUserParam")).
			Return("", errors.New("keycloak unavailable"))

		res, err := t.svc.RegisterUser(t.ctx, "test@example.com", "John", "Doe", "pass@1234A", &isAdmin)

		t.Error(err)
		t.Nil(res)
		t.repo.AssertExpectations(t.T())
		t.tx.AssertExpectations(t.T())
	})
}

// ─── ValidateAuthenticateUser ─────────────────────────────────────────────────

func (t *UserServiceTestSuite) TestValidateAuthenticateUser() {
	t.Run("it should return user when email is valid and account is active", func() {
		u, _ := entity.NewUser("active@example.com", "John", "Doe", nil)
		u.Status = entity.UserStatusActive
		t.repo.On("FindByEmail", t.ctx, "active@example.com").Return(u, nil)

		res, err := t.svc.ValidateAuthenticateUser(t.ctx, "active@example.com")

		t.NoError(err)
		t.NotNil(res)
		t.Equal("active@example.com", res.Email.String())
		t.repo.AssertExpectations(t.T())
	})

	t.Run("it should return error when email is not registered", func() {
		businessErr := pkgErrors.NewBusinessError("NOT_FOUND_", "USER_FBE", "User not found")
		t.repo.On("FindByEmail", t.ctx, "ghost@example.com").Return(nil, businessErr)

		res, err := t.svc.ValidateAuthenticateUser(t.ctx, "ghost@example.com")

		t.Error(err)
		t.Nil(res)
		t.Contains(err.Error(), "Email was not registered")
		t.repo.AssertExpectations(t.T())
	})

	t.Run("it should return error when account is suspended", func() {
		u, _ := entity.NewUser("suspended@example.com", "John", "Doe", nil)
		u.Status = entity.UserStatusSuspended
		t.repo.On("FindByEmail", t.ctx, "suspended@example.com").Return(u, nil)

		res, err := t.svc.ValidateAuthenticateUser(t.ctx, "suspended@example.com")

		t.Error(err)
		t.Nil(res)
		t.Contains(err.Error(), "suspended")
		t.repo.AssertExpectations(t.T())
	})

	t.Run("it should return error when account is inactive", func() {
		u, _ := entity.NewUser("inactive@example.com", "John", "Doe", nil)
		u.Status = user.UserStatusInactive
		t.repo.On("FindByEmail", t.ctx, "inactive@example.com").Return(u, nil)

		res, err := t.svc.ValidateAuthenticateUser(t.ctx, "inactive@example.com")

		t.Error(err)
		t.Nil(res)
		t.Contains(err.Error(), "inactive")
		t.repo.AssertExpectations(t.T())
	})
}

// ─── GenerateToken ────────────────────────────────────────────────────────────

func (t *UserServiceTestSuite) TestGenerateToken() {
	t.Run("it should generate access token and refresh token successfully", func() {
		t.ip.On("GetAccessToken", t.ctx, "test@example.com", "correctpass", "", constant.KeycloakGrantTypePasswordConst, constant.KeycloakScope).
			Return(&port.TokenResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			}, nil)

		acc, ref, err := t.svc.GenerateToken(t.ctx, "test@example.com", "correctpass", constant.KeycloakGrantTypePasswordConst)

		t.NoError(err)
		t.NotNil(acc)
		t.NotNil(ref)
		t.Equal("access-token", *acc)
		t.Equal("refresh-token", *ref)
		t.ip.AssertExpectations(t.T())
	})

	t.Run("it should return error when credentials are invalid", func() {
		t.ip.On("GetAccessToken", t.ctx, "test@example.com", "wrongpass", "", constant.KeycloakGrantTypePasswordConst, constant.KeycloakScope).
			Return(nil, pkgErrors.NewBusinessError("KCGT", "003", "Invalid user credentials"))

		acc, ref, err := t.svc.GenerateToken(t.ctx, "test@example.com", "wrongpass", constant.KeycloakGrantTypePasswordConst)

		t.Error(err)
		t.Nil(acc)
		t.Nil(ref)
		t.ip.AssertExpectations(t.T())
	})
}

// ─── RefreshToken ─────────────────────────────────────────────────────────────

func (t *UserServiceTestSuite) TestRefreshToken() {
	t.Run("it should refresh token and return user successfully", func() {
		t.ip.On("RefreshToken", t.ctx, "valid-refresh-token").Return(&port.TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
		}, nil)
		t.ip.On("GetUserInfo", t.ctx, "new-access-token").Return(&port.UserInfoResponse{
			Sub:   "uuid-1234",
			Email: "test@example.com",
		}, nil)
		u, _ := entity.NewUser("test@example.com", "John", "Doe", nil)
		u.SetKeycloakUUID("uuid-1234")
		t.repo.On("FindByKeycloakUUID", t.ctx, "uuid-1234").Return(u, nil)

		resUser, acc, ref, err := t.svc.RefreshToken(t.ctx, "valid-refresh-token")

		t.NoError(err)
		t.NotNil(resUser)
		t.Equal("new-access-token", *acc)
		t.Equal("new-refresh-token", *ref)
		t.ip.AssertExpectations(t.T())
		t.repo.AssertExpectations(t.T())
	})

	t.Run("it should return error when refresh token is invalid or expired", func() {
		t.ip.On("RefreshToken", t.ctx, "expired-token").
			Return(nil, pkgErrors.NewBusinessError("KCRT", "003", "Token is not active"))

		resUser, acc, ref, err := t.svc.RefreshToken(t.ctx, "expired-token")

		t.Error(err)
		t.Nil(resUser)
		t.Nil(acc)
		t.Nil(ref)
		t.ip.AssertExpectations(t.T())
	})

	t.Run("it should return error when user info cannot be fetched from identity provider", func() {
		t.ip.On("RefreshToken", t.ctx, "valid-refresh-token").Return(&port.TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
		}, nil)
		t.ip.On("GetUserInfo", t.ctx, "new-access-token").
			Return(nil, pkgErrors.NewBusinessError("KCGU", "002", "failed to fetch user info"))

		resUser, acc, ref, err := t.svc.RefreshToken(t.ctx, "valid-refresh-token")

		t.Error(err)
		t.Nil(resUser)
		t.Nil(acc)
		t.Nil(ref)
		t.ip.AssertExpectations(t.T())
	})
}

// ─── DeleteUser ───────────────────────────────────────────────────────────────

func (t *UserServiceTestSuite) TestDeleteUser() {
	t.Run("it should delete user successfully", func() {
		u := &entity.User{KeycloakUUID: "uuid-123"}
		u.SetID(1)
		t.tx.On("WithTransaction", t.ctx, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		t.repo.On("FindByID", t.ctx, uint(1)).Return(u, nil)
		t.ip.On("DeleteUser", t.ctx, "uuid-123").Return(nil)
		t.repo.On("DeleteByID", t.ctx, uint(1)).Return(nil)

		err := t.svc.DeleteUser(t.ctx, uint(1))

		t.NoError(err)
		t.repo.AssertExpectations(t.T())
		t.ip.AssertExpectations(t.T())
		t.tx.AssertExpectations(t.T())
	})

	t.Run("it should return error when user is not found", func() {
		t.tx.On("WithTransaction", t.ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(pkgErrors.NewBusinessError("NOT_FOUND_", "USER_BYID", "User not found"))
		t.repo.On("FindByID", t.ctx, uint(99)).
			Return(nil, pkgErrors.NewBusinessError("NOT_FOUND_", "USER_BYID", "User not found"))

		err := t.svc.DeleteUser(t.ctx, uint(99))

		t.Error(err)
		t.repo.AssertExpectations(t.T())
		t.tx.AssertExpectations(t.T())
	})

	t.Run("it should return error when identity provider fails to delete", func() {
		u := &entity.User{KeycloakUUID: "uuid-123"}
		u.SetID(1)
		t.tx.On("WithTransaction", t.ctx, mock.AnythingOfType("func(context.Context) error")).
			Return(pkgErrors.NewInfrastructureError("KCDU", "001", "keycloak error"))
		t.repo.On("FindByID", t.ctx, uint(1)).Return(u, nil)
		t.ip.On("DeleteUser", t.ctx, "uuid-123").
			Return(pkgErrors.NewInfrastructureError("KCDU", "001", "keycloak error"))

		err := t.svc.DeleteUser(t.ctx, uint(1))

		t.Error(err)
		t.repo.AssertExpectations(t.T())
		t.ip.AssertExpectations(t.T())
		t.tx.AssertExpectations(t.T())
	})
}

*/
