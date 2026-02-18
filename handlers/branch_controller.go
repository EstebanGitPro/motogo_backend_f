package handlers

import (
	"errors"

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

		// 3. Decode all encoded IDs from request
		decodedBrands, decodedDepartmentID, decodedCityID, decodedFranchiseID, decodeErr := h.decodeBranchRequestIDs(req, log)
		if decodeErr != nil {
			h.Response.Error(c, decodeErr.msgCode)
			return
		}

		// 4. Map DTO to domain model (using mapper in branch.go)
		branch := req.ToDomain(person.ID)
		branch.Brands = decodedBrands           // Override with decoded brand IDs
		branch.FranchiseID = decodedFranchiseID // Override with decoded franchise ID
		// Override location IDs with decoded values
		if branch.Location != nil {
			branch.Location.DepartmentID = decodedDepartmentID
			branch.Location.CityID = decodedCityID
		}

		// 5. Check if user provided coordinates (to determine geocoding status later)
		userProvidedCoords := branch.Location != nil &&
			branch.Location.Latitude != nil &&
			branch.Location.Longitude != nil

		// 6. Call interactor (returns geocodingSucceeded flag)
		savedBranch, geocodingSucceeded, err := h.BranchInteractor.RegisterBranch(c.Request.Context(), branch)
		if err != nil {
			log.Error(logger.LogBranchControllerRegError,
				"error", err,
				"branch_name", req.Name,
				"client_ip", c.ClientIP())
			h.mapRegisterBranchError(c, err)
			return
		}

		// 7. Encode ID for response
		encodedID, err := h.EncodeID(savedBranch.ID)
		if err != nil {
			h.HandleIDEncodingError(c, savedBranch.ID, err)
			return
		}

		// 8. Build HATEOAS links and set Location header
		baseURL := GetBaseURL(c)
		links := BuildBranchCreatedLinks(baseURL, encodedID)
		SetLocationHeader(c, baseURL, "branches", encodedID)

		// 9. Determine geocoding status for response
		var geocodingStatus GeocodingStatus
		switch {
		case userProvidedCoords:
			geocodingStatus = GeocodingStatusSkipped // User provided coords, no geocoding attempted
		case geocodingSucceeded:
			geocodingStatus = GeocodingStatusSuccess // Coords were auto-generated
		default:
			geocodingStatus = GeocodingStatusFailed // Geocoding was attempted but failed
		}

		// 10. Build response DTO and encode IDs
		response := NewBranchResponse(savedBranch, encodedID, geocodingStatus, links)
		h.encodeBranchResponseIDs(&response, savedBranch, log)

		log.Success(logger.LogBranchControllerRegSuccess,
			"branch_id", savedBranch.ID,
			"encoded_id", encodedID,
			"branch_name", savedBranch.Name,
			"representative_id", savedBranch.RepresentativeID,
			"geocoding_status", geocodingStatus,
			"client_ip", c.ClientIP())

		// 12. Send success response (201 Created)
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
			if errors.Is(err, domain.ErrBranchNotFound) {
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

		// 6. Build response DTO and encode IDs
		response := NewBranchResponse(branch, encodedID, "", links)
		h.encodeBranchResponseIDs(&response, branch, log)

		log.Success(logger.LogBranchControllerGetSuccess,
			"branch_id", branch.ID,
			"encoded_id", encodedID,
			"branch_name", branch.Name,
			"is_owner", isOwner,
			"client_ip", c.ClientIP())

		// 9. Send success response (200 OK)
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
			item, ok := h.buildBranchListItem(branch, baseURL, log)
			if !ok {
				continue
			}
			items = append(items, item)
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

// UpdateBranch handles PUT /branches/:id - updates branch information (HU60)
// @Summary Update a branch
// @Description Updates an existing branch. Only the owner can update.
// @Tags Branch
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Param request body RegisterBranchRequest true "Branch update data"
// @Success 200 {object} StandardResponse{data=BranchResponse}
// @Failure 400 {object} StandardResponse
// @Failure 401 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id} [put]
func (h *handler) UpdateBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBranchControllerUpdateRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode branch ID from URL
		encodedID := c.Param("id")
		branchID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogBranchControllerIDDecodeError, "encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 3. Bind request body
		var req RegisterBranchRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogBranchControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		// 3.1 Sanitize input
		req.Sanitize()

		log.Info(logger.LogBranchControllerProcessing,
			"branch_id", branchID,
			"branch_name", req.Name,
			"establishment_type", req.EstablishmentType)

		// 4. Decode all encoded IDs from request
		decodedBrands, decodedDepartmentID, decodedCityID, decodedFranchiseID, decodeErr := h.decodeBranchRequestIDs(req, log)
		if decodeErr != nil {
			h.Response.Error(c, decodeErr.msgCode)
			return
		}

		// 5. Map DTO to domain model
		branch := req.ToDomain(person.ID)
		branch.Brands = decodedBrands
		branch.FranchiseID = decodedFranchiseID // Override with decoded franchise ID
		if branch.Location != nil {
			branch.Location.DepartmentID = decodedDepartmentID
			branch.Location.CityID = decodedCityID
		}

		// 6. Check if user provided coordinates
		userProvidedCoords := branch.Location != nil &&
			branch.Location.Latitude != nil &&
			branch.Location.Longitude != nil

		// 7. Call interactor with ownership validation
		updatedBranch, geocodingSucceeded, err := h.BranchInteractor.UpdateBranch(c.Request.Context(), branchID, branch, person.ID)
		if err != nil {
			log.Error(logger.LogBranchControllerUpdateError, "error", err, "branch_id", branchID)
			h.mapUpdateBranchError(c, err)
			return
		}

		// 8. Encode ID for response
		responseEncodedID, err := h.EncodeID(updatedBranch.ID)
		if err != nil {
			h.HandleIDEncodingError(c, updatedBranch.ID, err)
			return
		}

		// 9. Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildBranchDetailLinks(baseURL, responseEncodedID, true)

		// 10. Determine geocoding status
		var geocodingStatus GeocodingStatus
		switch {
		case userProvidedCoords:
			geocodingStatus = GeocodingStatusSkipped
		case geocodingSucceeded:
			geocodingStatus = GeocodingStatusSuccess
		default:
			geocodingStatus = GeocodingStatusFailed
		}

		// 11. Build response DTO and encode IDs
		response := NewBranchResponse(updatedBranch, responseEncodedID, geocodingStatus, links)
		h.encodeBranchResponseIDs(&response, updatedBranch, log)

		log.Success(logger.LogBranchControllerUpdateSuccess,
			"branch_id", updatedBranch.ID,
			"branch_name", updatedBranch.Name,
			"geocoding_status", geocodingStatus)

		// 13. Send success response (200 OK)
		h.Response.SuccessWithData(c, domain.MsgBranchUpdated, response)
	}
}

