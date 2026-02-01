package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// CreateFranchiseRequest Tests
// ============================================

func TestCreateFranchiseRequest_ToFranchiseDomain(t *testing.T) {
	// Arrange
	description := "Test description"
	req := handlers.CreateFranchiseRequest{
		Name:        "Test Franchise",
		Description: &description,
		BranchIDs:   []string{"branch-1", "branch-2"},
	}

	// Act
	result := req.ToFranchiseDomain()

	// Assert
	assert.Equal(t, "Test Franchise", result.Name)
	assert.NotNil(t, result.Description)
	assert.Equal(t, "Test description", *result.Description)
}

func TestCreateFranchiseRequest_ToFranchiseDomain_NoDescription(t *testing.T) {
	// Arrange
	req := handlers.CreateFranchiseRequest{
		Name:        "Franchise Without Desc",
		Description: nil,
		BranchIDs:   []string{"branch-1"},
	}

	// Act
	result := req.ToFranchiseDomain()

	// Assert
	assert.Equal(t, "Franchise Without Desc", result.Name)
	assert.Nil(t, result.Description)
}

// ============================================
// UpdateFranchiseRequest Tests
// ============================================

func TestUpdateFranchiseRequest_ToFranchiseDomain(t *testing.T) {
	// Arrange
	description := "Updated description"
	req := handlers.UpdateFranchiseRequest{
		Name:        "Updated Franchise",
		Description: &description,
	}

	// Act
	result := req.ToFranchiseDomain("franchise-123")

	// Assert
	assert.Equal(t, "franchise-123", result.ID)
	assert.Equal(t, "Updated Franchise", result.Name)
	assert.NotNil(t, result.Description)
	assert.Equal(t, "Updated description", *result.Description)
}

func TestUpdateFranchiseRequest_ToFranchiseDomain_NilDescription(t *testing.T) {
	// Arrange
	req := handlers.UpdateFranchiseRequest{
		Name:        "Franchise No Desc",
		Description: nil,
	}

	// Act
	result := req.ToFranchiseDomain("franchise-456")

	// Assert
	assert.Equal(t, "franchise-456", result.ID)
	assert.Equal(t, "Franchise No Desc", result.Name)
	assert.Nil(t, result.Description)
}

// ============================================
// ToFranchiseResponse Tests
// ============================================

func TestToFranchiseResponse_Success(t *testing.T) {
	// Arrange
	description := "A great franchise"
	franchise := domain.Franchise{
		ID:          "franchise-uuid",
		Name:        "Super Franchise",
		Description: &description,
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/franchises/xyz", Method: "GET"},
	}

	// Act
	response := handlers.ToFranchiseResponse(&franchise, 5, "encoded-id", links)

	// Assert
	assert.Equal(t, "encoded-id", response.ID)
	assert.Equal(t, "Super Franchise", response.Name)
	assert.NotNil(t, response.Description)
	assert.Equal(t, "A great franchise", *response.Description)
	assert.Equal(t, 5, response.BranchCount)
	assert.Len(t, response.Links, 1)
}

func TestToFranchiseResponse_WithoutDescription(t *testing.T) {
	// Arrange
	franchise := domain.Franchise{
		ID:          "franchise-uuid",
		Name:        "Basic Franchise",
		Description: nil,
	}

	// Act
	response := handlers.ToFranchiseResponse(&franchise, 0, "encoded-id", nil)

	// Assert
	assert.Equal(t, "encoded-id", response.ID)
	assert.Equal(t, "Basic Franchise", response.Name)
	assert.Nil(t, response.Description)
	assert.Equal(t, 0, response.BranchCount)
}

// ============================================
// ToFranchiseListResponse Tests
// ============================================

func TestToFranchiseListResponse_Success(t *testing.T) {
	// Arrange
	desc1 := "First franchise"
	desc2 := "Second franchise"
	franchises := []domain.Franchise{
		{ID: "uuid-1", Name: "Franchise 1", Description: &desc1},
		{ID: "uuid-2", Name: "Franchise 2", Description: &desc2},
	}
	encodedIDs := []string{"enc-1", "enc-2"}
	baseURL := "http://localhost:8080"

	// Act
	response := handlers.ToFranchiseListResponse(franchises, encodedIDs, baseURL)

	// Assert
	assert.Equal(t, 2, response.Total)
	assert.Len(t, response.Franchises, 2)
	assert.Equal(t, "enc-1", response.Franchises[0].ID)
	assert.Equal(t, "Franchise 1", response.Franchises[0].Name)
	assert.Equal(t, "enc-2", response.Franchises[1].ID)
	assert.Equal(t, "Franchise 2", response.Franchises[1].Name)
	assert.NotEmpty(t, response.Links)
}

func TestToFranchiseListResponse_Empty(t *testing.T) {
	// Arrange
	franchises := []domain.Franchise{}
	encodedIDs := []string{}
	baseURL := "http://localhost:8080"

	// Act
	response := handlers.ToFranchiseListResponse(franchises, encodedIDs, baseURL)

	// Assert
	assert.Equal(t, 0, response.Total)
	assert.Empty(t, response.Franchises)
	assert.NotEmpty(t, response.Links)
}
