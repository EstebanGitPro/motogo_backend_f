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

// ============================================
// RegisterBranchRequest.Sanitize Tests
// ============================================

func TestRegisterBranchRequest_Sanitize(t *testing.T) {
	// Arrange
	franchiseID := "  franchise-123  "
	req := &handlers.RegisterBranchRequest{
		Name:              "  Taller Express  ",
		EstablishmentType: "  WORKSHOP  ",
		FranchiseID:       &franchiseID,
		Brands:            []string{"  Honda  ", "  Yamaha  "},
		Location: handlers.LocationDTO{
			DepartmentID:   "  dept-1  ",
			CityID:         "  city-1  ",
			CityName:       "  Medellín  ",
			DepartmentName: "  Antioquia  ",
			Address:        "  Calle 10 #20-30  ",
		},
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "Taller Express", req.Name)
	assert.Equal(t, "WORKSHOP", req.EstablishmentType)
	assert.Equal(t, "franchise-123", *req.FranchiseID)
	assert.Equal(t, "Honda", req.Brands[0])
	assert.Equal(t, "Yamaha", req.Brands[1])
	assert.Equal(t, "dept-1", req.Location.DepartmentID)
	assert.Equal(t, "city-1", req.Location.CityID)
	assert.Equal(t, "Medellín", req.Location.CityName)
	assert.Equal(t, "Antioquia", req.Location.DepartmentName)
	assert.Equal(t, "Calle 10 #20-30", req.Location.Address)
}

func TestRegisterBranchRequest_Sanitize_NilFields(t *testing.T) {
	// Arrange
	req := &handlers.RegisterBranchRequest{
		Name:              "Taller",
		EstablishmentType: "STORE",
		FranchiseID:       nil,
		Location: handlers.LocationDTO{
			DepartmentID: "dept-1",
			CityID:       "city-1",
			Address:      "Calle 1",
		},
	}

	// Act
	req.Sanitize()

	// Assert - should not panic
	assert.Equal(t, "Taller", req.Name)
	assert.Nil(t, req.FranchiseID)
}

// ============================================
// RegisterBranchRequest.ToDomain Tests
// ============================================

func TestRegisterBranchRequest_ToDomain_Complete(t *testing.T) {
	// Arrange
	franchiseID := "franchise-123"
	profileURL := "https://example.com/image.jpg"
	lat := 6.2518
	lng := -75.5636
	req := &handlers.RegisterBranchRequest{
		Name:              "Taller Express",
		EstablishmentType: "WORKSHOP",
		FranchiseID:       &franchiseID,
		ProfileImageURL:   &profileURL,
		Brands:            []string{"Honda", "Yamaha"},
		Location: handlers.LocationDTO{
			DepartmentID:   "dept-antioquia",
			CityID:         "city-medellin",
			CityName:       "Medellín",
			DepartmentName: "Antioquia",
			Address:        "Calle 10 #20-30",
			Latitude:       &lat,
			Longitude:      &lng,
		},
	}

	// Act
	result := req.ToDomain("rep-123")

	// Assert
	assert.Equal(t, "rep-123", result.RepresentativeID)
	assert.Equal(t, "Taller Express", result.Name)
	assert.Equal(t, "WORKSHOP", result.EstablishmentType)
	assert.Equal(t, "franchise-123", *result.FranchiseID)
	assert.Equal(t, profileURL, *result.ProfileImageURL)
	assert.Len(t, result.Brands, 2)
	assert.NotNil(t, result.Location)
	assert.Equal(t, "Medellín", result.Location.CityName)
	assert.Equal(t, 6.2518, *result.Location.Latitude)
	assert.Equal(t, -75.5636, *result.Location.Longitude)
}

func TestRegisterBranchRequest_ToDomain_MinimalFields(t *testing.T) {
	// Arrange
	req := &handlers.RegisterBranchRequest{
		Name:              "Store Básico",
		EstablishmentType: "STORE",
		Location: handlers.LocationDTO{
			DepartmentID: "dept-1",
			CityID:       "city-1",
			Address:      "Calle 1",
		},
	}

	// Act
	result := req.ToDomain("rep-456")

	// Assert
	assert.Equal(t, "rep-456", result.RepresentativeID)
	assert.Equal(t, "Store Básico", result.Name)
	assert.Nil(t, result.FranchiseID)
	assert.Nil(t, result.ProfileImageURL)
	assert.Empty(t, result.Brands)
	assert.NotNil(t, result.Location)
	assert.Nil(t, result.Location.Latitude)
}

