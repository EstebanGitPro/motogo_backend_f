package handlers

import (
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// GetMotorcycle handles GET /motorcycles/:id - retrieves motorcycle details (HU46)
// @Summary Get motorcycle by ID
// @Description Retrieves detailed information of a motorcycle including reference and brand info. Returns conditional HATEOAS links based on ownership (Richardson Maturity Level 3).
// @Tags Motorcycles
// @Produce json
// @Security BearerAuth
// @Param id path string true "Motorcycle ID (hashid encoded)"
// @Success 200 {object} StandardResponse{data=MotorcycleResponse} "Motorcycle retrieved successfully"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Motorcycle not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id} [get]
func (h *handler) GetMotorcycle() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerGetRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated user from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Get and decode motorcycle ID from URL
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogMotorcycleControllerIDDecodeError, "encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		log.Debug(logger.LogMotorcycleControllerGetByID, "encoded_id", encodedID, "motorcycle_id", motorcycleID)

		// 3. Call interactor to get motorcycle
		motorcycle, err := h.MotorcycleInteractor.GetMotorcycleByID(c.Request.Context(), motorcycleID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerGetError, "error", err, "motorcycle_id", motorcycleID, "client_ip", c.ClientIP())
			h.mapMotorcycleGetError(c, err)
			return
		}

		// 4. Validate ownership (returns 404 to not reveal existence - security by obscurity)
		log.Debug(logger.LogMotorcycleControllerOwnershipDebug, "person_id", safePersonID(person), "motorcycle_owner_id", motorcycle.OwnerID)
		if person == nil || person.ID != motorcycle.OwnerID {
			log.Warn(logger.LogMotorcycleControllerOwnershipDenied,
				"motorcycle_id", motorcycleID, "requested_by", safePersonID(person),
				"owner_id", motorcycle.OwnerID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 5. Build response DTO
		response := ToMotorcycleResponse(motorcycle)
		response.ID = encodedID
		h.encodeReferenceIDs(&response, motorcycle)

		// 6. Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildMotorcycleDetailLinks(baseURL, encodedID, true)

		log.Success(logger.LogMotorcycleControllerGetSuccess,
			"motorcycle_id", motorcycle.ID, "encoded_id", encodedID,
			"is_owner", true, "client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMotorcycleRetrieved, response)
	}
}

// RegisterMotorcycle handles POST /motorcycles - registers a new motorcycle (HU43)
// @Summary Register a new motorcycle
// @Description Registers a new motorcycle for the authenticated user. Validates that reference exists in catalog and license plate is unique. Returns HATEOAS links (Richardson Maturity Level 3).
// @Tags Motorcycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param motorcycle body RegisterMotorcycleRequest true "Motorcycle registration data"
// @Success 201 {object} StandardResponse{data=MotorcycleResponse} "Motorcycle registered successfully"
// @Failure 400 {object} StandardResponse "Bad request - invalid data or missing required fields"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Reference not found in catalog"
// @Failure 409 {object} StandardResponse "Conflict - duplicate license plate"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles [post]
func (h *handler) RegisterMotorcycle() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerRegRequest,
			"method", c.Request.Method, "path", c.Request.URL.Path, "client_ip", c.ClientIP())

		// 1. Get authenticated user from context
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists || person == nil {
			log.Warn(logger.LogMotorcycleControllerAuthError, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Parse request body
		var req RegisterMotorcycleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogMotorcycleControllerBindError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}
		req.Sanitize()

		// 3. Convert to domain object
		motorcycle := req.ToDomain(person.ID)

		// 4. Decode reference ID if provided
		if req.ReferenceID != nil && *req.ReferenceID != "" {
			referenceID, err := h.DecodeID(*req.ReferenceID)
			if err != nil {
				log.Warn(logger.LogMotorcycleControllerRefDecError, "encoded_id", *req.ReferenceID, "error", err, "client_ip", c.ClientIP())
				h.Response.Error(c, domain.MsgMotorcycleReferenceNotFound)
				return
			}
			motorcycle.ReferenceID = referenceID
		}

		log.Debug(logger.LogMotorcycleControllerRegBody,
			"license_plate", motorcycle.LicensePlate, "reference_id", motorcycle.ReferenceID, "owner_id", person.ID)

		// 5. Call interactor to register motorcycle
		createdMotorcycle, err := h.MotorcycleInteractor.RegisterMotorcycle(c.Request.Context(), motorcycle)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerRegError, "error", err, "license_plate", motorcycle.LicensePlate, "client_ip", c.ClientIP())
			h.mapMotorcycleRegError(c, err)
			return
		}

		// 6. Encode the new motorcycle ID for response
		encodedID, err := h.EncodeID(createdMotorcycle.ID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerIDEncError, "motorcycle_id", createdMotorcycle.ID, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 7. Build response DTO
		response := ToMotorcycleResponse(createdMotorcycle)
		response.ID = encodedID
		h.encodeReferenceIDs(&response, createdMotorcycle)

		// 8. Build HATEOAS links
		baseURL := GetBaseURL(c)
		response.Links = BuildMotorcycleDetailLinks(baseURL, encodedID, true)

		log.Success(logger.LogMotorcycleControllerRegSuccess,
			"motorcycle_id", createdMotorcycle.ID, "encoded_id", encodedID,
			"license_plate", createdMotorcycle.LicensePlate, "client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMotorcycleCreated, response)
	}
}

