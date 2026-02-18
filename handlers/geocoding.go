package handlers

// GeocodingTestRequest is the DTO for testing geocoding
type GeocodingTestRequest struct {
	Address        string `json:"address" binding:"required"`
	CityName       string `json:"city_name" binding:"required"`
	DepartmentName string `json:"department_name" binding:"required"`
}

// Sanitize trims whitespace from all string fields
func (r *GeocodingTestRequest) Sanitize() {
	r.Address = TrimString(r.Address)
	r.CityName = TrimString(r.CityName)
	r.DepartmentName = TrimString(r.DepartmentName)
}

// GeocodingTestResponse is the response for geocoding test
type GeocodingTestResponse struct {
	Geocoded         bool    `json:"geocoded"`
	Latitude         float64 `json:"latitude,omitempty"`
	Longitude        float64 `json:"longitude,omitempty"`
	FormattedAddress string  `json:"formatted_address,omitempty"`
	Confidence       float64 `json:"confidence,omitempty"`
	Error            string  `json:"error,omitempty"`
}
