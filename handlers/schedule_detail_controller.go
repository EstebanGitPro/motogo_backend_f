package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// ============================================
// Schedule Detail Controller Endpoints (HU6-9)
// ============================================

// CreateScheduleDetail handles POST /branches/:id/schedules/details (HU6)
// @Summary Create a time slot for a branch schedule
// @Description Creates a new time slot (schedule detail) for a specific day of the week
// @Tags Schedule Details
// @Accept json
// @Produce json
// @Param id path string true "Branch ID (encoded)"
// @Param body body CreateScheduleDetailRequest true "Schedule detail data"
// @Success 201 {object} StandardResponse
// @Failure 400 {object} StandardResponse
// @Failure 403 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Failure 409 {object} StandardResponse
// @Router /branches/{id}/schedules/details [post]
func (h *handler) CreateScheduleDetail(
	scheduleDetailInteractor *interactor.ScheduleDetailInteractor,
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
		var req CreateScheduleDetailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Warn(logger.LogScheduleDetailControllerBindError, "error", err)
			h.Response.Error(c, domain.MsgValBadFormat)
			return
		}
		req.Sanitize()

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

		// 5. Build domain object
		detail := domain.ScheduleDetail{
			ScheduleID:  schedule.ID,
			DayOfWeek:   &req.DayOfWeek,
			OpeningTime: req.OpeningTime,
			ClosingTime: req.ClosingTime,
			IsClosed:    req.IsClosed,
		}

		// 6. Create schedule detail
		createdDetail, err := scheduleDetailInteractor.CreateDetail(c.Request.Context(), detail, person.ID, branchID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerCreateError, "error", err, "schedule_id", schedule.ID)
			switch err {
			case domain.ErrScheduleNotFound:
				h.Response.Error(c, domain.MsgScheduleNotFound)
			case domain.ErrScheduleDetailInvalidDay:
				h.Response.Error(c, domain.MsgScheduleDetailInvalidDay)
			case domain.ErrScheduleDetailInvalidTime:
				h.Response.Error(c, domain.MsgScheduleDetailInvalidTime)
			case domain.ErrScheduleDetailTimeConflict:
				h.Response.Error(c, domain.MsgScheduleDetailTimeConflict)
			case domain.ErrBranchNotFound:
				h.Response.Error(c, domain.MsgBranchNotFound)
			case domain.ErrForbidden:
				h.Response.Error(c, domain.MsgForbidden)
			default:
				h.Response.Error(c, domain.MsgServerError)
			}
			return
		}

		// 7. Encode detail ID for response
		encodedDetailID, err := h.EncodeID(createdDetail.ID)
		if err != nil {
			h.HandleIDEncodingError(c, createdDetail.ID, err)
			return
		}

		encodedScheduleID, err := h.EncodeID(createdDetail.ScheduleID)
		if err != nil {
			h.HandleIDEncodingError(c, createdDetail.ScheduleID, err)
			return
		}

		// 8. Build HATEOAS response
		baseURL := GetBaseURL(c)
		links := BuildScheduleDetailLinks(baseURL, encodedBranchID, encodedDetailID)
		SetLocationHeader(c, baseURL, "schedule-details", encodedDetailID)

		response := NewScheduleDetailResponse(createdDetail, encodedDetailID, encodedScheduleID, links)

		log.Success(logger.LogScheduleDetailControllerCreateOK,
			"detail_id", createdDetail.ID,
			"schedule_id", schedule.ID,
			"day_of_week", req.DayOfWeek)

		h.Response.SuccessWithData(c, domain.MsgScheduleDetailCreated, response)
	}
}

// ListScheduleDetails handles GET /branches/:id/schedules/details (HU9)
// @Summary List time slots for a branch schedule
// @Description Retrieves all time slots (schedule details) for a branch
// @Tags Schedule Details
// @Produce json
// @Param id path string true "Branch ID (encoded)"
// @Success 200 {object} StandardResponse
// @Failure 404 {object} StandardResponse
// @Router /branches/{id}/schedules/details [get]
func (h *handler) ListScheduleDetails(
	scheduleDetailInteractor *interactor.ScheduleDetailInteractor,
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

		// 4. Get details for this schedule
		details, err := scheduleDetailInteractor.ListDetails(c.Request.Context(), schedule.ID)
		if err != nil {
			log.Error(logger.LogScheduleDetailControllerListError, "error", err, "schedule_id", schedule.ID)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// 5. Encode details for response
		encodedScheduleID, _ := h.EncodeID(schedule.ID)
		baseURL := GetBaseURL(c)

		var detailResponses []ScheduleDetailResponse
		for _, detail := range details {
			encodedDetailID, err := h.EncodeID(detail.ID)
			if err != nil {
				continue // Skip on encoding error
			}
			links := BuildScheduleDetailLinks(baseURL, encodedBranchID, encodedDetailID)
			detailResponses = append(detailResponses, NewScheduleDetailResponse(&detail, encodedDetailID, encodedScheduleID, links))
		}

		// 6. Build list response with HATEOAS
		listLinks := BuildScheduleDetailListLinks(baseURL, encodedBranchID)
		response := NewScheduleDetailListResponse(detailResponses, listLinks)

		log.Success(logger.LogScheduleDetailControllerListOK,
			"schedule_id", schedule.ID,
			"count", len(details))

		h.Response.SuccessWithData(c, domain.MsgScheduleDetailsListed, response)
	}
}
