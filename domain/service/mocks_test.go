package service

import (
	"auth-service/domain/aggregate"
	"auth-service/domain/valueobject"

	"auth-service/pkg/pagination"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *aggregate.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmailAndUsername(ctx context.Context, email, username string) (bool, error) {
	args := m.Called(ctx, email, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*aggregate.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*aggregate.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByKeycloakUUID(ctx context.Context, keycloakUUID string) (*aggregate.User, error) {
	args := m.Called(ctx, keycloakUUID)
	if args.Get(0) != nil {
		return args.Get(0).(*aggregate.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*aggregate.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*aggregate.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, u *aggregate.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) ExistEmailByUserID(ctx context.Context, id uint, email string) (bool, error) {
	args := m.Called(ctx, id, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) DeleteByID(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(ctx context.Context, req pagination.PaginationRequest) ([]*aggregate.User, int64, error) {
	args := m.Called(ctx, req)
	if args.Get(0) != nil {
		return args.Get(0).([]*aggregate.User), args.Get(1).(int64), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	args := m.Called(ctx, fn)
	if fn != nil {
		// Attempt to run the actual function for accurate testing
		err := fn(ctx)
		if err != nil {
			return err
		}
	}
	return args.Error(0)
}

type MockIdentityProvider struct {
	mock.Mock
}

func (m *MockIdentityProvider) CreateUser(ctx context.Context, param valueobject.RegisterUserParam) (string, error) {
	args := m.Called(ctx, param)
	return args.String(0), args.Error(1)
}

func (m *MockIdentityProvider) GetAccessToken(ctx context.Context, email, password, clientID, grantType, scope string) (*valueobject.TokenResponse, error) {
	args := m.Called(ctx, email, password, clientID, grantType, scope)
	if args.Get(0) != nil {
		return args.Get(0).(*valueobject.TokenResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockIdentityProvider) RefreshToken(ctx context.Context, refreshToken string) (*valueobject.TokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) != nil {
		return args.Get(0).(*valueobject.TokenResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockIdentityProvider) GetUserInfo(ctx context.Context, accessToken string) (*valueobject.UserInfoResponse, error) {
	args := m.Called(ctx, accessToken)
	if args.Get(0) != nil {
		return args.Get(0).(*valueobject.UserInfoResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockIdentityProvider) UpdateUser(ctx context.Context, keycloakUUID string, param valueobject.UpdateUserParam) error {
	args := m.Called(ctx, keycloakUUID, param)
	return args.Error(0)
}

func (m *MockIdentityProvider) DeleteUser(ctx context.Context, keycloakUUID string) error {
	args := m.Called(ctx, keycloakUUID)
	return args.Error(0)
}
