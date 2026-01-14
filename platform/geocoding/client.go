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

// Coordinates represents geographic coordinates from geocoding
type Coordinates struct {
	Latitude   float64
	Longitude  float64
	Confidence int // 1-10, higher is more precise
}

// Client interface for geocoding operations (output port)
type Client interface {
	Geocode(ctx context.Context, address, city, department string) (*Coordinates, error)
}

// openCageClient implements Client using OpenCage API
type openCageClient struct {
	apiKey      string
	baseURL     string
	httpClient  *http.Client
	countryCode string
	logger      logger.Logger
}

// OpenCageResponse represents the API response structure
type OpenCageResponse struct {
	Results []struct {
		Geometry struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"geometry"`
		Confidence int `json:"confidence"`
	} `json:"results"`
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Rate struct {
		Limit     int `json:"limit"`
		Remaining int `json:"remaining"`
	} `json:"rate"`
}

// NewOpenCageClient creates a new OpenCage geocoding client
func NewOpenCageClient(apiKey, baseURL, countryCode string, timeout time.Duration, log logger.Logger) Client {
	return &openCageClient{
		apiKey:      apiKey,
		baseURL:     baseURL,
		countryCode: countryCode,
		httpClient:  &http.Client{Timeout: timeout},
		logger:      log,
	}
}

// Geocode converts address components to coordinates
// Returns nil, nil if no results found (not an error)
// Returns nil, error if there was an API/network error
func (c *openCageClient) Geocode(ctx context.Context, address, city, department string) (*Coordinates, error) {
	// Build query: "Calle 123, Bogotá, Cundinamarca, Colombia"
	query := fmt.Sprintf("%s, %s, %s, Colombia", address, city, department)

	// Build URL with parameters
	reqURL := fmt.Sprintf("%s/json?q=%s&key=%s&countrycode=%s&limit=1&no_annotations=1",
		c.baseURL,
		url.QueryEscape(query),
		c.apiKey,
		c.countryCode,
	)

	c.logger.Debug(logger.LogGeocodingRequest, "query", query)

	// Create request with context for timeout/cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geocoding request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(logger.LogGeocodingError, "error", err)
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result OpenCageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.logger.Error(logger.LogGeocodingError, "error", "failed to decode response")
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	// Check API status
	if result.Status.Code != 200 {
		c.logger.Warn(logger.LogGeocodingError,
			"status_code", result.Status.Code,
			"message", result.Status.Message)
		return nil, fmt.Errorf("geocoding API error: %s (code: %d)", result.Status.Message, result.Status.Code)
	}

	// Log rate limit info for monitoring
	c.logger.Debug(logger.LogGeocodingRateLimit,
		"remaining", result.Rate.Remaining,
		"limit", result.Rate.Limit)

	// Check if we got results
	if len(result.Results) == 0 {
		c.logger.Warn(logger.LogGeocodingNoResults, "query", query)
		return nil, nil // No results, but not an error - address may be invalid
	}

	coords := &Coordinates{
		Latitude:   result.Results[0].Geometry.Lat,
		Longitude:  result.Results[0].Geometry.Lng,
		Confidence: result.Results[0].Confidence,
	}

	c.logger.Info(logger.LogGeocodingSuccess,
		"query", query,
		"lat", coords.Latitude,
		"lng", coords.Longitude,
		"confidence", coords.Confidence)

	return coords, nil
}
