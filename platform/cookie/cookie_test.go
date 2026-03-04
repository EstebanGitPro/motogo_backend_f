package cookie_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/platform/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const testProdDomain = "motogo.com"

func localCookieConfig() config.CookieConfig {
	return config.CookieConfig{
		Domain:   "",
		Secure:   false,
		SameSite: "lax",
	}
}

func prodCookieConfig() config.CookieConfig {
	return config.CookieConfig{
		Domain:   testProdDomain,
		Secure:   true,
		SameSite: "strict",
	}
}

func TestSetAccessTokenWritesCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	cookie.SetAccessToken(c, "test-access-token", localCookieConfig())

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "mg_access_token", cookies[0].Name)
	assert.Equal(t, "test-access-token", cookies[0].Value)
	assert.Equal(t, "/", cookies[0].Path)
	assert.True(t, cookies[0].HttpOnly)
	assert.False(t, cookies[0].Secure)
}

func TestSetRefreshTokenWritesCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	cookie.SetRefreshToken(c, "test-refresh-token", localCookieConfig())

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.Equal(t, "mg_refresh_token", cookies[0].Name)
	assert.Equal(t, "test-refresh-token", cookies[0].Value)
	assert.Equal(t, "/motogo/api/v1/auth", cookies[0].Path)
	assert.True(t, cookies[0].HttpOnly)
}

func TestSetAccessTokenSecureMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	cookie.SetAccessToken(c, "secure-token", prodCookieConfig())

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.True(t, cookies[0].Secure)
	assert.Equal(t, testProdDomain, cookies[0].Domain)
}

func TestSetRefreshTokenSecureMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	cookie.SetRefreshToken(c, "secure-refresh", prodCookieConfig())

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 1)
	assert.True(t, cookies[0].Secure)
	assert.Equal(t, testProdDomain, cookies[0].Domain)
}

func TestClearTokensSetsBothCookiesNegativeMaxAge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	cookie.ClearTokens(c, localCookieConfig())

	cookies := w.Result().Cookies()
	assert.Len(t, cookies, 2)

	assert.Equal(t, "mg_access_token", cookies[0].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.Equal(t, -1, cookies[0].MaxAge)

	assert.Equal(t, "mg_refresh_token", cookies[1].Name)
	assert.Equal(t, "", cookies[1].Value)
	assert.Equal(t, -1, cookies[1].MaxAge)
}

func TestGetAccessTokenSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.AddCookie(&http.Cookie{
		Name:  "mg_access_token",
		Value: "my-access-token",
	})

	token, err := cookie.GetAccessToken(c)
	assert.NoError(t, err)
	assert.Equal(t, "my-access-token", token)
}

func TestGetAccessTokenNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	token, err := cookie.GetAccessToken(c)
	assert.Error(t, err)
	assert.Equal(t, "", token)
}

func TestGetRefreshTokenSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/", nil)
	c.Request.AddCookie(&http.Cookie{
		Name:  "mg_refresh_token",
		Value: "my-refresh-token",
	})

	token, err := cookie.GetRefreshToken(c)
	assert.NoError(t, err)
	assert.Equal(t, "my-refresh-token", token)
}

func TestGetRefreshTokenNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	token, err := cookie.GetRefreshToken(c)
	assert.Error(t, err)
	assert.Equal(t, "", token)
}
