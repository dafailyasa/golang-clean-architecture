package entity

import (
	"auth-service/domain/valueobject"

	"testing"

	"github.com/stretchr/testify/suite"
)

type EntityTestSuite struct {
	suite.Suite
	isAdmin bool
}

func (suite *EntityTestSuite) SetupTest() {
	suite.isAdmin = false
}

func (suite *EntityTestSuite) TestNewUser() {
	suite.Run("Success", func() {
		got, err := NewUser("test@example.com", "John", "Doe", &suite.isAdmin)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal("test@example.com", got.Email.String())
		suite.Equal("John", got.FirstName)
		suite.Equal("Doe", got.LastName)
		suite.Equal(valueobject.UserStatusActive, got.Status)
		suite.Equal("John Doe", got.GetFullName())
		suite.NotEmpty(got.Username)
		suite.False(got.IsInActiveAndSuspend())
	})

	suite.Run("Invalid Email", func() {
		got, err := NewUser("invalid", "John", "Doe", &suite.isAdmin)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("Invalid First Name", func() {
		got, err := NewUser("test@example.com", "J", "Doe", &suite.isAdmin)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("Invalid Last Name", func() {
		got, err := NewUser("test@example.com", "John", "D", &suite.isAdmin)
		suite.Error(err)
		suite.Nil(got)
	})
}

func (suite *EntityTestSuite) TestUser_Rebuild() {
	u, _ := NewUser("test@example.com", "John", "Doe", &suite.isAdmin)

	suite.Run("Success Rebuild", func() {
		newIsAdmin := true
		err := u.Rebuild("Jane", "Smith", "jane@example.com", valueobject.UserStatusSuspended, &newIsAdmin)
		suite.NoError(err)
		suite.Equal("Jane", u.FirstName)
		suite.Equal("Smith", u.LastName)
		suite.Equal("jane@example.com", u.Email.String())
		suite.Equal(valueobject.UserStatusSuspended, u.Status)
		suite.True(u.IsInActiveAndSuspend())
	})

	suite.Run("Invalid Status", func() {
		err := u.Rebuild("Jane", "Smith", "jane@example.com", "unknown_status", nil)
		suite.Error(err)
	})

	suite.Run("Invalid Email", func() {
		err := u.Rebuild("Jane", "Smith", "invalid", valueobject.UserStatusActive, nil)
		suite.Error(err)
	})
}

func TestEntityTestSuite(t *testing.T) {
	suite.Run(t, new(EntityTestSuite))
}
