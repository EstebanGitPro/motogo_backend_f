package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockFirebaseClient is a mock implementation of firebase.FirebaseClient
type MockFirebaseClient struct {
	mock.Mock
}

// CreateCustomToken mocks the creation of a Firebase custom token
func (m *MockFirebaseClient) CreateCustomToken(ctx context.Context, uid string) (string, error) {
	args := m.Called(ctx, uid)
	return args.String(0), args.Error(1)
}

// CreateCustomTokenWithClaims mocks the creation of a Firebase custom token with claims
func (m *MockFirebaseClient) CreateCustomTokenWithClaims(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
	args := m.Called(ctx, uid, claims)
	return args.String(0), args.Error(1)
}
