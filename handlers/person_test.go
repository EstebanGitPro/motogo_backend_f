package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPersonRequest_ToDomain(t *testing.T) {
	// Arrange
	req := handlers.PersonRequest{
		IdentityNumber: "12345678",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		Email:          "john@example.com",
		PhoneNumber:    "3001234567",
		Password:       "securepass123",
		Role:           "customer",
	}

	// Act
	domain := req.ToDomain()

	// Assert
	assert.Equal(t, req.IdentityNumber, domain.IdentityNumber)
	assert.Equal(t, req.FirstName, domain.FirstName)
	assert.Equal(t, req.LastName, domain.LastName)
	assert.Equal(t, req.SecondLastName, domain.SecondLastName)
	assert.Equal(t, req.Email, domain.Email)
	assert.Equal(t, req.PhoneNumber, domain.PhoneNumber)
	assert.Equal(t, req.Password, domain.Password)
	assert.Equal(t, req.Role, string(domain.Role))
}

func TestLoginRequest_Binding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        map[string]string
		expectError bool
	}{
		{
			name: "valid request",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectError: false,
		},
		{
			name: "missing email",
			body: map[string]string{
				"password": "password123",
			},
			expectError: true,
		},
		{
			name: "missing password",
			body: map[string]string{
				"email": "test@example.com",
			},
			expectError: true,
		},
		{
			name: "invalid email format",
			body: map[string]string{
				"email":    "not-an-email",
				"password": "password123",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			var req handlers.LoginRequest
			err := c.ShouldBindJSON(&req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.body["email"], req.Email)
				assert.Equal(t, tt.body["password"], req.Password)
			}
		})
	}
}

func TestResendVerificationEmailRequest_Binding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        map[string]string
		expectError bool
	}{
		{
			name:        "valid request",
			body:        map[string]string{"email": "test@example.com"},
			expectError: false,
		},
		{
			name:        "missing email",
			body:        map[string]string{},
			expectError: true,
		},
		{
			name:        "invalid email",
			body:        map[string]string{"email": "invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("POST", "/resend", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			var req handlers.ResendVerificationEmailRequest
			err := c.ShouldBindJSON(&req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyEmailRequest_Binding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        map[string]string
		expectError bool
	}{
		{
			name:        "valid request",
			body:        map[string]string{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
			expectError: false,
		},
		{
			name:        "missing token",
			body:        map[string]string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("POST", "/verify", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			var req handlers.VerifyEmailRequest
			err := c.ShouldBindJSON(&req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResetPasswordWithTokenRequest_Binding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        map[string]string
		expectError bool
	}{
		{
			name: "valid request",
			body: map[string]string{
				"token":        "valid-token",
				"new_password": "newpassword123",
			},
			expectError: false,
		},
		{
			name: "missing token",
			body: map[string]string{
				"new_password": "newpassword123",
			},
			expectError: true,
		},
		{
			name: "missing password",
			body: map[string]string{
				"token": "valid-token",
			},
			expectError: true,
		},
		{
			name: "password too short",
			body: map[string]string{
				"token":        "valid-token",
				"new_password": "short",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest("POST", "/reset", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			var req handlers.ResetPasswordWithTokenRequest
			err := c.ShouldBindJSON(&req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPersonResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := handlers.PersonResponse{
		IdentityNumber: "12345678",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		Email:          "john@example.com",
		PhoneNumber:    "3001234567",
		Role:           "customer",
	}

	// Act
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "12345678", result["identity_number"])
	assert.Equal(t, "John", result["first_name"])
	assert.Equal(t, "Doe", result["last_name"])
	assert.Equal(t, "john@example.com", result["email"])
}

func TestLoginResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := handlers.LoginResponse{
		AccessToken: "access-token-123",
		ExpiresIn:   3600,
		TokenType:   "Bearer",
		Links:       []handlers.Link{{Rel: "self", Href: "/auth/me", Method: http.MethodGet}},
	}

	// Act
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "access-token-123", result["access_token"])
	assert.Equal(t, float64(3600), result["expires_in"])
	assert.Equal(t, "Bearer", result["token_type"])
}

func TestAuthMeResponse_JSONSerialization(t *testing.T) {
	// Arrange
	response := handlers.AuthMeResponse{
		ID:             "encoded-id-123",
		IdentityNumber: "12345678",
		Email:          "john@example.com",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		PhoneNumber:    "3001234567",
		Role:           "customer",
	}

	// Act
	data, err := json.Marshal(response)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "encoded-id-123", result["id"])
	assert.Equal(t, "john@example.com", result["email"])
	assert.Equal(t, "John", result["first_name"])
}

// ============================================
// PersonRequest.Sanitize Tests
// ============================================

func TestPersonRequest_Sanitize(t *testing.T) {
	// Arrange
	req := &handlers.PersonRequest{
		IdentityNumber: "  1234567890  ",
		FirstName:      "  Juan  ",
		LastName:       "  Pérez  ",
		SecondLastName: "  García\t",
		Email:          "  juan@example.com  ",
		PhoneNumber:    "  +57 300 1234567\n",
		Password:       "  password123  ", // Should NOT be trimmed
		Role:           "  USER  ",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "1234567890", req.IdentityNumber)
	assert.Equal(t, "Juan", req.FirstName)
	assert.Equal(t, "Pérez", req.LastName)
	assert.Equal(t, "García", req.SecondLastName)
	assert.Equal(t, "juan@example.com", req.Email)
	assert.Equal(t, "+57 300 1234567", req.PhoneNumber)
	assert.Equal(t, "  password123  ", req.Password) // Password NOT trimmed
	assert.Equal(t, "USER", req.Role)
}

// ============================================
// LoginRequest.Sanitize Tests
// ============================================

func TestLoginRequest_Sanitize(t *testing.T) {
	// Arrange
	req := &handlers.LoginRequest{
		Email:    "  user@example.com  ",
		Password: "  password123  ",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "user@example.com", req.Email)
	assert.Equal(t, "  password123  ", req.Password) // Password NOT trimmed
}

// ============================================
// UpdateProfileRequest.Sanitize Tests
// ============================================

func TestUpdateProfileRequest_Sanitize(t *testing.T) {
	// Arrange
	req := &handlers.UpdateProfileRequest{
		IdentityNumber: "  9876543210  ",
		FirstName:      "  María  ",
		LastName:       "  López\t",
		SecondLastName: "  Rodríguez\n",
		PhoneNumber:    "  +57 311 9876543  ",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "9876543210", req.IdentityNumber)
	assert.Equal(t, "María", req.FirstName)
	assert.Equal(t, "López", req.LastName)
	assert.Equal(t, "Rodríguez", req.SecondLastName)
	assert.Equal(t, "+57 311 9876543", req.PhoneNumber)
}

// ============================================
// ResendVerificationEmailRequest.Sanitize Tests
// ============================================

func TestResendVerificationEmailRequest_Sanitize(t *testing.T) {
	// Arrange
	req := &handlers.ResendVerificationEmailRequest{
		Email: "  user@domain.com  \n",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "user@domain.com", req.Email)
}

// ============================================
// PasswordResetRequest.Sanitize Tests
// ============================================

func TestPasswordResetRequest_Sanitize(t *testing.T) {
	// Arrange
	req := &handlers.PasswordResetRequest{
		Email: "\t  reset@example.com  \t",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "reset@example.com", req.Email)
}
