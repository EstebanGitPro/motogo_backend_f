package services

import (
	"context"
	"errors"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ─────────────────────────────────────────────
// Lightweight mocks — only the methods used by ratingService
// ─────────────────────────────────────────────

type stubRatingRepo struct{ mock.Mock }

func (m *stubRatingRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(output.Tx), args.Error(1)
}
func (m *stubRatingRepo) RateServiceItem(ctx context.Context, tx output.Tx, itemID string, rating int, comment *string, isOffensive bool) error {
	args := m.Called(ctx, tx, itemID, rating, comment, isOffensive)
	return args.Error(0)
}
func (m *stubRatingRepo) GetItemByID(ctx context.Context, itemID string) (*domain.CompletedServiceItem, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CompletedServiceItem), args.Error(1)
}
func (m *stubRatingRepo) GetReviewsByServiceID(ctx context.Context, serviceID string) (*domain.ServiceReviewSummary, error) {
	args := m.Called(ctx, serviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ServiceReviewSummary), args.Error(1)
}

type stubTx struct{ mock.Mock }

func (t *stubTx) Commit() error   { return t.Called().Error(0) }
func (t *stubTx) Rollback() error { return t.Called().Error(0) }

// ─────────────────────────────────────────────
// Tests — NewRatingService verified via coverage (constructor is trivial)
// Tests for each method via direct struct construction
// ─────────────────────────────────────────────

func TestRatingService_BeginTx_Success(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	mockTx := new(stubTx)
	rr.On("BeginTx", mock.Anything).Return(mockTx, nil)

	tx, err := svc.BeginTx(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, mockTx, tx)
}

func TestRatingService_BeginTx_Error(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	rr.On("BeginTx", mock.Anything).Return(nil, errors.New("tx fail"))

	tx, err := svc.BeginTx(context.Background())
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestRatingService_RateServiceItem_WithoutComment(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	mockTx := new(stubTx)
	rr.On("RateServiceItem", mock.Anything, mockTx, "item-1", 5, (*string)(nil), false).Return(nil)

	err := svc.RateServiceItem(context.Background(), mockTx, "item-1", 5, nil)
	assert.NoError(t, err)
	rr.AssertCalled(t, "RateServiceItem", mock.Anything, mockTx, "item-1", 5, (*string)(nil), false)
}

func TestRatingService_RateServiceItem_WithComment(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	mockTx := new(stubTx)
	comment := "Buen servicio"
	rr.On("RateServiceItem", mock.Anything, mockTx, "item-1", 4, &comment, false).Return(nil)

	err := svc.RateServiceItem(context.Background(), mockTx, "item-1", 4, &comment)
	assert.NoError(t, err)
}

func TestRatingService_RateServiceItem_EmptyComment(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	mockTx := new(stubTx)
	empty := ""
	rr.On("RateServiceItem", mock.Anything, mockTx, "item-1", 3, &empty, false).Return(nil)

	err := svc.RateServiceItem(context.Background(), mockTx, "item-1", 3, &empty)
	assert.NoError(t, err)
}

func TestRatingService_RateServiceItem_Error(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	mockTx := new(stubTx)
	rr.On("RateServiceItem", mock.Anything, mockTx, "item-1", 5, (*string)(nil), false).Return(errors.New("repo err"))

	err := svc.RateServiceItem(context.Background(), mockTx, "item-1", 5, nil)
	assert.Error(t, err)
}

func TestRatingService_GetItemByID_Success(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	item := &domain.CompletedServiceItem{ID: "item-1"}
	rr.On("GetItemByID", mock.Anything, "item-1").Return(item, nil)

	result, err := svc.GetItemByID(context.Background(), "item-1")
	assert.NoError(t, err)
	assert.Equal(t, item, result)
}

func TestRatingService_GetItemByID_Error(t *testing.T) {
	rr := new(stubRatingRepo)
	svc := &ratingService{ratingRepo: rr}

	rr.On("GetItemByID", mock.Anything, "item-bad").Return(nil, errors.New("not found"))

	result, err := svc.GetItemByID(context.Background(), "item-bad")
	assert.Error(t, err)
	assert.Nil(t, result)
}
