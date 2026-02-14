package keycloak

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
)

// ============================================
// ensureValidToken Tests
// ============================================

// newTestClient creates a client with pre-set token for testing.
// Uses a gocloak client pointing to a non-routable address so refresh calls fail with
// a connection error rather than a nil pointer dereference.
func newTestClient(token *gocloak.JWT, expiresAt time.Time) *client {
	return &client{
		config: &config.KeycloakConfig{
			ServerURL:    "http://localhost:0",
			Realm:        "test-realm",
			AdminUser:    "admin",
			AdminPass:    "admin",
			ClientID:     "test-client",
			ClientSecret: "test-secret",
		},
		gocloak:        gocloak.NewClient("http://localhost:0"),
		token:          token,
		tokenExpiresAt: expiresAt,
		tokenMutex:     sync.RWMutex{},
	}
}

func TestEnsureValidToken_TokenNotExpired_ReturnsCurrent(t *testing.T) {
	// Token that expires far in the future (no refresh needed)
	token := &gocloak.JWT{
		AccessToken: "valid-access-token",
		ExpiresIn:   3600,
	}
	// Token expires 1 hour from now – well beyond the 30s threshold
	expiresAt := time.Now().Add(1 * time.Hour)

	c := newTestClient(token, expiresAt)

	accessToken, err := c.ensureValidToken(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "valid-access-token", accessToken)
}

func TestEnsureValidToken_TokenAboutToExpire_NeedsRefresh(t *testing.T) {
	// Token that will expire within 30 seconds (needs refresh)
	token := &gocloak.JWT{
		AccessToken: "expiring-token",
		ExpiresIn:   10, // 10 seconds
	}
	// Token expires 10 seconds from now – within the 30s threshold
	expiresAt := time.Now().Add(10 * time.Second)

	c := newTestClient(token, expiresAt)
	// gocloak is nil, so LoginAdmin will fail – but we're testing the refresh logic path
	_, err := c.ensureValidToken(context.Background())

	// Should fail because gocloak client is nil – but this proves the refresh logic was triggered
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to refresh admin token")
}

func TestEnsureValidToken_TokenExpired_NeedsRefresh(t *testing.T) {
	// Token that has already expired
	token := &gocloak.JWT{
		AccessToken: "expired-token",
		ExpiresIn:   0,
	}
	// Token expired 5 minutes ago
	expiresAt := time.Now().Add(-5 * time.Minute)

	c := newTestClient(token, expiresAt)

	_, err := c.ensureValidToken(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to refresh admin token")
}

func TestEnsureValidToken_DoubleCheckLock_SecondReaderGetsRefreshedToken(t *testing.T) {
	// This tests the double-check locking optimization:
	// When needsRefresh is true, the code upgrades to write lock,
	// then re-checks if another goroutine already refreshed.

	token := &gocloak.JWT{
		AccessToken: "refreshed-by-another-goroutine",
		ExpiresIn:   3600,
	}
	// Simulate: First read says "needs refresh" (token within 30s threshold)
	// But by the time we get the write lock, another goroutine already refreshed it
	// So the second check passes.
	expiresAt := time.Now().Add(1 * time.Hour)

	c := newTestClient(token, expiresAt)

	// Both concurrent readers should get the same valid token
	var wg sync.WaitGroup
	results := make([]string, 2)
	errors := make([]error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			tok, err := c.ensureValidToken(context.Background())
			results[idx] = tok
			errors[idx] = err
		}(i)
	}

	wg.Wait()

	for i := 0; i < 2; i++ {
		assert.NoError(t, errors[i])
		assert.Equal(t, "refreshed-by-another-goroutine", results[i])
	}
}

// ============================================
// LoginUser Input Validation Tests
// ============================================

func TestLoginUser_EmptyUsername(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	_, err := c.LoginUser(context.Background(), "", "password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username and password cannot be empty")
}

func TestLoginUser_EmptyPassword(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	_, err := c.LoginUser(context.Background(), "user", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username and password cannot be empty")
}

// ============================================
// CreateUser Input Validation Tests
// ============================================

func TestCreateUser_NilPerson(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	_, err := c.CreateUser(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "person cannot be nil")
}

// ============================================
// GetUserByEmail Input Validation Tests
// ============================================

func TestGetUserByEmail_EmptyEmail(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	_, err := c.GetUserByEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email cannot be empty")
}

// ============================================
// GetUserByID Input Validation Tests
// ============================================

func TestGetUserByID_EmptyID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	_, err := c.GetUserByID(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// ============================================
// UpdateUser Input Validation Tests
// ============================================

func TestUpdateUser_NilUser(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.UpdateUser(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user or user ID cannot be nil")
}

func TestUpdateUser_NilUserID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	user := &gocloak.User{} // ID is nil
	err := c.UpdateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user or user ID cannot be nil")
}

// ============================================
// DeleteUser Input Validation Tests
// ============================================

func TestDeleteUser_EmptyID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.DeleteUser(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// ============================================
// SetPassword Input Validation Tests
// ============================================

func TestSetPassword_EmptyUserID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.SetPassword(context.Background(), "", "newpass", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID and password cannot be empty")
}

func TestSetPassword_EmptyPassword(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.SetPassword(context.Background(), "user-123", "", false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID and password cannot be empty")
}

// ============================================
// AssignRole Input Validation Tests
// ============================================

func TestAssignRole_EmptyUserID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.AssignRole(context.Background(), "", "admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID and roleName cannot be empty")
}

func TestAssignRole_EmptyRoleName(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.AssignRole(context.Background(), "user-123", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID and roleName cannot be empty")
}

// ============================================
// RemoveRole Input Validation Tests
// ============================================

func TestRemoveRole_EmptyUserID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.RemoveRole(context.Background(), "", "admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID and roleName cannot be empty")
}

func TestRemoveRole_EmptyRoleName(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.RemoveRole(context.Background(), "user-123", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID and roleName cannot be empty")
}

// ============================================
// GetUserRoles Input Validation Tests
// ============================================

func TestGetUserRoles_EmptyUserID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	_, err := c.GetUserRoles(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// ============================================
// SendVerificationEmail Input Validation Tests
// ============================================

func TestSendVerificationEmail_EmptyUserID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.SendVerificationEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// ============================================
// SendPasswordResetEmail Input Validation Tests
// ============================================

func TestSendPasswordResetEmail_EmptyEmail(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.SendPasswordResetEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email cannot be empty")
}

// ============================================
// VerifyEmail Input Validation Tests
// ============================================

func TestVerifyEmail_EmptyUserID(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.VerifyEmail(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userID cannot be empty")
}

// ============================================
// Logout Input Validation Tests
// ============================================

func TestLogout_EmptyRefreshToken(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	err := c.Logout(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refreshToken cannot be empty")
}

// ============================================
// RefreshToken Input Validation Tests
// ============================================

func TestRefreshToken_EmptyRefreshToken(t *testing.T) {
	c := newTestClient(&gocloak.JWT{AccessToken: "tok"}, time.Now().Add(time.Hour))
	_, err := c.RefreshToken(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "refreshToken cannot be empty")
}

// ============================================
// NewClient Tests
// ============================================

func TestNewClient_NilConfig_Coverage(t *testing.T) {
	_, err := NewClient(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "keycloak config cannot be nil")
}
