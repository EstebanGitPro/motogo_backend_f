package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/dto"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/EstebanGitPro/motogo-backend/middleware"
	"github.com/EstebanGitPro/motogo-backend/mocks"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestRegisterPerson_Integration_Success validates the full HTTP handler pipeline
// for the success path of RegisterPerson (POST /persons).
//
// It exercises: JSON binding → sanitization → Interactor.RegisterPerson
// (validate → tx → save DB → Keycloak create → set password → assign role →
// update keycloak ID → commit → send verification email) → ID encoding →
// 201 response.
func TestRegisterPerson_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	// Mock dependencies
	mockPersonService := new(mocks.MockPersonService)
	mockTx := new(mocks.MockTx)

	personInteractor := interactor.NewInteractor(mockPersonService)
	h := handlers.NewForTestWithPerson(personInteractor, nil, encoder, responseHandler)

	keycloakUID := "kc-new-user-123"

	// Mock: Step 1 validate (RegisterPerson on service)
	mockPersonService.On("RegisterPerson", mock.Anything, mock.AnythingOfType("domain.Person")).
		Return(&dto.RegistrationResult{}, nil)

	// Mock: Step 1.5 check consistent state
	mockPersonService.On("CheckAndCleanInconsistentState", mock.Anything, "juan@example.com").Return(nil)

	// Mock: Step 2 begin transaction
	mockPersonService.On("BeginTx", mock.Anything).Return(mockTx, nil)

	// Mock: Step 3 save to DB
	mockPersonService.On("SavePersonToDB", mock.Anything, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)

	// Mock: Step 4 create Keycloak user
	mockPersonService.On("CreateUserInKeycloak", mock.Anything, mock.AnythingOfType("*domain.Person")).
		Return(keycloakUID, nil)

	// Mock: Step 5 set password
	mockPersonService.On("SetUserPassword", mock.Anything, keycloakUID, "Password123!").Return(nil)

	// Mock: Step 6 assign role
	mockPersonService.On("AssignUserRole", mock.Anything, keycloakUID, "representative").Return(nil)

	// Mock: Step 7 update DB with keycloak ID
	mockPersonService.On("UpdatePersonKeycloakID", mock.Anything, mockTx, mock.AnythingOfType("string"), keycloakUID).Return(nil)

	// Mock: Step 8 commit transaction
	mockTx.On("Commit").Return(nil)

	// Mock: Step 9 send verification email (non-blocking)
	mockPersonService.On("SendVerificationEmail", mock.Anything, keycloakUID).Return(nil)

	// Request body
	reqBody := map[string]interface{}{
		"identity_number": "1234567890",
		"first_name":      "Juan",
		"last_name":       "Pérez",
		"email":           "juan@example.com",
		"phone_number":    "3001234567",
		"password":        "Password123!",
		"role":            "representative",
	}
	body, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/persons", h.RegisterPerson())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/persons", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_U_REG_EXI_00001", response["code"])

	mockPersonService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestLogin_Integration_Success validates the full HTTP handler pipeline
// for the success path of Login (POST /auth/login).
//
// It exercises: JSON binding → sanitization → Interactor.Login →
// token response mapping → HATEOAS links → 200 response.
func TestLogin_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockPersonService := new(mocks.MockPersonService)
	personInteractor := interactor.NewInteractor(mockPersonService)
	h := handlers.NewForTestWithPerson(personInteractor, nil, encoder, responseHandler)

	// Mock: Login returns Keycloak JWT
	mockPersonService.On("Login", mock.Anything, "juan@example.com", "Password123!").
		Return(&gocloak.JWT{
			AccessToken:  "access-token-abc",
			RefreshToken: "refresh-token-xyz",
			ExpiresIn:    300,
			TokenType:    "Bearer",
		}, nil)

	reqBody := map[string]interface{}{
		"email":    "juan@example.com",
		"password": "Password123!",
	}
	body, _ := json.Marshal(reqBody)

	router := gin.New()
	router.POST("/auth/login", h.Login())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_AUTH_LOGIN_SUCCESS_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "access-token-abc", data["access_token"])
	assert.Equal(t, "refresh-token-xyz", data["refresh_token"])
	assert.Equal(t, "Bearer", data["token_type"])
	assert.NotEmpty(t, data["_links"])

	mockPersonService.AssertExpectations(t)
}

