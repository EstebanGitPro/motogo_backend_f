package geocoding

import "context"

// Client defines the interface for geocoding operations
// All geocoding providers (Google, Mapbox) must implement this interface
type Client interface {
	// Geocode converts address components to geographic coordinates
	// Returns nil, nil if no results found (not an error)
	// Returns nil, error if there was an API/network error
	Geocode(ctx context.Context, address, city, department string) (*Coordinates, error)
}

// Coordinates represents geographic coordinates from geocoding
type Coordinates struct {
	Latitude   float64
	Longitude  float64
	Confidence int // 1-10 scale (10 = highest confidence)
}
