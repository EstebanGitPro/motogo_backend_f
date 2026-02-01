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

		// Encode department IDs before building response
		encodedDepts := make([]DepartmentResponse, 0, len(departments))
		for _, dept := range departments {
			encodedID, err := h.EncodeID(dept.ID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "department_id", dept.ID, "error", err)
				continue // Skip departments with encoding errors
			}
			encodedDepts = append(encodedDepts, DepartmentResponse{
				ID:   encodedID,
				Name: dept.Name,
			})
		}

		links := BuildDepartmentListLinks()
		response := DepartmentListResponse{
			Departments: encodedDepts,
			Links:       links,
		}

		log.Success(logger.LogLocationControllerGetDeptsOK,
			"count", len(encodedDepts),
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

		encodedDeptID := c.Param("id")

		log.Info(logger.LogLocationControllerGetCities,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"encoded_department_id", encodedDeptID,
			"client_ip", c.ClientIP())

		// Decode department ID from URL
		departmentID, err := h.DecodeID(encodedDeptID)
		if err != nil {
			log.Warn(logger.LogIDDecodeError, "encoded_id", encodedDeptID, "error", err)
			h.Response.Error(c, domain.MsgValIDInvalid)
			return
		}

		cities, err := h.LocationInteractor.GetCitiesByDepartment(c.Request.Context(), departmentID)
		if err != nil {
			log.Error(logger.LogLocationControllerGetCitiesErr, "error", err, "department_id", departmentID)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Encode city IDs before building response
		encodedCities := make([]CityResponse, 0, len(cities))
		for _, city := range cities {
			encodedCityID, err := h.EncodeID(city.ID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "city_id", city.ID, "error", err)
				continue // Skip cities with encoding errors
			}
			encodedCities = append(encodedCities, CityResponse{
				ID:   encodedCityID,
				Name: city.Name,
			})
		}

		// Use encoded department ID for HATEOAS links
		links := BuildCityListLinks(encodedDeptID)
		response := CityListResponse{
			Cities: encodedCities,
			Links:  links,
		}

		log.Success(logger.LogLocationControllerGetCitiesOK,
			"count", len(encodedCities),
			"department_id", departmentID,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgCitiesRetrieved, response)
	}
}