// DeleteBranch handles DELETE /branches/:id - deletes a branch (HU61)
// @Summary Delete a branch
// @Description Deletes an existing branch. Only the owner can delete. Cannot delete if has diagnostics or completed services.
// @Tags Branch
// @Produce json
// @Security BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Success 200 {object} StandardResponse
// @Failure 401 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 409 {object} StandardResponse
// @Router /branches/{id} [delete]
func (h *handler) DeleteBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBranchControllerDeleteRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode branch ID from URL
		encodedID := c.Param("id")
		branchID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogBranchControllerIDDecodeError, "encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		log.Debug(logger.LogBranchDeleteProcessing, "branch_id", branchID, "person_id", person.ID)

		// 3. Call interactor with ownership validation
		err = h.BranchInteractor.DeleteBranch(c.Request.Context(), branchID, person.ID)
		if err != nil {
			log.Error(logger.LogBranchControllerDeleteError, "error", err, "branch_id", branchID)
			switch {
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			case errors.Is(err, domain.ErrBranchCannotDelete):
				h.Response.Error(c, domain.MsgBranchHasAssoc)
			default:
				h.Response.Error(c, domain.MsgBranchCannotDelete)
			}
			return
		}

		// 4. Build HATEOAS links for after-delete actions
		baseURL := GetBaseURL(c)
		links := BuildBranchDeletedLinks(baseURL)

		// 5. Build minimal response
		response := struct {
			Links []Link `json:"_links"`
		}{
			Links: links,
		}

		log.Success(logger.LogBranchControllerDeleteSuccess,
			"branch_id", branchID,
			"person_id", person.ID,
			"client_ip", c.ClientIP())

		// 6. Send success response (200 OK)
		h.Response.SuccessWithData(c, domain.MsgBranchDeleted, response)
	}
}

