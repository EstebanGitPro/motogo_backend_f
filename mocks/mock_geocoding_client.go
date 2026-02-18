package mocks

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/platform/geocoding"
	"github.com/stretchr/testify/mock"
)

// MockGeocodingClient is a mock implementation of geocoding.Geocoder
type MockGeocodingClient struct {
	mock.Mock
}

func (m *MockGeocodingClient) Geocode(ctx context.Context, address, city, department string) (*geocoding.Coordinates, error) {
	args := m.Called(ctx, address, city, department)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*geocoding.Coordinates), args.Error(1)
}
