package geocoding

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewGoogleMapsClient Tests
// ============================================

func TestNewGoogleMapsClient_ReturnsClient(t *testing.T) {
	log := logger.NewSlogLogger()
	client := NewGoogleMapsClient("api-key", "https://maps.googleapis.com/maps/api/geocode", "co", time.Second*5, log)

	assert.NotNil(t, client)
}

// ============================================
// Geocode Tests with Mock Server
// ============================================

func TestGoogleMapsClient_Geocode_Success_Rooftop(t *testing.T) {
	response := GoogleMapsResponse{
		Status: "OK",
		Results: []GoogleMapsResult{
			{
				Geometry: GoogleMapsGeometry{
					Location:     GoogleMapsLocation{Lat: 4.7110, Lng: -74.0721},
					LocationType: "ROOFTOP",
				},
				FormattedAddress: "Calle 100 #15-20, Bogotá, Colombia",
				Types:            []string{"street_address"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "test-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, 4.7110, coords.Latitude)
	assert.Equal(t, -74.0721, coords.Longitude)
	assert.Equal(t, 10, coords.Confidence) // ROOFTOP = 10
}

func TestGoogleMapsClient_Geocode_Success_Approximate(t *testing.T) {
	response := GoogleMapsResponse{
		Status: "OK",
		Results: []GoogleMapsResult{
			{
				Geometry: GoogleMapsGeometry{
					Location:     GoogleMapsLocation{Lat: 6.2518, Lng: -75.5636},
					LocationType: "APPROXIMATE",
				},
				FormattedAddress: "Medellín, Colombia",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "test-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Centro", "Medellín", "Antioquia")

	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, 3, coords.Confidence) // APPROXIMATE = 3
}

func TestGoogleMapsClient_Geocode_ZeroResults(t *testing.T) {
	response := GoogleMapsResponse{
		Status:  "ZERO_RESULTS",
		Results: []GoogleMapsResult{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "test-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Dirección Inexistente", "Ciudad", "Departamento")

	assert.NoError(t, err) // Not an error, just no results
	assert.Nil(t, coords)
}

func TestGoogleMapsClient_Geocode_OverQueryLimit(t *testing.T) {
	response := GoogleMapsResponse{
		Status:  "OVER_QUERY_LIMIT",
		Results: []GoogleMapsResult{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "test-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Contains(t, err.Error(), "quota exceeded")
}

func TestGoogleMapsClient_Geocode_RequestDenied(t *testing.T) {
	response := GoogleMapsResponse{
		Status: "REQUEST_DENIED",
		Error:  &GoogleMapsError{Message: "Invalid API key"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "invalid-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Contains(t, err.Error(), "request denied")
}

func TestGoogleMapsClient_Geocode_InvalidRequest(t *testing.T) {
	response := GoogleMapsResponse{
		Status: "INVALID_REQUEST",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "test-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "", "", "")

	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Contains(t, err.Error(), "invalid request")
}

func TestGoogleMapsClient_Geocode_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "test-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.Error(t, err)
	assert.Nil(t, coords)
}

func TestGoogleMapsClient_Geocode_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := &googleMapsClient{
		apiKey:     "test-key",
		baseURL:    server.URL,
		region:     "co",
		httpClient: &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Contains(t, err.Error(), "decode")
}

// ============================================
// googleLocationTypeToConfidence Tests
// ============================================

func TestGoogleLocationTypeToConfidence(t *testing.T) {
	tests := []struct {
		locationType string
		expected     int
	}{
		{"ROOFTOP", 10},
		{"RANGE_INTERPOLATED", 8},
		{"GEOMETRIC_CENTER", 5},
		{"APPROXIMATE", 3},
		{"UNKNOWN", 1},
		{"", 1},
	}

	for _, tc := range tests {
		t.Run(tc.locationType, func(t *testing.T) {
			result := googleLocationTypeToConfidence(tc.locationType)
			assert.Equal(t, tc.expected, result)
		})
	}
}