// GetNearbyBranches handles GET /branches/nearby - search branches near location (HU89)
// @Summary Search nearby branches
// @Description Returns branches within the specified radius from given coordinates
// @Tags Branch
// @Produce json
// @Security BearerAuth
// @Param lat query number true "Latitude of the center point"
// @Param lng query number true "Longitude of the center point"
// @Param radius query number false "Search radius in kilometers (default: 5, max: 50)"
// @Param type query string false "Filter by establishment type (WORKSHOP, STORE, WORKSHOP_STORE)"
// @Success 200 {object} StandardResponse{data=NearbyBranchesResponse}
// @Failure 400 {object} StandardResponse
// @Failure 401 {object} StandardResponse
// @Router /branches/nearby [get]
func (h *handler) GetNearbyBranches() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogBranchControllerGetRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Parse and validate all nearby filters
		filters, filterErr := h.parseNearbyFilters(c, log)
		if filterErr != nil {
			h.Response.Error(c, filterErr.msgCode)
			return
		}
		lat, lng, radiusKm := filters.lat, filters.lng, filters.radiusKm
		establishmentType, brandID, displacementRange := filters.establishmentType, filters.brandID, filters.displacementRange

		log.Info(logger.LogBranchNearbySearch,
			"lat", lat,
			"lng", lng,
			"radius_km", radiusKm,
			"type", establishmentType,
			"brand", brandID,
			"displacement_range", displacementRange)

		// 6. Call interactor
		branches, err := h.BranchInteractor.GetBranchesNearby(c.Request.Context(), lat, lng, radiusKm, establishmentType, brandID, displacementRange)
		if err != nil {
			log.Error(logger.LogBranchNearbyError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 6. Build response with encoded IDs and HATEOAS links
		baseURL := GetBaseURL(c)
		items := h.buildNearbyBranchItems(branches, baseURL, log)

		// 7. Build collection response
		response := NearbyBranchesResponse{
			Branches: items,
			Count:    len(items),
			Center: struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			}{
				Latitude:  lat,
				Longitude: lng,
			},
			RadiusKm: radiusKm,
			Links:    BuildNearbyBranchesLinks(baseURL, lat, lng, radiusKm),
		}

		log.Success(logger.LogBranchControllerNearbyFound,
			"count", len(items),
			"radius_km", radiusKm,
			"client_ip", c.ClientIP())

		// 8. Send success response
		h.Response.SuccessWithData(c, domain.MsgBranchNearbyFound, response)
	}
}

// parseFloat converts string to float64
func parseFloat(s string) (float64, error) {
	var result float64
	var negative bool
	var i int

	if len(s) == 0 {
		return 0, domain.ErrInvalidBranchType
	}

	if s[0] == '-' {
		negative = true
		i = 1
	}

	var decimal bool
	var decimalDiv float64 = 10

	for ; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			if decimal {
				return 0, domain.ErrInvalidBranchType
			}
			decimal = true
			continue
		}
		if c < '0' || c > '9' {
			return 0, domain.ErrInvalidBranchType
		}
		digit := float64(c - '0')
		if decimal {
			result += digit / decimalDiv
			decimalDiv *= 10
		} else {
			result = result*10 + digit
		}
	}

	if negative {
		result = -result
	}
	return result, nil
}

