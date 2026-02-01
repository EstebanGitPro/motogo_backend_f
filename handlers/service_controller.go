package handlers

import (
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetServiceTypes handles GET /service-types - retrieves all service types (HU75)
// @Summary Get all service types
// @Description Retrieves the complete list of service types available in the system
// @Tags Services
// @Produce json
// @Success 200 {object} StandardResponse{data=ServiceTypeListResponse}
// @Failure 500 {object} StandardResponse
// @Router /service-types [get]
func (h *handler) GetServiceTypes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogServiceControllerGetTypes,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Call interactor to get all service types
		types := h.ServiceInteractor.GetServiceTypes(c.Request.Context())

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildServiceTypeListLinks(baseURL)

		// Build response DTO
		response := NewServiceTypeListResponse(types, links)

		log.Success(logger.LogServiceControllerGetTypesOK,
			"types_count", len(types),
			"client_ip", c.ClientIP())

		// Send success response
		h.Response.SuccessWithData(c, domain.MsgServiceTypesRetrieved, response)
	}
}

// GetServices handles GET /services - retrieves all services from catalog (HU63)
// @Summary Get all services
// @Description Retrieves the complete list of services available in the catalog, optionally filtered by type
// @Tags Services
// @Produce json
// @Param type query string false "Filter by service type (e.g., Maintenance, Repair)"
// @Success 200 {object} StandardResponse{data=ServiceListResponse}
// @Failure 500 {object} StandardResponse
// @Router /services [get]
func (h *handler) GetServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogServiceControllerGetAll,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Check for type filter
		filterType := c.Query("type")

		var services []domain.Service
		var err error

		if filterType != "" {
			// Validate service type
			if !domain.IsValidServiceType(filterType) {
				log.Warn(logger.LogServiceControllerInvalidType,
					"type", filterType,
					"client_ip", c.ClientIP())
				h.Response.Error(c, domain.MsgServiceInvalidType)
				return
			}
			// Filter by type
			services, err = h.ServiceInteractor.GetServicesByType(c.Request.Context(), domain.ServiceType(filterType))
		} else {
			// Get all services
			services, err = h.ServiceInteractor.GetAllServices(c.Request.Context())
		}

		if err != nil {
			log.Error(logger.LogServiceControllerGetAllError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildServiceListLinks(baseURL, filterType)

		// Build response DTO with encoded IDs
		response, encErr := NewServiceListResponseWithEncoder(services, links, h.IDEncoder)
		if encErr != nil {
			log.Error(logger.LogMessageIDEncodeError, "error", encErr)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		log.Success(logger.LogServiceControllerGetAllOK,
			"services_count", len(services),
			"filter_type", filterType,
			"client_ip", c.ClientIP())

		// Send success response
		h.Response.SuccessWithData(c, domain.MsgServicesRetrieved, response)
	}
}

// GetBranchServices handles GET /branches/:id/services - retrieves services for a specific branch
// @Summary Get services for a branch
// @Description Retrieves all services associated with a specific branch, including when they were added
// @Tags Branches
// @Produce json
// @Param id path string true "Branch ID"
// @Success 200 {object} StandardResponse{data=BranchServiceListResponse}
// @Failure 400 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /branches/{id}/services [get]
func (h *handler) GetBranchServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		branchID := c.Param("id")

		log.Info(logger.LogBranchServicesControllerGet,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"branch_id", branchID,
			"client_ip", c.ClientIP())

		// Decode branch ID if encoded
		decodedID, err := h.IDEncoder.Decode(branchID)
		if err != nil {
			log.Warn(logger.LogBranchServicesControllerInvalidID, "id", branchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// Get services for this branch
		services, err := h.ServiceInteractor.GetServicesByBranch(c.Request.Context(), decodedID)
		if err != nil {
			log.Error(logger.LogBranchServicesControllerGetError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := []Link{
			{Rel: "self", Href: baseURL + "/branches/" + branchID + "/services", Method: "GET"},
			{Rel: "branch", Href: baseURL + "/branches/" + branchID, Method: "GET"},
		}

		// Build response with encoded service IDs
		response, err := NewBranchServiceListResponseWithEncoder(services, links, h.IDEncoder)
		if err != nil {
			log.Error(logger.LogBranchServicesControllerGetError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		log.Success(logger.LogBranchServicesControllerGetOK,
			"branch_id", branchID,
			"services_count", len(services),
			"client_ip", c.ClientIP())

		// Send success response
		h.Response.SuccessWithData(c, domain.MsgServicesRetrieved, response)
	}
}

// AssociateBranchServicesRequest represents the request body for POST /branches/:id/services
type AssociateBranchServicesRequest struct {
	ServiceIDs []string `json:"service_ids" binding:"required,min=1"`
}

// AssociateBranchServicesResponse represents the response for POST /branches/:id/services
type AssociateBranchServicesResponse struct {
	AddedCount int    `json:"added_count"`
	Links      []Link `json:"_links"`
}

// AssociateBranchServices handles POST /branches/:id/services - associates services to a branch
// @Summary Associate services to a branch
// @Description Associates one or more services from the catalog to a specific branch
// @Tags Branches
// @Accept json
// @Produce json
// @Param id path string true "Branch ID"
// @Param body body AssociateBranchServicesRequest true "Service IDs to associate"
// @Success 201 {object} StandardResponse{data=AssociateBranchServicesResponse}
// @Failure 400 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 409 {object} StandardResponse
// @Router /branches/{id}/services [post]
func (h *handler) AssociateBranchServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		branchID := c.Param("id")
		log.Info(logger.LogBranchServicesControllerAssociate, "branch_id", branchID)

		// Decode branch ID
		decodedBranchID, err := h.IDEncoder.Decode(branchID)
		if err != nil {
			log.Warn(logger.LogBranchServicesControllerInvalidID, "id", branchID)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// Parse request body
		var req AssociateBranchServicesRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogBranchServicesControllerInvalidBody, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Decode all service IDs
		decodedServiceIDs := make([]string, len(req.ServiceIDs))
		for i, id := range req.ServiceIDs {
			decoded, err := h.IDEncoder.Decode(id)
			if err != nil {
				log.Warn(logger.LogBranchServicesControllerInvalidSvcID, "id", id)
				h.Response.Error(c, domain.MsgServiceNotFound)
				return
			}
			decodedServiceIDs[i] = decoded
		}

		// Validate service IDs exist
		if err := h.ServiceInteractor.ValidateServiceIDs(c.Request.Context(), decodedServiceIDs); err != nil {
			log.Warn(logger.LogBranchServicesControllerInvalidSvcs, "error", err)
			h.Response.Error(c, domain.MsgServiceNotFound)
			return
		}

		// Start transaction and associate
		tx, err := h.ServiceInteractor.BeginTx(c.Request.Context())
		if err != nil {
			log.Error(logger.LogBranchServicesControllerTxError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		err = h.ServiceInteractor.AssociateBranchServices(c.Request.Context(), tx, decodedBranchID, decodedServiceIDs)
		if err != nil {
			tx.Rollback()
			log.Error(logger.LogBranchServicesControllerAssociateErr, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			log.Error(logger.LogBranchServicesControllerCommitError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Build response
		baseURL := GetBaseURL(c)
		response := AssociateBranchServicesResponse{
			AddedCount: len(decodedServiceIDs),
			Links: []Link{
				{Rel: "branch_services", Href: baseURL + "/branches/" + branchID + "/services", Method: "GET"},
			},
		}

		log.Success(logger.LogBranchServicesControllerAssociateOK, "branch_id", branchID, "added_count", len(decodedServiceIDs))
		h.Response.SuccessWithData(c, domain.MsgServiceAssociated, response)
	}
}

// DissociateBranchService handles DELETE /branches/:id/services/:serviceId - removes a service from a branch
// @Summary Dissociate a service from a branch
// @Description Removes a service association from a specific branch
// @Tags Branches
// @Produce json
// @Param id path string true "Branch ID"
// @Param serviceId path string true "Service ID"
// @Success 200 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id}/services/{serviceId} [delete]
func (h *handler) DissociateBranchService() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		branchID := c.Param("id")
		serviceID := c.Param("serviceId")
		log.Info(logger.LogBranchServicesControllerDissociate, "branch_id", branchID, "service_id", serviceID)

		// Decode IDs
		decodedBranchID, err := h.IDEncoder.Decode(branchID)
		if err != nil {
			log.Warn(logger.LogBranchServicesControllerInvalidID, "id", branchID)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		decodedServiceID, err := h.IDEncoder.Decode(serviceID)
		if err != nil {
			log.Warn(logger.LogBranchServicesControllerInvalidSvcID, "id", serviceID)
			h.Response.Error(c, domain.MsgServiceNotFound)
			return
		}

		// Start transaction and dissociate
		tx, err := h.ServiceInteractor.BeginTx(c.Request.Context())
		if err != nil {
			log.Error(logger.LogBranchServicesControllerTxError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		err = h.ServiceInteractor.DissociateBranchService(c.Request.Context(), tx, decodedBranchID, decodedServiceID)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, domain.ErrServiceNotFound) {
				log.Warn(logger.LogBranchServicesControllerNotFound, "branch_id", branchID, "service_id", serviceID)
				h.Response.Error(c, domain.MsgServiceNotFound)
				return
			}
			log.Error(logger.LogBranchServicesControllerDisassocErr, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			log.Error(logger.LogBranchServicesControllerCommitError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		log.Success(logger.LogBranchServicesControllerDissociateOK, "branch_id", branchID, "service_id", serviceID)
		h.Response.Success(c, domain.MsgServiceDissociated)
	}
}

// UpdateService handles PUT /admin/services/:id - updates a service in the global catalog (HU68 - Admin only)
// @Summary Update a service
// @Description Updates an existing service in the global catalog. Requires ADMIN role.
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param body body UpdateServiceRequest true "Service data to update"
// @Success 200 {object} StandardResponse{data=ServiceDetailResponse}
// @Failure 400 {object} StandardResponse
// @Failure 401 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /admin/services/{id} [put]
func (h *handler) UpdateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		serviceID := c.Param("id")
		log.Info(logger.LogServiceControllerUpdate, "service_id", serviceID)

		// Decode service ID
		decodedID, err := h.IDEncoder.Decode(serviceID)
		if err != nil {
			log.Warn(logger.LogServiceControllerUpdateError, "id", serviceID, "error", err)
			h.Response.Error(c, domain.MsgServiceNotFound)
			return
		}

		// Parse request body
		var req UpdateServiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogServiceControllerUpdateError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Validate service type
		if !domain.IsValidServiceType(req.ServiceType) {
			log.Warn(logger.LogServiceControllerInvalidType, "type", req.ServiceType)
			h.Response.Error(c, domain.MsgServiceTypeInvalid)
			return
		}

		// Verify service exists
		existingService, err := h.ServiceInteractor.GetServiceByID(c.Request.Context(), decodedID)
		if err != nil {
			if errors.Is(err, domain.ErrServiceNotFound) {
				log.Warn(logger.LogServiceControllerUpdateError, "service_id", serviceID, "error", "not found")
				h.Response.Error(c, domain.MsgServiceResNotFound)
				return
			}
			log.Error(logger.LogServiceControllerUpdateError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Update service fields
		existingService.Name = req.Name
		existingService.Description = req.Description
		existingService.ServiceType = domain.ServiceType(req.ServiceType)

		// Update activation status if provided
		if req.IsActive != nil {
			existingService.IsActive = *req.IsActive
		}

		// Perform update
		if err := h.ServiceInteractor.UpdateService(c.Request.Context(), *existingService); err != nil {
			if errors.Is(err, domain.ErrServiceNotFound) {
				log.Warn(logger.LogServiceControllerUpdateError, "service_id", serviceID, "error", "update failed - not found")
				h.Response.Error(c, domain.MsgServiceResNotFound)
				return
			}
			log.Error(logger.LogServiceControllerUpdateError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := []Link{
			{Rel: "self", Href: baseURL + "/admin/services/" + serviceID, Method: "PUT"},
			{Rel: "services", Href: baseURL + "/services", Method: "GET"},
		}

		// Build response with encoded ID
		response, encErr := NewServiceDetailResponseWithEncoder(existingService, links, h.IDEncoder)
		if encErr != nil {
			log.Error(logger.LogMessageIDEncodeError, "error", encErr)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		log.Success(logger.LogServiceControllerUpdateOK, "service_id", serviceID)
		h.Response.SuccessWithData(c, domain.MsgServiceUpdated, response)
	}
}
