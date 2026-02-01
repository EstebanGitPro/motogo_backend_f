package jwt

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock JWT token with given claims
func createMockToken(claims map[string]interface{}) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))

	claimsJSON, _ := json.Marshal(claims)
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signature := base64.RawURLEncoding.EncodeToString([]byte("mock_signature"))

	return header + "." + payload + "." + signature
}

// ============================================
// NewTokenParser Tests
// ============================================

func TestNewTokenParser_ReturnsInstance(t *testing.T) {
	parser := NewTokenParser()
	assert.NotNil(t, parser)
}

// ============================================
// ExtractEmailFromToken Tests
// ============================================

func TestExtractEmailFromToken_InvalidFormat_TwoParts(t *testing.T) {
	parser := NewTokenParser()
	email, err := parser.ExtractEmailFromToken("header.payload")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTokenFormat, err)
	assert.Empty(t, email)
}

func TestExtractEmailFromToken_InvalidFormat_OnePart(t *testing.T) {
	parser := NewTokenParser()
	email, err := parser.ExtractEmailFromToken("invalidtoken")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTokenFormat, err)
	assert.Empty(t, email)
}

func TestExtractEmailFromToken_InvalidFormat_Empty(t *testing.T) {
	parser := NewTokenParser()
	email, err := parser.ExtractEmailFromToken("")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTokenFormat, err)
	assert.Empty(t, email)
}

func TestExtractEmailFromToken_InvalidBase64Payload(t *testing.T) {
	parser := NewTokenParser()
	// Invalid base64 in payload
	token := "header.!!!invalid-base64!!!.signature"
	email, err := parser.ExtractEmailFromToken(token)

	assert.Error(t, err)
	assert.Equal(t, ErrPayloadDecode, err)
	assert.Empty(t, email)
}

func TestExtractEmailFromToken_InvalidJSONPayload(t *testing.T) {
	parser := NewTokenParser()
	// Valid base64 but invalid JSON
	invalidJSON := base64.RawURLEncoding.EncodeToString([]byte("not valid json"))
	token := "header." + invalidJSON + ".signature"
	email, err := parser.ExtractEmailFromToken(token)

	assert.Error(t, err)
	assert.Equal(t, ErrClaimsParse, err)
	assert.Empty(t, email)
}

func TestExtractEmailFromToken_WithEmlField(t *testing.T) {
	parser := NewTokenParser()
	token := createMockToken(map[string]interface{}{
		"eml": "user@example.com",
		"sub": "user-uuid-123",
	})

	email, err := parser.ExtractEmailFromToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", email)
}

func TestExtractEmailFromToken_WithEmailField(t *testing.T) {
	parser := NewTokenParser()
	token := createMockToken(map[string]interface{}{
		"email": "user@example.com",
		"sub":   "user-uuid-123",
	})

	email, err := parser.ExtractEmailFromToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", email)
}

func TestExtractEmailFromToken_WithSubAsEmail(t *testing.T) {
	parser := NewTokenParser()
	token := createMockToken(map[string]interface{}{
		"sub": "user@example.com",
	})

	email, err := parser.ExtractEmailFromToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", email)
}

func TestExtractEmailFromToken_EmlTakesPrecedence(t *testing.T) {
	parser := NewTokenParser()
	// When both eml and email are present, eml takes precedence
	token := createMockToken(map[string]interface{}{
		"eml":   "primary@example.com",
		"email": "secondary@example.com",
	})

	email, err := parser.ExtractEmailFromToken(token)

	assert.NoError(t, err)
	assert.Equal(t, "primary@example.com", email)
}

func TestExtractEmailFromToken_NoEmailFields(t *testing.T) {
	parser := NewTokenParser()
	token := createMockToken(map[string]interface{}{
		"sub":  "user-uuid-123",
		"name": "John Doe",
	})

	email, err := parser.ExtractEmailFromToken(token)

	assert.Error(t, err)
	assert.Equal(t, ErrEmailNotFound, err)
	assert.Empty(t, email)
}

