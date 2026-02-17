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
// Schedule Exception Controller Endpoints (HU20-25)
// ============================================

// CreateScheduleException handles POST /branches/:id/schedules/exceptions (HU20)
func (h *handler) CreateScheduleException(
	exceptionInteractor *interactor.ScheduleExceptionInteractor,
	scheduleInteractor *interactor.ScheduleInteractor,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleDetailControllerCreateRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode branch ID
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerIDDecodeError, "encoded_id", encodedBranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 3. Parse request body
		var req CreateScheduleExceptionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogScheduleDetailControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Sanitize input
		req.Sanitize()

		// 4. Get schedule for this branch
		schedule, err := scheduleInteractor.GetScheduleByBranchID(c.Request.Context(), branchID, person.ID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerCreateError, "error", err, "branch_id", branchID)
			switch {
			case errors.Is(err, domain.ErrScheduleNotFound):
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 5. Parse exception dates (start date required, end date optional)
		// IMPORTANT: Use ParseInLocation with time.Local to match MySQL loc=Local config
		// Using time.Parse creates UTC dates which cause a 1-day shift when MySQL converts them
		exceptionStartDate, err := time.ParseInLocation(constants.DateFormat, req.ExceptionStartDate, time.Local)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerBindError, "error", err, "date", req.ExceptionStartDate)
			h.Response.Error(c, domain.MsgScheduleExceptionDatePast)
			return
		}

		// If end date not provided, use start date (single day exception)
		exceptionEndDate := exceptionStartDate
		if req.ExceptionEndDate != "" {
			exceptionEndDate, err = time.ParseInLocation(constants.DateFormat, req.ExceptionEndDate, time.Local)
			if err != nil {
				log.Warn(logger.LogScheduleDetailControllerBindError, "error", err, "date", req.ExceptionEndDate)
				h.Response.Error(c, domain.MsgScheduleExceptionDatePast)
				return
			}
		}

		// 6. Build domain object
		exception := domain.ScheduleDetail{
			ScheduleID:         schedule.ID,
			ExceptionStartDate: &exceptionStartDate,
			ExceptionEndDate:   &exceptionEndDate,
			OpeningTime:        req.OpeningTime,
			ClosingTime:        req.ClosingTime,
			IsClosed:           req.IsClosed,
		}

		// 7. Create exception
		createdException, err := exceptionInteractor.CreateException(c.Request.Context(), exception, person.ID, branchID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerCreateError, "error", err, "schedule_id", schedule.ID)
			switch {
			case errors.Is(err, domain.ErrScheduleNotFound):
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case errors.Is(err, domain.ErrScheduleExceptionDatePast):
				h.Response.Error(c, domain.MsgScheduleExceptionDatePast)
			case errors.Is(err, domain.ErrScheduleExceptionDateConflict):
				h.Response.Error(c, domain.MsgScheduleExceptionDateConflict)
			case errors.Is(err, domain.ErrScheduleExceptionInvalidTime):
				h.Response.Error(c, domain.MsgScheduleExceptionInvalidTime)
			case errors.Is(err, domain.ErrScheduleExceptionRedundant):
				h.Response.Error(c, domain.MsgScheduleExceptionRedundant)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 8. Encode IDs for response
		encodedExceptionID, err := h.EncodeID(createdException.ID)
		if err != nil {
			h.HandleIDEncodingError(c, createdException.ID, err)
			return
		}

		encodedScheduleID, err := h.EncodeID(createdException.ScheduleID)
		if err != nil {
			h.HandleIDEncodingError(c, createdException.ScheduleID, err)
			return
		}

		// 9. Build HATEOAS response
		baseURL := GetBaseURL(c)
		links := BuildScheduleExceptionLinks(baseURL, encodedBranchID, encodedExceptionID)
		SetLocationHeader(c, baseURL, "schedule-exceptions", encodedExceptionID)

		response := NewScheduleExceptionResponse(createdException, encodedExceptionID, encodedScheduleID, links)

		log.Success(logger.LogScheduleDetailControllerCreateOK,
			"exception_id", createdException.ID,
			"schedule_id", schedule.ID,
			"exception_start_date", req.ExceptionStartDate)

		h.Response.SuccessWithData(c, domain.MsgScheduleExceptionCreated, response)
	}
}

// ListScheduleExceptions handles GET /branches/:id/schedules/exceptions (HU23)
func (h *handler) ListScheduleExceptions(
	exceptionInteractor *interactor.ScheduleExceptionInteractor,
	scheduleInteractor *interactor.ScheduleInteractor,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleDetailControllerListRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Decode branch ID (no ownership check needed for public viewing)
		encodedBranchID := c.Param("id")
		branchID, err := h.DecodeID(encodedBranchID)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerIDDecodeError, "encoded_id", encodedBranchID, "error", err)
			h.Response.Error(c, domain.MsgBranchNotFound)
			return
		}

		// 2. Get schedule for this branch (public access - no ownership validation)
		schedule, err := scheduleInteractor.GetScheduleByBranchIDPublic(c.Request.Context(), branchID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerListError, "error", err, "branch_id", branchID)
			switch {
			case errors.Is(err, domain.ErrScheduleNotFound):
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case errors.Is(err, domain.ErrBranchNotFound):
				h.Response.Error(c, domain.MsgBranchNotFound)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 4. Get exceptions for this schedule
		exceptions, err := exceptionInteractor.ListExceptions(c.Request.Context(), schedule.ID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerListError, "error", err, "schedule_id", schedule.ID)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 5. Build response
		encodedScheduleID, _ := h.EncodeID(schedule.ID)
		baseURL := GetBaseURL(c)

		var exceptionResponses []ScheduleExceptionResponse
		for _, exception := range exceptions {
			encodedExceptionID, err := h.EncodeID(exception.ID)
			if err != nil {
				continue
			}
			links := BuildScheduleExceptionLinks(baseURL, encodedBranchID, encodedExceptionID)
			exceptionResponses = append(exceptionResponses, NewScheduleExceptionResponse(&exception, encodedExceptionID, encodedScheduleID, links))
		}

		listLinks := BuildScheduleExceptionListLinks(baseURL, encodedBranchID)
		response := ScheduleExceptionListResponse{
			Exceptions: exceptionResponses,
			Links:      listLinks,
		}

		log.Success(logger.LogScheduleDetailControllerListOK,
			"schedule_id", schedule.ID,
			"count", len(exceptions))

		h.Response.SuccessWithData(c, domain.MsgScheduleExceptionsListed, response)
	}
}

