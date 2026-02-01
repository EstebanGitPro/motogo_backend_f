package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewDepartmentListResponse Tests
// ============================================

func TestNewDepartmentListResponse_Success(t *testing.T) {
	// Arrange
	departments := []domain.Department{
		{ID: "dept-1", Name: "Antioquia"},
		{ID: "dept-2", Name: "Cundinamarca"},
		{ID: "dept-3", Name: "Valle del Cauca"},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/departments", Method: "GET"},
	}

	// Act
	response := handlers.NewDepartmentListResponse(departments, links)

	// Assert
	assert.Len(t, response.Departments, 3)
	assert.Equal(t, "dept-1", response.Departments[0].ID)
	assert.Equal(t, "Antioquia", response.Departments[0].Name)
	assert.Equal(t, "dept-2", response.Departments[1].ID)
	assert.Equal(t, "Cundinamarca", response.Departments[1].Name)
	assert.Equal(t, "dept-3", response.Departments[2].ID)
	assert.Equal(t, "Valle del Cauca", response.Departments[2].Name)
	assert.Len(t, response.Links, 1)
}

func TestNewDepartmentListResponse_Empty(t *testing.T) {
	// Arrange
	departments := []domain.Department{}
	links := []handlers.Link{
		{Rel: "self", Href: "/departments", Method: "GET"},
	}

	// Act
	response := handlers.NewDepartmentListResponse(departments, links)

	// Assert
	assert.Empty(t, response.Departments)
	assert.Len(t, response.Links, 1)
}

func TestNewDepartmentListResponse_NoLinks(t *testing.T) {
	// Arrange
	departments := []domain.Department{
		{ID: "dept-1", Name: "Antioquia"},
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewDepartmentListResponse(departments, links)

	// Assert
	assert.Len(t, response.Departments, 1)
	assert.Empty(t, response.Links)
}

func TestNewDepartmentListResponse_SingleDepartment(t *testing.T) {
	// Arrange
	departments := []domain.Department{
		{ID: "dept-bogota", Name: "Bogotá D.C."},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/departments", Method: "GET"},
		{Rel: "cities", Href: "/departments/dept-bogota/cities", Method: "GET"},
	}

	// Act
	response := handlers.NewDepartmentListResponse(departments, links)

	// Assert
	assert.Len(t, response.Departments, 1)
	assert.Equal(t, "dept-bogota", response.Departments[0].ID)
	assert.Equal(t, "Bogotá D.C.", response.Departments[0].Name)
	assert.Len(t, response.Links, 2)
}

// ============================================
// NewCityListResponse Tests
// ============================================

func TestNewCityListResponse_Success(t *testing.T) {
	// Arrange
	cities := []domain.City{
		{ID: "city-1", Name: "Medellín"},
		{ID: "city-2", Name: "Envigado"},
		{ID: "city-3", Name: "Bello"},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/departments/dept-1/cities", Method: "GET"},
		{Rel: "department", Href: "/departments/dept-1", Method: "GET"},
	}

	// Act
	response := handlers.NewCityListResponse(cities, links)

	// Assert
	assert.Len(t, response.Cities, 3)
	assert.Equal(t, "city-1", response.Cities[0].ID)
	assert.Equal(t, "Medellín", response.Cities[0].Name)
	assert.Equal(t, "city-2", response.Cities[1].ID)
	assert.Equal(t, "Envigado", response.Cities[1].Name)
	assert.Equal(t, "city-3", response.Cities[2].ID)
	assert.Equal(t, "Bello", response.Cities[2].Name)
	assert.Len(t, response.Links, 2)
}

func TestNewCityListResponse_Empty(t *testing.T) {
	// Arrange
	cities := []domain.City{}
	links := []handlers.Link{
		{Rel: "self", Href: "/departments/dept-1/cities", Method: "GET"},
	}

	// Act
	response := handlers.NewCityListResponse(cities, links)

	// Assert
	assert.Empty(t, response.Cities)
	assert.Len(t, response.Links, 1)
}

func TestNewCityListResponse_NoLinks(t *testing.T) {
	// Arrange
	cities := []domain.City{
		{ID: "city-1", Name: "Bogotá"},
	}
	links := []handlers.Link{}

	// Act
	response := handlers.NewCityListResponse(cities, links)

	// Assert
	assert.Len(t, response.Cities, 1)
	assert.Empty(t, response.Links)
}

func TestNewCityListResponse_SingleCity(t *testing.T) {
	// Arrange
	cities := []domain.City{
		{ID: "city-bogota", Name: "Bogotá"},
	}
	links := []handlers.Link{
		{Rel: "self", Href: "/departments/dept-cundinamarca/cities", Method: "GET"},
	}

	// Act
	response := handlers.NewCityListResponse(cities, links)

	// Assert
	assert.Len(t, response.Cities, 1)
	assert.Equal(t, "city-bogota", response.Cities[0].ID)
	assert.Equal(t, "Bogotá", response.Cities[0].Name)
	assert.Len(t, response.Links, 1)
}
