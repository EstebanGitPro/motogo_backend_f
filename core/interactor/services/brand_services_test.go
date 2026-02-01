package services

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBrandRepository is a mock for BrandRepository
type MockBrandRepository struct {
	mock.Mock
}

func (m *MockBrandRepository) GetAllBrands(ctx context.Context) ([]domain.Brand, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Brand), args.Error(1)
}

func (m *MockBrandRepository) ValidateBrandIDs(ctx context.Context, brandIDs []string) error {
	args := m.Called(ctx, brandIDs)
	return args.Error(0)
}

// ============================================
// NewBrandService Tests
// ============================================

func TestNewBrandService(t *testing.T) {
	mockRepo := new(MockBrandRepository)
	service := NewBrandService(mockRepo)
	assert.NotNil(t, service)
}

// ============================================
// GetAllBrands Tests
// ============================================

func TestGetAllBrands_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockBrandRepository)
	service := NewBrandService(mockRepo)

	expectedBrands := []domain.Brand{
		{ID: "brand-1", Name: "Honda"},
		{ID: "brand-2", Name: "Yamaha"},
		{ID: "brand-3", Name: "Suzuki"},
	}

	mockRepo.On("GetAllBrands", ctx).Return(expectedBrands, nil)

	// Act
	brands, err := service.GetAllBrands(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, brands, 3)
	assert.Equal(t, "Honda", brands[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAllBrands_Error(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockBrandRepository)
	service := NewBrandService(mockRepo)

	dbError := errors.New("database connection failed")
	mockRepo.On("GetAllBrands", ctx).Return(nil, dbError)

	// Act
	brands, err := service.GetAllBrands(ctx)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, brands)
	assert.Equal(t, dbError, err)

	mockRepo.AssertExpectations(t)
}

// ============================================
// ValidateBrandIDs Tests
// ============================================

func TestValidateBrandIDs_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockBrandRepository)
	service := NewBrandService(mockRepo)

	brandIDs := []string{"brand-1", "brand-2"}
	mockRepo.On("ValidateBrandIDs", ctx, brandIDs).Return(nil)

	// Act
	err := service.ValidateBrandIDs(ctx, brandIDs)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestValidateBrandIDs_EmptyList(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockBrandRepository)
	service := NewBrandService(mockRepo)

	// Act - empty list should return nil without calling repository
	err := service.ValidateBrandIDs(ctx, []string{})

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertNotCalled(t, "ValidateBrandIDs")
}

func TestValidateBrandIDs_InvalidBrand(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockRepo := new(MockBrandRepository)
	service := NewBrandService(mockRepo)

	brandIDs := []string{"invalid-brand"}
	mockRepo.On("ValidateBrandIDs", ctx, brandIDs).Return(domain.ErrBrandNotFound)

	// Act
	err := service.ValidateBrandIDs(ctx, brandIDs)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrBrandNotFound, err)

	mockRepo.AssertExpectations(t)
}
