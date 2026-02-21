package handlers

import (
	"errors"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/constants"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// ============================================
// Schedule Controller Endpoints (HU30-35, HU10)
// ============================================

// CreateBranchSchedule handles POST /branches/:id/schedules (HU30)
// @Summary Create a schedule for a branch
// @Description Creates a new schedule configuration for a branch
// @Tags Schedules
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Success 201 {object} StandardResponse
// @Failure 400 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 409 {object} StandardResponse
// @Router /branches/{id}/schedules [post]
func (h *handler) CreateBranchSchedule(scheduleInteractor *interactor.ScheduleInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleControllerCreateRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			log.Warn(logger.LogScheduleControllerIDDecodeError, "encoded_id", encodedBranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 3. Create schedule
		schedule, err := scheduleInteractor.CreateSchedule(c.Request.Context(), branchID, person.ID)
		if err != nil {
			log.Error(logger.LogScheduleControllerCreateError, "error", err, "branch_id", branchID)
			switch {
			case errors.Is(err, domain.ErrScheduleAlreadyExists):
				h.Response.Error(c, domain.MsgScheduleAlreadyExists)
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Encode schedule ID for response
		encodedScheduleID, err := h.EncodeID(schedule.ID)
		if err != nil {
			h.HandleIDEncodingError(c, schedule.ID, err)
			return
		}

		// 5. Build HATEOAS response
		baseURL := GetBaseURL(c)
		links := BuildScheduleLinks(baseURL, encodedBranchID, encodedScheduleID)
		SetLocationHeader(c, baseURL, "branches/"+encodedBranchID+"/schedules", "")

		response := NewScheduleResponse(schedule, encodedScheduleID, encodedBranchID, links)

		log.Success(logger.LogScheduleControllerCreateOK,
			"schedule_id", schedule.ID,
			"branch_id", branchID)

		h.Response.SuccessWithData(c, domain.MsgScheduleCreated, response)
	}
}

// GetBranchSchedule handles GET /branches/:id/schedules (HU32)
// @Summary Get schedule for a branch
// @Description Retrieves the schedule configuration for a branch
// @Tags Schedules
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Success 200 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id}/schedules [get]
func (h *handler) GetBranchSchedule(scheduleInteractor *interactor.ScheduleInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleControllerGetRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path)

		// 1. Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			log.Warn(logger.LogScheduleControllerIDDecodeError, "encoded_id", encodedBranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 2. Get schedule
		schedule, err := scheduleInteractor.GetScheduleByBranchIDPublic(c.Request.Context(), branchID)
		if err != nil {
			log.Error(logger.LogScheduleControllerGetError, "error", err, "branch_id", branchID)
			if errors.Is(err, domain.ErrScheduleNotFound) {
				h.Response.Error(c, domain.MsgScheduleNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 3. Encode schedule ID
		encodedScheduleID, err := h.EncodeID(schedule.ID)
		if err != nil {
			h.HandleIDEncodingError(c, schedule.ID, err)
			return
		}

		// 4. Build HATEOAS response
		baseURL := GetBaseURL(c)
		links := BuildScheduleLinks(baseURL, encodedBranchID, encodedScheduleID)

		response := NewScheduleResponse(schedule, encodedScheduleID, encodedBranchID, links)

		log.Success(logger.LogScheduleControllerGetOK, "schedule_id", schedule.ID, "branch_id", branchID)

		h.Response.SuccessWithData(c, domain.MsgScheduleRetrieved, response)
	}
}

// UpdateBranchSchedule handles PUT /branches/:id/schedules (HU31)
// @Summary Update schedule for a branch
// @Description Updates the schedule configuration including validity dates
// @Tags Schedules
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Param schedule body UpdateScheduleRequest true "Schedule update data"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id}/schedules [put]
func (h *handler) UpdateBranchSchedule(scheduleInteractor *interactor.ScheduleInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleControllerUpdateRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			log.Warn(logger.LogScheduleControllerIDDecodeError, "encoded_id", encodedBranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 3. Parse request body
		var req UpdateScheduleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogScheduleControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgServerError)
			return
		}
		req.Sanitize()

		// 4. Get existing schedule
		schedule, err := h.fetchSchedulePublic(c, scheduleInteractor, branchID)
		if err != nil {
			return
		}

		// 5. Apply updates and parse dates
		if req.Active != nil {
			schedule.Active = *req.Active
		}
		if parseErr := parseScheduleDates(req, schedule, log); parseErr != nil {
			h.Response.Error(c, parseErr.msgCode)
			return
		}

		// 6. Update schedule
		if err := scheduleInteractor.UpdateSchedule(c.Request.Context(), *schedule, person.ID); err != nil {
			log.Error(logger.LogScheduleControllerUpdateError, "error", err, "schedule_id", schedule.ID)
			h.mapScheduleUpdateError(c, err)
			return
		}

		// 7. Encode schedule ID
		encodedScheduleID, err := h.EncodeID(schedule.ID)
		if err != nil {
			h.HandleIDEncodingError(c, schedule.ID, err)
			return
		}

		// 8. Build HATEOAS response
		baseURL := GetBaseURL(c)
		links := BuildScheduleLinks(baseURL, encodedBranchID, encodedScheduleID)

		response := NewScheduleResponse(schedule, encodedScheduleID, encodedBranchID, links)

		log.Success(logger.LogScheduleControllerUpdateOK, "schedule_id", schedule.ID, "branch_id", branchID)

		h.Response.SuccessWithData(c, domain.MsgScheduleUpdated, response)
	}
}