// ============================================
// Branch controller helpers (extracted to reduce cognitive complexity)
// ============================================

// buildBranchListItem encodes all IDs for a single branch and builds its list response item.
// Returns the response item and false if the branch should be skipped (encoding error on branch ID).
func (h *handler) buildBranchListItem(branch domain.Branch, baseURL string, log logger.Logger) (BranchListItemResponse, bool) {
	encodedID, err := h.EncodeID(branch.ID)
	if err != nil {
		log.Warn(logger.LogIDEncodeError, "branch_id", branch.ID, "error", err)
		return BranchListItemResponse{}, false
	}

	// Encode franchise_id if present
	var encodedFranchiseID *string
	if branch.FranchiseID != nil && *branch.FranchiseID != "" {
		encoded, err := h.EncodeID(*branch.FranchiseID)
		if err != nil {
			log.Warn(logger.LogIDEncodeError, "franchise_id", *branch.FranchiseID, "error", err)
		} else {
			encodedFranchiseID = &encoded
		}
	}

	// Encode brand IDs for response
	encodedBrands := make([]string, 0, len(branch.Brands))
	for _, brandID := range branch.Brands {
		encodedBrand, err := h.EncodeID(brandID)
		if err != nil {
			log.Warn(logger.LogIDEncodeError, "brand_id", brandID, "error", err)
			continue
		}
		encodedBrands = append(encodedBrands, encodedBrand)
	}

	// Owner always sees full links since this is "my branches"
	itemLinks := BuildBranchDetailLinks(baseURL, encodedID, true)
	item := NewBranchListItemResponse(branch, encodedID, encodedFranchiseID, itemLinks)
	item.Brands = encodedBrands // Override with encoded brand IDs

	// Encode location IDs if present
	if item.Location != nil && branch.Location != nil {
		if encodedDeptID, err := h.EncodeID(branch.Location.DepartmentID); err == nil {
			item.Location.DepartmentID = encodedDeptID
		}
		if encodedCityID, err := h.EncodeID(branch.Location.CityID); err == nil {
			item.Location.CityID = encodedCityID
		}
	}

	return item, true
}

// branchDecodeError wraps a message code for decode/parse errors.
type branchDecodeError struct {
	msgCode string
}

// decodeBranchRequestIDs decodes all encoded IDs from a branch register/update request.
// Returns decoded brands, department ID, city ID, franchise ID, and an error if any decoding fails.
func (h *handler) decodeBranchRequestIDs(req RegisterBranchRequest, log logger.Logger) (
	brands []string, departmentID string, cityID string, franchiseID *string, err *branchDecodeError,
) {
	// Decode brand IDs
	brands = make([]string, 0, len(req.Brands))
	for _, encodedBrandID := range req.Brands {
		decoded, decErr := h.DecodeID(encodedBrandID)
		if decErr != nil {
			log.Warn(logger.LogBranchControllerIDDecodeError, "encoded_id", encodedBrandID, "error", decErr)
			return nil, "", "", nil, &branchDecodeError{msgCode: domain.MsgBrandNotFound}
		}
		brands = append(brands, decoded)
	}

	// Decode location IDs
	departmentID, decErr := h.DecodeID(req.Location.DepartmentID)
	if decErr != nil {
		log.Warn(logger.LogBranchControllerIDDecodeError, "encoded_id", req.Location.DepartmentID, "error", decErr)
		return nil, "", "", nil, &branchDecodeError{msgCode: domain.MsgValIDInvalid}
	}
	cityID, decErr = h.DecodeID(req.Location.CityID)
	if decErr != nil {
		log.Warn(logger.LogBranchControllerIDDecodeError, "encoded_id", req.Location.CityID, "error", decErr)
		return nil, "", "", nil, &branchDecodeError{msgCode: domain.MsgValIDInvalid}
	}

	// Decode franchise_id if provided
	if req.FranchiseID != nil && *req.FranchiseID != "" {
		decoded, decErr := h.DecodeID(*req.FranchiseID)
		if decErr != nil {
			log.Warn(logger.LogBranchControllerIDDecodeError, "encoded_id", *req.FranchiseID, "error", decErr)
			return nil, "", "", nil, &branchDecodeError{msgCode: domain.MsgValIDInvalid}
		}
		franchiseID = &decoded
	}

	return brands, departmentID, cityID, franchiseID, nil
}

