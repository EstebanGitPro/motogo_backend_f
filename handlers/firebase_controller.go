package handlers

import (
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	"github.com/gin-gonic/gin"
)

// FirebaseTokenResponse represents the response for GET /auth/firebase-token
type FirebaseTokenResponse struct {
	FirebaseToken string `json:"firebase_token"`
}

// GetFirebaseToken generates a Firebase custom token for authenticated users
// @Summary Get Firebase custom token
// @Description Generates a Firebase custom token using the authenticated user's Keycloak ID
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse{data=FirebaseTokenResponse}
// @Failure 401 {object} StandardResponse
// @Failure 500 {object} StandardResponse
// @Router /auth/firebase-token [get]
func (h *handler) GetFirebaseToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := middleware.GetRequestID(c)
		log := Logger.WithTraceID(traceID)

		log.Info("Firebase token request received",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP())

		// Get authenticated person from context (set by auth middleware)
		person, exists := middleware.GetAuthenticatedUser(c)
		if !exists || person == nil {
			log.Warn("Unauthenticated request for Firebase token")
			h.Response.Error(c, domain.MsgUnauthorized)
			return
		}

		// Use Keycloak user ID as Firebase UID
		keycloakUID := person.KeycloakUserID
		if keycloakUID == "" {
			log.Error(logger.LogKeycloakUserNotFound, "person_id", person.ID)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		// Generate Firebase custom token
		token, err := h.FirebaseClient.CreateCustomToken(c.Request.Context(), keycloakUID)
		if err != nil {
			log.Error("Error generating Firebase token", "error", err, "keycloak_uid", keycloakUID)
			h.Response.Error(c, domain.MsgServerError)
			return
		}

		log.Success("Firebase token generated",
			"keycloak_uid", keycloakUID,
			"client_ip", c.ClientIP())

		response := FirebaseTokenResponse{
			FirebaseToken: token,
		}

		h.Response.SuccessWithData(c, domain.MsgOpSuccess, response)
	}
}
