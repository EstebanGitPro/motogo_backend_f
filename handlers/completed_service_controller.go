package handlers

import (
	"errors"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// ==========================================
// DTOs (Request / Response)
// ==========================================

// CreateCompletedServiceRequest represents the POST /completed-services request body.
type CreateCompletedServiceRequest struct {
	BranchID            string   `json:"branch_id" binding:"required"`
	MotorcycleID        string   `json:"motorcycle_id" binding:"required"`
	DiagnosticID        *string  `json:"diagnostic_id,omitempty"`
	ServiceIDs          []string `json:"service_ids" binding:"required,min=1"`
	QuotedPrice         *float64 `json:"quoted_price,omitempty"`
	FinalPrice          *float64 `json:"final_price,omitempty"`
	RepresentativeNotes *string  `json:"representative_notes,omitempty"`
}

// Sanitize trims whitespace from string fields.
func (r *CreateCompletedServiceRequest) Sanitize() {
	r.BranchID = strings.TrimSpace(r.BranchID)
	r.MotorcycleID = strings.TrimSpace(r.MotorcycleID)
	if r.DiagnosticID != nil {
		trimmed := strings.TrimSpace(*r.DiagnosticID)
		r.DiagnosticID = &trimmed
	}

	if r.RepresentativeNotes != nil {
		trimmed := strings.TrimSpace(*r.RepresentativeNotes)
		r.RepresentativeNotes = &trimmed
	}
	for i := range r.ServiceIDs {
		r.ServiceIDs[i] = strings.TrimSpace(r.ServiceIDs[i])
	}
}

// CompletedServiceResponse represents the API response for a completed service.
type CompletedServiceResponse struct {
	ID                  string                         `json:"id"`
	BranchID            string                         `json:"branch_id"`
	BranchName          *string                        `json:"branch_name,omitempty"`
	MotorcycleID        string                         `json:"motorcycle_id"`
	DiagnosticID        *string                        `json:"diagnostic_id,omitempty"`
	Status              string                         `json:"status"`
	RequestDate         string                         `json:"request_date"`
	QuotedPrice         *float64                       `json:"quoted_price,omitempty"`
	FinalPrice          *float64                       `json:"final_price,omitempty"`
	RepresentativeNotes *string                        `json:"representative_notes,omitempty"`
	Services            []CompletedServiceItemResponse `json:"services,omitempty"`
}

// CompletedServiceItemResponse represents a single service item in the response.
type CompletedServiceItemResponse struct {
	ID        string `json:"id"`
	ServiceID string `json:"service_id"`
}

// ToCompletedServiceResponse converts a domain CompletedService to an API response.
func ToCompletedServiceResponse(cs *domain.CompletedService) CompletedServiceResponse {
	resp := CompletedServiceResponse{
		ID:                  cs.ID,
		BranchID:            cs.BranchID,
		BranchName:          cs.BranchName,
		MotorcycleID:        cs.MotorcycleID,
		DiagnosticID:        cs.DiagnosticID,
		Status:              string(cs.Status),
		RequestDate:         cs.RequestDate.Format(constants.DateFormat),
		QuotedPrice:         cs.QuotedPrice,
		FinalPrice:          cs.FinalPrice,
		RepresentativeNotes: cs.RepresentativeNotes,
	}

	if cs.Services != nil {
		resp.Services = make([]CompletedServiceItemResponse, len(cs.Services))
		for i, item := range cs.Services {
			resp.Services[i] = CompletedServiceItemResponse{
				ID:        item.ID,
				ServiceID: item.ServiceID,
			}
		}
	}
	return resp
}

// ==========================================
// Handler Methods
// ==========================================

// RegisterCompletedService handles POST /completed-services - register a performed service (HU64)
// @Summary Register performed service
// @Description Registers a new performed service for a motorcycle at a branch. Creates the service record, pivot items, and initial status history.
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param completedService body CreateCompletedServiceRequest true "Completed service data"
// @Success 201 {object} StandardResponse{data=CompletedServiceResponse} "Service registered successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Branch or motorcycle not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /completed-services [post]
func (h *handler) RegisterCompletedService() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogCSControllerCreateRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Parse request body
		var request CreateCompletedServiceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			Logger.Error(logger.LogCSControllerCreateError, "bind_error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Sanitize input
		request.Sanitize()

		// Step 3: Decode obfuscated IDs
		branchID, err := h.DecodeID(request.BranchID)
		if err != nil {
			Logger.Error(logger.LogCSControllerCreateError, "branch_decode_error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		motorcycleID, err := h.DecodeID(request.MotorcycleID)
		if err != nil {
			Logger.Error(logger.LogCSControllerCreateError, "motorcycle_decode_error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Decode optional diagnostic ID
		var diagnosticID *string
		if request.DiagnosticID != nil && *request.DiagnosticID != "" {
			decoded, err := h.DecodeID(*request.DiagnosticID)
			if err != nil {
				Logger.Error(logger.LogCSControllerCreateError, "diagnostic_decode_error", err)
				h.Response.Error(c, domain.MsgDiagnosticNotFound)
				return
			}
			diagnosticID = &decoded
		}

		// Decode service IDs
		decodedServiceIDs := make([]string, len(request.ServiceIDs))
		for i, encodedID := range request.ServiceIDs {
			decoded, err := h.DecodeID(encodedID)
			if err != nil {
				Logger.Error(logger.LogCSControllerCreateError, "service_decode_error", err, "index", i)
				h.Response.Error(c, domain.MsgInvalidBranchServices)
				return
			}
			decodedServiceIDs[i] = decoded
		}

		// Step 4: Build domain entity
		cs := &domain.CompletedService{
			BranchID:            branchID,
			MotorcycleID:        motorcycleID,
			DiagnosticID:        diagnosticID,
			QuotedPrice:         request.QuotedPrice,
			FinalPrice:          request.FinalPrice,
			RepresentativeNotes: request.RepresentativeNotes,
		}

		// Step 5: Register through interactor
		result, err := h.CompletedServiceInteractor.RegisterCompletedService(
			c.Request.Context(),
			cs,
			decodedServiceIDs,
			person.ID,
		)
		if err != nil {
			Logger.Error(logger.LogCSControllerCreateError, "error", err)

			switch {
			case errors.Is(err, domain.ErrInvalidBranchServices):
				h.Response.Error(c, domain.MsgInvalidBranchServices)
			case errors.Is(err, domain.ErrDiagnosticNotForMotorcycle):
				h.Response.Error(c, domain.MsgDiagnosticNotForMotorcycle)
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			case errors.Is(err, domain.ErrActiveServiceExists):
				h.Response.Error(c, domain.MsgActiveServiceExists)
			default:
				h.Response.Error(c, domain.MsgCompletedServiceCannotSave)
			}
			return
		}

		Logger.Info(logger.LogCSControllerCreateSuccess, "id", result.ID, "trace_id", traceID)

		// Step 6: Build response with encoded IDs
		response := ToCompletedServiceResponse(result)
		encodedCSID, _ := h.EncodeID(result.ID)
		response.ID = encodedCSID
		response.BranchID = request.BranchID         // return obfuscated branch ID client sent
		response.MotorcycleID = request.MotorcycleID // return obfuscated motorcycle ID client sent

		// Encode service item IDs
		for i := range response.Services {
			encodedItemID, _ := h.EncodeID(result.Services[i].ID)
			response.Services[i].ID = encodedItemID
			response.Services[i].ServiceID = request.ServiceIDs[i] // return obfuscated service IDs
		}

		h.Response.SuccessWithData(c, domain.MsgCompletedServiceCreated, response)
	}
}

// GetCompletedServicesByBranch handles GET /branches/:id/completed-services - list services for a branch (HU64)
// @Summary List completed services for a branch
// @Description Lists all completed services registered at a branch
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=[]CompletedServiceResponse} "Completed services list"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Branch not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /branches/{id}/completed-services [get]
func (h *handler) GetCompletedServicesByBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogCSControllerListRequest, "trace_id", traceID)

		// Step 1: Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			Logger.Error(logger.LogCSControllerListError, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// Step 2: Get services through interactor
		services, err := h.CompletedServiceInteractor.GetCompletedServicesByBranch(
			c.Request.Context(),
			branchID,
		)
		if err != nil {
			Logger.Error(logger.LogCSControllerListError, "error", err)
			h.Response.Error(c, domain.MsgCompletedServiceNotFound)
			return
		}

		Logger.Info(logger.LogCSControllerListSuccess, "count", len(services), "trace_id", traceID)

		// Step 3: Build response with encoded IDs
		responses := make([]CompletedServiceResponse, len(services))
		for i := range services {
			responses[i] = ToCompletedServiceResponse(&services[i])
			encodedID, _ := h.EncodeID(services[i].ID)
			responses[i].ID = encodedID
			responses[i].BranchID = encodedBranchID
			encodedMotoID, _ := h.EncodeID(services[i].MotorcycleID)
			responses[i].MotorcycleID = encodedMotoID
		}

		h.Response.SuccessWithData(c, domain.MsgCompletedServicesListed, responses)
	}
}