// ListMotorcycles handles GET /motorcycles - lists all motorcycles for authenticated user (HU47)
// @Summary List user's motorcycles
// @Description Retrieves all motorcycles owned by the authenticated user. Returns HATEOAS links for each motorcycle (Richardson Maturity Level 3).
// @Tags Motorcycles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse{data=[]MotorcycleResponse} "Motorcycles listed successfully"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles [get]
func (h *handler) ListMotorcycles() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerListRequest,
			"method", c.Request.Method, "path", c.Request.URL.Path, "client_ip", c.ClientIP())

		// 1. Get authenticated user from context
		person, _ := middleware.GetAuthenticatedUser(c)
		if person == nil {
			log.Warn(logger.LogMotorcycleControllerNoAuthUser, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// 2. Get all motorcycles for this owner
		motorcycles, err := h.MotorcycleInteractor.GetMotorcyclesByOwner(c.Request.Context(), person.ID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerListError, "error", err, "owner_id", person.ID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 3. Build response DTOs with HATEOAS links
		baseURL := GetBaseURL(c)
		responses := h.buildMotorcycleListResponses(motorcycles, baseURL, log)

		log.Success(logger.LogMotorcycleControllerListSuccess,
			"owner_id", person.ID, "count", len(responses), "client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMotorcyclesListed, responses)
	}
}

// LookupMotorcycleByPlate handles GET /motorcycles/lookup - lookup motorcycle by plate (HU47)
// @Summary Lookup motorcycle by license plate
// @Description Retrieves motorcycle information by license plate. Accessible by representatives (workshops). Returns motorcycle details for service purposes. Enforces branch-specific permissions: only diagnostics from authorized branches are returned.
// @Tags Motorcycles
// @Produce json
// @Security BearerAuth
// @Param plate query string true "License plate to lookup (exact match)"
// @Success 200 {object} StandardResponse{data=MotorcycleLookupResponse} "Motorcycle found successfully"
// @Failure 400 {object} StandardResponse "Bad request - missing plate parameter"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} StandardResponse "Forbidden - no permission to view this motorcycle's diagnostics"
// @Failure 404 {object} StandardResponse "Motorcycle not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/lookup [get]
func (h *handler) LookupMotorcycleByPlate() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerPlateRequest,
			"method", c.Request.Method, "path", c.Request.URL.Path, "client_ip", c.ClientIP())

		// 1. Get authenticated representative from context
		person, _ := middleware.GetAuthenticatedUser(c)
		if person == nil {
			log.Warn(logger.LogMotorcycleControllerNoAuthUser, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Get plate from query parameter
		plate := c.Query("plate")
		if plate == "" {
			log.Warn(logger.LogMotorcycleControllerMissingPlate, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgMissingPlateParam)
			return
		}

		log.Debug(logger.LogMotorcycleControllerPlateDebug, "license_plate", plate)

		// 3. Get representative's branches
		repBranches, err := h.BranchInteractor.GetBranchesByRepresentative(c.Request.Context(), person.ID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerRepBranchErr, "error", err, "person_id", person.ID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		repBranchSet := make(map[string]domain.Branch, len(repBranches))
		for _, b := range repBranches {
			repBranchSet[b.ID] = b
		}

		// 4. Call interactor to get motorcycle by plate
		motorcycle, err := h.MotorcycleInteractor.GetMotorcycleByLicensePlate(c.Request.Context(), plate)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerPlateError, "error", err, "license_plate", plate, "client_ip", c.ClientIP())
			h.mapMotorcycleGetError(c, err)
			return
		}

		// 5. Get active permissions for this motorcycle
		permissions, err := h.MotorcycleInteractor.LookupPermissions(c.Request.Context(), motorcycle.ID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerPermErr, "error", err, "motorcycle_id", motorcycle.ID, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 6. Intersect: find branches where rep has access AND motorcycle has active permission
		permittedBranches, permittedBranchIDs := h.buildPermittedBranches(permissions, repBranchSet, log, c.ClientIP())
		if len(permittedBranches) == 0 {
			log.Warn(logger.LogMotorcycleControllerNoPerm,
				"person_id", person.ID, "motorcycle_id", motorcycle.ID,
				"license_plate", plate, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgMotorcycleNoPermission)
			return
		}

		// 7. Build response DTO
		response := ToMotorcycleLookupResponse(motorcycle)
		encodedID, err := h.EncodeID(motorcycle.ID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerIDEncError, "motorcycle_id", motorcycle.ID, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}
		response.ID = encodedID

		// 8. Hydrate sub-resources
		response.Diagnostics = h.fetchAndEncodeDiagnostics(c, motorcycle.ID, encodedID, permittedBranchIDs, log)
		response.Evidence = h.fetchAndEncodeEvidence(c, motorcycle.ID, encodedID, log)
		response.CompletedServices = h.fetchAndEncodeCompletedServices(c, motorcycle.ID, encodedID, log)
		response.PermittedBranches = permittedBranches

		log.Success(logger.LogMotorcycleControllerPlateSuccess,
			"motorcycle_id", motorcycle.ID, "license_plate", plate,
			"diagnostics_count", len(response.Diagnostics),
			"permitted_branches_count", len(permittedBranches), "client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMotorcycleRetrieved, response)
	}
}

// UpdateMotorcycle handles PUT /motorcycles/:id - update motorcycle (HU44)
// @Summary Update motorcycle information
// @Description Updates motorcycle details. Only the owner can update their motorcycle. Returns 404 for non-owners (security by obscurity). License plate cannot be changed.
// @Tags Motorcycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Motorcycle ID (hashid encoded)"
// @Param motorcycle body UpdateMotorcycleRequest true "Motorcycle update data"
// @Success 200 {object} StandardResponse{data=MotorcycleResponse} "Motorcycle updated successfully"
// @Failure 400 {object} StandardResponse "Bad request - invalid data"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id} [put]
func (h *handler) UpdateMotorcycle() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerUpdateRequest,
			"method", c.Request.Method, "path", c.Request.URL.Path, "client_ip", c.ClientIP())

		// 1. Get authenticated user from context
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists || person == nil {
			log.Warn(logger.LogMotorcycleControllerAuthError, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Get and decode motorcycle ID from URL
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogMotorcycleControllerIDDecodeError, "encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 3. Parse request body
		var req UpdateMotorcycleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogMotorcycleControllerBindError, "error", err, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgMotorcycleCannotUpdate)
			return
		}
		req.Sanitize()

		log.Debug(logger.LogMotorcycleControllerUpdateDebug,
			"motorcycle_id", motorcycleID, "encoded_id", encodedID, "owner_id", person.ID)

		// 4. Call interactor to update motorcycle
		updates := req.ToDomain()
		motorcycle, err := h.MotorcycleInteractor.UpdateMotorcycle(c.Request.Context(), motorcycleID, person.ID, updates)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerUpdateError, "error", err, "motorcycle_id", motorcycleID, "client_ip", c.ClientIP())
			h.mapMotorcycleUpdateError(c, err)
			return
		}

		// 5. Build response DTO with HATEOAS links
		response := ToMotorcycleResponse(motorcycle)
		response.ID = encodedID
		h.encodeReferenceIDs(&response, motorcycle)

		baseURL := GetBaseURL(c)
		response.Links = BuildMotorcycleDetailLinks(baseURL, encodedID, true)

		log.Success(logger.LogMotorcycleControllerUpdateSuccess,
			"motorcycle_id", motorcycle.ID, "encoded_id", encodedID, "client_ip", c.ClientIP())

		h.Response.SuccessWithData(c, domain.MsgMotorcycleUpdated, response)
	}
}

