package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/EstebanGitPro/motogo-backend/tools/idencoder"
	"github.com/stretchr/testify/assert"
)

// Helper to create a test encoder
func createBrandTestEncoder() *idencoder.HashidsEncoder {
	config := idencoder.Config{
		Secret:    "test-salt-for-brand-tests",
		MinLength: 8,
	}
	encoder, _ := idencoder.NewHashidsEncoder(config, logger.NewSlogLogger())
	return encoder
}

func TestNewBrandListResponse_Success(t *testing.T) {
	// Arrange
	encoder := createBrandTestEncoder()
	brands := []domain.Brand{
		{ID: "f6a7b8c9-1111-4000-8000-000000000001", Name: "Honda"},
		{ID: "f6a7b8c9-2222-4000-8000-000000000002", Name: "Yamaha"},
		{ID: "f6a7b8c9-3333-4000-8000-000000000003", Name: "Suzuki"},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/brands", Method: "GET"},
	}

	// Act
	response := handlers.NewBrandListResponse(brands, links, encoder)

	// Assert
	assert.Len(t, response.Brands, 3)
	// IDs should be encoded (not the original UUIDs)
	assert.NotEqual(t, brands[0].ID, response.Brands[0].ID)
	assert.NotEqual(t, brands[1].ID, response.Brands[1].ID)
	assert.NotEqual(t, brands[2].ID, response.Brands[2].ID)
	// Names should match
	assert.Equal(t, "Honda", response.Brands[0].Name)
	assert.Equal(t, "Yamaha", response.Brands[1].Name)
	assert.Equal(t, "Suzuki", response.Brands[2].Name)
	assert.Len(t, response.Links, 1)
	assert.Equal(t, "self", response.Links[0].Rel)
}

func TestNewBrandListResponse_Empty(t *testing.T) {
	// Arrange
	encoder := createBrandTestEncoder()
	brands := []domain.Brand{}
	links := []handlers.Link{
		{Rel: "self", Href: "/brands", Method: "GET"},
	}

	// Act
	response := handlers.NewBrandListResponse(brands, links, encoder)

	// Assert
	assert.Empty(t, response.Brands)
	assert.Len(t, response.Links, 1)
}

func TestNewBrandListResponse_NoLinks(t *testing.T) {
	// Arrange
	encoder := createBrandTestEncoder()
	brands := []domain.Brand{
		{ID: "f6a7b8c9-4444-4000-8000-000000000004", Name: "Honda"},
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewBrandListResponse(brands, links, encoder)

	// Assert
	assert.Len(t, response.Brands, 1)
	assert.Empty(t, response.Links)
}

func TestNewBrandListResponse_SingleBrand(t *testing.T) {
	// Arrange
	encoder := createBrandTestEncoder()
	brands := []domain.Brand{
		{ID: "f6a7b8c9-5555-4000-8000-000000000005", Name: "Honda"},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/brands", Method: "GET"},
		{Rel: "references", Href: "/references", Method: "GET"},
	}

	// Act
	response := handlers.NewBrandListResponse(brands, links, encoder)

	// Assert
	assert.Len(t, response.Brands, 1)
	assert.NotEqual(t, brands[0].ID, response.Brands[0].ID) // Should be encoded
	assert.Equal(t, "Honda", response.Brands[0].Name)
	assert.Len(t, response.Links, 2)
}
