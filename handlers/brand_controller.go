package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetBrands handles GET /brands - retrieves all brands from the catalog
// @Summary Get all motorcycle brands
// @Description Retrieves the complete list of motorcycle brands available in the system
// @Tags Brands
// @Produce json
// @Success 200 {object} StandardResponse{data=BrandListResponse}
// @Failure 500 {object} StandardResponse
// @Router /brands [get]
func (h *handler) GetBrands() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBrandInteractorGetAll,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Call interactor to get all brands
		brands, err := h.BrandInteractor.GetAllBrands(c.Request.Context())
		if err != nil {
			log.Error(logger.LogBrandInteractorGetAllError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildBrandListLinks(baseURL)

		// Build response DTO
		response := NewBrandListResponse(brands, links, h.IDEncoder)

		log.Success(logger.LogBrandInteractorGetAllOK,
			"brands_count", len(brands),
			"client_ip", c.ClientIP())

		// Send success response
		h.Response.SuccessWithData(c, domain.MsgBrandsRetrieved, response)
	}
}
