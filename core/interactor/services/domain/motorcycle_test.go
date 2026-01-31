package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// Motorcycle Tests
// ============================================

func TestMotorcycle_Structure(t *testing.T) {
	year := 2022
	mileage := 15000
	notes := "Well maintained"

	motorcycle := Motorcycle{
		ID:             "moto-001",
		LicensePlate:   "ABC123",
		ReferenceID:    "ref-001",
		OwnerID:        "owner-001",
		Year:           &year,
		CurrentMileage: &mileage,
		OwnerNotes:     &notes,
	}

	assert.Equal(t, "moto-001", motorcycle.ID)
	assert.Equal(t, "ABC123", motorcycle.LicensePlate)
	assert.Equal(t, "ref-001", motorcycle.ReferenceID)
	assert.Equal(t, "owner-001", motorcycle.OwnerID)
	assert.NotNil(t, motorcycle.Year)
	assert.Equal(t, 2022, *motorcycle.Year)
	assert.NotNil(t, motorcycle.CurrentMileage)
	assert.Equal(t, 15000, *motorcycle.CurrentMileage)
	assert.NotNil(t, motorcycle.OwnerNotes)
	assert.Equal(t, "Well maintained", *motorcycle.OwnerNotes)
}

func TestMotorcycle_NullableFields(t *testing.T) {
	motorcycle := Motorcycle{
		ID:           "moto-002",
		LicensePlate: "XYZ789",
		OwnerID:      "owner-002",
	}

	assert.Empty(t, motorcycle.ReferenceID)
	assert.Nil(t, motorcycle.Year)
	assert.Nil(t, motorcycle.CurrentMileage)
	assert.Nil(t, motorcycle.OwnerNotes)
	assert.Nil(t, motorcycle.Reference)
}

func TestMotorcycle_WithReference(t *testing.T) {
	ref := &MotorcycleReference{
		ID:                 "ref-001",
		BrandID:            "brand-001",
		BrandName:          "Yamaha",
		Model:              "MT-07",
		Category:           "Naked",
		EngineDisplacement: 689,
	}

	motorcycle := Motorcycle{
		ID:           "moto-003",
		LicensePlate: "REF123",
		OwnerID:      "owner-003",
		ReferenceID:  "ref-001",
		Reference:    ref,
	}

	assert.NotNil(t, motorcycle.Reference)
	assert.Equal(t, "Yamaha", motorcycle.Reference.BrandName)
	assert.Equal(t, "MT-07", motorcycle.Reference.Model)
	assert.Equal(t, 689, motorcycle.Reference.EngineDisplacement)
}

// ============================================
// MotorcycleReference Tests
// ============================================

func TestMotorcycleReference_Structure(t *testing.T) {
	ref := MotorcycleReference{
		ID:                 "ref-001",
		BrandID:            "brand-001",
		BrandName:          "Honda",
		Model:              "CB 190R",
		Category:           "Sport",
		EngineDisplacement: 184,
	}

	assert.Equal(t, "ref-001", ref.ID)
	assert.Equal(t, "brand-001", ref.BrandID)
	assert.Equal(t, "Honda", ref.BrandName)
	assert.Equal(t, "CB 190R", ref.Model)
	assert.Equal(t, "Sport", ref.Category)
	assert.Equal(t, 184, ref.EngineDisplacement)
}
