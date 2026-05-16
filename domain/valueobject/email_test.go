package valueobject

import (
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
		status Status
		want   bool
	}{
		{"Active", UserStatusActive, true},
		{"Inactive", UserStatusInactive, true},
		{"Suspended", UserStatusSuspended, true},
		{"Invalid status", "unknown", false},
		{"Empty status", "", false},
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
