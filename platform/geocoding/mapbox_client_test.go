package geocoding

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============================================
// NewMapboxClient Tests
// ============================================

func TestNewMapboxClient_ReturnsClient(t *testing.T) {
	client := NewMapboxClient("access-token", "https://api.mapbox.com/search/geocode/v6", "co", time.Second*5)

	assert.NotNil(t, client)
}

// ============================================
// Mapbox Geocode Tests with Mock Server
// ============================================

func TestMapboxClient_Geocode_Success_ExactConfidence(t *testing.T) {
	response := MapboxResponse{
		Type: "FeatureCollection",
		Features: []MapboxFeature{
			{
				Type: "Feature",
				Geometry: MapboxGeometry{
					Type:        "Point",
					Coordinates: []float64{-74.0721, 4.7110},
				},
				Properties: MapboxProperties{
					FullAddress: "Calle 100 #15-20, Bogotá, Colombia",
					Coordinates: MapboxCoordinates{
						Longitude: -74.0721,
						Latitude:  4.7110,
						Accuracy:  "point",
					},
					MatchCode: MapboxMatchCode{
						Confidence: "exact",
						Street:     "matched",
						Place:      "matched",
						Region:     "matched",
						Country:    "matched",
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &mapboxClient{
		accessToken: "test-token",
		baseURL:     server.URL,
		countryCode: "co",
		httpClient:  &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, 4.7110, coords.Latitude)
	assert.Equal(t, -74.0721, coords.Longitude)
	assert.Equal(t, 10, coords.Confidence) // exact = 10
}

func TestMapboxClient_Geocode_Success_LowConfidence(t *testing.T) {
	response := MapboxResponse{
		Type: "FeatureCollection",
		Features: []MapboxFeature{
			{
				Type: "Feature",
				Geometry: MapboxGeometry{
					Type:        "Point",
					Coordinates: []float64{-75.5636, 6.2518},
				},
				Properties: MapboxProperties{
					Coordinates: MapboxCoordinates{
						Longitude: -75.5636,
						Latitude:  6.2518,
					},
					MatchCode: MapboxMatchCode{
						Confidence: "low",
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &mapboxClient{
		accessToken: "test-token",
		baseURL:     server.URL,
		countryCode: "co",
		httpClient:  &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Dirección Vaga", "Medellín", "Antioquia")

	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, 3, coords.Confidence) // low = 3
}

func TestMapboxClient_Geocode_FallbackToGeometryCoordinates(t *testing.T) {
	// Test case where properties.coordinates is empty but geometry has coordinates
	response := MapboxResponse{
		Type: "FeatureCollection",
		Features: []MapboxFeature{
			{
				Type: "Feature",
				Geometry: MapboxGeometry{
					Type:        "Point",
					Coordinates: []float64{-74.5, 5.0},
				},
				Properties: MapboxProperties{
					FullAddress: "Some address",
					Coordinates: MapboxCoordinates{}, // Empty
					MatchCode: MapboxMatchCode{
						Confidence: "medium",
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &mapboxClient{
		accessToken: "test-token",
		baseURL:     server.URL,
		countryCode: "co",
		httpClient:  &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Address", "City", "Department")

	assert.NoError(t, err)
	assert.NotNil(t, coords)
	assert.Equal(t, 5.0, coords.Latitude)
	assert.Equal(t, -74.5, coords.Longitude)
}

func TestMapboxClient_Geocode_NoResults(t *testing.T) {
	response := MapboxResponse{
		Type:     "FeatureCollection",
		Features: []MapboxFeature{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &mapboxClient{
		accessToken: "test-token",
		baseURL:     server.URL,
		countryCode: "co",
		httpClient:  &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Dirección Inexistente", "Ciudad", "Departamento")

	assert.NoError(t, err) // No error, just no results
	assert.Nil(t, coords)
}

func TestMapboxClient_Geocode_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := &mapboxClient{
		accessToken: "invalid-token",
		baseURL:     server.URL,
		countryCode: "co",
		httpClient:  &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Contains(t, err.Error(), "HTTP 401")
}

func TestMapboxClient_Geocode_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := &mapboxClient{
		accessToken: "test-token",
		baseURL:     server.URL,
		countryCode: "co",
		httpClient:  &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Calle 100", "Bogotá", "Cundinamarca")

	assert.Error(t, err)
	assert.Nil(t, coords)
	assert.Contains(t, err.Error(), "decode")
}

func TestMapboxClient_Geocode_NoCoordinatesInResponse(t *testing.T) {
	response := MapboxResponse{
		Type: "FeatureCollection",
		Features: []MapboxFeature{
			{
				Type:     "Feature",
				Geometry: MapboxGeometry{}, // Empty geometry
				Properties: MapboxProperties{
					Coordinates: MapboxCoordinates{}, // Empty coordinates
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &mapboxClient{
		accessToken: "test-token",
		baseURL:     server.URL,
		countryCode: "co",
		httpClient:  &http.Client{Timeout: time.Second * 5},
	}

	coords, err := client.Geocode(context.Background(), "Address", "City", "Department")

	assert.NoError(t, err)
	assert.Nil(t, coords) // No coordinates found, returns nil
}

// ============================================
// mapboxConfidenceToScale Tests
// ============================================

func TestMapboxConfidenceToScale(t *testing.T) {
	tests := []struct {
		confidence string
		expected   int
	}{
		{"exact", 10},
		{"high", 8},
		{"medium", 5},
		{"low", 3},
		{"unknown", 1},
		{"", 1},
	}

	for _, tc := range tests {
		t.Run(tc.confidence, func(t *testing.T) {
			result := mapboxConfidenceToScale(tc.confidence)
			assert.Equal(t, tc.expected, result)
		})
	}
}
