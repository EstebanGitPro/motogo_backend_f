package handlers

import (
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// CreateDiagnostic handles POST /motorcycles/:id/diagnostics - creates a diagnostic (HU11)
// @Summary Create motorcycle diagnostic
// @Description Creates a diagnostic request for a motorcycle.
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param diagnostic body CreateDiagnosticRequest true "Diagnostic data"
// @Success 201 {object} StandardResponse{data=DiagnosticResponse} "Diagnostic created successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle or branch not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/diagnostics [post]
func (h *handler) CreateDiagnostic() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogDiagnosticControllerCreateRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode motorcycle ID
		encodedMotorcycleID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedMotorcycleID)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerCreateError, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Step 3: Parse request body
		var request CreateDiagnosticRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			Logger.Error(logger.LogDiagnosticControllerCreateError, "bind_error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Sanitize input
		request.Sanitize()

		// Step 4: Decode branch ID
		branchID, err := h.DecodeID(request.BranchID)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerCreateError, "branch_decode_error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// Step 5: Create diagnostic through interactor
		diagnostic, err := h.DiagnosticInteractor.RegisterDiagnostic(
			c.Request.Context(),
			motorcycleID,
			branchID,
			user.ID,
			request.ProblemDescription,
		)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerCreateError, "error", err)

			switch {
			case errors.Is(err, domain.ErrMotorcycleNotFound):
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			default:
				h.Response.Error(c, domain.MsgDiagnosticCannotSave)
			}
			return
		}

		Logger.Info(logger.LogDiagnosticControllerCreateSuccess, "id", diagnostic.ID, "trace_id", traceID)

		// Step 6: Build response with encoded IDs
		response := ToDiagnosticResponse(diagnostic)
		encodedDiagID, _ := h.EncodeID(diagnostic.ID)
		response.ID = encodedDiagID
		response.MotorcycleID = encodedMotorcycleID
		response.BranchID = request.BranchID // return the obfuscated branch ID the client sent

		h.Response.SuccessWithData(c, domain.MsgDiagnosticCreated, response)
	}
}

// ListDiagnostics handles GET /motorcycles/:id/diagnostics - lists diagnostics for a motorcycle (HU14)
// @Summary List motorcycle diagnostics
// @Description Lists all diagnostics for a motorcycle owned by authenticated user
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=[]DiagnosticResponse} "Diagnostics list"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/diagnostics [get]
func (h *handler) ListDiagnostics() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogDiagnosticControllerListRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode motorcycle ID
		encodedMotorcycleID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedMotorcycleID)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerListError, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Step 3: Get diagnostics through interactor
		diagnostics, err := h.DiagnosticInteractor.ListDiagnosticsByMotorcycle(
			c.Request.Context(),
			motorcycleID,
			user.ID,
		)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerListError, "error", err)

			if errors.Is(err, domain.ErrMotorcycleNotFound) {
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			} else {
				h.Response.Error(c, domain.MsgDiagnosticNotFound)
			}
			return
		}

		Logger.Info(logger.LogDiagnosticControllerListSuccess, "count", len(diagnostics), "trace_id", traceID)

		// Step 4: Build response with encoded IDs
		responses := ToDiagnosticResponseList(diagnostics)
		for i := range responses {
			encodedDiagID, _ := h.EncodeID(diagnostics[i].ID)
			responses[i].ID = encodedDiagID
			responses[i].MotorcycleID = encodedMotorcycleID
			encodedBranchID, _ := h.EncodeID(diagnostics[i].BranchID)
			responses[i].BranchID = encodedBranchID
		}

		h.Response.SuccessWithData(c, domain.MsgDiagnosticsListed, responses)
	}
}

// GetDiagnostic handles GET /motorcycles/:id/diagnostics/:diagnosticId - gets a diagnostic (HU14)
// @Summary Get motorcycle diagnostic
// @Description Gets a specific diagnostic with evidence for a motorcycle owned by authenticated user
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param diagnosticId path string true "Diagnostic ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=DiagnosticResponse} "Diagnostic found"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Diagnostic not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/diagnostics/{diagnosticId} [get]
func (h *handler) GetDiagnostic() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogDiagnosticControllerGetRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode diagnostic ID
		encodedDiagnosticID := c.Param("diagnosticId")
		diagnosticID, err := h.DecodeID(encodedDiagnosticID)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerGetError, "error", err)
			h.Response.Error(c, domain.MsgDiagnosticNotFound)
			return
		}

		encodedMotorcycleID := c.Param("id")

		// Step 3: Get diagnostic through interactor
		diagnostic, err := h.DiagnosticInteractor.GetDiagnosticByID(
			c.Request.Context(),
			diagnosticID,
			user.ID,
		)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerGetError, "error", err)
			h.Response.Error(c, domain.MsgDiagnosticNotFound)
			return
		}

		Logger.Info(logger.LogDiagnosticControllerGetSuccess, "id", diagnosticID, "trace_id", traceID)

		// Step 4: Build response with encoded IDs
		response := ToDiagnosticResponse(diagnostic)
		response.ID = encodedDiagnosticID
		response.MotorcycleID = encodedMotorcycleID
		encodedBranchID, _ := h.EncodeID(diagnostic.BranchID)
		response.BranchID = encodedBranchID

		h.Response.SuccessWithData(c, domain.MsgDiagnosticRetrieved, response)
	}
}