// DeleteMotorcycle handles DELETE /motorcycles/:id - soft delete motorcycle (HU45)
// @Summary Delete a motorcycle
// @Description Soft deletes a motorcycle. Only the owner can delete their motorcycle. Returns 404 for non-owners (security by obscurity).
// @Tags Motorcycles
// @Produce json
// @Security BearerAuth
// @Param id path string true "Motorcycle ID (hashid encoded)"
// @Success 200 {object} StandardResponse "Motorcycle deleted successfully with HATEOAS links"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 404 {object} StandardResponse "Motorcycle not found or not owner"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycles/{id} [delete]
func (h *handler) DeleteMotorcycle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerDeleteRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated user from context
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists || person == nil {
			log.Warn(logger.LogMotorcycleControllerAuthError, "client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Get and decode motorcycle ID from URL
		encodedID := c.Param("id")
		motorcycleID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogMotorcycleControllerIDDecodeError,
				"encoded_id", encodedID,
				"error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		log.Debug(logger.LogMotorcycleControllerDeleteRequest,
			"motorcycle_id", motorcycleID,
			"encoded_id", encodedID,
			"owner_id", person.ID)

		// 3. Call interactor to delete motorcycle
		err = h.MotorcycleInteractor.DeleteMotorcycle(c.Request.Context(), motorcycleID, person.ID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerDeleteError,
				"error", err,
				"motorcycle_id", motorcycleID,
				"client_ip", c.ClientIP())
			switch {
			case errors.Is(err, domain.ErrMotorcycleNotFound):
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			case errors.Is(err, domain.ErrMotorcycleCannotDelete):
				h.Response.Error(c, domain.MsgMotorcycleCannotDelete)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogMotorcycleControllerDeleteSuccess,
			"motorcycle_id", motorcycleID,
			"encoded_id", encodedID,
			"client_ip", c.ClientIP())

		// 4. Build HATEOAS links for next actions (Richardson Level 3)
		baseURL := GetBaseURL(c)
		links := BuildMotorcycleDeletedLinks(baseURL)

		// 5. Send success response (200 OK with HATEOAS links)
		response := map[string]interface{}{
			"_links": links,
		}
		h.Response.SuccessWithData(c, domain.MsgMotorcycleDeleted, response)
	}
}

