package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// Branch Domain Tests
// ============================================

func TestBranch_SetID(t *testing.T) {
	branch := &Branch{}
	branch.SetID()

	assert.NotEmpty(t, branch.ID)
	assert.Len(t, branch.ID, 36) // UUID format
}

func TestBranch_IsValidEstablishmentType_Valid(t *testing.T) {
	tests := []struct {
		name          string
		establishment string
		expected      bool
	}{
		{"Workshop", EstablishmentTypeWorkshop, true},
		{"Store", EstablishmentTypeStore, true},
		{"WorkshopStore", EstablishmentTypeWorkshopStore, true},
		{"Invalid", "INVALID_TYPE", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branch := &Branch{EstablishmentType: tt.establishment}
			result := branch.IsValidEstablishmentType()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidEstablishmentType_Standalone(t *testing.T) {
	assert.True(t, IsValidEstablishmentType(EstablishmentTypeWorkshop))
	assert.True(t, IsValidEstablishmentType(EstablishmentTypeStore))
	assert.True(t, IsValidEstablishmentType(EstablishmentTypeWorkshopStore))
	assert.False(t, IsValidEstablishmentType("INVALID"))
	assert.False(t, IsValidEstablishmentType(""))
}

func TestBranch_ToLogger(t *testing.T) {
	branch := &Branch{
		ID:                "branch-123",
		Name:              "Taller Norte",
		EstablishmentType: EstablishmentTypeWorkshop,
		RepresentativeID:  "rep-456",
	}

	result := branch.ToLogger()

	assert.Len(t, result, 4)
	assert.Contains(t, result[0], "id:branch-123")
	assert.Contains(t, result[1], "name:Taller Norte")
	assert.Contains(t, result[2], "type:WORKSHOP")
	assert.Contains(t, result[3], "representative_id:rep-456")
}

func TestGetAllEstablishmentTypes(t *testing.T) {
	types := GetAllEstablishmentTypes()

	assert.Len(t, types, 3)
	assert.Equal(t, EstablishmentTypeWorkshop, types[0].Code)
	assert.Equal(t, EstablishmentTypeStore, types[1].Code)
	assert.Equal(t, EstablishmentTypeWorkshopStore, types[2].Code)
}

func TestGetEstablishmentTypeLabel(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{EstablishmentTypeWorkshop, "Taller"},
		{EstablishmentTypeStore, "Tienda"},
		{EstablishmentTypeWorkshopStore, "Taller y Tienda"},
		{"UNKNOWN", "UNKNOWN"}, // Falls back to code
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := GetEstablishmentTypeLabel(tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}
