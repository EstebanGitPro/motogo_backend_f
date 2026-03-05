package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/stretchr/testify/assert"
)

// helper to create RatingInteractor with fresh mock service
func setupRatingInteractor() (*interactor.RatingInteractor, *mocks.MockRatingService) {
	svc := new(mocks.MockRatingService)
	ri := interactor.NewRatingInteractor(svc)
	return ri, svc
}

// ============================================
// RateServiceItem — success
// ============================================

func TestRateServiceItem_Success(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()
	mockTx := new(mocks.MockTx)

	comment := "Great service!"
	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil, // not yet rated
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("RateServiceItem", ctx, mockTx, "item-1", 5, &comment).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	err := ri.RateServiceItem(ctx, "item-1", 5, &comment)

	assert.NoError(t, err)
	svc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestRateServiceItem_Success_NoComment(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()
	mockTx := new(mocks.MockTx)

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("RateServiceItem", ctx, mockTx, "item-1", 4, (*string)(nil)).Return(nil)
	mockTx.On("Commit").Return(nil)
	mockTx.On("Rollback").Return(nil).Maybe()

	err := ri.RateServiceItem(ctx, "item-1", 4, nil)

	assert.NoError(t, err)
	svc.AssertExpectations(t)
}

// ============================================
// RateServiceItem — validation errors
// ============================================

func TestRateServiceItem_ItemNotFound(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	svc.On("GetItemByID", ctx, "item-bad").Return(nil, errors.New("not found"))

	err := ri.RateServiceItem(ctx, "item-bad", 5, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceItemNotFound, err)
}

func TestRateServiceItem_AlreadyRated(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	existingRating := 3
	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             &existingRating,
	}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)

	err := ri.RateServiceItem(ctx, "item-1", 5, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceItemAlreadyRated, err)
}

func TestRateServiceItem_ServiceNotFinalized(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusInProgress}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)

	err := ri.RateServiceItem(ctx, "item-1", 5, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceNotFinalized, err)
}

func TestRateServiceItem_ParentServiceNotFound(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-bad",
		Rating:             nil,
	}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-bad").Return(nil, errors.New("not found"))

	err := ri.RateServiceItem(ctx, "item-1", 5, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrServiceItemNotFound, err)
}

func TestRateServiceItem_InvalidRating_TooLow(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)

	err := ri.RateServiceItem(ctx, "item-1", 0, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidRatingRange, err)
}

func TestRateServiceItem_InvalidRating_TooHigh(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)

	err := ri.RateServiceItem(ctx, "item-1", 6, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidRatingRange, err)
}

// ============================================
// RateServiceItem — transaction errors
// ============================================

func TestRateServiceItem_BeginTxError(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(nil, errors.New("tx error"))

	err := ri.RateServiceItem(ctx, "item-1", 5, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrRatingCannotSave, err)
}

func TestRateServiceItem_SaveError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()
	mockTx := new(mocks.MockTx)

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("RateServiceItem", ctx, mockTx, "item-1", 5, (*string)(nil)).Return(errors.New("save error"))
	mockTx.On("Rollback").Return(nil)

	err := ri.RateServiceItem(ctx, "item-1", 5, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrRatingCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

func TestRateServiceItem_CommitError_RollsBack(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()
	mockTx := new(mocks.MockTx)

	item := &domain.CompletedServiceItem{
		ID:                 "item-1",
		CompletedServiceID: "cs-1",
		Rating:             nil,
	}
	cs := &domain.CompletedService{ID: "cs-1", Status: domain.ServiceStatusCompleted}

	svc.On("GetItemByID", ctx, "item-1").Return(item, nil)
	svc.On("GetCompletedServiceByID", ctx, "cs-1").Return(cs, nil)
	svc.On("BeginTx", ctx).Return(mockTx, nil)
	svc.On("RateServiceItem", ctx, mockTx, "item-1", 5, (*string)(nil)).Return(nil)
	mockTx.On("Commit").Return(errors.New("commit error"))
	mockTx.On("Rollback").Return(nil)

	err := ri.RateServiceItem(ctx, "item-1", 5, nil)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrRatingCannotSave, err)
	mockTx.AssertCalled(t, "Rollback")
}

// ============================================
// GetServiceReviews — success
// ============================================

func TestGetServiceReviews_Success(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	comment := "Excellent"
	model := "Honda CB 160F"
	summary := &domain.ServiceReviewSummary{
		ServiceID:     "svc-1",
		ServiceName:   "Oil Change",
		AverageRating: 4.5,
		TotalReviews:  2,
		Breakdown:     map[int]int{5: 1, 4: 1, 3: 0, 2: 0, 1: 0},
		Reviews: []domain.ServiceReview{
			{ReviewerName: "Carlos M.", Rating: 5, Comment: &comment, MotorcycleModel: &model},
			{ReviewerName: "Ana P.", Rating: 4, Comment: nil, MotorcycleModel: nil},
		},
	}

	svc.On("GetReviewsByServiceID", ctx, "svc-1").Return(summary, nil)

	result, err := ri.GetServiceReviews(ctx, "svc-1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "svc-1", result.ServiceID)
	assert.Equal(t, 4.5, result.AverageRating)
	assert.Equal(t, 2, result.TotalReviews)
	assert.Len(t, result.Reviews, 2)
	svc.AssertExpectations(t)
}

func TestGetServiceReviews_ServiceError(t *testing.T) {
	ctx := context.Background()
	ri, svc := setupRatingInteractor()

	svc.On("GetReviewsByServiceID", ctx, "svc-bad").Return(nil, errors.New("db error"))

	result, err := ri.GetServiceReviews(ctx, "svc-bad")

	assert.Error(t, err)
	assert.Nil(t, result)
	svc.AssertExpectations(t)
}
