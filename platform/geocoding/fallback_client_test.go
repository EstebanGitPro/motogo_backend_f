package geocoding

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of geocoding.Geocoder for testing
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Geocode(ctx context.Context, address, city, department string) (*Coordinates, error) {
	args := m.Called(ctx, address, city, department)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Coordinates), args.Error(1)
}

// ============================================
// NewFallbackClient Tests
// ============================================

func TestNewFallbackClient_ReturnsClient(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	// Act
	client := NewFallbackClient(primary, fallback)

	// Assert
	assert.NotNil(t, client)
	assert.Implements(t, (*Geocoder)(nil), client)
}

// ============================================
// Geocode Tests - Primary Success
// ============================================

func TestFallbackClient_Geocode_PrimarySuccess(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	expectedCoords := &Coordinates{
		Latitude:   4.7110,
		Longitude:  -74.0721,
		Confidence: 9,
	}

	primary.On("Geocode", mock.Anything, "Calle 100", "Bogotá", "Cundinamarca").
		Return(expectedCoords, nil)

	client := NewFallbackClient(primary, fallback)

	// Act
	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, expectedCoords.Latitude, coords.Latitude)
	assert.Equal(t, expectedCoords.Longitude, coords.Longitude)
	primary.AssertExpectations(t)
	fallback.AssertNotCalled(t, "Geocode")
}

func TestFallbackClient_Geocode_PrimaryReturnsNilCoords(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	primary.On("Geocode", mock.Anything, "Unknown Address", "Unknown", "Unknown").
		Return(nil, nil) // No results found, not an error

	client := NewFallbackClient(primary, fallback)

	// Act
	coords, err := client.Geocode(context.Background(), "Unknown Address", "Unknown", "Unknown")

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, coords)
	primary.AssertExpectations(t)
	fallback.AssertNotCalled(t, "Geocode")
}

// ============================================
// Geocode Tests - Quota Error Triggers Fallback
// ============================================

func TestFallbackClient_Geocode_QuotaError_FallbackSuccess(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	quotaErr := errors.New("OVER_QUERY_LIMIT: Quota exceeded")
	fallbackCoords := &Coordinates{
		Latitude:   4.7110,
		Longitude:  -74.0721,
		Confidence: 8,
	}

	primary.On("Geocode", mock.Anything, "Calle 100", "Bogotá", "Cundinamarca").
		Return(nil, quotaErr)
	fallback.On("Geocode", mock.Anything, "Calle 100", "Bogotá", "Cundinamarca").
		Return(fallbackCoords, nil)

	client := NewFallbackClient(primary, fallback)

	// Act
	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, fallbackCoords.Latitude, coords.Latitude)
	primary.AssertExpectations(t)
	fallback.AssertExpectations(t)
}

func TestFallbackClient_Geocode_QuotaError_FallbackAlsoFails(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	quotaErr := errors.New("quota exceeded")
	fallbackErr := errors.New("too many requests")

	primary.On("Geocode", mock.Anything, "Calle 100", "Bogotá", "Cundinamarca").
		Return(nil, quotaErr)
	fallback.On("Geocode", mock.Anything, "Calle 100", "Bogotá", "Cundinamarca").
		Return(nil, fallbackErr)

	client := NewFallbackClient(primary, fallback)

	// Act
	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Equal(t, fallbackErr, err)
	primary.AssertExpectations(t)
	fallback.AssertExpectations(t)
}

func TestFallbackClient_Geocode_RateLimitError_TriggersFallback(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	rateLimitErr := errors.New("rate limit exceeded")
	fallbackCoords := &Coordinates{Latitude: 4.7, Longitude: -74.0, Confidence: 7}

	primary.On("Geocode", mock.Anything, "Address", "City", "Dept").
		Return(nil, rateLimitErr)
	fallback.On("Geocode", mock.Anything, "Address", "City", "Dept").
		Return(fallbackCoords, nil)

	client := NewFallbackClient(primary, fallback)

	// Act
	coords, err := client.Geocode(context.Background(), "Address", "City", "Dept")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, coords)
	primary.AssertExpectations(t)
	fallback.AssertExpectations(t)
}

// ============================================
// Geocode Tests - Non-Quota Error
// ============================================

func TestFallbackClient_Geocode_NonQuotaError_DoesNotTriggerFallback(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	networkErr := errors.New("network connection failed")

	primary.On("Geocode", mock.Anything, "Calle 100", "Bogotá", "Cundinamarca").
		Return(nil, networkErr)

	client := NewFallbackClient(primary, fallback)

	// Act
	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Equal(t, networkErr, err)
	primary.AssertExpectations(t)
	fallback.AssertNotCalled(t, "Geocode")
}

func TestFallbackClient_Geocode_InvalidAddressError_DoesNotTriggerFallback(t *testing.T) {
	// Arrange
	primary := new(MockClient)
	fallback := new(MockClient)

	invalidErr := errors.New("ZERO_RESULTS: Address not found")

	primary.On("Geocode", mock.Anything, "XYZ123", "ABC", "DEF").
		Return(nil, invalidErr)

	client := NewFallbackClient(primary, fallback)

	// Act
	coords, err := client.Geocode(context.Background(), "XYZ123", "ABC", "DEF")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, coords)
	fallback.AssertNotCalled(t, "Geocode")
}

// ============================================
// isQuotaError Tests
// ============================================

func TestIsQuotaError_NilError(t *testing.T) {
	assert.False(t, isQuotaError(nil))
}

func TestIsQuotaError_QuotaKeyword(t *testing.T) {
	err := errors.New("API quota exceeded for today")
	assert.True(t, isQuotaError(err))
}

func TestIsQuotaError_OverQueryLimit(t *testing.T) {
	err := errors.New("OVER_QUERY_LIMIT: You have exceeded your daily request quota")
	assert.True(t, isQuotaError(err))
}

func TestIsQuotaError_OverDailyLimit(t *testing.T) {
	err := errors.New("over_daily_limit error")
	assert.True(t, isQuotaError(err))
}

func TestIsQuotaError_RateLimit(t *testing.T) {
	err := errors.New("rate limit exceeded, try again later")
	assert.True(t, isQuotaError(err))
}

func TestIsQuotaError_TooManyRequests(t *testing.T) {
	err := errors.New("Too many requests, please wait")
	assert.True(t, isQuotaError(err))
}

func TestIsQuotaError_CaseInsensitive(t *testing.T) {
	err := errors.New("QUOTA EXCEEDED")
	assert.True(t, isQuotaError(err))
}

func TestIsQuotaError_NonQuotaError(t *testing.T) {
	testCases := []struct {
		name    string
		errMsg  string
		isQuota bool
	}{
		{"network error", "network connection timeout", false},
		{"invalid API key", "invalid API key provided", false},
		{"server error", "internal server error", false},
		{"not found", "address not found", false},
		{"empty error", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := errors.New(tc.errMsg)
			assert.Equal(t, tc.isQuota, isQuotaError(err))
		})
	}
}
