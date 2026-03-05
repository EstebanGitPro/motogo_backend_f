package cookie

import (
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/gin-gonic/gin"
)

const (
	// AccessTokenCookie es el nombre de la cookie que almacena el access token.
	// Se envía en todas las peticiones (Path=/).
	AccessTokenCookie = "mg_access_token"

	// RefreshTokenCookie es el nombre de la cookie que almacena el refresh token.
	// Solo viaja a endpoints de auth (Path=/motogo/api/v1/auth).
	RefreshTokenCookie = "mg_refresh_token"

	// refreshTokenPath restringe el refresh cookie a endpoints de autenticación.
	refreshTokenPath = "/motogo/api/v1/auth"

	// accessTokenMaxAge es el tiempo de vida del access token cookie (15 minutos).
	accessTokenMaxAge = 15 * 60

	// refreshTokenMaxAge es el tiempo de vida del refresh token cookie (7 días).
	refreshTokenMaxAge = 7 * 24 * 60 * 60
)

// SetAccessToken escribe la cookie mg_access_token en el response.
// Path=/ para que se envíe en todas las peticiones.
func SetAccessToken(c *gin.Context, token string, cfg config.CookieConfig) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    token,
		Path:     "/",
		Domain:   cfg.Domain,
		MaxAge:   accessTokenMaxAge,
		Secure:   cfg.Secure,
		HttpOnly: true,
		SameSite: cfg.ParseSameSite(),
	})
}

// SetRefreshToken escribe la cookie mg_refresh_token en el response.
// Path restringido a /motogo/api/v1/auth para minimizar exposición.
func SetRefreshToken(c *gin.Context, token string, cfg config.CookieConfig) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    token,
		Path:     refreshTokenPath,
		Domain:   cfg.Domain,
		MaxAge:   refreshTokenMaxAge,
		Secure:   cfg.Secure,
		HttpOnly: true,
		SameSite: cfg.ParseSameSite(),
	})
}

// ClearTokens limpia ambas cookies de autenticación (MaxAge=-1).
// Útil para logout o invalidación de sesión.
func ClearTokens(c *gin.Context, cfg config.CookieConfig) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		Domain:   cfg.Domain,
		MaxAge:   -1,
		Secure:   cfg.Secure,
		HttpOnly: true,
		SameSite: cfg.ParseSameSite(),
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     refreshTokenPath,
		Domain:   cfg.Domain,
		MaxAge:   -1,
		Secure:   cfg.Secure,
		HttpOnly: true,
		SameSite: cfg.ParseSameSite(),
	})
}

// GetAccessToken lee el valor de la cookie mg_access_token de la request.
// Retorna cadena vacía y error si no existe.
func GetAccessToken(c *gin.Context) (string, error) {
	cookie, err := c.Cookie(AccessTokenCookie)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

// GetRefreshToken lee el valor de la cookie mg_refresh_token de la request.
// Retorna cadena vacía y error si no existe.
func GetRefreshToken(c *gin.Context) (string, error) {
	cookie, err := c.Cookie(RefreshTokenCookie)
	if err != nil {
		return "", err
	}
	return cookie, nil
}
