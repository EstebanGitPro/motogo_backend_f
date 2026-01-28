package handlers

import (
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
		// Create logger with trace ID for this request
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
			log.Warn(logger.LogMotorcycleControllerIDDecodeError,
				"encoded_id", encodedID,
				"error", err)
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		log.Debug(logger.LogMotorcycleControllerGetByID,
			"encoded_id", encodedID,
			"motorcycle_id", motorcycleID)

		// 3. Call interactor to get motorcycle
		motorcycle, err := h.MotorcycleInteractor.GetMotorcycleByID(c.Request.Context(), motorcycleID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerGetError,
				"error", err,
				"motorcycle_id", motorcycleID,
				"client_ip", c.ClientIP())
			if err == domain.ErrMotorcycleNotFound {
				h.Response.Error(c, domain.MsgMotorcycleNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Validate ownership - only owner can view their motorcycle (security)
		// Returns 404 to not reveal existence to non-owners (security by obscurity)
		if person == nil || person.ID != motorcycle.OwnerID {
			log.Warn(logger.LogMotorcycleControllerOwnershipDenied,
				"motorcycle_id", motorcycleID,
				"requested_by", func() string {
					if person != nil {
						return person.ID
					}
					return "anonymous"
				}(),
				"owner_id", motorcycle.OwnerID,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgMotorcycleNotFound)
			return
		}

		// 5. Build response DTO
		response := ToMotorcycleResponse(motorcycle)

		// 6. Build HATEOAS links (owner sees edit/delete, others only see self)
		baseURL := GetBaseURL(c)
		response.Links = BuildMotorcycleDetailLinks(baseURL, encodedID, true) // Owner validated above

		log.Success(logger.LogMotorcycleControllerGetSuccess,
			"motorcycle_id", motorcycle.ID,
			"encoded_id", encodedID,
			"is_owner", true,
			"client_ip", c.ClientIP())

		// 7. Send success response (200 OK)
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
		// Create logger with trace ID for this request
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogMotorcycleControllerRegRequest,
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

		// 2. Parse request body
		var req RegisterMotorcycleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogMotorcycleControllerBindError,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 3. Convert to domain object
		motorcycle := req.ToDomain(person.ID)

		// 4. Decode reference ID if provided (optional until Release 11)
		if req.ReferenceID != nil && *req.ReferenceID != "" {
			referenceID, err := h.DecodeID(*req.ReferenceID)
			if err != nil {
				log.Warn(logger.LogMotorcycleControllerRefDecError,
					"encoded_id", *req.ReferenceID,
					"error", err,
					"client_ip", c.ClientIP())
				h.Response.Error(c, domain.MsgMotorcycleReferenceNotFound)
				return
			}
			motorcycle.ReferenceID = referenceID
		}

		log.Debug(logger.LogMotorcycleControllerRegBody,
			"license_plate", motorcycle.LicensePlate,
			"reference_id", motorcycle.ReferenceID,
			"owner_id", person.ID)

		// 5. Call interactor to register motorcycle
		createdMotorcycle, err := h.MotorcycleInteractor.RegisterMotorcycle(c.Request.Context(), motorcycle)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerRegError,
				"error", err,
				"license_plate", motorcycle.LicensePlate,
				"client_ip", c.ClientIP())

			switch err {
			case domain.ErrReferenceNotFound:
				h.Response.Error(c, domain.MsgMotorcycleReferenceNotFound)
			case domain.ErrDuplicateLicensePlate:
				h.Response.Error(c, domain.MsgDuplicateLicensePlate)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 6. Encode the new motorcycle ID for response
		encodedID, err := h.EncodeID(createdMotorcycle.ID)
		if err != nil {
			log.Error(logger.LogMotorcycleControllerIDEncError,
				"motorcycle_id", createdMotorcycle.ID,
				"error", err,
				"client_ip", c.ClientIP())
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 7. Build response DTO (owner always sees all links on creation)
		response := ToMotorcycleResponse(createdMotorcycle)
		response.ID = encodedID

		// 8. Build HATEOAS links (Richardson Maturity Level 3)
		baseURL := GetBaseURL(c)
		response.Links = BuildMotorcycleDetailLinks(baseURL, encodedID, true)

		log.Success(logger.LogMotorcycleControllerRegSuccess,
			"motorcycle_id", createdMotorcycle.ID,
			"encoded_id", encodedID,
			"license_plate", createdMotorcycle.LicensePlate,
			"client_ip", c.ClientIP())

		// 9. Send success response (201 Created)
		h.Response.SuccessWithData(c, domain.MsgMotorcycleCreated, response)
	}
}
