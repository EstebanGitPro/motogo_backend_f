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

// mapboxClient implements Client using Mapbox Search API
type mapboxClient struct {
	accessToken string
	baseURL     string
	httpClient  *http.Client
	countryCode string
}

// MapboxResponse represents the Mapbox Search API response (GeoJSON FeatureCollection)
type MapboxResponse struct {
	Type     string          `json:"type"` // "FeatureCollection"
	Features []MapboxFeature `json:"features"`
}

// MapboxFeature represents a single geocoding result
type MapboxFeature struct {
	Type       string           `json:"type"` // "Feature"
	Geometry   MapboxGeometry   `json:"geometry"`
	Properties MapboxProperties `json:"properties"`
}

// MapboxGeometry contains the geographic coordinates
type MapboxGeometry struct {
	Type        string    `json:"type"`        // "Point"
	Coordinates []float64 `json:"coordinates"` // [longitude, latitude]
}

// MapboxProperties contains detailed information about the result
type MapboxProperties struct {
	FullAddress string            `json:"full_address"`
	Coordinates MapboxCoordinates `json:"coordinates"`
	MatchCode   MapboxMatchCode   `json:"match_code"`
}

// MapboxCoordinates contains precise coordinate information
type MapboxCoordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Accuracy  string  `json:"accuracy"` // "point", "street", "place", etc.
}

// MapboxMatchCode indicates the quality of the geocoding match
type MapboxMatchCode struct {
	Confidence string `json:"confidence"` // "exact", "high", "medium", "low"
	Street     string `json:"street"`     // "matched", "unmatched", "plausible"
	Place      string `json:"place"`      // "matched", "unmatched"
	Region     string `json:"region"`     // "matched", "unmatched"
	Country    string `json:"country"`    // "matched", "inferred"
}

// NewMapboxClient creates a new Mapbox geocoding client
func NewMapboxClient(accessToken, baseURL, countryCode string, timeout time.Duration) Client {
	return &mapboxClient{
		accessToken: accessToken,
		baseURL:     baseURL,
		countryCode: countryCode,
		httpClient:  &http.Client{Timeout: timeout},
	}
}

// Geocode converts address components to coordinates using Mapbox Search API
// Returns nil, nil if no results found (not an error)
// Returns nil, error if there was an API/network error
func (c *mapboxClient) Geocode(ctx context.Context, address, city, department string) (*Coordinates, error) {
	// Build query: "Carrera 46 40 80, Rionegro, Antioquia, Colombia"
	query := fmt.Sprintf("%s, %s, %s, Colombia", address, city, department)

	// Build URL with parameters for Mapbox Search API v6 forward geocoding
	reqURL := fmt.Sprintf("%s/forward?q=%s&access_token=%s&country=%s&limit=1",
		c.baseURL,
		url.QueryEscape(query),
		c.accessToken,
		c.countryCode,
	)

	log.Debug(logger.LogGeocodingRequest, "query", query, "provider", "mapbox")

	// Create request with context for timeout/cancellation
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create geocoding request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(logger.LogGeocodingError, "error", err, "provider", "mapbox")
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }() // Body close error intentionally ignored

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		log.Error(logger.LogGeocodingError,
			"status_code", resp.StatusCode,
			"provider", "mapbox")
		return nil, fmt.Errorf("mapbox API error: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var result MapboxResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(logger.LogGeocodingError, "error", "failed to decode response", "provider", "mapbox")
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	// Check if we got results
	if len(result.Features) == 0 {
		log.Warn(logger.LogGeocodingNoResults, "query", query, "provider", "mapbox")
		return nil, nil // No results, but not an error - address may be invalid
	}

	feature := result.Features[0]

	// Extract coordinates - prefer properties.coordinates for precision
	var lat, lng float64
	if feature.Properties.Coordinates.Latitude != 0 || feature.Properties.Coordinates.Longitude != 0 {
		lat = feature.Properties.Coordinates.Latitude
		lng = feature.Properties.Coordinates.Longitude
	} else if len(feature.Geometry.Coordinates) >= 2 {
		// Fallback to geometry coordinates [longitude, latitude]
		lng = feature.Geometry.Coordinates[0]
		lat = feature.Geometry.Coordinates[1]
	} else {
		log.Warn(logger.LogGeocodingNoResults, "query", query, "reason", "no coordinates in response")
		return nil, nil
	}

	// Convert Mapbox confidence to our 1-10 scale
	confidence := mapboxConfidenceToScale(feature.Properties.MatchCode.Confidence)

	coords := &Coordinates{
		Latitude:   lat,
		Longitude:  lng,
		Confidence: confidence,
	}

	log.Info(logger.LogGeocodingSuccess,
		"query", query,
		"lat", coords.Latitude,
		"lng", coords.Longitude,
		"confidence", coords.Confidence,
		"mapbox_confidence", feature.Properties.MatchCode.Confidence,
		"accuracy", feature.Properties.Coordinates.Accuracy,
		"provider", "mapbox")

	return coords, nil
}

// mapboxConfidenceToScale converts Mapbox confidence levels to our 1-10 scale
func mapboxConfidenceToScale(confidence string) int {
	switch confidence {
	case "exact":
		return 10
	case "high":
		return 8
	case "medium":
		return 5
	case "low":
		return 3
	default:
		return 1 // Unknown confidence
	}
}