// GetMotorcycleReferences handles GET /motorcycle-references - lists motorcycle reference catalog (HU50)
// @Summary List motorcycle reference catalog
// @Description Retrieves all motorcycle references (brands/models) for selection during motorcycle registration.
// @Tags Motorcycles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse{data=[]MotorcycleReferenceCatalogItem} "References retrieved successfully"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /motorcycle-references [get]
func (h *handler) GetMotorcycleReferences() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerRefsRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Call interactor to get all references
		references, err := h.MotorcycleInteractor.GetMotorcycleReferences(c.Request.Context())
		if err != nil {
			log.Error(logger.LogMotorcycleControllerRefsError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 2. Build response with encoded IDs
		var responseItems []MotorcycleReferenceCatalogItem
		for _, ref := range references {
			encodedID, err := h.EncodeID(ref.ID)
			if err != nil {
				log.Warn(logger.LogMotorcycleControllerIDEncError,
					"reference_id", ref.ID,
					"error", err)
				continue // Skip if encoding fails
			}

			encodedBrandID, err := h.EncodeID(ref.BrandID)
			if err != nil {
				log.Warn(logger.LogMotorcycleControllerIDEncError,
					"brand_id", ref.BrandID,
					"error", err)
				continue
			}

			responseItems = append(responseItems, MotorcycleReferenceCatalogItem{
				ID:                 encodedID,
				BrandID:            encodedBrandID,
				BrandName:          ref.BrandName,
				Model:              ref.Model,
				Category:           ref.Category,
				EngineDisplacement: ref.EngineDisplacement,
			})
		}

		log.Success(logger.LogMotorcycleControllerRefsSuccess,
			"count", len(responseItems),
			"client_ip", c.ClientIP())

		// 3. Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildMotorcycleReferencesLinks(baseURL)

		// 4. Build final response with data and links
		response := map[string]interface{}{
			"references": responseItems,
			"_links":     links,
		}

		// 5. Send success response (200 OK)
		h.Response.SuccessWithData(c, domain.MsgMotorcycleReferencesListed, response)
	}
}