// UpdateDiagnostic handles PUT /motorcycles/:id/diagnostics/:diagnosticId - updates a diagnostic (HU12)
// @Summary Update motorcycle diagnostic
// @Description Updates a diagnostic. Only the owner can update.
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param diagnosticId path string true "Diagnostic ID (obfuscated)"
// @Param diagnostic body UpdateDiagnosticRequest true "Updated diagnostic data"
// @Success 200 {object} StandardResponse{data=DiagnosticResponse} "Diagnostic updated"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Diagnostic not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/diagnostics/{diagnosticId} [put]
func (h *handler) UpdateDiagnostic() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogDiagnosticControllerUpdateRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode diagnostic ID
		encodedDiagnosticID := c.Param("diagnosticId")
		diagnosticID, err := h.DecodeID(encodedDiagnosticID)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerUpdateError, "error", err)
			h.Response.Error(c, domain.MsgDiagnosticNotFound)
			return
		}

		encodedMotorcycleID := c.Param("id")

		// Step 3: Parse request body
		var request UpdateDiagnosticRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			Logger.Error(logger.LogDiagnosticControllerUpdateError, "bind_error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Sanitize input
		request.Sanitize()

		// Step 4: Update diagnostic through interactor
		diagnostic, err := h.DiagnosticInteractor.UpdateDiagnostic(
			c.Request.Context(),
			diagnosticID,
			user.ID,
			request.ToDomain(),
		)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerUpdateError, "error", err)

			if errors.Is(err, domain.ErrDiagnosticNotFound) {
				h.Response.Error(c, domain.MsgDiagnosticNotFound)
			} else {
				h.Response.Error(c, domain.MsgDiagnosticCannotUpdate)
			}
			return
		}

		Logger.Info(logger.LogDiagnosticControllerUpdateSuccess, "id", diagnosticID, "trace_id", traceID)

		// Step 5: Build response with encoded IDs
		response := ToDiagnosticResponse(diagnostic)
		response.ID = encodedDiagnosticID
		response.MotorcycleID = encodedMotorcycleID
		encodedBranchID, _ := h.EncodeID(diagnostic.BranchID)
		response.BranchID = encodedBranchID

		h.Response.SuccessWithData(c, domain.MsgDiagnosticUpdated, response)
	}
}

// DeleteDiagnostic handles DELETE /motorcycles/:id/diagnostics/:diagnosticId - deletes a diagnostic (HU13)
// @Summary Delete motorcycle diagnostic
// @Description Deletes a diagnostic. Only the owner can delete.
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param diagnosticId path string true "Diagnostic ID (obfuscated)"
// @Success 200 {object} StandardResponse "Diagnostic deleted"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Diagnostic not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/diagnostics/{diagnosticId} [delete]
func (h *handler) DeleteDiagnostic() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogDiagnosticControllerDeleteRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode diagnostic ID
		encodedDiagnosticID := c.Param("diagnosticId")
		diagnosticID, err := h.DecodeID(encodedDiagnosticID)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerDeleteError, "error", err)
			h.Response.Error(c, domain.MsgDiagnosticNotFound)
			return
		}

		// Step 3: Delete diagnostic through interactor
		err = h.DiagnosticInteractor.DeleteDiagnostic(
			c.Request.Context(),
			diagnosticID,
			user.ID,
		)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerDeleteError, "error", err)

			if errors.Is(err, domain.ErrDiagnosticNotFound) {
				h.Response.Error(c, domain.MsgDiagnosticNotFound)
			} else {
				h.Response.Error(c, domain.MsgDiagnosticCannotDelete)
			}
			return
		}

		Logger.Info(logger.LogDiagnosticControllerDeleteSuccess, "id", diagnosticID, "trace_id", traceID)

		h.Response.Success(c, domain.MsgDiagnosticDeleted)
	}
}

// SetDiagnosticSolution handles PATCH /diagnostics/:id/solution - sets solution for a diagnostic (representative)
// @Summary Set diagnostic solution
// @Description Sets the possible solution for a diagnostic. Used by workshop representatives.
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Diagnostic ID (obfuscated)"
// @Param solution body SetSolutionRequest true "Solution data"
// @Success 200 {object} StandardResponse "Solution saved"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Diagnostic not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /diagnostics/{id}/solution [patch]
func (h *handler) SetDiagnosticSolution() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogDiagnosticControllerSetSolutionRequest, "trace_id", traceID)

		// Step 1: Decode diagnostic ID
		encodedDiagnosticID := c.Param("id")
		diagnosticID, err := h.DecodeID(encodedDiagnosticID)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerSetSolutionError, "error", err)
			h.Response.Error(c, domain.MsgDiagnosticNotFound)
			return
		}

		// Step 2: Parse request body
		var request SetSolutionRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			Logger.Error(logger.LogDiagnosticControllerSetSolutionError, "bind_error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Sanitize input
		request.Sanitize()

		// Step 3: Set solution through interactor (no ownership check)
		err = h.DiagnosticInteractor.SetSolution(
			c.Request.Context(),
			diagnosticID,
			request.PossibleSolution,
		)
		if err != nil {
			Logger.Error(logger.LogDiagnosticControllerSetSolutionError, "error", err)

			if errors.Is(err, domain.ErrDiagnosticNotFound) {
				h.Response.Error(c, domain.MsgDiagnosticNotFound)
			} else {
				h.Response.Error(c, domain.MsgDiagnosticCannotUpdate)
			}
			return
		}

		Logger.Info(logger.LogDiagnosticControllerSetSolutionSuccess, "id", diagnosticID, "trace_id", traceID)

		h.Response.Success(c, domain.MsgDiagnosticUpdated)
	}
}