// TestUpdateProfile_Integration_Success validates the full HTTP handler pipeline
// for the success path of UpdateProfile (PUT /persons/me).
//
// It exercises: auth context extraction → JSON binding → field merging →
// Interactor.UpdateProfile (tx → update → commit) → ID encoding →
// HATEOAS links → 200 response.
func TestUpdateProfile_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockPersonService := new(mocks.MockPersonService)
	mockTx := new(mocks.MockTx)
	personInteractor := interactor.NewInteractor(mockPersonService)
	h := handlers.NewForTestWithPerson(personInteractor, nil, encoder, responseHandler)

	// Existing person in context
	person := &domain.Person{
		ID:             "a1111111-1111-4000-8000-111111111111",
		IdentityNumber: "1234567890",
		Email:          "juan@example.com",
		FirstName:      "Juan",
		LastName:       "Pérez",
		Role:           domain.RoleRepresentative,
		KeycloakUserID: "kc-user-abc",
	}

	// Mock: transaction lifecycle
	mockPersonService.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockPersonService.On("UpdatePersonProfile", mock.Anything, mockTx, mock.AnythingOfType("domain.Person")).Return(nil)
	mockTx.On("Commit").Return(nil)

	// Update only first name and phone
	reqBody := map[string]interface{}{
		"first_name":   "Carlos",
		"phone_number": "3109876543",
	}
	body, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/persons/me", func(c *gin.Context) {
		c.Set("authenticated_user", person)
	}, h.UpdateProfile())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/persons/me", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_P_UPD_EXI_00002", response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, "Carlos", data["first_name"])
	assert.Equal(t, "3109876543", data["phone_number"])
	assert.NotEmpty(t, data["_links"])

	mockPersonService.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

// TestChangePassword_Integration_Success validates the full HTTP handler pipeline
// for the success path of ChangePassword (PUT /persons/me/password).
//
// It exercises: auth context extraction → JSON binding → Interactor.ChangePassword →
// HATEOAS links → 200 response.
func TestChangePassword_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockPersonService := new(mocks.MockPersonService)
	personInteractor := interactor.NewInteractor(mockPersonService)
	h := handlers.NewForTestWithPerson(personInteractor, nil, encoder, responseHandler)

	person := &domain.Person{
		ID:             "a1111111-1111-4000-8000-111111111111",
		KeycloakUserID: "kc-user-abc",
	}

	// Mock: ChangePassword succeeds
	mockPersonService.On("ChangePassword", mock.Anything, "kc-user-abc", "OldPassword123!", "NewPassword456!").Return(nil)

	reqBody := map[string]interface{}{
		"current_password": "OldPassword123!",
		"new_password":     "NewPassword456!",
	}
	body, _ := json.Marshal(reqBody)

	router := gin.New()
	router.PUT("/persons/me/password", func(c *gin.Context) {
		c.Set("authenticated_user", person)
	}, h.ChangePassword())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/persons/me/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_P_CHANGE_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.NotEmpty(t, data["_links"])

	mockPersonService.AssertExpectations(t)
}

// TestDeleteSelf_Integration_Success validates the full HTTP handler pipeline
// for the success path of DeleteSelf (DELETE /persons/me).
//
// It exercises: auth context → branch check (no branches) →
// Interactor.DeleteKeycloakUser → Interactor.DeletePersonFromDB → 200 response.
func TestDeleteSelf_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockPersonService := new(mocks.MockPersonService)
	mockBranchService := new(mocks.MockBranchService)

	personInteractor := interactor.NewInteractor(mockPersonService)
	branchInteractor := interactor.NewBranchInteractor(mockBranchService)
	h := handlers.NewForTestWithPerson(personInteractor, branchInteractor, encoder, responseHandler)

	person := &domain.Person{
		ID:             "a1111111-1111-4000-8000-111111111111",
		Email:          "juan@example.com",
		KeycloakUserID: "kc-user-abc",
	}

	// Mock: no branches owned by this person
	mockBranchService.On("GetBranchesByRepresentative", mock.Anything, person.ID).
		Return([]domain.Branch{}, nil)

	// Mock: delete from Keycloak
	mockPersonService.On("RollbackKeycloakUser", mock.Anything, person.KeycloakUserID).Return(nil)

	// Mock: delete from DB
	mockPersonService.On("RollbackPerson", mock.Anything, person.ID).Return(nil)

	router := gin.New()
	router.DELETE("/persons/me", func(c *gin.Context) {
		c.Set("authenticated_user", person)
	}, h.DeleteSelf())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/persons/me", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_P_DEL_EXI_00001", response["code"])

	mockPersonService.AssertExpectations(t)
	mockBranchService.AssertExpectations(t)
}

// TestGetPublicContact_Integration_Success validates the full HTTP handler pipeline
// for the success path of GetPublicContact (GET /persons/:id/contact).
//
// It exercises: ID decoding → Interactor.GetPublicContact → phone number response →
// HATEOAS links → 200 response.
func TestGetPublicContact_Integration_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	encoder := createTestEncoder()
	msgCache := createTestMessageCache()
	responseHandler := middleware.NewResponseHandler(msgCache)

	mockPersonService := new(mocks.MockPersonService)
	personInteractor := interactor.NewInteractor(mockPersonService)
	h := handlers.NewForTestWithPerson(personInteractor, nil, encoder, responseHandler)

	personID := "a1111111-1111-4000-8000-111111111111"
	encodedID, _ := encoder.Encode(personID)

	// Mock: get person by ID
	mockPersonService.On("GetPersonByID", mock.Anything, personID).
		Return(&domain.Person{
			ID:          personID,
			PhoneNumber: "3001234567",
		}, nil)

	router := gin.New()
	router.GET("/persons/:id/contact", h.GetPublicContact())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/persons/"+encodedID+"/contact", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "MOD_P_CONTACT_EXI_00001", response["code"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "3001234567", data["phone_number"])
	assert.NotEmpty(t, data["_links"])

	mockPersonService.AssertExpectations(t)
}
