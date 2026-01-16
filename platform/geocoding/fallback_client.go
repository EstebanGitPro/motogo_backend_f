package geocoding

import (
	"context"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var log logger.Logger = logger.NewSlogLogger()

// fallbackClient implements Client with automatic failover between providers
// Primary: Google Maps (more accurate, 10k free/month)
// Fallback: Mapbox (100k free/month)
type fallbackClient struct {
	primary  Client
	fallback Client
}

// NewFallbackClient creates a geocoding client with automatic failover
// When the primary provider returns a quota error, it automatically tries the fallback
func NewFallbackClient(primary, fallback Client) Client {
	return &fallbackClient{
		primary:  primary,
		fallback: fallback,
	}
}

// Geocode attempts geocoding with primary provider, falls back on quota errors
func (c *fallbackClient) Geocode(ctx context.Context, address, city, department string) (*Coordinates, error) {
	// Try primary provider first (Google Maps)
	coords, err := c.primary.Geocode(ctx, address, city, department)

	// If successful or no results, return as-is
	if err == nil {
		return coords, nil
	}

	// Check if this is a quota-related error that warrants fallback
	if isQuotaError(err) {
		log.Warn("geocoding_primary_quota_exceeded",
			"error", err.Error(),
			"action", "falling_back_to_secondary_provider")

		// Try fallback provider (Mapbox)
		fallbackCoords, fallbackErr := c.fallback.Geocode(ctx, address, city, department)
		if fallbackErr != nil {
			log.Error("geocoding_fallback_also_failed",
				"primary_error", err.Error(),
				"fallback_error", fallbackErr.Error())
			// Return the fallback error but could also return primary
			return nil, fallbackErr
		}

		log.Info("geocoding_fallback_succeeded",
			"lat", fallbackCoords.Latitude,
			"lng", fallbackCoords.Longitude,
			"provider", "fallback")

		return fallbackCoords, nil
	}

	// For non-quota errors, just return the primary error
	return nil, err
}

// isQuotaError checks if the error indicates a quota/rate limit issue
func isQuotaError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())

	// Google Maps quota errors
	if strings.Contains(errMsg, "quota") ||
		strings.Contains(errMsg, "over_query_limit") ||
		strings.Contains(errMsg, "over_daily_limit") ||
		strings.Contains(errMsg, "rate limit") {
		return true
	}

	// Mapbox quota errors (in case we flip the order)
	if strings.Contains(errMsg, "rate limit") ||
		strings.Contains(errMsg, "too many requests") {
		return true
	}

	return false
}
