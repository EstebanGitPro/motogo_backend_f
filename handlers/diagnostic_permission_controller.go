package handlers

import (
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// GrantDiagnosticPermission handles POST /motorcycles/:id/permissions - grant diagnostic access to a branch
// @Summary Grant diagnostic permission
// @Description Allows a motorcycle owner to grant a specific branch permission to view diagnostic details.
// @Tags Diagnostic Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Motorcycle ID (hashid encoded)"
// @Param permission body GrantPermissionRequest true "Branch to grant permission"
// @Success 201 {object} StandardResponse{data=DiagnosticPermissionResponse} "Permission granted"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/permissions [post]
func (h *handler) GrantDiagnosticPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogDiagPermControllerGrantRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated user
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists || person == nil {
			log.Warn(logger.LogDiagPermControllerGrantRequest, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Decode motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogDiagPermControllerGrantError,
				"encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 3. Parse request body
		var req GrantPermissionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogDiagPermControllerGrantError,
				"error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgPermissionCannotSave)
			return
		}
		req.Sanitize()

		// 4. Decode branch ID
		branchID, err := h.DecodeID(req.BranchID)
		if err != nil {
			log.Warn(logger.LogDiagPermControllerGrantError,
				"encoded_branch_id", req.BranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 5. Resolve active value (default true for backward compat)
		active := true
		if req.Active != nil {
			active = *req.Active
		}

		// 6. Call interactor
		permission, err := h.MotorcycleInteractor.GrantDiagnosticPermission(
			c.Request.Context(), motorcycleID, branchID, person.ID, active)
		if err != nil {
			log.Error(logger.LogDiagPermControllerGrantError,
				"error", err,
				"motorcycle_id", motorcycleID,
				"branch_id", branchID,
				"client_ip", c.ClientIP())
			switch {
			case errors.Is(err, domain.ErrMotorcycleNotFound):
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			case errors.Is(err, domain.ErrPermissionCannotSave):
				h.Response.Error(c, domain.MsgPermissionCannotSave)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 6. Build response
		response := ToDiagnosticPermissionResponse(permission)
		// Encode IDs for response
		response.ID = "" // Not exposed externally
		response.MotorcycleID = encodedID
		response.BranchID = req.BranchID // Already encoded from client

		// 7. HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildPermissionLinks(baseURL, encodedID)

		log.Success(logger.LogDiagPermControllerGrantRequest,
			"motorcycle_id", motorcycleID,
			"branch_id", branchID,
			"client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgPermissionGranted, response)
	}
}

// RevokeDiagnosticPermission handles DELETE /motorcycles/:id/permissions/:branchId - revoke diagnostic access
// @Summary Revoke diagnostic permission
// @Description Allows a motorcycle owner to revoke a branch's permission to view diagnostic details.
// @Tags Diagnostic Permissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Motorcycle ID (hashid encoded)"
// @Param branchId path string true "Branch ID (hashid encoded)"
// @Success 200 {object} StandardResponse "Permission revoked"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/permissions/{branchId} [delete]
func (h *handler) RevokeDiagnosticPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogDiagPermControllerRevokeRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated user
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists || person == nil {
			log.Warn(logger.LogDiagPermControllerRevokeRequest, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Decode motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogDiagPermControllerRevokeError,
				"encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 3. Decode branch ID
		encodedBranchID := c.Param("branchId")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			log.Warn(logger.LogDiagPermControllerRevokeError,
				"encoded_branch_id", encodedBranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 4. Call interactor
		err = h.MotorcycleInteractor.RevokeDiagnosticPermission(
			c.Request.Context(), motorcycleID, branchID, person.ID)
		if err != nil {
			log.Error(logger.LogDiagPermControllerRevokeError,
				"error", err,
				"motorcycle_id", motorcycleID,
				"branch_id", branchID,
				"client_ip", c.ClientIP())
			switch {
			case errors.Is(err, domain.ErrMotorcycleNotFound):
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			case errors.Is(err, domain.ErrPermissionNotFound):
				h.Response.Error(c, domain.MsgPermissionNotFound)
			case errors.Is(err, domain.ErrPermissionCannotDelete):
				h.Response.Error(c, domain.MsgPermissionCannotDelete)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogDiagPermControllerRevokeRequest,
			"motorcycle_id", motorcycleID,
			"branch_id", branchID,
			"client_ip", c.ClientIP())

		// 5. HATEOAS links for next actions
		baseURL := GetBaseURL(c)
		links := BuildPermissionLinks(baseURL, encodedID)
		response := map[string]interface{}{
			"_links": links,
		}

		h.Response.SuccessWithData(c, domain.MsgPermissionRevoked, response)
	}
}

// ListDiagnosticPermissions handles GET /motorcycles/:id/permissions - list all permissions for a motorcycle
// @Summary List diagnostic permissions
// @Description Lists all active diagnostic permissions for a motorcycle. Only the owner can view permissions.
// @Tags Diagnostic Permissions
// @Produce json
// @Security BearerAuth
// @Param id path string true "Motorcycle ID (hashid encoded)"
// @Success 200 {object} StandardResponse{data=[]DiagnosticPermissionResponse} "Permissions listed"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/permissions [get]
func (h *handler) ListDiagnosticPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogDiagPermControllerListRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated user
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists || person == nil {
			log.Warn(logger.LogDiagPermControllerListRequest, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Decode motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogDiagPermControllerListError,
				"encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 3. Call interactor
		permissions, err := h.MotorcycleInteractor.ListDiagnosticPermissions(
			c.Request.Context(), motorcycleID, person.ID)
		if err != nil {
			log.Error(logger.LogDiagPermControllerListError,
				"error", err,
				"motorcycle_id", motorcycleID,
				"client_ip", c.ClientIP())
			if errors.Is(err, domain.ErrMotorcycleNotFound) {
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Build response DTOs with encoded IDs
		var responses []DiagnosticPermissionResponse
		for _, p := range permissions {
			resp := ToDiagnosticPermissionResponse(&p)
			resp.ID = "" // Internal ID not exposed
			resp.MotorcycleID = encodedID
			if encodedBranchID, err := h.EncodeID(p.BranchID); err == nil {
				resp.BranchID = encodedBranchID
			}
			responses = append(responses, resp)
		}

		log.Success(logger.LogDiagPermControllerListRequest,
			"motorcycle_id", motorcycleID,
			"count", len(responses),
			"client_ip", c.ClientIP())

		// 5. HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildPermissionLinks(baseURL, encodedID)

		result := map[string]interface{}{
			"permissions": responses,
			"_links":      links,
		}

		h.Response.SuccessWithData(c, domain.MsgPermissionsListed, result)
	}
}