// GetBrandLines handles GET /admin/brands/:brandId/lines - lists motorcycle lines for a brand (HU40)
// @Summary List motorcycle lines for a specific brand
// @Description Retrieves all motorcycle references (models/lines) associated with a specific brand. This endpoint requires ADMIN role and is used for catalog management.
// @Tags Admin Brands
// @Produce json
// @Security BearerAuth
// @Param brandId path string true "Brand ID (hashid encoded)"
// @Success 200 {object} StandardResponse{data=[]MotorcycleReferenceCatalogItem} "Lines retrieved successfully"
// @Failure 401 {object} StandardResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} StandardResponse "Forbidden - requires ADMIN role"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /admin/brands/{brandId}/lines [get]
func (h *handler) GetBrandLines() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerBrandLinesRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get and decode brand ID from URL
		encodedBrandID := c.Param("brandId")
		brandID, err := h.DecodeID(encodedBrandID)
		if err != nil {
			log.Warn(logger.LogMotorcycleControllerBrandLinesError,
				"encoded_brand_id", encodedBrandID,
				"error", err)
			// Return empty list for invalid brand ID
			h.Response.SuccessWithData(c, domain.MsgBrandLinesRetrieved, map[string]interface{}{
				"lines":  []BrandLineItem{},
				"_links": BuildBrandLinesLinks(GetBaseURL(c), encodedBrandID),
			})
			return
		}

		log.Debug(logger.LogMotorcycleControllerBrandIDDecoded, "encoded", encodedBrandID, "decoded", brandID)

		// 2. Call interactor to get references by brand
		references, err := h.MotorcycleInteractor.GetReferencesByBrandID(c.Request.Context(), brandID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerBrandLinesError,
				"error", err,
				"brand_id", brandID,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 3. Build simplified response (only brand_name and model)
		var responseItems []BrandLineItem
		for _, ref := range references {
			responseItems = append(responseItems, BrandLineItem{
				BrandName: ref.BrandName,
				Model:     ref.Model,
			})
		}

		log.Success(logger.LogMotorcycleControllerBrandLinesSuccess,
			"brand_id", brandID,
			"count", len(responseItems),
			"client_ip", c.ClientIP())

		// 4. Build HATEOAS links
		baseURL := GetBaseURL(c)
		links := BuildBrandLinesLinks(baseURL, encodedBrandID)

		// 5. Build final response with data and links
		response := map[string]interface{}{
			"lines":  responseItems,
			"_links": links,
		}

		// 6. Send success response (200 OK)
		h.Response.SuccessWithData(c, domain.MsgBrandLinesRetrieved, response)
	}
}

// ============================================
// Motorcycle Controller Helpers (extracted to reduce cognitive complexity)
// ============================================

// safePersonID returns the person ID or a fallback string for logging.
func safePersonID(person *domain.Person) string {
	if person != nil {
		return person.ID
	}
	return "anonymous"
}

// encodeReferenceIDs encodes Reference.ID and Reference.BrandID in a MotorcycleResponse.
func (h *handler) encodeReferenceIDs(response *MotorcycleResponse, motorcycle *domain.Motorcycle) {
	if response.Reference == nil || motorcycle.Reference == nil {
		return
	}
	if encodedRefID, err := h.EncodeID(motorcycle.Reference.ID); err == nil {
		response.Reference.ID = encodedRefID
	}
	if encodedBrandID, err := h.EncodeID(motorcycle.Reference.BrandID); err == nil {
		response.Reference.BrandID = encodedBrandID
	}
}

