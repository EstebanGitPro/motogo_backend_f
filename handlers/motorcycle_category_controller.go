package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetMotorcycleCategories handles GET /motorcycle-categories - lists distinct motorcycle categories (HU41)
// @Summary List motorcycle categories
// @Description Retrieves all distinct motorcycle categories with line counts. Public endpoint (no authentication required). Returns HATEOAS links for category drill-down (Richardson Maturity Level 3).
// @Tags Catalogs
// @Produce json
// @Success 200 {object} StandardResponse{data=CategoryListResponse} "Categories retrieved successfully"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycle-categories [get]
func (h *handler) GetMotorcycleCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleCatControllerRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		categories, err := h.MotorcycleInteractor.GetDistinctCategories(c.Request.Context())
		if err != nil {
			log.Error(logger.LogMotorcycleCatControllerError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		baseURL := GetBaseURL(c)
		response := NewCategoryListResponse(categories, baseURL)

		log.Success(logger.LogMotorcycleCatControllerSuccess,
			"count", len(categories),
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMotorcycleCategoriesRetrieved, response)
	}
}

// GetCategoryLines handles GET /motorcycle-categories/:categoryName/lines - lists lines for a category (HU41)
// @Summary List motorcycle lines by category
// @Description Retrieves all motorcycle lines (models) for a specific category with brand and engine displacement. Public endpoint (no authentication required). Returns HATEOAS links (Richardson Maturity Level 3).
// @Tags Catalogs
// @Produce json
// @Param categoryName path string true "Category name (e.g., Sport, Scooter, Urban)"
// @Success 200 {object} StandardResponse{data=CategoryLinesResponse} "Lines retrieved successfully"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycle-categories/{categoryName}/lines [get]
func (h *handler) GetCategoryLines() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		categoryName := c.Param("categoryName")

		log.Info(logger.LogMotorcycleCatLinesControllerRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"category", categoryName,
			"client_ip", c.ClientIP())

		lines, err := h.MotorcycleInteractor.GetLinesByCategory(c.Request.Context(), categoryName)
		if err != nil {
			log.Error(logger.LogMotorcycleCatLinesControllerError, "error", err, "category", categoryName, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		baseURL := GetBaseURL(c)
		response := NewCategoryLinesResponse(categoryName, lines, baseURL)

		log.Success(logger.LogMotorcycleCatLinesControllerSuccess,
			"category", categoryName,
			"count", len(lines),
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgCategoryLinesRetrieved, response)
	}
}

// GetEngineDisplacements handles GET /engine-displacements - lists distinct engine displacement ranges (HU49)
// @Summary List engine displacement ranges
// @Description Retrieves all distinct engine displacement values with reference counts. Public endpoint (no authentication required). Returns HATEOAS links (Richardson Maturity Level 3).
// @Tags Catalogs
// @Produce json
// @Success 200 {object} StandardResponse{data=DisplacementListResponse} "Displacements retrieved successfully"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /engine-displacements [get]
func (h *handler) GetEngineDisplacements() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleDispControllerRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		displacements, err := h.MotorcycleInteractor.GetDistinctDisplacements(c.Request.Context())
		if err != nil {
			log.Error(logger.LogMotorcycleDispControllerError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		baseURL := GetBaseURL(c)
		response := NewDisplacementListResponse(displacements, baseURL)

		log.Success(logger.LogMotorcycleDispControllerSuccess,
			"count", len(displacements),
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgEngineDisplacementsRetrieved, response)
	}
}

// GetRatingRanges handles GET /rating-ranges - lists valid rating range values (HU48)
// @Summary List rating ranges
// @Description Retrieves all valid rating values (1-5) with labels. Public endpoint (no authentication required). Returns HATEOAS links (Richardson Maturity Level 3).
// @Tags Catalogs
// @Produce json
// @Success 200 {object} StandardResponse{data=RatingRangeListResponse} "Rating ranges retrieved successfully"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /rating-ranges [get]
func (h *handler) GetRatingRanges() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogRatingRangeControllerRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		ranges, err := h.MotorcycleInteractor.GetRatingRanges(c.Request.Context())
		if err != nil {
			log.Error(logger.LogRatingRangeControllerError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		baseURL := GetBaseURL(c)
		response := NewRatingRangeListResponse(ranges, baseURL)

		log.Success(logger.LogRatingRangeControllerSuccess,
			"count", len(ranges),
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgRatingRangesRetrieved, response)
	}
}