// encodeBranchResponseIDs encodes brand IDs, franchise ID, and location IDs in the response.
func (h *handler) encodeBranchResponseIDs(response *BranchResponse, branch *domain.Branch, log logger.Logger) {
	// Encode brand IDs
	encodedBrands := make([]string, 0, len(branch.Brands))
	for _, brandID := range branch.Brands {
		encodedBrand, err := h.EncodeID(brandID)
		if err != nil {
			log.Warn(logger.LogIDEncodeError, "brand_id", brandID, "error", err)
			continue
		}
		encodedBrands = append(encodedBrands, encodedBrand)
	}
	response.Brands = encodedBrands

	// Encode franchise_id if present
	if branch.FranchiseID != nil && *branch.FranchiseID != "" {
		if encodedFranchiseID, err := h.EncodeID(*branch.FranchiseID); err == nil {
			response.FranchiseID = &encodedFranchiseID
		} else {
			log.Warn(logger.LogIDEncodeError, "franchise_id", *branch.FranchiseID, "error", err)
		}
	}

	// Encode location IDs if present
	if response.Location != nil && branch.Location != nil {
		if encodedDeptID, err := h.EncodeID(branch.Location.DepartmentID); err == nil {
			response.Location.DepartmentID = encodedDeptID
		}
		if encodedCityID, err := h.EncodeID(branch.Location.CityID); err == nil {
			response.Location.CityID = encodedCityID
		}
	}
}

// mapRegisterBranchError maps RegisterBranch interactor errors to response messages.
func (h *handler) mapRegisterBranchError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidBranchType):
		h.Response.Error(c, domain.MsgBranchInvalidType)
	case errors.Is(err, domain.ErrDuplicateBranchName):
		h.Response.Error(c, domain.MsgBranchDuplicateName)
	case errors.Is(err, domain.ErrBrandNotFound):
		h.Response.Error(c, domain.MsgBrandNotFound)
	case errors.Is(err, domain.ErrDuplicateAddress):
		h.Response.Error(c, domain.MsgDuplicateAddress)
	default:
		h.Response.Error(c, domain.MsgBranchCannotSave)
	}
}

// mapUpdateBranchError maps UpdateBranch interactor errors to response messages.
func (h *handler) mapUpdateBranchError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrBranchNotFound):
		h.Response.Error(c, domain.MsgBranchNotFound)
	case errors.Is(err, domain.ErrForbidden):
		h.Response.Error(c, domain.MsgForbidden)
	case errors.Is(err, domain.ErrInvalidBranchType):
		h.Response.Error(c, domain.MsgBranchInvalidType)
	case errors.Is(err, domain.ErrBrandNotFound):
		h.Response.Error(c, domain.MsgBrandNotFound)
	default:
		h.Response.Error(c, domain.MsgBranchCannotUpdate)
	}
}

// buildNearbyBranchItems encodes all IDs for nearby branches and builds response items.
func (h *handler) buildNearbyBranchItems(branches []domain.NearbyBranch, baseURL string, log logger.Logger) []NearbyBranchResponse {
	items := make([]NearbyBranchResponse, 0, len(branches))
	for _, branch := range branches {
		encodedID, err := h.EncodeID(branch.ID)
		if err != nil {
			log.Warn(logger.LogIDEncodeError, "branch_id", branch.ID, "error", err)
			continue
		}
		resp := NewNearbyBranchResponse(branch, encodedID, baseURL)

		// Encode brand IDs
		if len(branch.Brands) > 0 {
			encodedBrands := make([]string, 0, len(branch.Brands))
			for _, brandID := range branch.Brands {
				encodedBrand, err := h.EncodeID(brandID)
				if err != nil {
					log.Warn(logger.LogIDEncodeError, "brand_id", brandID, "error", err)
					continue
				}
				encodedBrands = append(encodedBrands, encodedBrand)
			}
			resp.Brands = encodedBrands
		}

		items = append(items, resp)
	}
	return items
}

