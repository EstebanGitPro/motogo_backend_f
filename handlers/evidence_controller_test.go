package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreateEvidence_Integration_Success validates the full HTTP pipeline
// for POST /motorcycles/:id/evidence (HU16).
//
// Exercises: auth → ID decoding → JSON binding → sanitize → interactor →
// ID encoding → HATEOAS → 201 response.
func TestCreateEvidence_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockEvidenceInteractor := new(mocks.MockEvidenceInteractor)

	h := handlers.NewForTest(
		nil, nil, nil,
		mockEvidenceInteractor,
		nil,
		encoder,
		responseHandler,
	)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	evidenceUUID := "a3333333-3333-4000-8000-333333333333"
	angle := domain.EvidenceAngle("FRONTAL")

	encodedMotorcycleID, err := encoder.Encode(motorcycleUUID)
	assert.NoError(t, err)
	assert.NotEmpty(t, encodedMotorcycleID)

	createdEvidence := &domain.MotorcycleEvidence{
		ID:           evidenceUUID,
		MotorcycleID: motorcycleUUID,
		ImageURL:     "https://firebasestorage.googleapis.com/v0/b/test/evidence.jpg",
		Angle:        &angle,
		CreatedAt:    time.Now(),
	}

	mockEvidenceInteractor.On("CreateEvidence",
		mock.Anything,
		motorcycleUUID,
		ownerID,
		mock.AnythingOfType("*domain.MotorcycleEvidence"),
	).Return(createdEvidence, nil)

	reqBody := map[string]interface{}{
		"image_url": "  https://firebasestorage.googleapis.com/v0/b/test/evidence.jpg  ",
		"angle":     "FRONTAL",
	}
	bodyJSON, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("authenticated_user", &domain.Person{
			ID:   ownerID,
			Role: "USUARIO",
		})
		c.Next()
	})
	router.POST("/motorcycles/:id/evidence", h.CreateEvidence())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/motorcycles/"+encodedMotorcycleID+"/evidence", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	assert.NotEmpty(t, data["id"])
	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])
	assert.Equal(t, "FRONTAL", data["angle"])
	assert.NotEmpty(t, data["image_url"])
	assert.NotEmpty(t, data["created_at"])

	links, ok := data["_links"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, links)

	mockEvidenceInteractor.AssertExpectations(t)
}