// fetchSchedulePublic retrieves a schedule by branch ID for public access.
// Responds with an error and returns nil on failure.
func (h *handler) fetchSchedulePublic(c *gin.Context, si *interactor.ScheduleInteractor, branchID string) (*domain.BranchSchedule, error) {
	schedule, err := si.GetScheduleByBranchIDPublic(c.Request.Context(), branchID)
	if err != nil {
		if errors.Is(err, domain.ErrScheduleNotFound) {
			h.Response.Error(c, domain.MsgScheduleNotFound)
		} else {
			h.Response.Error(c, domain.MsgServerError)
		}
		return nil, err
	}
	return schedule, nil
}

// mapScheduleUpdateError maps schedule update/delete errors to API responses.
func (h *handler) mapScheduleUpdateError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrScheduleNotFound):
		h.Response.Error(c, domain.MsgScheduleNotFound)
	case errors.Is(err, domain.ErrForbidden):
		h.Response.Error(c, domain.MsgForbidden)
	default:
		h.Response.Error(c, domain.MsgServerError)
	}
}

// @Summary Delete schedule for a branch
// @Description Removes the schedule configuration for a branch
// @Tags Schedules
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Success 200 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id}/schedules [delete]
func (h *handler) DeleteBranchSchedule(scheduleInteractor *interactor.ScheduleInteractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleControllerDeleteRequest, "client_ip", c.ClientIP())

		// 1. Get authenticated person
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 3. Get schedule first to get its ID
		schedule, err := scheduleInteractor.GetScheduleByBranchIDPublic(c.Request.Context(), branchID)
		if err != nil {
			if errors.Is(err, domain.ErrScheduleNotFound) {
				h.Response.Error(c, domain.MsgScheduleNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Delete schedule
		if err := scheduleInteractor.DeleteSchedule(c.Request.Context(), schedule.ID, person.ID); err != nil {
			log.Error(logger.LogScheduleControllerDeleteError, "error", err, "schedule_id", schedule.ID)
			switch {
			case errors.Is(err, domain.ErrScheduleNotFound):
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 5. Build response with next actions
		baseURL := GetBaseURL(c)
		response := NewScheduleDeleteResponse(baseURL, encodedBranchID)

		log.Success(logger.LogScheduleControllerDeleteOK, "schedule_id", schedule.ID)

		h.Response.SuccessWithData(c, domain.MsgScheduleDeleted, response)
	}
}

// ActivateBranchSchedule handles PUT /branches/:id/schedules/activate (HU34)
// @Summary Activate schedule for a branch
// @Description Activates the schedule configuration for a branch
// @Tags Schedules
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Success 200 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id}/schedules/activate [put]
func (h *handler) ActivateBranchSchedule(scheduleInteractor *interactor.ScheduleInteractor) gin.HandlerFunc {
	return h.toggleScheduleActive(scheduleInteractor, true)
}

// DeactivateBranchSchedule handles PUT /branches/:id/schedules/deactivate (HU35)
// @Summary Deactivate schedule for a branch
// @Description Deactivates the schedule configuration for a branch
// @Tags Schedules
// @Produce json
// @Security     BearerAuth
// @Param id path string true "Branch ID (encoded)"
// @Success 200 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id}/schedules/deactivate [put]
func (h *handler) DeactivateBranchSchedule(scheduleInteractor *interactor.ScheduleInteractor) gin.HandlerFunc {
	return h.toggleScheduleActive(scheduleInteractor, false)
}

// toggleScheduleActive is a shared helper for activate/deactivate schedule handlers.
func (h *handler) toggleScheduleActive(scheduleInteractor *interactor.ScheduleInteractor, activate bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		reqLogMsg := logger.LogScheduleControllerActivateReq
		errLogMsg := logger.LogScheduleControllerActivateError
		okLogMsg := logger.LogScheduleControllerActivateOK
		successMsg := domain.MsgScheduleActivated
		if !activate {
			reqLogMsg = logger.LogScheduleControllerDeactivateReq
			errLogMsg = logger.LogScheduleControllerDeactivateErr
			okLogMsg = logger.LogScheduleControllerDeactivateOK
			successMsg = domain.MsgScheduleDeactivated
		}

		log.Info(reqLogMsg, "client_ip", c.ClientIP())

		// 1. Get authenticated person
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 3. Get schedule to get its ID
		schedule, err := scheduleInteractor.GetScheduleByBranchIDPublic(c.Request.Context(), branchID)
		if err != nil {
			if errors.Is(err, domain.ErrScheduleNotFound) {
				h.Response.Error(c, domain.MsgScheduleNotFound)
			} else {
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Toggle schedule active status
		var toggleErr error
		if activate {
			toggleErr = scheduleInteractor.ActivateSchedule(c.Request.Context(), schedule.ID, person.ID)
		} else {
			toggleErr = scheduleInteractor.DeactivateSchedule(c.Request.Context(), schedule.ID, person.ID)
		}
		if toggleErr != nil {
			log.Error(errLogMsg, "error", toggleErr, "schedule_id", schedule.ID)
			switch {
			case errors.Is(toggleErr, domain.ErrScheduleNotFound):
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case errors.Is(toggleErr, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 5. Encode schedule ID
		encodedScheduleID, _ := h.EncodeID(schedule.ID)

		// 6. Build HATEOAS response
		baseURL := GetBaseURL(c)
		links := BuildScheduleLinks(baseURL, encodedBranchID, encodedScheduleID)

		schedule.Active = activate
		response := NewScheduleResponse(schedule, encodedScheduleID, encodedBranchID, links)

		log.Success(okLogMsg, "schedule_id", schedule.ID)

		h.Response.SuccessWithData(c, successMsg, response)
	}
}

// GetDaysOfWeek handles GET /schedules/days (HU10)
// @Summary Get catalog of days of week
// @Description Returns the catalog of available days of week for schedules
// @Tags Schedules
// @Produce json
// @Security     BearerAuth
// @Success 200 {object} StandardResponse
// @Router /schedules/days [get]
func (h *handler) GetDaysOfWeek() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleControllerGetDaysReq)

		// Use domain catalog
		days := domain.GetDayCatalog()
		baseURL := GetBaseURL(c)
		response := NewDaysCatalogResponse(days, baseURL)

		log.Success(logger.LogScheduleControllerGetDaysOK, "count", len(days))

		h.Response.SuccessWithData(c, domain.MsgDaysCatalogRetrieved, response)
	}
}

// ============================================
// HATEOAS Link Builders
// ============================================

// BuildScheduleLinks generates HATEOAS links for a schedule resource
func BuildScheduleLinks(baseURL, encodedBranchID, encodedScheduleID string) []Link {
	branchScheduleURL := BuildResourceURL(baseURL, "branches", encodedBranchID) + pathSchedules
	daysCatalogURL := baseURL + pathScheduleDays

	return []Link{
		{Rel: "self", Href: branchScheduleURL, Method: "GET"},
		{Rel: "delete", Href: branchScheduleURL, Method: "DELETE"},
		{Rel: "activate", Href: branchScheduleURL + "/activate", Method: "PUT"},
		{Rel: "deactivate", Href: branchScheduleURL + "/deactivate", Method: "PUT"},
		{Rel: "details-list", Href: branchScheduleURL + pathDetails, Method: "GET"},    // HU6-9: Schedule Details
		{Rel: "details-create", Href: branchScheduleURL + pathDetails, Method: "POST"}, // HU6: Create Detail
		{Rel: "branch", Href: BuildResourceURL(baseURL, "branches", encodedBranchID), Method: "GET"},
		{Rel: relDaysCatalog, Href: daysCatalogURL, Method: "GET"},
	}
}

// BuildScheduleDetailLinks generates HATEOAS links for a schedule detail resource (HU6-9)
// Note: GET /schedule-details/:id does not exist as a standalone route,
// so we omit the self-GET link. Update and delete are the primary actions.
func BuildScheduleDetailLinks(baseURL, encodedBranchID, encodedDetailID string) []Link {
	branchScheduleURL := BuildResourceURL(baseURL, "branches", encodedBranchID) + pathSchedules
	detailURL := BuildResourceURL(baseURL, "schedule-details", encodedDetailID)

	return []Link{
		{Rel: "update", Href: detailURL, Method: "PUT"},
		{Rel: "delete", Href: detailURL, Method: "DELETE"},
		{Rel: "schedule", Href: branchScheduleURL, Method: "GET"},
		{Rel: "details-list", Href: branchScheduleURL + pathDetails, Method: "GET"},
		{Rel: relDaysCatalog, Href: baseURL + pathScheduleDays, Method: "GET"},
	}
}

// BuildScheduleDetailListLinks generates HATEOAS links for the schedule details list (HU9)
func BuildScheduleDetailListLinks(baseURL, encodedBranchID string) []Link {
	branchScheduleURL := BuildResourceURL(baseURL, "branches", encodedBranchID) + pathSchedules

	return []Link{
		{Rel: "self", Href: branchScheduleURL + pathDetails, Method: "GET"},
		{Rel: "create", Href: branchScheduleURL + pathDetails, Method: "POST"},
		{Rel: "schedule", Href: branchScheduleURL, Method: "GET"},
		{Rel: "branch", Href: BuildResourceURL(baseURL, "branches", encodedBranchID), Method: "GET"},
		{Rel: relDaysCatalog, Href: baseURL + pathScheduleDays, Method: "GET"},
	}
}

// scheduleParseError wraps a message code for schedule date parse/validation errors.
type scheduleParseError struct {
	msgCode string
}

// parseScheduleDates parses start_date and end_date from the request and validates the date range.
func parseScheduleDates(req UpdateScheduleRequest, schedule *domain.BranchSchedule, log logger.Logger) *scheduleParseError {
	if req.StartDate != nil {
		parsed, err := time.Parse(constants.DateFormat, *req.StartDate)
		if err != nil {
			log.Warn(logger.LogScheduleControllerDateParseError, "field", "start_date", "error", err)
			return &scheduleParseError{msgCode: domain.MsgScheduleInvalidDateFormat}
		}
		schedule.StartDate = parsed
	}
	if req.EndDate != nil {
		parsed, err := time.Parse(constants.DateFormat, *req.EndDate)
		if err != nil {
			log.Warn(logger.LogScheduleControllerDateParseError, "field", "end_date", "error", err)
			return &scheduleParseError{msgCode: domain.MsgScheduleInvalidDateFormat}
		}
		schedule.EndDate = &parsed
	}

	// Validate date range
	if schedule.EndDate != nil && schedule.EndDate.Before(schedule.StartDate) {
		log.Warn(logger.LogScheduleControllerDateValidationError, "start", schedule.StartDate, "end", schedule.EndDate)
		return &scheduleParseError{msgCode: domain.MsgScheduleInvalidDateRange}
	}

	return nil
}