// UpdateScheduleException handles PUT /schedule-exceptions/:id (HU21)
func (h *handler) UpdateScheduleException(
	exceptionInteractor *interactor.ScheduleExceptionInteractor,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleDetailControllerUpdateRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode exception ID
		encodedExceptionID := c.Param("id")
		exceptionID, err := h.DecodeID(encodedExceptionID)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerIDDecodeError, "encoded_id", encodedExceptionID, "error", err)
			h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			return
		}

		// 3. Parse request body
		var req UpdateScheduleExceptionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogScheduleDetailControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}

		// Sanitize input
		req.Sanitize()

		// 4. Build domain object
		exception := domain.ScheduleDetail{
			ID:          exceptionID,
			OpeningTime: req.OpeningTime,
			ClosingTime: req.ClosingTime,
			IsClosed:    req.IsClosed,
		}

		// 5. Update exception
		if err := exceptionInteractor.UpdateException(c.Request.Context(), exception, person.ID); err != nil {
			log.Error(logger.LogScheduleDetailControllerUpdateError, "error", err, "exception_id", exceptionID)
			switch {
			case errors.Is(err, domain.ErrScheduleExceptionNotFound):
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case errors.Is(err, domain.ErrScheduleExceptionInvalidTime):
				h.Response.Error(c, domain.MsgScheduleExceptionInvalidTime)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogScheduleDetailControllerUpdateOK, "exception_id", exceptionID)
		h.Response.Success(c, domain.MsgScheduleExceptionUpdated)
	}
}

// DeleteScheduleException handles DELETE /schedule-exceptions/:id (HU22)
func (h *handler) DeleteScheduleException(
	exceptionInteractor *interactor.ScheduleExceptionInteractor,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleDetailControllerDeleteRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode exception ID
		encodedExceptionID := c.Param("id")
		exceptionID, err := h.DecodeID(encodedExceptionID)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerIDDecodeError, "encoded_id", encodedExceptionID, "error", err)
			h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			return
		}

		// 3. Delete exception
		if err := exceptionInteractor.DeleteException(c.Request.Context(), exceptionID, person.ID); err != nil {
			log.Error(logger.LogScheduleDetailControllerDeleteError, "error", err, "exception_id", exceptionID)
			switch {
			case errors.Is(err, domain.ErrScheduleExceptionNotFound):
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogScheduleDetailControllerDeleteOK, "exception_id", exceptionID)
		h.Response.Success(c, domain.MsgScheduleExceptionDeleted)
	}
}

// ActivateScheduleException handles PUT /schedule-exceptions/:id/activate (HU24)
func (h *handler) ActivateScheduleException(
	exceptionInteractor *interactor.ScheduleExceptionInteractor,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleDetailControllerUpdateRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode exception ID
		encodedExceptionID := c.Param("id")
		exceptionID, err := h.DecodeID(encodedExceptionID)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerIDDecodeError, "encoded_id", encodedExceptionID, "error", err)
			h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			return
		}

		// 3. Activate exception
		if err := exceptionInteractor.ActivateException(c.Request.Context(), exceptionID, person.ID); err != nil {
			log.Error(logger.LogScheduleDetailControllerUpdateError, "error", err, "exception_id", exceptionID)
			switch {
			case errors.Is(err, domain.ErrScheduleExceptionNotFound):
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogScheduleDetailControllerUpdateOK, "exception_id", exceptionID)
		h.Response.Success(c, domain.MsgScheduleExceptionActivated)
	}
}

// DeactivateScheduleException handles PUT /schedule-exceptions/:id/deactivate (HU25)
func (h *handler) DeactivateScheduleException(
	exceptionInteractor *interactor.ScheduleExceptionInteractor,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info(logger.LogScheduleDetailControllerUpdateRequest,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// 1. Get authenticated person from context
		person, _ := middleware.GetAuthenticatedUser(c)

		// 2. Decode exception ID
		encodedExceptionID := c.Param("id")
		exceptionID, err := h.DecodeID(encodedExceptionID)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerIDDecodeError, "encoded_id", encodedExceptionID, "error", err)
			h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			return
		}

		// 3. Deactivate exception
		if err := exceptionInteractor.DeactivateException(c.Request.Context(), exceptionID, person.ID); err != nil {
			log.Error(logger.LogScheduleDetailControllerUpdateError, "error", err, "exception_id", exceptionID)
			switch {
			case errors.Is(err, domain.ErrScheduleExceptionNotFound):
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case errors.Is(err, domain.ErrForbidden):
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		log.Success(logger.LogScheduleDetailControllerUpdateOK, "exception_id", exceptionID)
		h.Response.Success(c, domain.MsgScheduleExceptionDeactivated)
	}
}