// GetCompletedServicesByMotorcycle handles GET /motorcycles/:id/completed-services (HU64)
// @Summary List completed services for a motorcycle
// @Description Lists all completed services associated with a motorcycle. Only the owner can access.
// @Tags Completed Services
// @Produce json
// @Security BearerAuth
// @Param id path string true "Motorcycle ID (hashid encoded)"
// @Success 200 {object} StandardResponse{data=[]CompletedServiceResponse} "Completed services list"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/completed-services [get]
func (h *handler) GetCompletedServicesByMotorcycle() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogCSControllerListByMotoReq,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated user from context
		person, _ := middleware.GetAuthenticatedUser(c)
		if person == nil {
			log.Warn(logger.LogMotorcycleControllerNoAuthUser, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// 2. Decode motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogMotorcycleControllerIDDecodeError,
				"encoded_id", encodedID,
				"error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 3. Validate ownership
		motorcycle, err := h.MotorcycleInteractor.GetMotorcycleByID(c.Request.Context(), motorcycleID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerGetError,
				"error", err,
				"motorcycle_id", motorcycleID)
			if errors.Is(err, domain.ErrMotorcycleNotFound) {
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		if person.ID != motorcycle.OwnerID {
			log.Warn(logger.LogMotorcycleControllerOwnershipDenied,
				"motorcycle_id", motorcycleID,
				"requested_by", person.ID,
				"owner_id", motorcycle.OwnerID)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 4. Get completed services
		services, err := h.CompletedServiceInteractor.GetCompletedServicesByMotorcycle(
			c.Request.Context(),
			motorcycleID,
		)
		if err != nil {
			log.Error(logger.LogCSControllerListByMotoError,
				"error", err,
				"motorcycle_id", motorcycleID)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 5. Build response with encoded IDs
		responses := make([]CompletedServiceResponse, len(services))
		for i := range services {
			responses[i] = ToCompletedServiceResponse(&services[i])
			encodedCSID, _ := h.EncodeID(services[i].ID)
			responses[i].ID = encodedCSID
			responses[i].MotorcycleID = encodedID
			encodedBranchID, _ := h.EncodeID(services[i].BranchID)
			responses[i].BranchID = encodedBranchID
		}

		log.Success(logger.LogCSControllerListByMotoSuccess,
			"motorcycle_id", motorcycleID,
			"count", len(responses),
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgCompletedServicesListed, responses)
	}
}

// DeleteCompletedService handles DELETE /completed-services/:id - delete a completed service (HU65)
// @Summary Delete a completed service
// @Description Deletes a completed service record and its associated items and status history (via DB cascade).
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Completed service ID (obfuscated)"
// @Success 200 {object} StandardResponse "Service deleted successfully"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Completed service not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /completed-services/{id} [delete]
func (h *handler) DeleteCompletedService() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogCSControllerDeleteRequest, "trace_id", traceID)

		// Step 1: Decode completed service ID
		encodedID := c.Param("id")
		serviceID, err := h.DecodeID(encodedID)
		if err != nil {
			Logger.Error(logger.LogCSControllerDeleteError, "decode_error", err)
			h.Response.Error(c, domain.MsgCompletedServiceNotFound)
			return
		}

		// Step 2: Delete through interactor
		if err := h.CompletedServiceInteractor.DeleteCompletedService(
			c.Request.Context(),
			serviceID,
		); err != nil {
			Logger.Error(logger.LogCSControllerDeleteError, "error", err)

			switch {
			case errors.Is(err, domain.ErrCompletedServiceNotFound):
				h.Response.Error(c, domain.MsgCompletedServiceNotFound)
			default:
				h.Response.Error(c, domain.MsgCompletedServiceDeleteError)
			}
			return
		}

		Logger.Info(logger.LogCSControllerDeleteSuccess, "id", serviceID, "trace_id", traceID)

		h.Response.Success(c, domain.MsgCompletedServiceDeleted)
	}
}

