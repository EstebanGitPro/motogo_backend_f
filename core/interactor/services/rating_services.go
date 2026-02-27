package services

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/tools/moderation"
)

type ratingService struct {
	ratingRepo output.RatingRepository
	csRepo     output.CompletedServiceRepository // read-only: parent status validation
}

// NewRatingService creates a new RatingService instance
func NewRatingService(
	ratingRepo output.RatingRepository,
	csRepo output.CompletedServiceRepository,
) *ratingService {
	return &ratingService{
		ratingRepo: ratingRepo,
		csRepo:     csRepo,
	}
}

// BeginTx starts a new database transaction via the rating repository
func (s *ratingService) BeginTx(ctx context.Context) (output.Tx, error) {
	return s.ratingRepo.BeginTx(ctx)
}

// RateServiceItem rates a service item with moderation check (RELEASE_14 / HU48)
func (s *ratingService) RateServiceItem(ctx context.Context, tx output.Tx, itemID string, rating int, comment *string) error {
	isOffensive := false
	if comment != nil && *comment != "" {
		isOffensive = moderation.IsOffensive(*comment)
	}
	if err := s.ratingRepo.RateServiceItem(ctx, tx, itemID, rating, comment, isOffensive); err != nil {
		return err
	}
	if isOffensive {
		return domain.ErrOffensiveCommentFiltered
	}
	return nil
}

// GetItemByID retrieves a single completed service item by ID (RELEASE_14 / HU48)
func (s *ratingService) GetItemByID(ctx context.Context, itemID string) (*domain.CompletedServiceItem, error) {
	return s.ratingRepo.GetItemByID(ctx, itemID)
}

// GetCompletedServiceByID retrieves a completed service for parent status validation
func (s *ratingService) GetCompletedServiceByID(ctx context.Context, serviceID string) (*domain.CompletedService, error) {
	return s.csRepo.GetByID(ctx, serviceID)
}

// GetReviewsByServiceID retrieves aggregated reviews for a service type (RELEASE_14 / HU48)
func (s *ratingService) GetReviewsByServiceID(ctx context.Context, serviceID string) (*domain.ServiceReviewSummary, error) {
	return s.ratingRepo.GetReviewsByServiceID(ctx, serviceID)
}
