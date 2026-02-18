package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// googleMapsClient implements Client using Google Maps Geocoding API
type googleMapsClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	region     string // "co" for Colombia - biases results
}

// GoogleMapsResponse represents the Google Maps Geocoding API response
type GoogleMapsResponse struct {
	Results []GoogleMapsResult `json:"results"`
	Status  string             `json:"status"` // "OK", "ZERO_RESULTS", "OVER_QUERY_LIMIT", etc.
	Error   *GoogleMapsError   `json:"error_message,omitempty"`
}

// GoogleMapsResult represents a single geocoding result
type GoogleMapsResult struct {
	Geometry         GoogleMapsGeometry `json:"geometry"`
	FormattedAddress string             `json:"formatted_address"`
	Types            []string           `json:"types"`
}

// GoogleMapsGeometry contains location data
type GoogleMapsGeometry struct {
	Location     GoogleMapsLocation `json:"location"`
	LocationType string             `json:"location_type"` // "ROOFTOP", "RANGE_INTERPOLATED", "GEOMETRIC_CENTER", "APPROXIMATE"
}

// GoogleMapsLocation contains the coordinates
type GoogleMapsLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// GoogleMapsError represents an error from the API
type GoogleMapsError struct {
	Message string `json:"message"`
}

// NewGoogleMapsClient creates a new Google Maps geocoding client
func NewGoogleMapsClient(apiKey, baseURL, region string, timeout time.Duration, log logger.Logger) Geocoder {
	return &googleMapsClient{
		apiKey:     apiKey,
		baseURL:    baseURL,
		region:     region,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// Geocode converts address components to coordinates using Google Maps Geocoding API
// Returns nil, nil if no results found (not an error)
// Returns nil, error if there was an API/network error
func (c *googleMapsClient) Geocode(ctx context.Context, address, city, department string) (*Coordinates, error) {
	// Build query: "Carrera 46 40 80, Rionegro, Antioquia, Colombia"
	query := fmt.Sprintf("%s, %s, %s, Colombia", address, city, department)

	// Build URL with parameters for Google Maps Geocoding API
	// region=co biases results to Colombia
	reqURL := fmt.Sprintf("%s/json?address=%s&key=%s&region=%s",
		c.baseURL,
		url.QueryEscape(query),
		c.apiKey,
		c.region,
	)

	log.Debug(logger.LogGeocodingRequest, "query", query, "provider", "google")

	// Create request with context for timeout/cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geocoding request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(logger.LogGeocodingError, "error", err, "provider", "google")
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }() // Body close error intentionally ignored

	// Parse response
	var result GoogleMapsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(logger.LogGeocodingError, "error", "failed to decode response", "provider", "google")
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	// Check API status
	switch result.Status {
	case "OK":
		// Success, continue processing
	case "ZERO_RESULTS":
		log.Warn(logger.LogGeocodingNoResults, "query", query, "provider", "google")
		return nil, nil // No results, but not an error
	case "OVER_QUERY_LIMIT", "OVER_DAILY_LIMIT":
		log.Error(logger.LogGeocodingError, "status", result.Status, "provider", "google")
		return nil, fmt.Errorf("google maps API quota exceeded: %s", result.Status)
	case "REQUEST_DENIED":
		log.Error(logger.LogGeocodingError, "status", result.Status, "error", result.Error, "provider", "google")
		return nil, fmt.Errorf("google maps API request denied: check API key configuration")
	case "INVALID_REQUEST":
		log.Error(logger.LogGeocodingError, "status", result.Status, "provider", "google")
		return nil, fmt.Errorf("google maps API invalid request")
	default:
		log.Error(logger.LogGeocodingError, "status", result.Status, "provider", "google")
		return nil, fmt.Errorf("google maps API error: %s", result.Status)
	}

	// Check if we got results
	if len(result.Results) == 0 {
		log.Warn(logger.LogGeocodingNoResults, "query", query, "provider", "google")
		return nil, nil
	}

	firstResult := result.Results[0]

	// Convert Google's location_type to our confidence scale (1-10)
	confidence := googleLocationTypeToConfidence(firstResult.Geometry.LocationType)

	coords := &Coordinates{
		Latitude:   firstResult.Geometry.Location.Lat,
		Longitude:  firstResult.Geometry.Location.Lng,
		Confidence: confidence,
	}

	log.Info(logger.LogGeocodingSuccess,
		"query", query,
		"lat", coords.Latitude,
		"lng", coords.Longitude,
		"confidence", coords.Confidence,
		"location_type", firstResult.Geometry.LocationType,
		"formatted_address", firstResult.FormattedAddress,
		"provider", "google")

	return coords, nil
}

// googleLocationTypeToConfidence converts Google's location_type to our 1-10 confidence scale
// https://developers.google.com/maps/documentation/geocoding/requests-geocoding#results
func googleLocationTypeToConfidence(locationType string) int {
	switch locationType {
	case "ROOFTOP":
		return 10 // Exact location, highest precision
	case "RANGE_INTERPOLATED":
		return 8 // Interpolated between two precise points
	case "GEOMETRIC_CENTER":
		return 5 // Center of a region (street, neighborhood)
	case "APPROXIMATE":
		return 3 // Approximate location
	default:
		return 1 // Unknown
	}
}
