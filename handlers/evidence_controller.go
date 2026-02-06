package handlers

import (
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// CreateEvidence handles POST /motorcycles/:id/evidence - creates evidence for a motorcycle (HU16)
// @Summary Upload photographic evidence
// @Description Upload photographic evidence for a motorcycle. The image must already be uploaded to Firebase Storage.
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param evidence body CreateEvidenceRequest true "Evidence data"
// @Success 201 {object} StandardResponse{data=EvidenceResponse} "Evidence created successfully"
// @Failure 400 {object} StandardResponse "Bad request - invalid image URL"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 409 {object} StandardResponse "Evidence limit exceeded"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/evidence [post]
func (h *handler) CreateEvidence() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogEvidenceControllerCreateRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode and validate motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			Logger.Error(logger.LogBranchControllerIDDecodeError, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Step 3: Parse request body
		var request CreateEvidenceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			Logger.Error(logger.LogBranchControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Sanitize input
		request.Sanitize()

		// Step 4: Create evidence through interactor
		evidence, err := h.EvidenceInteractor.CreateEvidence(
			c.Request.Context(),
			motorcycleID,
			user.ID,
			request.ToDomain(),
		)
		if err != nil {
			Logger.Error(logger.LogEvidenceControllerCreateError, "error", err)

			switch {
			case errors.Is(err, domain.ErrMotorcycleNotFound):
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			case errors.Is(err, domain.ErrInvalidEvidenceURL):
				h.Response.Error(c, domain.MsgInvalidEvidenceURL)
			case errors.Is(err, domain.ErrEvidenceLimitExceeded):
				h.Response.Error(c, domain.MsgEvidenceLimitExceeded)
			default:
				h.Response.Error(c, domain.MsgEvidenceCannotSave)
			}
			return
		}

		Logger.Info(logger.LogEvidenceControllerCreateSuccess, "id", evidence.ID, "trace_id", traceID)

		// Step 5: Build response with HATEOAS links
		baseURL := GetBaseURL(c)
		encodedEvidenceID, _ := h.EncodeID(evidence.ID)
		response := ToEvidenceResponse(evidence)
		response.ID = encodedEvidenceID
		response.MotorcycleID = encodedID
		response.Links = BuildEvidenceDetailLinks(baseURL, encodedID, encodedEvidenceID, true)

		h.Response.SuccessWithData(c, domain.MsgEvidenceCreated, response)
	}
}

// ListEvidence handles GET /motorcycles/:id/evidence - lists all evidence for a motorcycle (HU18)
// @Summary List motorcycle evidence
// @Description Lists all photographic evidence for a motorcycle owned by authenticated user
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=[]EvidenceResponse} "Evidence list"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/evidence [get]
func (h *handler) ListEvidence() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogEvidenceControllerListRequest, "trace_id", traceID)

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
			Logger.Error(logger.LogBranchControllerIDDecodeError, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Step 3: Get evidence list through interactor
		evidences, err := h.EvidenceInteractor.ListEvidenceByMotorcycle(
			c.Request.Context(),
			motorcycleID,
			user.ID,
		)
		if err != nil {
			Logger.Error(logger.LogEvidenceControllerListError, "error", err)

			if errors.Is(err, domain.ErrMotorcycleNotFound) {
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			} else {
				h.Response.Error(c, domain.MsgEvidenceNotFound)
			}
			return
		}

		Logger.Info(logger.LogEvidenceControllerListSuccess, "count", len(evidences), "trace_id", traceID)

		// Step 4: Build response with HATEOAS links
		baseURL := GetBaseURL(c)
		responses := make([]EvidenceResponse, len(evidences))
		for i, e := range evidences {
			encodedEvidenceID, _ := h.EncodeID(e.ID)
			responses[i] = ToEvidenceResponse(&e)
			responses[i].ID = encodedEvidenceID
			responses[i].MotorcycleID = encodedMotorcycleID
			responses[i].Links = BuildEvidenceDetailLinks(baseURL, encodedMotorcycleID, encodedEvidenceID, true)
		}

		// Include list-level links
		listResponse := map[string]interface{}{
			"items":  responses,
			"_links": BuildEvidenceListLinks(baseURL, encodedMotorcycleID),
		}

		h.Response.SuccessWithData(c, domain.MsgEvidencesListed, listResponse)
	}
}

