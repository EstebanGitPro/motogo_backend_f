package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetDepartments handles GET /departments
// @Summary Get all departments
// @Description Retrieves all Colombian departments for location selection
// @Tags Catalogs
// @Produce json
// @Success 200 {object} StandardResponse{data=DepartmentListResponse}
// @Failure 500 {object} StandardResponse
// @Router /departments [get]
func (h *handler) GetDepartments() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogLocationControllerGetDepts,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		departments, err := h.LocationInteractor.GetAllDepartments(c.Request.Context())
		if err != nil {
			log.Error(logger.LogLocationControllerGetDeptsError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		links := BuildDepartmentListLinks()
		response := NewDepartmentListResponse(departments, links)

		log.Success(logger.LogLocationControllerGetDeptsOK,
			"count", len(departments),
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgDepartmentsRetrieved, response)
	}
}

// GetCitiesByDepartment handles GET /departments/:id/cities
// @Summary Get cities by department
// @Description Retrieves all cities for a specific department
// @Tags Catalogs
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} StandardResponse{data=CityListResponse}
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /departments/{id}/cities [get]
func (h *handler) GetCitiesByDepartment() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		departmentID := c.Param("id")

		log.Info(logger.LogLocationControllerGetCities,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"department_id", departmentID,
			"client_ip", c.ClientIP())

		cities, err := h.LocationInteractor.GetCitiesByDepartment(c.Request.Context(), departmentID)
		if err != nil {
			log.Error(logger.LogLocationControllerGetCitiesErr, "error", err, "department_id", departmentID)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		links := BuildCityListLinks(departmentID)
		response := NewCityListResponse(cities, links)

		log.Success(logger.LogLocationControllerGetCitiesOK,
			"count", len(cities),
			"department_id", departmentID,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgCitiesRetrieved, response)
	}
}