func TestExtractEmailFromToken_EmptyEmailFields(t *testing.T) {
	parser := NewTokenParser()
	token := createMockToken(map[string]interface{}{
		"eml":   "",
		"email": "",
		"sub":   "not-an-email",
	})

	email, err := parser.ExtractEmailFromToken(token)

	assert.Error(t, err)
	assert.Equal(t, ErrEmailNotFound, err)
	assert.Empty(t, email)
}

// ============================================
// ExtractClaimsFromToken Tests
// ============================================

func TestExtractClaimsFromToken_Success(t *testing.T) {
	parser := NewTokenParser()
	expectedClaims := map[string]interface{}{
		"sub":   "user-uuid-123",
		"email": "user@example.com",
		"roles": []interface{}{"user", "admin"},
	}
	token := createMockToken(expectedClaims)

	claims, err := parser.ExtractClaimsFromToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "user-uuid-123", claims["sub"])
	assert.Equal(t, "user@example.com", claims["email"])
}

func TestExtractClaimsFromToken_InvalidFormat(t *testing.T) {
	parser := NewTokenParser()
	claims, err := parser.ExtractClaimsFromToken("invalid.token")

	assert.Error(t, err)
	assert.Equal(t, ErrInvalidTokenFormat, err)
	assert.Nil(t, claims)
}

func TestExtractClaimsFromToken_InvalidPayload(t *testing.T) {
	parser := NewTokenParser()
	invalidPayload := base64.RawURLEncoding.EncodeToString([]byte("{invalid}"))
	token := "header." + invalidPayload + ".signature"

	claims, err := parser.ExtractClaimsFromToken(token)

	assert.Error(t, err)
	assert.Equal(t, ErrClaimsParse, err)
	assert.Nil(t, claims)
}

func TestExtractClaimsFromToken_InvalidBase64(t *testing.T) {
	parser := NewTokenParser()
	// Invalid base64 in payload
	token := "header.!!!invalid-base64!!!.signature"

	claims, err := parser.ExtractClaimsFromToken(token)

	assert.Error(t, err)
	assert.Equal(t, ErrPayloadDecode, err)
	assert.Nil(t, claims)
}

// ============================================
// base64URLDecode Tests
// ============================================

func TestBase64URLDecode_NoPadding(t *testing.T) {
	// "hello" in base64url without padding
	encoded := "aGVsbG8"
	decoded, err := base64URLDecode(encoded)

	assert.NoError(t, err)
	assert.Equal(t, "hello", string(decoded))
}

func TestBase64URLDecode_WithOnePadding(t *testing.T) {
	// String that needs one = padding
	encoded := "aGVsbG9X" // "helloW"
	decoded, err := base64URLDecode(encoded)

	assert.NoError(t, err)
	assert.Equal(t, "helloW", string(decoded))
}

func TestBase64URLDecode_WithTwoPadding(t *testing.T) {
	// String that needs two == padding
	encoded := "aGk" // "hi"
	decoded, err := base64URLDecode(encoded)

	assert.NoError(t, err)
	assert.Equal(t, "hi", string(decoded))
}

// ============================================
// isValidEmail Tests
// ============================================

func TestIsValidEmail_ValidEmail(t *testing.T) {
	testCases := []string{
		"user@example.com",
		"test.user@domain.org",
		"name+tag@company.co",
		"a@b.c",
	}

	for _, email := range testCases {
		t.Run(email, func(t *testing.T) {
			assert.True(t, isValidEmail(email))
		})
	}
}

func TestIsValidEmail_InvalidEmail(t *testing.T) {
	testCases := []struct {
		name  string
		email string
	}{
		{"no at symbol", "userexample.com"},
		{"starts with at", "@example.com"},
		{"ends with at", "user@"},
		{"double at", "user@@example.com"},
		{"empty string", ""},
		{"only at", "@"},
		{"at at start", "@user"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.False(t, isValidEmail(tc.email))
		})
	}
}
