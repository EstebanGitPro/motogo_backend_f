package handlers

import (
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// ==========================================
// Private Helpers
// ==========================================

// decodedRegisterIDs holds all decoded IDs needed for RegisterCompletedService.
type decodedRegisterIDs struct {
	BranchID     string
	MotorcycleID string
	DiagnosticID *string
	ServiceIDs   []string
}

// decodeRegisterIDs decodes all obfuscated IDs from a CreateCompletedServiceRequest.
// Returns nil and sends an error response if any ID fails to decode.
func (h *handler) decodeRegisterIDs(c *gin.Context, request *CreateCompletedServiceRequest) *decodedRegisterIDs {
	branchID, err := h.DecodeID(request.BranchID)
	if err != nil {
		Logger.Error(logger.LogCSControllerCreateError, "branch_decode_error", err)
		h.Response.Error(c, domain.MsgBranchNotFound)
		return nil
	}

	motorcycleID, err := h.DecodeID(request.MotorcycleID)
	if err != nil {
		Logger.Error(logger.LogCSControllerCreateError, "motorcycle_decode_error", err)
		h.Response.Error(c, domain.MsgMotorcycleNotFound)
		return nil
	}

	var diagnosticID *string
	if request.DiagnosticID != nil && *request.DiagnosticID != "" {
		decoded, err := h.DecodeID(*request.DiagnosticID)
		if err != nil {
			Logger.Error(logger.LogCSControllerCreateError, "diagnostic_decode_error", err)
			h.Response.Error(c, domain.MsgDiagnosticNotFound)
			return nil
		}
		diagnosticID = &decoded
	}

	serviceIDs, err := h.decodeServiceIDs(request.ServiceIDs)
	if err != nil {
		Logger.Error(logger.LogCSControllerCreateError, "service_decode_error", err)
		h.Response.Error(c, domain.MsgInvalidBranchServices)
		return nil
	}

	return &decodedRegisterIDs{
		BranchID:     branchID,
		MotorcycleID: motorcycleID,
		DiagnosticID: diagnosticID,
		ServiceIDs:   serviceIDs,
	}
}

// decodeServiceIDs decodes a slice of obfuscated service IDs.
func (h *handler) decodeServiceIDs(encodedIDs []string) ([]string, error) {
	decodedIDs := make([]string, len(encodedIDs))
	for i, encodedID := range encodedIDs {
		decoded, err := h.DecodeID(encodedID)
		if err != nil {
			return nil, err
		}
		decodedIDs[i] = decoded
	}
	return decodedIDs, nil
}

// mapRegisterCSError maps domain errors from RegisterCompletedService to API responses.
func (h *handler) mapRegisterCSError(c *gin.Context, err error) {
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
}

// validateMotorcycleOwnership decodes and validates motorcycle ownership, returning the decoded ID.
func (h *handler) validateMotorcycleOwnership(c *gin.Context, log logger.Logger, encodedID string, personID string) (string, error) {
	motorcycleID, err := h.DecodeID(encodedID)
	if err != nil {
		log.Warn(logger.LogMotorcycleControllerIDDecodeError, "encoded_id", encodedID, "error", err)
		h.Response.Error(c, domain.MsgMotorcycleNotFound)
		return "", err
	}

	motorcycle, err := h.MotorcycleInteractor.GetMotorcycleByID(c.Request.Context(), motorcycleID)
	if err != nil {
		log.Error(logger.LogMotorcycleControllerGetError, "error", err, "motorcycle_id", motorcycleID)
		if errors.Is(err, domain.ErrMotorcycleNotFound) {
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
		} else {
			h.Response.Error(c, domain.MsgServerError)
		}
		return "", err
	}

	if personID != motorcycle.OwnerID {
		log.Warn(logger.LogMotorcycleControllerOwnershipDenied,
			"motorcycle_id", motorcycleID,
			"requested_by", personID,
			"owner_id", motorcycle.OwnerID)
		h.Response.Error(c, domain.MsgMotorcycleNotFound)
		return "", domain.ErrMotorcycleNotFound
	}

	return motorcycleID, nil
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

		// Step 3: Decode all obfuscated IDs
		ids := h.decodeRegisterIDs(c, &request)
		if ids == nil {
			return // error response already sent
		}

		// Step 4: Build domain entity
		cs := &domain.CompletedService{
			BranchID:            ids.BranchID,
			MotorcycleID:        ids.MotorcycleID,
			DiagnosticID:        ids.DiagnosticID,
			QuotedPrice:         request.QuotedPrice,
			FinalPrice:          request.FinalPrice,
			RepresentativeNotes: request.RepresentativeNotes,
		}

		// Step 5: Register through interactor
		result, err := h.CompletedServiceInteractor.RegisterCompletedService(
			c.Request.Context(),
			cs,
			ids.ServiceIDs,
			person.ID,
		)
		if err != nil {
			Logger.Error(logger.LogCSControllerCreateError, "error", err)
			h.mapRegisterCSError(c, err)
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

		// 2. Decode motorcycle ID and validate ownership
		encodedID := c.Param("id")
		motorcycleID, err := h.validateMotorcycleOwnership(c, log, encodedID, person.ID)
		if err != nil {
			return
		}

		// 3. Get completed services
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

		// 4. Build response with encoded IDs
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
