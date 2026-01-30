package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// Helper functions
// ============================================

func intPtrMoto(i int) *int {
	return &i
}

func stringPtrMoto(s string) *string {
	return &s
}

// ============================================
// ToMotorcycleResponse Tests
// ============================================

func TestToMotorcycleResponse_WithAllFields(t *testing.T) {
	// Arrange
	motorcycle := &domain.Motorcycle{
		ID:             "moto-123",
		LicensePlate:   "ABC123",
		Year:           intPtrMoto(2023),
		CurrentMileage: intPtrMoto(5000),
		OwnerNotes:     stringPtrMoto("Siempre con aceite sintético"),
		Reference: &domain.MotorcycleReference{
			ID:                 "ref-honda",
			BrandID:            "brand-honda",
			BrandName:          "Honda",
			Model:              "CB 190R",
			Category:           "Sport",
			EngineDisplacement: 190,
		},
	}

	// Act
	response := handlers.ToMotorcycleResponse(motorcycle)

	// Assert
	assert.Equal(t, "moto-123", response.ID)
	assert.Equal(t, "ABC123", response.LicensePlate)
	assert.Equal(t, 2023, *response.Year)
	assert.Equal(t, 5000, *response.CurrentMileage)
	assert.Equal(t, "Siempre con aceite sintético", *response.OwnerNotes)
	assert.NotNil(t, response.Reference)
	assert.Equal(t, "ref-honda", response.Reference.ID)
	assert.Equal(t, "brand-honda", response.Reference.BrandID)
	assert.Equal(t, "Honda", response.Reference.BrandName)
	assert.Equal(t, "CB 190R", response.Reference.Model)
	assert.Equal(t, "Sport", response.Reference.Category)
	assert.Equal(t, 190, response.Reference.EngineDisplacement)
}

func TestToMotorcycleResponse_WithoutOptionalFields(t *testing.T) {
	// Arrange
	motorcycle := &domain.Motorcycle{
		ID:           "moto-456",
		LicensePlate: "XYZ789",
	}

	// Act
	response := handlers.ToMotorcycleResponse(motorcycle)

	// Assert
	assert.Equal(t, "moto-456", response.ID)
	assert.Equal(t, "XYZ789", response.LicensePlate)
	assert.Nil(t, response.Year)
	assert.Nil(t, response.CurrentMileage)
	assert.Nil(t, response.OwnerNotes)
	assert.Nil(t, response.Reference)
}

func TestToMotorcycleResponse_WithReferenceWithoutOptionalFields(t *testing.T) {
	// Arrange
	motorcycle := &domain.Motorcycle{
		ID:           "moto-789",
		LicensePlate: "DEF456",
		Year:         intPtrMoto(2020),
		Reference: &domain.MotorcycleReference{
			ID:        "ref-yamaha",
			BrandName: "Yamaha",
			Model:     "FZ 25",
		},
	}

	// Act
	response := handlers.ToMotorcycleResponse(motorcycle)

	// Assert
	assert.Equal(t, "moto-789", response.ID)
	assert.Equal(t, 2020, *response.Year)
	assert.Nil(t, response.CurrentMileage)
	assert.NotNil(t, response.Reference)
	assert.Equal(t, "FZ 25", response.Reference.Model)
}

// ============================================
// ToMotorcycleLookupResponse Tests (Workshop View - HU47)
// ============================================

func TestToMotorcycleLookupResponse_ExcludesOwnerNotes(t *testing.T) {
	// Arrange - motorcycle with owner notes
	motorcycle := &domain.Motorcycle{
		ID:             "moto-123",
		LicensePlate:   "ABC123",
		Year:           intPtrMoto(2023),
		CurrentMileage: intPtrMoto(5000),
		OwnerNotes:     stringPtrMoto("Notas privadas del dueño"),
		Reference: &domain.MotorcycleReference{
			ID:                 "ref-honda",
			BrandID:            "brand-honda",
			BrandName:          "Honda",
			Model:              "CB 190R",
			Category:           "Sport",
			EngineDisplacement: 190,
		},
	}

	// Act
	response := handlers.ToMotorcycleLookupResponse(motorcycle)

	// Assert - owner_notes should NOT be in the lookup response (it's private)
	assert.Equal(t, "moto-123", response.ID)
	assert.Equal(t, "ABC123", response.LicensePlate)
	assert.Equal(t, 2023, *response.Year)
	assert.Equal(t, 5000, *response.CurrentMileage)
	// Reference should NOT include brand_id (not needed for workshops)
	assert.NotNil(t, response.Reference)
	assert.Equal(t, "Honda", response.Reference.BrandName)
	assert.Equal(t, "CB 190R", response.Reference.Model)
	assert.Equal(t, "Sport", response.Reference.Category)
	assert.Equal(t, 190, response.Reference.EngineDisplacement)
}

func TestToMotorcycleLookupResponse_WithoutReference(t *testing.T) {
	// Arrange
	motorcycle := &domain.Motorcycle{
		ID:           "moto-456",
		LicensePlate: "XYZ789",
	}

	// Act
	response := handlers.ToMotorcycleLookupResponse(motorcycle)

	// Assert
	assert.Equal(t, "moto-456", response.ID)
	assert.Nil(t, response.Reference)
}

// ============================================
// RegisterMotorcycleRequest.ToDomain Tests
// ============================================

