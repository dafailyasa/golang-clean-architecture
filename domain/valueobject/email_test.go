package valueobject

import (
	"auth-service/domain/enum"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValueObjectsTestSuite struct {
	suite.Suite
}

func (suite *ValueObjectsTestSuite) TestNewEmail() {
	tests := []struct {
		name    string
		email   string
		wantErr bool
		wantVal string
	}{
		{"Valid email", "test@example.com", false, "test@example.com"},
		{"Valid email uppercase", "Test@Example.Com", false, "test@example.com"},
		{"Valid email with spaces", " test@example.com ", false, "test@example.com"},
		{"Invalid email missing @", "testexample.com", true, ""},
		{"Invalid email missing domain", "test@", true, ""},
		{"Empty email", "", true, ""},
		{"Empty string with spaces", "   ", true, ""},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := NewEmail(tt.email)
			if tt.wantErr {
				suite.Error(err)
				suite.Nil(got)
			} else {
				suite.NoError(err)
				suite.NotNil(got)
				suite.Equal(tt.wantVal, got.String())
			}
		})
	}
}

func (suite *ValueObjectsTestSuite) TestStatus_IsValid() {
	tests := []struct {
		name   string
		status enum.Status
		want   bool
	}{
		{"Active", enum.UserStatusActive, true},
		{"Inactive", enum.UserStatusInactive, true},
		{"Suspended", enum.UserStatusSuspended, true},
		{"Invalid status", enum.Status("unknown"), false},
		{"Empty status", enum.Status(""), false},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got := tt.status.IsValid()
			suite.Equal(tt.want, got)
		})
	}
}

func TestValueObjectsTestSuite(t *testing.T) {
	suite.Run(t, new(ValueObjectsTestSuite))
}
