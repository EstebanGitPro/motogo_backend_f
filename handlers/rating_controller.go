package handlers

import (
	"errors"

	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// RateServiceItem handles POST /completed-services/:id/items/:itemId/rating (RELEASE_14 / HU48)
// @Summary Rate a completed service item
// @Description Submits a rating (1-5) and optional comment for a completed service item. Comments are moderated.
// @Tags Ratings
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Completed service ID (obfuscated)"
// @Param itemId path string true "Service item ID (obfuscated)"
// @Param body body RateServiceItemRequest true "Rating data"
// @Success 200 {object} StandardResponse "Rating created"
// @Failure 400 {object} StandardResponse "Invalid rating or already rated"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Service item not found"
// @Failure 409 {object} StandardResponse "Service not finalized"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /completed-services/{id}/items/{itemId}/rating [post]
func (h *handler) RateServiceItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		Logger.Info(logger.LogCSControllerRateRequest)

		// Step 1: Decode completed service ID (validate obfuscation)
		encodedServiceID := c.Param("id")
		if _, err := h.DecodeID(encodedServiceID); err != nil {
			Logger.Warn(logger.LogCSControllerRateError, "error", err)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Step 2: Decode item ID
		encodedItemID := c.Param("itemId")
		itemID, err := h.DecodeID(encodedItemID)
		if err != nil {
			Logger.Warn(logger.LogCSControllerRateError, "error", err)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Step 3: Parse request body
		var req RateServiceItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Logger.Warn(logger.LogCSControllerRateError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Sanitize
		req.Sanitize()

		// Step 4: Delegate to rating interactor
		if err := h.RatingInteractor.RateServiceItem(
			c.Request.Context(),
			itemID,
			req.Rating,
			req.Comment,
		); err != nil {
			// Offensive sentinel: rating saved, but comment was filtered
			if errors.Is(err, domain.ErrOffensiveCommentFiltered) {
				Logger.Success(logger.LogCSControllerRateSuccess, "offensive_filtered", true)
				h.Response.Success(c, domain.MsgRatingCreatedOffensive)
				return
			}

			Logger.Error(logger.LogCSControllerRateError, "error", err)

			switch {
			case errors.Is(err, domain.ErrServiceItemNotFound):
				h.Response.Error(c, domain.MsgServiceItemNotFound)
			case errors.Is(err, domain.ErrServiceItemAlreadyRated):
				h.Response.Error(c, domain.MsgServiceItemAlreadyRated)
			case errors.Is(err, domain.ErrInvalidRatingRange):
				h.Response.Error(c, domain.MsgInvalidRatingRange)
			case errors.Is(err, domain.ErrServiceNotFinalized):
				h.Response.Error(c, domain.MsgServiceNotFinalized)
			default:
				h.Response.Error(c, domain.MsgRatingCannotSave)
			}
			return
		}

		Logger.Success(logger.LogCSControllerRateSuccess)
		h.Response.Success(c, domain.MsgRatingCreated)
	}
}

// GetServiceReviews handles GET /branches/:id/services/:serviceId/reviews (RELEASE_14 / HU48)
// @Summary Get reviews for a service type at a specific branch
// @Description Returns aggregated reviews (average, breakdown, individual reviews) for a service type scoped to a branch
// @Tags Ratings
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (obfuscated)"
// @Param serviceId path string true "Service ID (obfuscated)"
// @Success 200 {object} StandardResponse "Reviews retrieved"
// @Failure 400 {object} StandardResponse "Invalid service ID or branch ID"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /branches/{id}/services/{serviceId}/reviews [get]
func (h *handler) GetServiceReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		Logger.Info(logger.LogCSControllerReviewsRequest)

		// Step 1: Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			Logger.Warn(logger.LogCSControllerReviewsError, "error", err)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Step 2: Decode service ID
		encodedServiceID := c.Param("serviceId")
		serviceID, err := h.DecodeID(encodedServiceID)
		if err != nil {
			Logger.Warn(logger.LogCSControllerReviewsError, "error", err)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		// Step 3: Delegate to interactor
		summary, err := h.RatingInteractor.GetServiceReviews(c.Request.Context(), serviceID, branchID)
		if err != nil {
			Logger.Error(logger.LogCSControllerReviewsError, "error", err)
			h.Response.Error(c, domain.MsgRatingCannotSave)
			return
		}

		// Step 4: Map domain to DTO
		reviewDTOs := make([]ServiceReviewItemDTO, 0, len(summary.Reviews))
		for _, r := range summary.Reviews {
			reviewDTOs = append(reviewDTOs, ServiceReviewItemDTO{
				ReviewerName:    r.ReviewerName,
				Rating:          r.Rating,
				Comment:         r.Comment,
				RatedAt:         r.RatedAt.Format("2006-01-02T15:04:05Z"),
				MotorcycleModel: r.MotorcycleModel,
			})
		}

		response := ServiceReviewResponse{
			ServiceID:     summary.ServiceID,
			ServiceName:   summary.ServiceName,
			AverageRating: summary.AverageRating,
			TotalReviews:  summary.TotalReviews,
			Breakdown: ReviewBreakdownDTO{
				Star5: summary.Breakdown[5],
				Star4: summary.Breakdown[4],
				Star3: summary.Breakdown[3],
				Star2: summary.Breakdown[2],
				Star1: summary.Breakdown[1],
			},
			Reviews: reviewDTOs,
		}

		Logger.Success(logger.LogCSControllerReviewsSuccess)
		h.Response.SuccessWithData(c, domain.MsgRatingsRetrieved, response)
	}
}
