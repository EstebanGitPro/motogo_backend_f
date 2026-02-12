package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockJWKSValidator is a mock implementation of output.JWTValidator
type MockJWKSValidator struct {
	mock.Mock
}

// ValidateToken mocks the JWT token validation
func (m *MockJWKSValidator) ValidateToken(tokenString string) (map[string]interface{}, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}
