package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFloat_ValidPositiveInteger(t *testing.T) {
	result, err := parseFloat("123")
	assert.NoError(t, err)
	assert.Equal(t, 123.0, result)
}

func TestParseFloat_ValidNegativeInteger(t *testing.T) {
	result, err := parseFloat("-456")
	assert.NoError(t, err)
	assert.Equal(t, -456.0, result)
}

func TestParseFloat_ValidDecimal(t *testing.T) {
	result, err := parseFloat("3.14159")
	assert.NoError(t, err)
	assert.InDelta(t, 3.14159, result, 0.00001)
}

func TestParseFloat_ValidNegativeDecimal(t *testing.T) {
	result, err := parseFloat("-9.81")
	assert.NoError(t, err)
	assert.InDelta(t, -9.81, result, 0.001)
}

func TestParseFloat_ValidLatitude(t *testing.T) {
	result, err := parseFloat("4.6097102")
	assert.NoError(t, err)
	assert.InDelta(t, 4.6097102, result, 0.0000001)
}

func TestParseFloat_ValidLongitude(t *testing.T) {
	result, err := parseFloat("-74.0817500")
	assert.NoError(t, err)
	assert.InDelta(t, -74.0817500, result, 0.0000001)
}

func TestParseFloat_EmptyString(t *testing.T) {
	_, err := parseFloat("")
	assert.Error(t, err)
}

func TestParseFloat_InvalidCharacters(t *testing.T) {
	_, err := parseFloat("abc")
	assert.Error(t, err)
}

func TestParseFloat_MultipleDots(t *testing.T) {
	_, err := parseFloat("1.2.3")
	assert.Error(t, err)
}

func TestParseFloat_Zero(t *testing.T) {
	result, err := parseFloat("0")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, result)
}

func TestParseFloat_ZeroDecimal(t *testing.T) {
	result, err := parseFloat("0.0")
	assert.NoError(t, err)
	assert.Equal(t, 0.0, result)
}

func TestParseFloat_LeadingDecimal(t *testing.T) {
	result, err := parseFloat(".5")
	assert.NoError(t, err)
	assert.InDelta(t, 0.5, result, 0.01)
}

// GeocodingTestRequest.Sanitize tests
func TestGeocodingTestRequest_Sanitize(t *testing.T) {
	req := GeocodingTestRequest{
		Address:        "  Calle 100 No 15-20  ",
		CityName:       "  Bogotá  ",
		DepartmentName: "  Cundinamarca  ",
	}

	req.Sanitize()

	assert.Equal(t, "Calle 100 No 15-20", req.Address)
	assert.Equal(t, "Bogotá", req.CityName)
	assert.Equal(t, "Cundinamarca", req.DepartmentName)
}

func TestGeocodingTestRequest_Sanitize_NoExtraSpace(t *testing.T) {
	req := GeocodingTestRequest{
		Address:        "Calle 50",
		CityName:       "Medellín",
		DepartmentName: "Antioquia",
	}

	req.Sanitize()

	assert.Equal(t, "Calle 50", req.Address)
	assert.Equal(t, "Medellín", req.CityName)
	assert.Equal(t, "Antioquia", req.DepartmentName)
}
