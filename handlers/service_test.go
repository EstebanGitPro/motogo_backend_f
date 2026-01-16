package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/stretchr/testify/assert"
)

// Note: createTestEncoder is defined in handler_test.go
func createServiceTestEncoder() *idencoder.HashidsEncoder {
	encoder, _ := idencoder.NewHashidsEncoder(idencoder.Config{
		Secret:    "test-secret-key-for-testing",
		MinLength: 8,
	}, nil)
	return encoder
}

func TestNewServiceListResponse_Success(t *testing.T) {
	// Arrange
	services := []domain.Service{
		{
			ID:          "uuid-1",
			Name:        "Cambio de aceite",
			Description: "Cambio de aceite de motor",
			ServiceType: domain.ServiceTypeMaintenance,
		},
		{
			ID:          "uuid-2",
			Name:        "Reparación de motor",
			Description: "Reparación completa",
			ServiceType: domain.ServiceTypeRepair,
		},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/services", Method: "GET"},
	}

	// Act
	response := handlers.NewServiceListResponse(services, links)

	// Assert
	assert.Len(t, response.Services, 2)
	assert.Equal(t, "uuid-1", response.Services[0].ID)
	assert.Equal(t, "Cambio de aceite", response.Services[0].Name)
	assert.Equal(t, "Mantenimiento", response.Services[0].ServiceType)
	assert.Len(t, response.Links, 1)
}

func TestNewServiceListResponseWithEncoder_Success(t *testing.T) {
	// Arrange
	encoder := createServiceTestEncoder()
	services := []domain.Service{
		{
			ID:          "a1234567-89ab-cdef-0123-456789abcdef",
			Name:        "Cambio de aceite",
			Description: "Cambio de aceite de motor",
			ServiceType: domain.ServiceTypeMaintenance,
		},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/services", Method: "GET"},
	}

	// Act
	response, err := handlers.NewServiceListResponseWithEncoder(services, links, encoder)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, response.Services, 1)
	assert.NotEqual(t, "a1234567-89ab-cdef-0123-456789abcdef", response.Services[0].ID) // ID should be encoded
	assert.NotEmpty(t, response.Services[0].ID)
	assert.Equal(t, "Cambio de aceite", response.Services[0].Name)
}

func TestNewServiceListResponseWithEncoder_InvalidUUID(t *testing.T) {
	// Arrange
	encoder := createServiceTestEncoder()
	services := []domain.Service{
		{
			ID:          "invalid-uuid", // Invalid UUID format
			Name:        "Test Service",
			ServiceType: domain.ServiceTypeMaintenance,
		},
	}
	links := []handlers.Link{}

	// Act
	response, err := handlers.NewServiceListResponseWithEncoder(services, links, encoder)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, response.Services)
}

func TestNewBranchServiceListResponse_Success(t *testing.T) {
	// Arrange
	services := []domain.BranchServiceInfo{
		{
			Service: domain.Service{
				ID:          "uuid-1",
				Name:        "Cambio de aceite",
				Description: "Cambio de aceite de motor",
				ServiceType: domain.ServiceTypeMaintenance,
			},
			AddedAt: "2026-01-15T10:30:00-05:00",
			Active:  true,
		},
		{
			Service: domain.Service{
				ID:          "uuid-2",
				Name:        "Reparación de motor",
				Description: "Reparación completa",
				ServiceType: domain.ServiceTypeRepair,
			},
			AddedAt: "2026-01-16T14:00:00-05:00",
			Active:  true,
		},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/branches/xyz/services", Method: "GET"},
		{Rel: "branch", Href: "/branches/xyz", Method: "GET"},
	}

	// Act
	response := handlers.NewBranchServiceListResponse(services, links)

	// Assert
	assert.Len(t, response.Services, 2)
	assert.Equal(t, "uuid-1", response.Services[0].ID)
	assert.Equal(t, "2026-01-15T10:30:00-05:00", response.Services[0].AddedAt)
	assert.True(t, response.Services[0].Active)
	assert.Len(t, response.Links, 2)
}

func TestNewBranchServiceListResponseWithEncoder_Success(t *testing.T) {
	// Arrange
	encoder := createServiceTestEncoder()
	services := []domain.BranchServiceInfo{
		{
			Service: domain.Service{
				ID:          "a1234567-89ab-cdef-0123-456789abcdef",
				Name:        "Cambio de aceite",
				Description: "Cambio de aceite de motor",
				ServiceType: domain.ServiceTypeMaintenance,
			},
			AddedAt: "2026-01-15T10:30:00-05:00",
			Active:  true,
		},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/branches/xyz/services", Method: "GET"},
	}

	// Act
	response, err := handlers.NewBranchServiceListResponseWithEncoder(services, links, encoder)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, response.Services, 1)
	assert.NotEqual(t, "a1234567-89ab-cdef-0123-456789abcdef", response.Services[0].ID) // ID should be encoded
	assert.NotEmpty(t, response.Services[0].ID)
	assert.Equal(t, "Cambio de aceite", response.Services[0].Name)
	assert.Equal(t, "2026-01-15T10:30:00-05:00", response.Services[0].AddedAt)
	assert.True(t, response.Services[0].Active)
}

func TestNewBranchServiceListResponseWithEncoder_InvalidUUID(t *testing.T) {
	// Arrange
	encoder := createServiceTestEncoder()
	services := []domain.BranchServiceInfo{
		{
			Service: domain.Service{
				ID:          "invalid-uuid", // Invalid UUID format
				Name:        "Test Service",
				ServiceType: domain.ServiceTypeMaintenance,
			},
			AddedAt: "2026-01-15T10:30:00-05:00",
			Active:  true,
		},
	}
	links := []handlers.Link{}

	// Act
	response, err := handlers.NewBranchServiceListResponseWithEncoder(services, links, encoder)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, response.Services)
}

func TestNewBranchServiceListResponseWithEncoder_EmptyList(t *testing.T) {
	// Arrange
	encoder := createServiceTestEncoder()
	services := []domain.BranchServiceInfo{}
	links := []handlers.Link{
		{Rel: "self", Href: "/branches/xyz/services", Method: "GET"},
	}

	// Act
	response, err := handlers.NewBranchServiceListResponseWithEncoder(services, links, encoder)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, response.Services)
	assert.Len(t, response.Links, 1)
}

func TestNewServiceTypeListResponse_Success(t *testing.T) {
	// Arrange
	types := []domain.ServiceType{
		domain.ServiceTypeMaintenance,
		domain.ServiceTypeRepair,
		domain.ServiceTypeTires,
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/service-types", Method: "GET"},
	}

	// Act
	response := handlers.NewServiceTypeListResponse(types, links)

	// Assert
	assert.Len(t, response.Types, 3)
	assert.Equal(t, "Mantenimiento", response.Types[0].Value)
	assert.Equal(t, "Reparación", response.Types[1].Value)
	assert.Equal(t, "Llantas", response.Types[2].Value)
	assert.Len(t, response.Links, 1)
}
