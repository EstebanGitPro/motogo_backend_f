package handlers

import (
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/gin-gonic/gin"
)

// GeocodingTestRequest is the DTO for testing geocoding
type GeocodingTestRequest struct {
	Address        string `json:"address" binding:"required"`
	CityName       string `json:"city_name" binding:"required"`
	DepartmentName string `json:"department_name" binding:"required"`
}

// Sanitize trims whitespace from all string fields
func (r *GeocodingTestRequest) Sanitize() {
	r.Address = TrimString(r.Address)
	r.CityName = TrimString(r.CityName)
	r.DepartmentName = TrimString(r.DepartmentName)
}

// GeocodingTestResponse is the response for geocoding test
type GeocodingTestResponse struct {
	Geocoded         bool    `json:"geocoded"`
	Latitude         float64 `json:"latitude,omitempty"`
	Longitude        float64 `json:"longitude,omitempty"`
	FormattedAddress string  `json:"formatted_address,omitempty"`
	Confidence       float64 `json:"confidence,omitempty"`
	Error            string  `json:"error,omitempty"`
}

// TestGeocoding handles POST /location/geocode
// @Summary Test geocoding service
// @Description Test address geocoding without creating a branch
// @Tags Dev Tools
// @Accept json
// @Produce json
// @Param request body GeocodingTestRequest true "Address to geocode"
// @Success 200 {object} middleware.APIResponse{data=GeocodingTestResponse}
// @Failure 400 {object} middleware.APIResponse
// @Router /location/geocode [post]
func (h *handler) TestGeocoding() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info("Geocoding test request received",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		var req GeocodingTestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn("Invalid geocoding test request", "error", err)
			c.JSON(http.StatusBadRequest, middleware.APIResponse{
				Success: false,
				Code:    "ERR_INVALID_REQUEST",
				Message: "Campos requeridos: address, city_name, department_name",
			})
			return
		}

		// Sanitize input
		req.Sanitize()

		// Create a location for geocoding
		location := &domain.Location{
			Address:        req.Address,
			CityName:       req.CityName,
			DepartmentName: req.DepartmentName,
		}

		// Attempt geocoding via BranchInteractor
		geocoded, err := h.BranchInteractor.GeocodeLocation(c.Request.Context(), location)

		response := GeocodingTestResponse{
			Geocoded: geocoded,
		}

		if geocoded && location.Latitude != nil && location.Longitude != nil {
			response.Latitude = *location.Latitude
			response.Longitude = *location.Longitude
			response.FormattedAddress = req.Address + ", " + req.CityName + ", " + req.DepartmentName + ", Colombia"
			response.Confidence = 0.8 // OpenCage confidence aproximado
			log.Success("Geocoding test successful",
				"latitude", response.Latitude,
				"longitude", response.Longitude)
		} else {
			response.Error = "No se pudieron obtener las coordenadas"
			if err != nil {
				response.Error = err.Error()
			}
			log.Warn("Geocoding test failed", "error", response.Error)
		}

		c.JSON(http.StatusOK, middleware.APIResponse{
			Success: true,
			Code:    "GEOCODING_TEST_COMPLETE",
			Message: "Prueba de geocodificación completada",
			Data:    response,
		})
	}
}
