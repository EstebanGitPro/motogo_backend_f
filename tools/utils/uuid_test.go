package utils_test

import (
	"testing"

	utils "github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerate_CreatesValidUUID(t *testing.T) {
	// Act
	id := utils.Generate()

	// Assert
	assert.NotEmpty(t, id)
	assert.True(t, utils.IsValid(id), "Generated UUID should be valid")
	assert.Equal(t, 36, len(id), "UUID should be 36 characters long (with hyphens)")
}

func TestGenerate_CreatesUniqueUUIDs(t *testing.T) {
	// Act
	id1 := utils.Generate()
	id2 := utils.Generate()

	// Assert
	assert.NotEqual(t, id1, id2, "Generated UUIDs should be unique")
}

func TestIsValid_ValidUUID(t *testing.T) {
	// Arrange
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	result := utils.IsValid(validUUID)

	// Assert
	assert.True(t, result)
}

func TestIsValid_InvalidUUID(t *testing.T) {
	testCases := []struct {
		name string
		uuid string
	}{
		{"Empty string", ""},
		{"Invalid format", "not-a-uuid"},
		{"Too short", "550e8400"},
		{"Invalid characters", "zzz-invalid-uuid-zzz"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := utils.IsValid(tc.uuid)

			// Assert
			assert.False(t, result, "UUID '%s' should be invalid", tc.uuid)
		})
	}
}

func TestParse_ValidUUID(t *testing.T) {
	// Arrange
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	parsed, err := utils.Parse(validUUID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
	assert.Equal(t, validUUID, parsed.String())
}

func TestParse_InvalidUUID(t *testing.T) {
	// Arrange
	invalidUUID := "invalid-uuid"

	// Act
	parsed, err := utils.Parse(invalidUUID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", parsed.String())
}

func TestMustParse_ValidUUID(t *testing.T) {
	// Arrange
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	// Act
	parsed := utils.MustParse(validUUID)

	// Assert
	assert.NotNil(t, parsed)
	assert.Equal(t, validUUID, parsed.String())
}

func TestMustParse_InvalidUUID_Panics(t *testing.T) {
	// Arrange
	invalidUUID := "invalid-uuid"

	// Assert that it panics
	assert.Panics(t, func() {
		utils.MustParse(invalidUUID)
	}, "MustParse should panic with invalid UUID")
}
