package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterBranch handles POST /branches - creates a new branch (HU59)
// @Summary Register a new branch
// @Description Creates a new branch (sede) for the authenticated representative
// @Tags Branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RegisterBranchRequest true "Branch registration data"
// @Success 201 {object} StandardResponse{data=BranchResponse}
// @Failure 400 {object} StandardResponse
// @Failure 401 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 409 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /branches [post]
func (h *handler) RegisterBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBranchControllerRegRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		// Auth and role validation are handled by RequireAuth + RequireRole middlewares
		person, _ := middleware.GetAuthenticatedUser(c)

		log.Debug(logger.LogBranchControllerUserAuth,
			"person_id", person.ID,
			"role", person.Role)

		// 2. Bind request
		var req RegisterBranchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogBranchControllerBindError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		// 2.1 Sanitize input (trim whitespace)
		req.Sanitize()

		log.Info(logger.LogBranchControllerProcessing,
			"branch_name", req.Name,
			"establishment_type", req.EstablishmentType,
			"brands_count", len(req.Brands))

		// 3. Map DTO to domain model (using mapper in branch.go)
		branch := req.ToDomain(person.ID)

		// 4. Check if user provided coordinates (to determine geocoding status later)
		userProvidedCoords := branch.Location != nil &&
			branch.Location.Latitude != nil &&
			branch.Location.Longitude != nil

		// 5. Call interactor (returns geocodingSucceeded flag)
		savedBranch, geocodingSucceeded, err := h.BranchInteractor.RegisterBranch(c.Request.Context(), branch)
		if err != nil {
			log.Error(logger.LogBranchControllerRegError,
				"error", err,
				"branch_name", req.Name,
				"client_ip", c.ClientIP())
			switch err {
			case domain.ErrInvalidBranchType:
				h.Response.Error(c, domain.MsgBranchInvalidType)
			case domain.ErrBrandNotFound:
				h.Response.Error(c, domain.MsgBrandNotFound)
			case domain.ErrDuplicateAddress:
				h.Response.Error(c, domain.MsgDuplicateAddress)
			default:
				h.Response.Error(c, domain.MsgBranchCannotSave)
			}
			return
		}

		// 6. Encode ID for response
		encodedID, err := h.EncodeID(savedBranch.ID)
		if err != nil {
			h.HandleIDEncodingError(c, savedBranch.ID, err)
			return
		}

		// 7. Build HATEOAS links and set Location header
		baseURL := GetBaseURL(c)
		links := BuildBranchCreatedLinks(baseURL, encodedID)
		SetLocationHeader(c, baseURL, "branches", encodedID)

		// 8. Determine geocoding status for response
		var geocodingStatus GeocodingStatus
		if userProvidedCoords {
			geocodingStatus = GeocodingStatusSkipped // User provided coords, no geocoding attempted
		} else if geocodingSucceeded {
			geocodingStatus = GeocodingStatusSuccess // Coords were auto-generated
		} else {
			geocodingStatus = GeocodingStatusFailed // Geocoding was attempted but failed
		}

		// 9. Build response DTO (using mapper in branch.go)
		response := NewBranchResponse(savedBranch, encodedID, geocodingStatus, links)

		log.Success(logger.LogBranchControllerRegSuccess,
			"branch_id", savedBranch.ID,
			"encoded_id", encodedID,
			"branch_name", savedBranch.Name,
			"representative_id", savedBranch.RepresentativeID,
			"geocoding_status", geocodingStatus,
			"client_ip", c.ClientIP())

		// 10. Send success response (201 Created)
		h.Response.SuccessWithData(c, domain.MsgBranchRegistered, response)
	}
}