// UpdateEvidence handles PUT /motorcycles/:id/evidence/:evidenceId - updates evidence (HU17)
// @Summary Update motorcycle evidence
// @Description Updates photographic evidence for a motorcycle. Only the owner can update.
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param evidenceId path string true "Evidence ID (obfuscated)"
// @Param evidence body CreateEvidenceRequest true "Updated evidence data"
// @Success 200 {object} StandardResponse{data=EvidenceResponse} "Evidence updated successfully"
// @Failure 400 {object} StandardResponse "Bad request - invalid data"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Evidence not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/evidence/{evidenceId} [put]
func (h *handler) UpdateEvidence() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogEvidenceControllerUpdateRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode evidence ID
		encodedEvidenceID := c.Param("evidenceId")
		evidenceID, err := h.DecodeID(encodedEvidenceID)
		if err != nil {
			Logger.Error(logger.LogBranchControllerIDDecodeError, "error", err)
			h.Response.Error(c, domain.MsgEvidenceNotFound)
			return
		}

		// Get motorcycle ID for HATEOAS links
		encodedMotorcycleID := c.Param("id")

		// Step 3: Parse request body
		var request CreateEvidenceRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			Logger.Error(logger.LogBranchControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Sanitize input
		request.Sanitize()

		// Step 4: Update evidence through interactor
		evidence, err := h.EvidenceInteractor.UpdateEvidence(
			c.Request.Context(),
			evidenceID,
			user.ID,
			request.ToDomain(),
		)
		if err != nil {
			Logger.Error(logger.LogEvidenceControllerUpdateError, "error", err)

			switch {
			case errors.Is(err, domain.ErrEvidenceNotFound):
				h.Response.Error(c, domain.MsgEvidenceNotFound)
			case errors.Is(err, domain.ErrInvalidEvidenceURL):
				h.Response.Error(c, domain.MsgInvalidEvidenceURL)
			default:
				h.Response.Error(c, domain.MsgEvidenceCannotUpdate)
			}
			return
		}

		Logger.Info(logger.LogEvidenceControllerUpdateSuccess, "id", evidenceID, "trace_id", traceID)

		// Step 5: Build response with HATEOAS links
		baseURL := GetBaseURL(c)
		response := ToEvidenceResponse(evidence)
		response.ID = encodedEvidenceID
		response.MotorcycleID = encodedMotorcycleID
		response.Links = BuildEvidenceDetailLinks(baseURL, encodedMotorcycleID, encodedEvidenceID, true)

		h.Response.SuccessWithData(c, domain.MsgEvidenceUpdated, response)
	}
}

// DeleteEvidence handles DELETE /motorcycles/:id/evidence/:evidenceId - deletes evidence (HU19)
// @Summary Delete motorcycle evidence
// @Description Deletes photographic evidence for a motorcycle. Only the owner can delete.
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param evidenceId path string true "Evidence ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=map[string]interface{}} "Evidence deleted, includes navigation links"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Evidence not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/evidence/{evidenceId} [delete]
func (h *handler) DeleteEvidence() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)

		Logger.Info(logger.LogEvidenceControllerDeleteRequest, "trace_id", traceID)

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogBranchControllerUserUnauth)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode evidence ID
		encodedEvidenceID := c.Param("evidenceId")
		evidenceID, err := h.DecodeID(encodedEvidenceID)
		if err != nil {
			Logger.Error(logger.LogBranchControllerIDDecodeError, "error", err)
			h.Response.Error(c, domain.MsgEvidenceNotFound)
			return
		}

		// Note: motorcycle ID is passed for HATEOAS links but not used in delete
		encodedMotorcycleID := c.Param("id")

		// Step 3: Delete evidence through interactor
		err = h.EvidenceInteractor.DeleteEvidence(
			c.Request.Context(),
			evidenceID,
			user.ID,
		)
		if err != nil {
			Logger.Error(logger.LogEvidenceControllerDeleteError, "error", err)

			if errors.Is(err, domain.ErrEvidenceNotFound) {
				h.Response.Error(c, domain.MsgEvidenceNotFound)
			} else {
				h.Response.Error(c, domain.MsgEvidenceCannotDelete)
			}
			return
		}

		Logger.Info(logger.LogEvidenceControllerDeleteSuccess, "id", evidenceID, "trace_id", traceID)

		// Step 4: Return HATEOAS navigation links after deletion
		baseURL := GetBaseURL(c)
		links := BuildEvidenceDeletedLinks(baseURL, encodedMotorcycleID)
		response := map[string]interface{}{
			"_links": links,
		}

		h.Response.SuccessWithData(c, domain.MsgEvidenceDeleted, response)
	}
}