// buildMotorcycleListResponses builds encoded MotorcycleResponse list from domain objects.
func (h *handler) buildMotorcycleListResponses(motorcycles []domain.Motorcycle, baseURL string, log logger.Logger) []MotorcycleResponse {
	var responses []MotorcycleResponse
	for _, moto := range motorcycles {
		encodedID, err := h.EncodeID(moto.ID)
		if err != nil {
			log.Warn(logger.LogMotorcycleControllerIDEncError, "motorcycle_id", moto.ID, "error", err)
			continue
		}
		response := ToMotorcycleResponse(&moto)
		response.ID = encodedID
		h.encodeReferenceIDs(&response, &moto)
		response.Links = BuildMotorcycleDetailLinks(baseURL, encodedID, true)
		responses = append(responses, response)
	}
	return responses
}

// mapMotorcycleGetError maps motorcycle get/lookup errors to API responses.
func (h *handler) mapMotorcycleGetError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrMotorcycleNotFound) {
		h.Response.Error(c, domain.MsgMotorcycleNotFound)
	} else {
		h.Response.Error(c, domain.MsgServerError)
	}
}

// mapMotorcycleRegError maps motorcycle registration errors to API responses.
func (h *handler) mapMotorcycleRegError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrReferenceNotFound):
		h.Response.Error(c, domain.MsgMotorcycleReferenceNotFound)
	case errors.Is(err, domain.ErrReferenceRequired):
		h.Response.Error(c, domain.MsgReferenceRequired)
	case errors.Is(err, domain.ErrDuplicateLicensePlate):
		h.Response.Error(c, domain.MsgDuplicateLicensePlate)
	case errors.Is(err, domain.ErrMotorcycleCannotSave):
		h.Response.Error(c, domain.MsgMotorcycleCannotSave)
	default:
		h.Response.Error(c, domain.MsgServerError)
	}
}

// mapMotorcycleUpdateError maps motorcycle update errors to API responses.
func (h *handler) mapMotorcycleUpdateError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrMotorcycleNotFound):
		h.Response.Error(c, domain.MsgMotorcycleNotFound)
	case errors.Is(err, domain.ErrReferenceNotFound):
		h.Response.Error(c, domain.MsgMotorcycleReferenceNotFound)
	case errors.Is(err, domain.ErrMotorcycleCannotUpdate):
		h.Response.Error(c, domain.MsgMotorcycleCannotUpdate)
	default:
		h.Response.Error(c, domain.MsgServerError)
	}
}

// buildPermittedBranches intersects active permissions with the representative's branches,
// encoding each branch ID. Returns the permitted branches list and a set of raw IDs.
func (h *handler) buildPermittedBranches(
	permissions []domain.DiagnosticPermission,
	repBranchSet map[string]domain.Branch,
	log logger.Logger,
	clientIP string,
) ([]PermittedBranchInfo, map[string]bool) {
	var permittedBranches []PermittedBranchInfo
	permittedBranchIDs := make(map[string]bool)

	for _, perm := range permissions {
		if !perm.Active {
			continue
		}
		if branch, ok := repBranchSet[perm.BranchID]; ok {
			encodedBranchID, encErr := h.EncodeID(branch.ID)
			if encErr != nil {
				log.Error(logger.LogMotorcycleControllerIDEncError,
					"branch_id", branch.ID,
					"error", encErr,
					"client_ip", clientIP)
				continue
			}
			permittedBranches = append(permittedBranches, PermittedBranchInfo{
				ID:   encodedBranchID,
				Name: branch.Name,
			})
			permittedBranchIDs[perm.BranchID] = true
		}
	}

	return permittedBranches, permittedBranchIDs
}

