package handlers_test

import (
	"bytes"
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

// TestUpdateProfileImage_Integration_Success validates the full HTTP pipeline
// for PUT /motorcycles/:id/profile-image (HU36/37).
//
// Exercises: auth → ID decoding → JSON binding → sanitize → interactor →
// HATEOAS → 200 response.
func TestUpdateProfileImage_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockMotorcycleInteractor := new(mocks.MockMotorcycleInteractor)

	h := handlers.NewForTest(
		nil, nil,
		mockMotorcycleInteractor,
		nil, nil,
		encoder,
		responseHandler,
	)

	ownerID := "a1111111-1111-4000-8000-111111111111"
	motorcycleUUID := "a2222222-2222-4000-8000-222222222222"
	imageURL := "https://firebasestorage.googleapis.com/v0/b/test/profile.jpg"

	encodedMotorcycleID, err := encoder.Encode(motorcycleUUID)
	assert.NoError(t, err)
	assert.NotEmpty(t, encodedMotorcycleID)

	updatedMotorcycle := &domain.Motorcycle{
		ID:              motorcycleUUID,
		LicensePlate:    "ABC123",
		OwnerID:         ownerID,
		ProfileImageURL: &imageURL,
	}

	mockMotorcycleInteractor.On("UpdateMotorcycle",
		mock.Anything,
		motorcycleUUID,
		ownerID,
		mock.AnythingOfType("*domain.Motorcycle"),
	).Return(updatedMotorcycle, nil)

	reqBody := map[string]interface{}{
		"image_url": "  https://firebasestorage.googleapis.com/v0/b/test/profile.jpg  ",
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
	router.PUT("/motorcycles/:id/profile-image", h.UpdateProfileImage())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/motorcycles/"+encodedMotorcycleID+"/profile-image", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, encodedMotorcycleID, data["motorcycle_id"])
	assert.Equal(t, imageURL, data["profile_image_url"])

	links, ok := data["_links"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, links)

	mockMotorcycleInteractor.AssertExpectations(t)
}