// ============================================
// LocationDTO.Sanitize Tests
// ============================================

func TestLocationDTO_Sanitize(t *testing.T) {
	// Arrange
	lat := 6.2518
	lng := -75.5636
	loc := &handlers.LocationDTO{
		DepartmentID:   "  dept-1  ",
		CityID:         "  city-1  ",
		CityName:       "\n Medellín \t",
		DepartmentName: "  Antioquia  ",
		Address:        "  Calle 10 #20-30  ",
		Latitude:       &lat,
		Longitude:      &lng,
	}

	// Act
	loc.Sanitize()

	// Assert
	assert.Equal(t, "dept-1", loc.DepartmentID)
	assert.Equal(t, "city-1", loc.CityID)
	assert.Equal(t, "Medellín", loc.CityName)
	assert.Equal(t, "Antioquia", loc.DepartmentName)
	assert.Equal(t, "Calle 10 #20-30", loc.Address)
	assert.Equal(t, 6.2518, *loc.Latitude)
}

// ============================================
// NewBranchResponse Tests
// ============================================

func TestNewBranchResponse_WithAllFields(t *testing.T) {
	// Arrange
	franchiseID := "franchise-123"
	profileURL := "https://example.com/image.jpg"
	lat := 6.2518
	lng := -75.5636
	branch := &domain.Branch{
		ID:                "branch-123",
		Name:              "Taller Express",
		EstablishmentType: "WORKSHOP",
		Status:            "ACTIVE",
		FranchiseID:       &franchiseID,
		ProfileImageURL:   &profileURL,
		Brands:            []string{"Honda", "Yamaha"},
		Location: &domain.Location{
			DepartmentID: "dept-1",
			CityID:       "city-1",
			Address:      "Calle 10 #20-30",
			Latitude:     &lat,
			Longitude:    &lng,
		},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/branches/enc-123", Method: "GET"},
	}

	// Act
	response := handlers.NewBranchResponse(branch, "enc-123", handlers.GeocodingStatusSuccess, links)

	// Assert
	assert.Equal(t, "enc-123", response.ID)
	assert.Equal(t, "Taller Express", response.Name)
	assert.Equal(t, "WORKSHOP", response.EstablishmentType)
	assert.Equal(t, "Taller", response.EstablishmentTypeLabel)
	assert.Equal(t, "ACTIVE", response.Status)
	assert.Equal(t, "franchise-123", *response.FranchiseID)
	assert.Equal(t, profileURL, *response.ProfileImageURL)
	assert.Len(t, response.Brands, 2)
	assert.Equal(t, handlers.GeocodingStatusSuccess, response.GeocodingStatus)
	assert.NotNil(t, response.Location)
	assert.Equal(t, 6.2518, *response.Location.Latitude)
	assert.Len(t, response.Links, 1)
}