// StatusTransitionResponse represents a single status transition in the API response
type StatusTransitionResponse struct {
	ID             string  `json:"id"`
	PreviousStatus *string `json:"previous_status"`
	NewStatus      string  `json:"new_status"`
	CreatedBy      string  `json:"created_by"`
	CreatedAt      string  `json:"created_at"`
}

// GetCompletedServiceTransitions handles GET /completed-services/:id/transitions (HU73)
// @Summary Get status transition history for a completed service
// @Description Retrieves the full status transition history for a completed service
// @Tags Completed Services
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Completed service ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=[]StatusTransitionResponse} "Transitions list"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Completed service not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /completed-services/{id}/transitions [get]
func (h *handler) GetCompletedServiceTransitions() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogCSControllerTransRequest, "trace_id", traceID)

		// Step 1: Decode completed service ID
		encodedID := c.Param("id")
		serviceID, err := h.DecodeID(encodedID)
		if err != nil {
			Logger.Error(logger.LogCSControllerTransError, "decode_error", err)
			h.Response.Error(c, domain.MsgCompletedServiceNotFound)
			return
		}

		// Step 2: Get transitions through interactor
		history, err := h.CompletedServiceInteractor.GetStatusHistory(
			c.Request.Context(),
			serviceID,
		)
		if err != nil {
			Logger.Error(logger.LogCSControllerTransError, "error", err)
			if errors.Is(err, domain.ErrCompletedServiceNotFound) {
				h.Response.Error(c, domain.MsgCompletedServiceNotFound)
			} else {
				h.Response.Error(c, domain.MsgStatusHistoryError)
			}
			return
		}

		// Step 3: Build response
		responses := make([]StatusTransitionResponse, len(history))
		for i := range history {
			responses[i] = StatusTransitionResponse{
				ID:        history[i].ID,
				NewStatus: string(history[i].NewStatus),
				CreatedBy: history[i].CreatedBy,
				CreatedAt: history[i].CreatedAt.Format(constants.DateFormat),
			}
			if history[i].PreviousStatus != nil {
				ps := string(*history[i].PreviousStatus)
				responses[i].PreviousStatus = &ps
			}
		}

		Logger.Info(logger.LogCSControllerTransSuccess,
			"service_id", serviceID,
			"count", len(responses),
			"trace_id", traceID)

		h.Response.SuccessWithData(c, domain.MsgStatusHistoryRetrieved, responses)
	}
}