// nearbyFilters holds parsed query parameters for GetNearbyBranches.
type nearbyFilters struct {
	lat               float64
	lng               float64
	radiusKm          float64
	establishmentType string
	brandID           string
	displacementRange string
}

// parseNearbyFilters parses and validates all query parameters for GetNearbyBranches.
func (h *handler) parseNearbyFilters(c *gin.Context, log logger.Logger) (*nearbyFilters, *branchDecodeError) {
	// Parse latitude (required)
	latStr := c.Query("lat")
	if latStr == "" {
		log.Warn(logger.LogBranchNearbyMissingLat, "client_ip", c.ClientIP())
		return nil, &branchDecodeError{msgCode: domain.MsgValLatitudeRequired}
	}
	lat, err := parseFloat(latStr)
	if err != nil || lat < -90 || lat > 90 {
		log.Warn(logger.LogBranchNearbyInvalidLat, "lat", latStr, "error", err)
		return nil, &branchDecodeError{msgCode: domain.MsgValLatitudeInvalid}
	}

	// Parse longitude (required)
	lngStr := c.Query("lng")
	if lngStr == "" {
		log.Warn(logger.LogBranchNearbyMissingLng, "client_ip", c.ClientIP())
		return nil, &branchDecodeError{msgCode: domain.MsgValLongitudeRequired}
	}
	lng, err := parseFloat(lngStr)
	if err != nil || lng < -180 || lng > 180 {
		log.Warn(logger.LogBranchNearbyInvalidLng, "lng", lngStr, "error", err)
		return nil, &branchDecodeError{msgCode: domain.MsgValLongitudeInvalid}
	}

	// Parse radius (optional, default 5km, max 50km)
	radiusKm := 5.0
	if radiusStr := c.Query("radius"); radiusStr != "" {
		radius, err := parseFloat(radiusStr)
		if err != nil || radius <= 0 || radius > 50 {
			log.Warn(logger.LogBranchNearbyInvalidRadius, "radius", radiusStr, "error", err)
			return nil, &branchDecodeError{msgCode: domain.MsgValRadiusInvalid}
		}
		radiusKm = radius
	}

	// Parse type (optional)
	establishmentType := c.Query("type")
	if establishmentType != "" {
		if !domain.IsValidEstablishmentType(domain.EstablishmentType(establishmentType)) {
			log.Warn(logger.LogBranchNearbyInvalidType, "type", establishmentType)
			return nil, &branchDecodeError{msgCode: domain.MsgBranchInvalidType}
		}
	}

	// Parse optional brand filter (extracted to reduce cognitive complexity)
	brandID := h.parseOptionalBrandFilter(c.Query("brand"), log)

	return &nearbyFilters{
		lat:               lat,
		lng:               lng,
		radiusKm:          radiusKm,
		establishmentType: establishmentType,
		brandID:           brandID,
		displacementRange: c.Query("displacement_range"),
	}, nil
}

// parseOptionalBrandFilter decodes an optional brand filter parameter.
// Returns the decoded brand ID, or empty string if not provided or invalid.
func (h *handler) parseOptionalBrandFilter(brandEncoded string, log logger.Logger) string {
	if brandEncoded == "" {
		return ""
	}
	decoded, err := h.DecodeID(brandEncoded)
	if err != nil {
		log.Warn(logger.LogBranchControllerIDDecodeError, "brand_filter", brandEncoded, "error", err)
		return "" // Invalid brand ID, ignore filter
	}
	return decoded
}
