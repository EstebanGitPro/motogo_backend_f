package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// parsePropertiesParam Tests
// ============================================

func TestParsePropertiesParam_SingleValue(t *testing.T) {
	result := parsePropertiesParam("name")

	assert.Len(t, result, 1)
	assert.Equal(t, "name", result[0])
}

func TestParsePropertiesParam_MultipleValues(t *testing.T) {
	result := parsePropertiesParam("name, email, phone")

	assert.Len(t, result, 3)
	assert.Equal(t, "name", result[0])
	assert.Equal(t, "email", result[1])
	assert.Equal(t, "phone", result[2])
}

func TestParsePropertiesParam_WithQuotes(t *testing.T) {
	result := parsePropertiesParam("'name', \"email\"")

	assert.Len(t, result, 2)
	assert.Equal(t, "name", result[0])
	assert.Equal(t, "email", result[1])
}

func TestParsePropertiesParam_EmptyString(t *testing.T) {
	result := parsePropertiesParam("")

	assert.Empty(t, result)
}

func TestParsePropertiesParam_InterfaceType(t *testing.T) {
	// Simulates what happens when a jsonschema param value is an interface{}
	var value interface{} = "latitude, longitude"
	result := parsePropertiesParam(value)

	assert.Len(t, result, 2)
	assert.Equal(t, "latitude", result[0])
	assert.Equal(t, "longitude", result[1])
}

func TestParsePropertiesParam_WhitespaceOnly(t *testing.T) {
	result := parsePropertiesParam("  ,  ,  ")

	assert.Empty(t, result)
}