func TestNewBranchResponse_WithoutLocation(t *testing.T) {
	// Arrange
	branch := &domain.Branch{
		ID:                "branch-456",
		Name:              "Store Sin Ubicación",
		EstablishmentType: "STORE",
		Status:            "PENDING",
		Location:          nil,
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewBranchResponse(branch, "enc-456", handlers.GeocodingStatusFailed, links)

	// Assert
	assert.Equal(t, "enc-456", response.ID)
	assert.Equal(t, "Tienda", response.EstablishmentTypeLabel)
	assert.Equal(t, handlers.GeocodingStatusFailed, response.GeocodingStatus)
	assert.Nil(t, response.Location)
}

// ============================================
// NewBranchListItemResponse Tests
// ============================================

func TestNewBranchListItemResponse_WithFranchise(t *testing.T) {
	// Arrange
	franchiseID := "franchise-123"
	branch := domain.Branch{
		ID:                "branch-123",
		Name:              "Taller Afiliado",
		EstablishmentType: "WORKSHOP",
		Status:            "ACTIVE",
		FranchiseID:       &franchiseID,
	}
	encodedFranchiseID := "enc-franchise"
	links := []handlers.Link{
		{Rel: "self", Href: "/branches/enc-123", Method: "GET"},
	}

	// Act
	response := handlers.NewBranchListItemResponse(branch, "enc-123", &encodedFranchiseID, links)

	// Assert
	assert.Equal(t, "enc-123", response.ID)
	assert.Equal(t, "Taller Afiliado", response.Name)
	assert.Equal(t, "enc-franchise", *response.FranchiseID)
	assert.Len(t, response.Links, 1)
}

func TestNewBranchListItemResponse_WithoutFranchise(t *testing.T) {
	// Arrange
	lat := 4.711
	lng := -74.072
	branch := domain.Branch{
		ID:                "branch-456",
		Name:              "Taller Independiente",
		EstablishmentType: "WORKSHOP",
		Status:            "ACTIVE",
		FranchiseID:       nil,
		Location: &domain.Location{
			DepartmentID: "dept-cundinamarca",
			CityID:       "city-bogota",
			Address:      "Carrera 7 #45-67",
			Latitude:     &lat,
			Longitude:    &lng,
		},
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewBranchListItemResponse(branch, "enc-456", nil, links)

	// Assert
	assert.Equal(t, "enc-456", response.ID)
	assert.Nil(t, response.FranchiseID)
	assert.NotNil(t, response.Location)
	assert.Equal(t, "Carrera 7 #45-67", response.Location.Address)
}

// ============================================
// NewNearbyBranchResponse Tests (HU89)
// ============================================

func TestNewNearbyBranchResponse_Success(t *testing.T) {
	// Arrange
	profileURL := "https://example.com/taller.jpg"
	branch := domain.NearbyBranch{
		ID:                "branch-789",
		Name:              "Taller Cercano",
		EstablishmentType: "WORKSHOP",
		ProfileImageURL:   &profileURL,
		DistanceKm:        1.5,
		Location: &domain.NearbyLocation{
			Address:        "Calle 50 #30-20",
			CityName:       "Medellín",
			DepartmentName: "Antioquia",
			Latitude:       6.2518,
			Longitude:      -75.5636,
		},
	}

	// Act
	response := handlers.NewNearbyBranchResponse(branch, "enc-789", "https://api.motogo.com")

	// Assert
	assert.Equal(t, "enc-789", response.ID)
	assert.Equal(t, "Taller Cercano", response.Name)
	assert.Equal(t, "WORKSHOP", response.EstablishmentType)
	assert.Equal(t, "Taller", response.EstablishmentTypeLabel)
	assert.Equal(t, profileURL, *response.ProfileImageURL)
	assert.Equal(t, 1.5, response.DistanceKm)
	assert.Equal(t, "Calle 50 #30-20", response.Address)
	assert.Equal(t, "Medellín", response.CityName)
	assert.Equal(t, "Antioquia", response.DepartmentName)
	assert.Equal(t, 6.2518, response.Latitude)
	assert.Equal(t, -75.5636, response.Longitude)
	assert.NotEmpty(t, response.Links)
}

func TestNewNearbyBranchResponse_WithoutLocation(t *testing.T) {
	// Arrange
	branch := domain.NearbyBranch{
		ID:                "branch-no-loc",
		Name:              "Taller Sin Ubicación",
		EstablishmentType: "STORE",
		DistanceKm:        2.0,
		Location:          nil,
	}

	// Act
	response := handlers.NewNearbyBranchResponse(branch, "enc-no-loc", "https://api.motogo.com")

	// Assert
	assert.Equal(t, "enc-no-loc", response.ID)
	assert.Equal(t, "Tienda", response.EstablishmentTypeLabel)
	assert.Equal(t, 2.0, response.DistanceKm)
	assert.Empty(t, response.Address)
	assert.Empty(t, response.CityName)
	assert.Equal(t, float64(0), response.Latitude)
}

// ============================================
// GeocodingStatus Constants Tests
// ============================================

func TestGeocodingStatus_Constants(t *testing.T) {
	assert.Equal(t, handlers.GeocodingStatus("SUCCESS"), handlers.GeocodingStatusSuccess)
	assert.Equal(t, handlers.GeocodingStatus("FAILED"), handlers.GeocodingStatusFailed)
	assert.Equal(t, handlers.GeocodingStatus("SKIPPED"), handlers.GeocodingStatusSkipped)
}
