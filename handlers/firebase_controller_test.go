package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestGetFirebaseToken_Integration_Success validates the full HTTP handler pipeline
// for the success path of GetFirebaseToken (GET /auth/firebase-token).
//
// It exercises: auth context extraction → Keycloak UID retrieval →
// FirebaseClient.CreateCustomToken → response serialization → 200 response.
func TestGetFirebaseToken_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Mock Firebase client
	mockFirebase := new(mocks.MockFirebaseClient)

	// Handler with Firebase client injected
	h := handlers.NewForTestWithFirebase(mockFirebase, encoder, responseHandler)

	// Test person data
	keycloakUID := "kc-user-abc-123"
	person := &domain.Person{
		ID:             "a1111111-1111-4000-8000-111111111111",
		Email:          "juan@example.com",
		KeycloakUserID: keycloakUID,
	}

	// Mock: Firebase generates a custom token
	expectedToken := "firebase-custom-token-abc123"
	mockFirebase.On("CreateCustomToken", mock.Anything, keycloakUID).Return(expectedToken, nil)

	// Setup router
	router := gin.New()
	router.GET("/auth/firebase-token", func(c *gin.Context) {
		c.Set("authenticated_user", person)
	}, h.GetFirebaseToken())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/firebase-token", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "GEN_OPE_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, expectedToken, data["firebase_token"])

	mockFirebase.AssertExpectations(t)
}
