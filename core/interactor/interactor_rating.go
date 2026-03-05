package interactor

import (
	"context"
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/input"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// RatingInteractor handles rating use cases (RELEASE_14 / HU48)
type RatingInteractor struct {
	service input.RatingService
}

// NewRatingInteractor creates a new RatingInteractor
func NewRatingInteractor(service input.RatingService) *RatingInteractor {
	return &RatingInteractor{service: service}
}

// RateServiceItem orchestrates the rating of a completed service item (RELEASE_14 / HU48)
func (i *RatingInteractor) RateServiceItem(ctx context.Context, itemID string, rating int, comment *string) error {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorRateStart, "item_id", itemID)

	// STEP 1: Verify item exists
	item, err := i.service.GetItemByID(ctx, itemID)
	if err != nil {
		log.Warn(logger.LogCSInteractorRateNotFound, "item_id", itemID, "error", err)
		return domain.ErrServiceItemNotFound
	}

	// STEP 2: Verify item has not been rated yet
	if item.Rating != nil {
		log.Warn(logger.LogCSInteractorRateAlready, "item_id", itemID)
		return domain.ErrServiceItemAlreadyRated
	}

	// STEP 3: Verify parent service is FINALIZADO
	cs, csErr := i.service.GetCompletedServiceByID(ctx, item.CompletedServiceID)
	if csErr != nil {
		log.Warn(logger.LogCSInteractorRateNotFound, "completed_service_id", item.CompletedServiceID, "error", csErr)
		return domain.ErrServiceItemNotFound
	}
	if cs.Status != domain.ServiceStatusCompleted {
		log.Warn(logger.LogCSInteractorRateNotFinal, "status", cs.Status, "item_id", itemID)
		return domain.ErrServiceNotFinalized
	}

	// STEP 4: Validate rating range (1-5)
	if rating < 1 || rating > 5 {
		log.Warn(logger.LogCSInteractorRateInvalid, "rating", rating, "item_id", itemID)
		return domain.ErrInvalidRatingRange
	}

	// STEP 5: Begin transaction
	tx, err := i.service.BeginTx(ctx)
	if err != nil {
		log.Error(logger.LogCSInteractorRateTxErr, "error", err)
		return domain.ErrRatingCannotSave
	}

	defer func() { _ = tx.Rollback() }()

	// STEP 6: Delegate to service (handles moderation + persistence)
	wasOffensive := false
	if err = i.service.RateServiceItem(ctx, tx, itemID, rating, comment); err != nil {
		// Offensive sentinel is NOT a real error — the data was saved successfully
		if errors.Is(err, domain.ErrOffensiveCommentFiltered) {
			wasOffensive = true
		} else {
			log.Error(logger.LogCSInteractorRateUpdErr, "item_id", itemID, "error", err)
			return domain.ErrRatingCannotSave
		}
	}

	// STEP 7: Commit
	if err = tx.Commit(); err != nil {
		log.Error(logger.LogCSInteractorRateCommErr, "error", err)
		return domain.ErrRatingCannotSave
	}

	log.Success(logger.LogCSInteractorRateSuccess, "item_id", itemID)

	// Propagate offensive sentinel to handler (if applicable)
	if wasOffensive {
		return domain.ErrOffensiveCommentFiltered
	}
	return nil
}

// GetServiceReviews retrieves aggregated reviews for a service type (RELEASE_14 / HU48)
func (i *RatingInteractor) GetServiceReviews(ctx context.Context, serviceID string) (*domain.ServiceReviewSummary, error) {
	traceID := middleware.GetTraceIDFromContext(ctx)
	log := log.WithTraceID(traceID)

	log.Info(logger.LogCSInteractorReviewsStart, "service_id", serviceID)

	summary, err := i.service.GetReviewsByServiceID(ctx, serviceID)
	if err != nil {
		log.Error(logger.LogCSInteractorReviewsError, "service_id", serviceID, "error", err)
		return nil, err
	}

	log.Success(logger.LogCSInteractorReviewsSuccess, "service_id", serviceID, "total_reviews", summary.TotalReviews)
	return summary, nil
}