// fetchAndEncodeDiagnostics fetches diagnostics for a motorcycle, filters by permitted branches,
// and encodes all IDs in the response. Returns nil on error (non-fatal).
func (h *handler) fetchAndEncodeDiagnostics(
	c *gin.Context,
	motorcycleID, encodedMotorcycleID string,
	permittedBranchIDs map[string]bool,
	log logger.Logger,
) []DiagnosticResponse {
	if h.DiagnosticInteractor == nil {
		return nil
	}

	diagnostics, err := h.DiagnosticInteractor.ListDiagnosticsByMotorcycleID(c.Request.Context(), motorcycleID)
	if err != nil {
		log.Warn(logger.LogDiagnosticControllerListError,
			"error", err,
			"motorcycle_id", motorcycleID,
			"client_ip", c.ClientIP())
		return nil
	}

	// Filter diagnostics to only include those from permitted branches
	var filteredDiagnostics []domain.Diagnostic
	for _, d := range diagnostics {
		if permittedBranchIDs[d.BranchID] {
			filteredDiagnostics = append(filteredDiagnostics, d)
		}
	}
	if len(filteredDiagnostics) == 0 {
		return nil
	}

	diagnosticResponses := ToDiagnosticResponseList(filteredDiagnostics)

	// Encode IDs (extracted to reduce cognitive complexity)
	h.encodeDiagnosticIDs(diagnosticResponses, filteredDiagnostics, encodedMotorcycleID)

	return diagnosticResponses
}

// encodeDiagnosticIDs encodes all IDs (diagnostic, branch, motorcycle, evidence) in a diagnostic response list.
func (h *handler) encodeDiagnosticIDs(responses []DiagnosticResponse, diagnostics []domain.Diagnostic, encodedMotorcycleID string) {
	for i := range responses {
		if encDiagID, err := h.EncodeID(diagnostics[i].ID); err == nil {
			responses[i].ID = encDiagID
		}
		if encBranchID, err := h.EncodeID(diagnostics[i].BranchID); err == nil {
			responses[i].BranchID = encBranchID
		}
		responses[i].MotorcycleID = encodedMotorcycleID

		for j := range responses[i].Evidence {
			if encEvidID, err := h.EncodeID(diagnostics[i].Evidence[j].ID); err == nil {
				responses[i].Evidence[j].ID = encEvidID
			}
		}
	}
}

// fetchAndEncodeEvidence fetches motorcycle evidence and encodes all IDs.
// Returns nil on error (non-fatal).
func (h *handler) fetchAndEncodeEvidence(
	c *gin.Context,
	motorcycleID, encodedMotorcycleID string,
	log logger.Logger,
) []EvidenceResponse {
	if h.EvidenceInteractor == nil {
		return nil
	}

	evidences, err := h.EvidenceInteractor.LookupEvidence(c.Request.Context(), motorcycleID)
	if err != nil {
		log.Warn(logger.LogEvidenceInteractorLookupError,
			"error", err,
			"motorcycle_id", motorcycleID,
			"client_ip", c.ClientIP())
		return nil
	}
	if len(evidences) == 0 {
		return nil
	}

	evidenceResponses := ToEvidenceResponseList(evidences)
	for i := range evidenceResponses {
		if encID, err := h.EncodeID(evidences[i].ID); err == nil {
			evidenceResponses[i].ID = encID
		}
		evidenceResponses[i].MotorcycleID = encodedMotorcycleID
	}

	return evidenceResponses
}

// fetchAndEncodeCompletedServices fetches completed services for a motorcycle and encodes all IDs.
// Returns nil on error (non-fatal, same pattern as fetchAndEncodeDiagnostics).
func (h *handler) fetchAndEncodeCompletedServices(
	c *gin.Context,
	motorcycleID, encodedMotorcycleID string,
	log logger.Logger,
) []CompletedServiceResponse {
	if h.CompletedServiceInteractor == nil {
		return nil
	}

	services, err := h.CompletedServiceInteractor.GetCompletedServicesByMotorcycle(
		c.Request.Context(),
		motorcycleID,
	)
	if err != nil {
		log.Warn(logger.LogCSControllerListByMotoError,
			"motorcycle_id", motorcycleID,
			"error", err)
		return nil
	}

	if len(services) == 0 {
		return nil
	}

	responses := make([]CompletedServiceResponse, 0, len(services))
	for _, svc := range services {
		resp := ToCompletedServiceResponse(&svc)
		encodedID, err := h.EncodeID(svc.ID)
		if err == nil {
			resp.ID = encodedID
		}
		resp.MotorcycleID = encodedMotorcycleID
		encodedBranchID, err := h.EncodeID(svc.BranchID)
		if err == nil {
			resp.BranchID = encodedBranchID
		}
		responses = append(responses, resp)
	}

	return responses
}