func TestRegisterMotorcycleRequest_ToDomain_WithAllFields(t *testing.T) {
	// Arrange
	refID := "ref-honda"
	req := &handlers.RegisterMotorcycleRequest{
		LicensePlate:   "ABC123",
		ReferenceID:    &refID,
		Year:           intPtrMoto(2023),
		CurrentMileage: intPtrMoto(0),
		OwnerNotes:     stringPtrMoto("Mi primera moto"),
	}

	// Act
	result := req.ToDomain("owner-123")

	// Assert
	assert.Equal(t, "ABC123", result.LicensePlate)
	assert.Equal(t, "owner-123", result.OwnerID)
	assert.Equal(t, "ref-honda", result.ReferenceID)
	assert.Equal(t, 2023, *result.Year)
	assert.Equal(t, 0, *result.CurrentMileage)
	assert.Equal(t, "Mi primera moto", *result.OwnerNotes)
}

func TestRegisterMotorcycleRequest_ToDomain_OnlyRequired(t *testing.T) {
	// Arrange
	req := &handlers.RegisterMotorcycleRequest{
		LicensePlate: "XYZ789",
	}

	// Act
	result := req.ToDomain("owner-456")

	// Assert
	assert.Equal(t, "XYZ789", result.LicensePlate)
	assert.Equal(t, "owner-456", result.OwnerID)
	assert.Empty(t, result.ReferenceID) // No reference
	assert.Nil(t, result.Year)
	assert.Nil(t, result.CurrentMileage)
	assert.Nil(t, result.OwnerNotes)
}

// ============================================
// UpdateMotorcycleRequest.ToDomain Tests
// ============================================

func TestUpdateMotorcycleRequest_ToDomain_WithAllFields(t *testing.T) {
	// Arrange
	refID := "ref-yamaha"
	req := &handlers.UpdateMotorcycleRequest{
		ReferenceID:    &refID,
		Year:           intPtrMoto(2022),
		CurrentMileage: intPtrMoto(10000),
		OwnerNotes:     stringPtrMoto("Actualizado"),
	}

	// Act
	result := req.ToDomain()

	// Assert
	assert.Equal(t, "ref-yamaha", result.ReferenceID)
	assert.Equal(t, 2022, *result.Year)
	assert.Equal(t, 10000, *result.CurrentMileage)
	assert.Equal(t, "Actualizado", *result.OwnerNotes)
}

func TestUpdateMotorcycleRequest_ToDomain_PartialUpdate(t *testing.T) {
	// Arrange - only updating mileage
	req := &handlers.UpdateMotorcycleRequest{
		CurrentMileage: intPtrMoto(15000),
	}

	// Act
	result := req.ToDomain()

	// Assert
	assert.Empty(t, result.ReferenceID)
	assert.Nil(t, result.Year)
	assert.Equal(t, 15000, *result.CurrentMileage)
	assert.Nil(t, result.OwnerNotes)
}

// ============================================
// ToMotorcycleResponseList Tests
// ============================================

func TestToMotorcycleResponseList_Success(t *testing.T) {
	// Arrange
	motorcycles := []domain.Motorcycle{
		{ID: "moto-1", LicensePlate: "ABC001"},
		{ID: "moto-2", LicensePlate: "ABC002", Year: intPtrMoto(2021)},
		{ID: "moto-3", LicensePlate: "ABC003"},
	}

	// Act
	result := handlers.ToMotorcycleResponseList(motorcycles)

	// Assert
	assert.Len(t, result, 3)
	assert.Equal(t, "moto-1", result[0].ID)
	assert.Equal(t, "ABC001", result[0].LicensePlate)
	assert.Equal(t, "moto-2", result[1].ID)
	assert.Equal(t, 2021, *result[1].Year)
	assert.Equal(t, "moto-3", result[2].ID)
}

func TestToMotorcycleResponseList_Empty(t *testing.T) {
	// Arrange
	motorcycles := []domain.Motorcycle{}

	// Act
	result := handlers.ToMotorcycleResponseList(motorcycles)

	// Assert
	assert.Empty(t, result)
}

// ============================================
// ToMotorcycleReferenceCatalogList Tests (HU50)
// ============================================

func TestToMotorcycleReferenceCatalogList_Success(t *testing.T) {
	// Arrange
	refs := []domain.MotorcycleReference{
		{
			ID:                 "ref-1",
			BrandID:            "brand-honda",
			BrandName:          "Honda",
			Model:              "CB 190R",
			Category:           "Sport",
			EngineDisplacement: 190,
		},
		{
			ID:                 "ref-2",
			BrandID:            "brand-yamaha",
			BrandName:          "Yamaha",
			Model:              "FZ 25",
			Category:           "Naked",
			EngineDisplacement: 250,
		},
	}

	// Act
	result := handlers.ToMotorcycleReferenceCatalogList(refs)

	// Assert
	assert.Len(t, result, 2)
	assert.Equal(t, "ref-1", result[0].ID)
	assert.Equal(t, "brand-honda", result[0].BrandID)
	assert.Equal(t, "Honda", result[0].BrandName)
	assert.Equal(t, "CB 190R", result[0].Model)
	assert.Equal(t, "Sport", result[0].Category)
	assert.Equal(t, 190, result[0].EngineDisplacement)
	assert.Equal(t, "Yamaha", result[1].BrandName)
	assert.Equal(t, "FZ 25", result[1].Model)
}

func TestToMotorcycleReferenceCatalogList_Empty(t *testing.T) {
	// Arrange
	refs := []domain.MotorcycleReference{}

	// Act
	result := handlers.ToMotorcycleReferenceCatalogList(refs)

	// Assert
	assert.Empty(t, result)
}
