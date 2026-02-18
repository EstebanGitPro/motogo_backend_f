package handlers

import (
	"errors"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// RegisterFranchise handles POST /franchises (HU26)
// @Summary Create a new franchise
// @Description Creates a new franchise and associates specified branches
// @Tags Franchises
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param request body CreateFranchiseRequest true "Franchise data with branch IDs"
// @Success 201 {object} StandardResponse
// @Failure 400 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 409 {object} StandardResponse
// @Router /franchises [post]
func (h *handler) RegisterFranchise(franchiseInteractor *interactor.FranchiseInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogFranchiseControllerRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Bind request
		var req CreateFranchiseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogFranchiseControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		// Sanitize input
		req.Sanitize()

		log.Info(logger.LogFranchiseControllerProcessing,
			"franchise_name", req.Name,
			"branch_count", len(req.BranchIDs))

		// 3. Decode branch IDs (they come encoded from frontend)
		branchIDs := make([]string, len(req.BranchIDs))
		for i, encodedID := range req.BranchIDs {
			decoded, err := h.DecodeID(encodedID)
			if err != nil {
				log.Warn(logger.LogFranchiseControllerIDDecodeError, "encoded_id", encodedID, "error", err)
				h.Response.Error(c, domain.MsgBranchNotFound)
				return
			}
			branchIDs[i] = decoded
		}

		// 4. Create franchise with branches
		franchise := req.ToFranchiseDomain()
		createdFranchise, err := franchiseInteractor.CreateFranchiseWithBranches(c.Request.Context(), franchise, branchIDs, person.ID)
		if err != nil {
			log.Error(logger.LogFranchiseControllerCreateError, "error", err, "franchise_name", req.Name)
			switch {
			case errors.Is(err, domain.ErrFranchiseDuplicateName):
				h.Response.Error(c, domain.MsgFranchiseDuplicateName)
			case errors.Is(err, domain.ErrFranchiseNoBranches):
				h.Response.Error(c, domain.MsgFranchiseNoBranches)
			case errors.Is(err, domain.ErrFranchiseBranchNotOwned):
				h.Response.Error(c, domain.MsgFranchiseBranchNotOwned)
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 5. Encode franchise ID for response
		encodedID, err := h.EncodeID(createdFranchise.ID)
		if err != nil {
			h.HandleIDEncodingError(c, createdFranchise.ID, err)
			return
		}

		// 6. Build HATEOAS links and set Location header
		baseURL := GetBaseURL(c)
		links := BuildFranchiseLinks(baseURL, encodedID)
		SetLocationHeader(c, baseURL, "franchises", encodedID)

		// 7. Build response DTO
		response := ToFranchiseResponse(createdFranchise, len(req.BranchIDs), encodedID, links)

		log.Success(logger.LogFranchiseControllerCreateSuccess,
			"franchise_id", createdFranchise.ID,
			"encoded_id", encodedID,
			"franchise_name", createdFranchise.Name,
			"branch_count", len(req.BranchIDs))

		// 8. Send success response (201 Created)
		h.Response.SuccessWithData(c, domain.MsgFranchiseCreated, response)
	}
}

// ListFranchises handles GET /franchises (HU29 - list)
// @Summary List franchises for the authenticated representative
// @Tags Franchises
// @Produce json
// @Security     BearerAuth
// @Success 200 {object} StandardResponse
// @Router /franchises [get]
func (h *handler) ListFranchises(franchiseInteractor *interactor.FranchiseInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogFranchiseControllerListRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Get franchises for this representative
		franchises, err := franchiseInteractor.GetFranchisesByRepresentative(c.Request.Context(), person.ID)
		if err != nil {
			log.Error(logger.LogFranchiseControllerListError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 3. Encode IDs for response
		baseURL := GetBaseURL(c)
		encodedIDs := make([]string, len(franchises))
		for i, f := range franchises {
			encoded, err := h.EncodeID(f.ID)
			if err != nil {
				log.Warn(logger.LogIDEncodeError, "franchise_id", f.ID, "error", err)
				continue
			}
			encodedIDs[i] = encoded
		}

		// 4. Build response
		response := ToFranchiseListResponse(franchises, encodedIDs, baseURL)

		log.Success(logger.LogFranchiseControllerListSuccess, "count", len(franchises))

		h.Response.SuccessWithData(c, domain.MsgFranchisesListed, response)
	}
}

// GetFranchise handles GET /franchises/:id (HU29 - detail)
// @Summary Get franchise by ID
// @Tags Franchises
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Franchise ID"
// @Success 200 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /franchises/{id} [get]
func (h *handler) GetFranchise(franchiseInteractor *interactor.FranchiseInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// 1. Decode franchise ID
		encodedID := c.Param("id")
		franchiseID, err := h.DecodeID(encodedID)
		if err != nil {
			log.Warn(logger.LogFranchiseControllerIDDecodeError, "encoded_id", encodedID, "error", err)
			h.Response.Error(c, domain.MsgFranchiseNotFound)
			return
		}

		// 2. Get franchise
		franchise, err := franchiseInteractor.GetFranchiseByID(c.Request.Context(), franchiseID)
		if err != nil {
			log.Error(logger.LogFranchiseControllerGetError, "error", err, "franchise_id", franchiseID)
			if errors.Is(err, domain.ErrFranchiseNotFound) {
				h.Response.Error(c, domain.MsgFranchiseNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 3. Build response with HATEOAS
		baseURL := GetBaseURL(c)
		links := BuildFranchiseLinks(baseURL, encodedID)
		response := ToFranchiseResponse(franchise, 0, encodedID, links)

		log.Success(logger.LogFranchiseControllerGetSuccess, "franchise_id", franchiseID, "name", franchise.Name)

		h.Response.SuccessWithData(c, domain.MsgFranchiseFound, response)
	}
}

// UpdateFranchise handles PUT /franchises/:id (HU27)
// @Summary Update a franchise
// @Tags Franchises
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Franchise ID"
// @Param request body UpdateFranchiseRequest true "Updated franchise data"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /franchises/{id} [put]
func (h *handler) UpdateFranchise(franchiseInteractor *interactor.FranchiseInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// 1. Get authenticated person
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode franchise ID
		encodedID := c.Param("id")
		franchiseID, err := h.DecodeID(encodedID)
		if err != nil {
			h.Response.Error(c, domain.MsgFranchiseNotFound)
			return
		}

		// 3. Bind request
		var req UpdateFranchiseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogFranchiseControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		// Sanitize input
		req.Sanitize()

		// 4. Update franchise
		franchise := req.ToFranchiseDomain(franchiseID)
		if err := franchiseInteractor.UpdateFranchise(c.Request.Context(), franchise, person.ID); err != nil {
			log.Error(logger.LogFranchiseControllerUpdateError, "error", err, "franchise_id", franchiseID)
			switch {
			case errors.Is(err, domain.ErrFranchiseNotFound):
				h.Response.Error(c, domain.MsgFranchiseNotFound)
			case errors.Is(err, domain.ErrFranchiseDuplicateName):
				h.Response.Error(c, domain.MsgFranchiseDuplicateName)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 5. Build response
		baseURL := GetBaseURL(c)
		links := BuildFranchiseLinks(baseURL, encodedID)
		response := struct {
			Links []Link `json:"_links"`
		}{Links: links}

		log.Success(logger.LogFranchiseControllerUpdateSuccess, "franchise_id", franchiseID, "name", req.Name)

		h.Response.SuccessWithData(c, domain.MsgFranchiseUpdated, response)
	}
}

// DeleteFranchise handles DELETE /franchises/:id (HU28)
// @Summary Delete a franchise
// @Description Deletes a franchise and disassociates all branches
// @Tags Franchises
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Franchise ID"
// @Success 200 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /franchises/{id} [delete]
func (h *handler) DeleteFranchise(franchiseInteractor *interactor.FranchiseInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		// 1. Get authenticated person
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode franchise ID
		encodedID := c.Param("id")
		franchiseID, err := h.DecodeID(encodedID)
		if err != nil {
			h.Response.Error(c, domain.MsgFranchiseNotFound)
			return
		}

		// 3. Delete franchise
		if err := franchiseInteractor.DeleteFranchise(c.Request.Context(), franchiseID, person.ID); err != nil {
			log.Error(logger.LogFranchiseControllerDeleteError, "error", err, "franchise_id", franchiseID)
			if errors.Is(err, domain.ErrFranchiseNotFound) {
				h.Response.Error(c, domain.MsgFranchiseNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Build response with next actions
		baseURL := GetBaseURL(c)
		response := struct {
			Links []Link `json:"_links"`
		}{
			Links: []Link{
				{Rel: "list", Href: BuildCollectionURL(baseURL, "franchises"), Method: "GET"},
				{Rel: "create", Href: BuildCollectionURL(baseURL, "franchises"), Method: "POST"},
			},
		}

		log.Success(logger.LogFranchiseControllerDeleteSuccess, "franchise_id", franchiseID)

		h.Response.SuccessWithData(c, domain.MsgFranchiseDeleted, response)
	}
}

// BuildFranchiseLinks generates HATEOAS links for a franchise
func BuildFranchiseLinks(baseURL, franchiseEncodedID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "franchises", franchiseEncodedID)
	collectionURL := BuildCollectionURL(baseURL, "franchises")

	return []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
		{Href: resourceURL, Rel: "update", Method: "PUT"},
		{Href: resourceURL, Rel: "delete", Method: "DELETE"},
		{Href: collectionURL, Rel: "list", Method: "GET"},
		{Href: fmt.Sprintf("%s?franchise_id=%s", BuildCollectionURL(baseURL, "branches"), franchiseEncodedID), Rel: "branches", Method: "GET"},
		{Href: fmt.Sprintf("%s/branches", resourceURL), Rel: "add-branch", Method: "POST"},
	}
}

// AddBranchToFranchise handles POST /franchises/:id/branches
// @Summary Add a branch to a franchise
// @Description Associates an existing branch to a franchise
// @Tags Franchises
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Franchise ID (encoded)"
// @Param request body AddBranchToFranchiseRequest true "Branch ID to add"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /franchises/{id}/branches [post]
func (h *handler) AddBranchToFranchise(franchiseInteractor *interactor.FranchiseInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogFranchiseControllerAddBranchRequest, "client_ip", c.ClientIP())

		// 1. Get authenticated person
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode franchise ID
		encodedFranchiseID := c.Param("id")
		franchiseID, err := h.DecodeID(encodedFranchiseID)
		if err != nil {
			log.Warn(logger.LogIDDecodeError, "encoded_id", encodedFranchiseID, "error", err)
			h.Response.Error(c, domain.MsgFranchiseNotFound)
			return
		}

		// 3. Bind request
		var req AddBranchToFranchiseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error(logger.LogBranchControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValJSONInvalid)
			return
		}

		// Sanitize input
		req.Sanitize()

		// 4. Decode branch ID
		branchID, err := h.DecodeID(req.BranchID)
		if err != nil {
			log.Warn(logger.LogIDDecodeError, "encoded_id", req.BranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 5. Call interactor
		if err := franchiseInteractor.AddBranchToFranchise(c.Request.Context(), franchiseID, branchID, person.ID); err != nil {
			log.Error(logger.LogFranchiseControllerAddBranchError, "error", err, "franchise_id", franchiseID, "branch_id", branchID)
			switch {
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			case errors.Is(err, domain.ErrFranchiseBranchNotOwned):
				h.Response.Error(c, domain.MsgFranchiseBranchNotOwned)
			case errors.Is(err, domain.ErrFranchiseNotFound):
				h.Response.Error(c, domain.MsgFranchiseNotFound)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 6. Build response
		baseURL := GetBaseURL(c)
		links := BuildFranchiseLinks(baseURL, encodedFranchiseID)

		log.Success(logger.LogFranchiseControllerAddBranchSuccess, "franchise_id", franchiseID, "branch_id", branchID)

		h.Response.SuccessWithData(c, domain.MsgFranchiseBranchAdded, struct {
			Links []Link `json:"_links"`
		}{Links: links})
	}
}

// RemoveBranchFromFranchise handles DELETE /franchises/:id/branches/:branchId
// @Summary Remove a branch from a franchise
// @Description Dissociates a branch from a franchise
// @Tags Franchises
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Franchise ID (encoded)"
// @Param branchId path string true "Branch ID (encoded)"
// @Success 200 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /franchises/{id}/branches/{branchId} [delete]
func (h *handler) RemoveBranchFromFranchise(franchiseInteractor *interactor.FranchiseInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogFranchiseControllerRemBranchRequest, "client_ip", c.ClientIP())

		// 1. Get authenticated person
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode franchise ID
		encodedFranchiseID := c.Param("id")
		franchiseID, err := h.DecodeID(encodedFranchiseID)
		if err != nil {
			log.Warn(logger.LogIDDecodeError, "encoded_id", encodedFranchiseID, "error", err)
			h.Response.Error(c, domain.MsgFranchiseNotFound)
			return
		}

		// 3. Decode branch ID
		encodedBranchID := c.Param("branchId")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			log.Warn(logger.LogIDDecodeError, "encoded_id", encodedBranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 4. Call interactor
		if err := franchiseInteractor.RemoveBranchFromFranchise(c.Request.Context(), franchiseID, branchID, person.ID); err != nil {
			log.Error(logger.LogFranchiseControllerRemBranchError, "error", err, "franchise_id", franchiseID, "branch_id", branchID)
			switch {
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			case errors.Is(err, domain.ErrFranchiseMinBranches):
				h.Response.Error(c, domain.MsgFranchiseMinBranches)
			case errors.Is(err, domain.ErrFranchiseNotFound):
				h.Response.Error(c, domain.MsgFranchiseNotFound)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 5. Build response
		baseURL := GetBaseURL(c)
		links := BuildFranchiseLinks(baseURL, encodedFranchiseID)

		log.Success(logger.LogFranchiseControllerRemBranchSuccess, "franchise_id", franchiseID, "branch_id", branchID)

		h.Response.SuccessWithData(c, domain.MsgFranchiseBranchRemoved, struct {
			Links []Link `json:"_links"`
		}{Links: links})
	}
}
