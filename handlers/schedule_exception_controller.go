package handlers

import (
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// ============================================
// Schedule Exception DTOs (HU20-25)
// ============================================

// CreateScheduleExceptionRequest represents the request body for creating an exception
type CreateScheduleExceptionRequest struct {
	ExceptionDate string  `json:"exception_date" binding:"required"` // YYYY-MM-DD
	OpeningTime   *string `json:"opening_time"`                      // HH:mm
	ClosingTime   *string `json:"closing_time"`                      // HH:mm
	IsClosed      bool    `json:"is_closed"`
}

// UpdateScheduleExceptionRequest represents the request body for updating an exception
type UpdateScheduleExceptionRequest struct {
	OpeningTime *string `json:"opening_time"` // HH:mm
	ClosingTime *string `json:"closing_time"` // HH:mm
	IsClosed    bool    `json:"is_closed"`
}

// ScheduleExceptionResponse represents the response for a schedule exception
type ScheduleExceptionResponse struct {
	ID                     string `json:"id"`
	ScheduleID             string `json:"schedule_id"`
	ExceptionDate          string `json:"exception_date"`
	ExceptionDateFormatted string `json:"exception_date_formatted"`
	DayName                string `json:"day_name"`
	OpeningTime            string `json:"opening_time,omitempty"`
	ClosingTime            string `json:"closing_time,omitempty"`
	IsClosed               bool   `json:"is_closed"`
	Active                 bool   `json:"active"`
	Links                  []Link `json:"_links"`
}

// ScheduleExceptionListResponse represents the list response
type ScheduleExceptionListResponse struct {
	Exceptions []ScheduleExceptionResponse `json:"exceptions"`
	Links      []Link                      `json:"_links"`
}

// NewScheduleExceptionResponse creates a response from domain entity
func NewScheduleExceptionResponse(
	exception *domain.ScheduleDetail,
	encodedExceptionID, encodedScheduleID string,
	links []Link,
) ScheduleExceptionResponse {
	response := ScheduleExceptionResponse{
		ID:         encodedExceptionID,
		ScheduleID: encodedScheduleID,
		IsClosed:   exception.IsClosed,
		Active:     exception.Active,
		Links:      links,
	}

	if exception.ExceptionDate != nil {
		response.ExceptionDate = exception.ExceptionDate.Format("2006-01-02")
		response.ExceptionDateFormatted = formatDateSpanish(*exception.ExceptionDate)
		response.DayName = getDayNameSpanish(*exception.ExceptionDate)
	}

	if exception.OpeningTime != nil {
		response.OpeningTime = *exception.OpeningTime
	}
	if exception.ClosingTime != nil {
		response.ClosingTime = *exception.ClosingTime
	}

	return response
}

// formatDateSpanish formats a date in Spanish (e.g., "24 de Diciembre, 2026")
func formatDateSpanish(date time.Time) string {
	months := []string{
		"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
		"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
	}
	return date.Format("2") + " de " + months[date.Month()-1] + ", " + date.Format("2006")
}

// getDayNameSpanish returns the Spanish day name for a date
func getDayNameSpanish(date time.Time) string {
	days := []string{
		"Domingo", "Lunes", "Martes", "Miércoles",
		"Jueves", "Viernes", "Sábado",
	}
	return days[date.Weekday()]
}

// BuildScheduleExceptionLinks builds HATEOAS links for a schedule exception
func BuildScheduleExceptionLinks(baseURL, encodedBranchID, encodedExceptionID string) []Link {
	return []Link{
		{Rel: "self", Href: baseURL + "/schedule-exceptions/" + encodedExceptionID, Method: "GET"},
		{Rel: "update", Href: baseURL + "/schedule-exceptions/" + encodedExceptionID, Method: "PUT"},
		{Rel: "delete", Href: baseURL + "/schedule-exceptions/" + encodedExceptionID, Method: "DELETE"},
		{Rel: "activate", Href: baseURL + "/schedule-exceptions/" + encodedExceptionID + "/activate", Method: "PUT"},
		{Rel: "deactivate", Href: baseURL + "/schedule-exceptions/" + encodedExceptionID + "/deactivate", Method: "PUT"},
		{Rel: "branch", Href: baseURL + "/branches/" + encodedBranchID, Method: "GET"},
	}
}

// BuildScheduleExceptionListLinks builds HATEOAS links for exception list
func BuildScheduleExceptionListLinks(baseURL, encodedBranchID string) []Link {
	return []Link{
		{Rel: "self", Href: baseURL + "/branches/" + encodedBranchID + "/schedules/exceptions", Method: "GET"},
		{Rel: "create", Href: baseURL + "/branches/" + encodedBranchID + "/schedules/exceptions", Method: "POST"},
		{Rel: "branch", Href: baseURL + "/branches/" + encodedBranchID, Method: "GET"},
	}
}

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

		// 4. Get schedule for this branch
		schedule, err := scheduleInteractor.GetScheduleByBranchID(c.Request.Context(), branchID, person.ID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerCreateError, "error", err, "branch_id", branchID)
			switch err {
			case domain.ErrScheduleNotFound:
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case domain.ErrBranchNotFound:
				h.Response.Error(c, domain.MsgBranchNotFound)
			case domain.ErrForbidden:
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 5. Parse exception date
		exceptionDate, err := time.Parse("2006-01-02", req.ExceptionDate)
		if err != nil {
			log.Warn(logger.LogScheduleDetailControllerBindError, "error", err, "date", req.ExceptionDate)
			h.Response.Error(c, domain.MsgScheduleExceptionDatePast)
			return
		}

		// 6. Build domain object
		exception := domain.ScheduleDetail{
			ScheduleID:    schedule.ID,
			ExceptionDate: &exceptionDate,
			OpeningTime:   req.OpeningTime,
			ClosingTime:   req.ClosingTime,
			IsClosed:      req.IsClosed,
		}

		// 7. Create exception
		createdException, err := exceptionInteractor.CreateException(c.Request.Context(), exception, person.ID, branchID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerCreateError, "error", err, "schedule_id", schedule.ID)
			switch err {
			case domain.ErrScheduleNotFound:
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case domain.ErrScheduleExceptionDatePast:
				h.Response.Error(c, domain.MsgScheduleExceptionDatePast)
			case domain.ErrScheduleExceptionDateConflict:
				h.Response.Error(c, domain.MsgScheduleExceptionDateConflict)
			case domain.ErrScheduleExceptionInvalidTime:
				h.Response.Error(c, domain.MsgScheduleExceptionInvalidTime)
			case domain.ErrForbidden:
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
			"exception_date", req.ExceptionDate)

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

		// 3. Get schedule for this branch
		schedule, err := scheduleInteractor.GetScheduleByBranchID(c.Request.Context(), branchID, person.ID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerListError, "error", err, "branch_id", branchID)
			switch err {
			case domain.ErrScheduleNotFound:
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case domain.ErrBranchNotFound:
				h.Response.Error(c, domain.MsgBranchNotFound)
			case domain.ErrForbidden:
				h.Response.Error(c, domain.MsgForbidden)
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
			switch err {
			case domain.ErrScheduleExceptionNotFound:
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case domain.ErrScheduleExceptionInvalidTime:
				h.Response.Error(c, domain.MsgScheduleExceptionInvalidTime)
			case domain.ErrForbidden:
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
			switch err {
			case domain.ErrScheduleExceptionNotFound:
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case domain.ErrForbidden:
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
			switch err {
			case domain.ErrScheduleExceptionNotFound:
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case domain.ErrForbidden:
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
			switch err {
			case domain.ErrScheduleExceptionNotFound:
				h.Response.Error(c, domain.MsgScheduleExceptionNotFound)
			case domain.ErrForbidden:
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
