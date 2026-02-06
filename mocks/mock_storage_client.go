package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockStorageClient is a mock implementation of output.StorageClient
type MockStorageClient struct {
	mock.Mock
}

// DeleteStorageFile mocks the DeleteStorageFile method
func (m *MockStorageClient) DeleteStorageFile(ctx context.Context, fileURL string) error {
	args := m.Called(ctx, fileURL)
	return args.Error(0)
}