// UpdateStatusRequest represents the request body for status updates
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateCompletedServiceStatus handles PATCH /completed-services/:id/status (HU74)
// @Summary Update the status of a completed service
// @Description Updates the status of a completed service, validating the transition is allowed
// @Tags Completed Services
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Completed service ID (obfuscated)"
// @Param body body UpdateStatusRequest true "New status"
// @Success 200 {object} StandardResponse "Status updated successfully"
// @Failure 400 {object} StandardResponse "Invalid status transition"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Completed service not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /completed-services/{id}/status [patch]
func (h *handler) UpdateCompletedServiceStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogCSControllerStatusRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		person, _ := middleware.GetAuthenticatedUser(c)
		if person == nil {
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode completed service ID
		encodedID := c.Param("id")
		serviceID, err := h.DecodeID(encodedID)
		if err != nil {
			Logger.Error(logger.LogCSControllerStatusError, "decode_error", err)
			h.Response.Error(c, domain.MsgCompletedServiceNotFound)
			return
		}

		// Step 3: Parse request body
		var req UpdateStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Logger.Error(logger.LogCSControllerStatusError, "bind_error", err)
			h.Response.Error(c, domain.MsgStatusTransitionError)
			return
		}

		// Step 4: Validate status is a known value
		if !domain.IsValidServiceStatus(req.Status) {
			Logger.Warn(logger.LogCSControllerStatusError, "invalid_status", req.Status)
			h.Response.Error(c, domain.MsgStatusTransitionError)
			return
		}

		// Step 5: Transition status through interactor
		if err := h.CompletedServiceInteractor.TransitionStatus(
			c.Request.Context(),
			serviceID,
			req.Status,
			person.ID,
		); err != nil {
			Logger.Error(logger.LogCSControllerStatusError, "error", err)

			switch {
			case errors.Is(err, domain.ErrCompletedServiceNotFound):
				h.Response.Error(c, domain.MsgCompletedServiceNotFound)
			case errors.Is(err, domain.ErrInvalidStatusTransition):
				h.Response.Error(c, domain.MsgStatusTransitionError)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		Logger.Info(logger.LogCSControllerStatusSuccess,
			"service_id", serviceID,
			"new_status", req.Status,
			"trace_id", traceID)

		h.Response.Success(c, domain.MsgStatusTransitionSuccess)
	}
}

// ServiceStatusItem represents a single status option for the frontend
type ServiceStatusItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// GetServiceStatuses handles GET /completed-services/statuses (HU15)
// @Summary Get all available service statuses
// @Description Returns all valid service statuses as an enum list for frontend selectors
func (h *handler) GetServiceStatuses() gin.HandlerFunc {
	return func(c *gin.Context) {
		statuses := domain.AllServiceStatuses()
		items := make([]ServiceStatusItem, len(statuses))
		for i, s := range statuses {
			items[i] = ServiceStatusItem{
				Value: string(s),
				Label: domain.ServiceStatusLabels[s],
			}
		}

		h.Response.SuccessWithData(c, domain.MsgStatusHistoryRetrieved, items)
	}
}
