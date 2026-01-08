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
			case domain.ErrDuplicateBranchName:
				h.Response.Error(c, domain.MsgBranchDuplicateName)
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

// GetBranch handles GET /branches/:id - retrieves branch details (HU62)
// @Summary Get branch by ID
// @Description Retrieves detailed information of a branch. Returns conditional HATEOAS links based on ownership.
// @Tags Branch
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID (hashid encoded)"
// @Success 200 {object} StandardResponse{data=BranchResponse}
// @Failure 401 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id} [get]
func (h *handler) GetBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBranchControllerGetRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated user from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Get and decode branch ID from URL
		encodedID := c.Param("id")
		branchID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogBranchControllerIDDecodeError,
				"encoded_id", encodedID,
				"error", err)
			// Use "Branch Not Found" for invalid IDs - more user-friendly
			// For end users, malformed ID and non-existent branch should look the same
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		log.Debug(logger.LogBranchControllerGetByID,
			"encoded_id", encodedID,
			"branch_id", branchID)

		// 3. Call interactor to get branch
		branch, err := h.BranchInteractor.GetBranchByID(c.Request.Context(), branchID)
		if err != nil {
			log.Error(logger.LogBranchControllerGetError,
				"error", err,
				"branch_id", branchID,
				"client_ip", c.ClientIP())
			if err == domain.ErrBranchNotFound {
				h.Response.Error(c, domain.MsgBranchNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Determine if current user is the owner (for HATEOAS links)
		isOwner := person != nil && person.ID == branch.RepresentativeID

		// 5. Build HATEOAS links (owner sees edit/delete, others only see self)
		baseURL := GetBaseURL(c)
		links := BuildBranchDetailLinks(baseURL, encodedID, isOwner)

		// 6. Build response DTO (no geocoding status for query)
		response := NewBranchResponse(branch, encodedID, "", links)

		log.Success(logger.LogBranchControllerGetSuccess,
			"branch_id", branch.ID,
			"encoded_id", encodedID,
			"branch_name", branch.Name,
			"is_owner", isOwner,
			"client_ip", c.ClientIP())

		// 7. Send success response (200 OK)
		h.Response.SuccessWithData(c, domain.MsgBranchFound, response)
	}
}

// GetBranchTypes handles GET /branch-types - returns all establishment types (HU76)
// @Summary Get all branch types
// @Description Returns all valid establishment types for branches (public catalog)
// @Tags Branch
// @Produce json
// @Success 200 {object} StandardResponse{data=[]domain.EstablishmentTypeInfo}
// @Router /branch-types [get]
func (h *handler) GetBranchTypes() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBranchControllerGetTypes, "client_ip", c.ClientIP())

		// Get all establishment types from domain
		types := domain.GetAllEstablishmentTypes()

		// Build response with HATEOAS links
		baseURL := GetBaseURL(c)
		response := struct {
			Types []domain.EstablishmentTypeInfo `json:"types"`
			Links []Link                         `json:"_links"`
		}{
			Types: types,
			Links: BuildBranchTypesLinks(baseURL),
		}

		log.Success(logger.LogBranchControllerGetTypesOK, "count", len(types))

		h.Response.SuccessWithData(c, domain.MsgBranchTypesFound, response)
	}
}

// ListBranches handles GET /branches - lists branches for authenticated representative (HU62)
// @Summary List my branches
// @Description Returns all branches owned by the authenticated representative, including location and brands
// @Tags Branch
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse{data=[]BranchListItemResponse}
// @Failure 401 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /branches [get]
func (h *handler) ListBranches() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBranchControllerListRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context (auth + role validated by middleware)
		person, _ := middleware.GetAuthenticatedUser(c)

		log.Debug(logger.LogBranchControllerUserAuth,
			"person_id", person.ID,
			"role", person.Role)

		// 2. Call interactor to get branches for this representative
		branches, err := h.BranchInteractor.GetBranchesByRepresentative(c.Request.Context(), person.ID)
		if err != nil {
			log.Error(logger.LogBranchControllerListError,
				"error", err,
				"person_id", person.ID,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 3. Build response items with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		items := make([]BranchListItemResponse, 0, len(branches))
		for _, branch := range branches {
			encodedID, err := h.EncodeID(branch.ID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "branch_id", branch.ID, "error", err)
				continue // Skip branches with encoding errors
			}
			// Owner always sees full links since this is "my branches"
			itemLinks := BuildBranchDetailLinks(baseURL, encodedID, true)
			items = append(items, NewBranchListItemResponse(branch, encodedID, itemLinks))
		}

		// 4. Build collection response with HATEOAS
		response := struct {
			Branches []BranchListItemResponse `json:"branches"`
			Count    int                      `json:"count"`
			Links    []Link                   `json:"_links"`
		}{
			Branches: items,
			Count:    len(items),
			Links:    BuildBranchListLinks(baseURL),
		}

		log.Success(logger.LogBranchControllerListSuccess,
			"person_id", person.ID,
			"count", len(items),
			"client_ip", c.ClientIP())

		// 5. Send success response (200 OK)
		h.Response.SuccessWithData(c, domain.MsgBranchListFound, response)
	}
}
