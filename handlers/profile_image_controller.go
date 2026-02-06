package handlers

import (
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// UpdateProfileImage handles PUT /motorcycles/:id/profile-image - Add or update image (HU36/37)
// @Summary Update motorcycle profile image
// @Description Adds or updates the profile image URL for a motorcycle. Only the owner can update.
// @Accept json
// @Produce json
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Param image body ProfileImageRequest true "Profile image URL"
// @Success 200 {object} StandardResponse{data=ProfileImageResponse} "Image updated successfully"
// @Failure 400 {object} StandardResponse "Invalid request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/profile-image [put]
func (h *handler) UpdateProfileImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)
		Logger.Info(logger.LogMotorcycleControllerUpdateRequest, "trace_id", traceID, "action", "update_profile_image")

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogMotorcycleControllerAuthError)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			Logger.Error(logger.LogMotorcycleControllerIDDecodeError, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Step 3: Parse request body
		var req ProfileImageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Logger.Error(logger.LogMotorcycleControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Step 4: Update via interactor
		updates := req.ToDomain()
		motorcycle, err := h.MotorcycleInteractor.UpdateMotorcycle(c.Request.Context(), motorcycleID, user.ID, &updates)
		if err != nil {
			Logger.Error(logger.LogMotorcycleControllerUpdateError, "error", err)

			switch {
			case errors.Is(err, domain.ErrMotorcycleNotFound):
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			default:
				h.Response.Error(c, domain.MsgMotorcycleCannotUpdate)
			}
			return
		}

		Logger.Info(logger.LogMotorcycleControllerUpdateSuccess, "id", motorcycle.ID, "trace_id", traceID)

		// Step 5: Build response with HATEOAS
		baseURL := GetBaseURL(c)
		response := ToProfileImageResponse(encodedID, motorcycle.ProfileImageURL)
		response.Links = BuildProfileImageLinks(baseURL, encodedID, motorcycle.ProfileImageURL != nil)

		h.Response.SuccessWithData(c, domain.MsgProfileImageUpdated, response)
	}
}

// GetProfileImage handles GET /motorcycles/:id/profile-image - Get image URL (HU38)
// @Summary Get motorcycle profile image
// @Description Gets the profile image URL for a motorcycle. Only the owner can view.
// @Accept json
// @Produce json
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=ProfileImageResponse} "Image retrieved successfully"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/profile-image [get]
func (h *handler) GetProfileImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)
		Logger.Info(logger.LogMotorcycleControllerGetRequest, "trace_id", traceID, "action", "get_profile_image")

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogMotorcycleControllerAuthError)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			Logger.Error(logger.LogMotorcycleControllerIDDecodeError, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Step 3: Get motorcycle via interactor (validates ownership internally)
		motorcycle, err := h.MotorcycleInteractor.GetMotorcycleByID(c.Request.Context(), motorcycleID)
		if err != nil {
			Logger.Error(logger.LogMotorcycleControllerGetError, "error", err)

			if errors.Is(err, domain.ErrMotorcycleNotFound) {
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// Step 4: Validate ownership
		if motorcycle.OwnerID != user.ID {
			Logger.Warn(logger.LogMotorcycleControllerOwnershipDenied, "motorcycle_owner", motorcycle.OwnerID, "requester", user.ID)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		Logger.Info(logger.LogMotorcycleControllerGetSuccess, "id", motorcycle.ID, "trace_id", traceID)

		// Step 5: Build response with HATEOAS
		baseURL := GetBaseURL(c)
		response := ToProfileImageResponse(encodedID, motorcycle.ProfileImageURL)
		response.Links = BuildProfileImageLinks(baseURL, encodedID, motorcycle.ProfileImageURL != nil)

		h.Response.SuccessWithData(c, domain.MsgProfileImageGet, response)
	}
}

// DeleteProfileImage handles DELETE /motorcycles/:id/profile-image - Remove image (HU39)
// @Summary Delete motorcycle profile image
// @Description Removes the profile image from a motorcycle. Only the owner can delete.
// @Accept json
// @Produce json
// @Param id path string true "Motorcycle ID (obfuscated)"
// @Success 200 {object} StandardResponse{data=ProfileImageResponse} "Image deleted successfully"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id}/profile-image [delete]
func (h *handler) DeleteProfileImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetTraceIDFromContext(c)
		Logger.Info(logger.LogMotorcycleControllerDeleteRequest, "trace_id", traceID, "action", "delete_profile_image")

		// Step 1: Get authenticated user
		user, exists := middleware.GetAuthenticatedUser(c)
		if !exists {
			Logger.Warn(logger.LogMotorcycleControllerAuthError)
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Step 2: Decode motorcycle ID
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			Logger.Error(logger.LogMotorcycleControllerIDDecodeError, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// Step 3: Delete profile image via interactor (handles Storage + DB cleanup)
		err = h.MotorcycleInteractor.DeleteProfileImage(c.Request.Context(), motorcycleID, user.ID)
		if err != nil {
			Logger.Error(logger.LogMotorcycleControllerUpdateError, "error", err)

			switch {
			case errors.Is(err, domain.ErrMotorcycleNotFound):
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			default:
				h.Response.Error(c, domain.MsgMotorcycleCannotUpdate)
			}
			return
		}

		Logger.Info(logger.LogMotorcycleControllerDeleteSuccess, "id", motorcycleID, "trace_id", traceID)

		// Step 4: Build response with HATEOAS
		baseURL := GetBaseURL(c)
		response := ToProfileImageResponse(encodedID, nil)
		response.Links = BuildProfileImageLinks(baseURL, encodedID, false)

		h.Response.SuccessWithData(c, domain.MsgProfileImageDeleted, response)
	}
}
