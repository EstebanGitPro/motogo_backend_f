package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type HATEOASResource struct {
	Links []Link `json:"_links"`
}

// GetBaseURL extrae la URL base de la petición (scheme + host)
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// SetLocationHeader establece el Location header con la URL del recurso
func SetLocationHeader(c *gin.Context, baseURL, resource, resourceID string) {
	locationURL := BuildResourceURL(baseURL, resource, resourceID)
	c.Header("Location", locationURL)
}

// BuildResourceURL construye una URL completa para un recurso
// Ejemplo: BuildResourceURL(baseURL, "accounts", encodedID) → "http://host/motogo/api/v1/accounts/xyz"
func BuildResourceURL(baseURL, resource, resourceID string) string {
	return fmt.Sprintf("%s/motogo/api/v1/%s/%s", baseURL, resource, resourceID)
}

// BuildCollectionURL construye una URL completa para una colección
// Ejemplo: BuildCollectionURL(baseURL, "accounts") → "http://host/motogo/api/v1/accounts"
func BuildCollectionURL(baseURL, resource string) string {
	return fmt.Sprintf("%s/motogo/api/v1/%s", baseURL, resource)
}

// BuildResourceLinks construye links HATEOAS genéricos para un recurso
// resource: nombre del recurso (ej: "accounts", "transactions")
// resourceID: ID ya ofuscado del recurso
func BuildResourceLinks(baseURL, resource, resourceID string) []Link {
	resourceURL := BuildResourceURL(baseURL, resource, resourceID)
	collectionURL := BuildCollectionURL(baseURL, resource)

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "collection",
			Method: "GET",
		},
	}
}

// BuildAccountLinks construye links específicos para cuentas (wrapper para compatibilidad)
func BuildAccountLinks(baseURL string, accountID string) []Link {
	return BuildResourceLinks(baseURL, "accounts", accountID)
}

// BuildMessageLinks construye links HATEOAS para un mensaje específico
func BuildMessageLinks(baseURL string, messageID string) []Link {
	return BuildResourceLinks(baseURL, "messages", messageID)
}

// BuildMessageCreatedLinks construye links para un mensaje recién creado
func BuildMessageCreatedLinks(baseURL string, messageID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "messages", messageID)
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// BuildMessageUpdatedLinks construye links para un mensaje actualizado
func BuildMessageUpdatedLinks(baseURL string, messageID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "messages", messageID)
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   resourceURL,
			Rel:    "delete",
			Method: "DELETE",
		},
		{
			Href:   collectionURL,
			Rel:    "list",
			Method: "GET",
		},
	}
}

// BuildMessageListLinks construye links para la lista de mensajes
func BuildMessageListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "messages")

	return []Link{
		{
			Href:   collectionURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   collectionURL,
			Rel:    "create",
			Method: "POST",
		},
	}
}

// ========================================
// Personalized HATEOAS Functions
// These functions return only relevant, implemented links for each endpoint context
// ========================================

// BuildLoginLinks constructs HATEOAS links for successful login response
// Only returns links that are actually implemented and relevant after login
func BuildLoginLinks(baseURL string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/auth/me", baseURL),
			Rel:    "profile",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/auth/password-reset", baseURL),
			Rel:    "password-reset",
			Method: "POST",
		},
	}
}

// BuildAuthMeLinks constructs HATEOAS links for authenticated user profile (/auth/me)
// Returns links specific to the authenticated user context
func BuildAuthMeLinks(baseURL string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/me", baseURL),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/me/password", baseURL),
			Rel:    "change-password",
			Method: "PUT",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/auth/login", baseURL),
			Rel:    "login",
			Method: "POST",
		},
	}
}

// BuildAccountCreatedLinks constructs HATEOAS links for newly created account
// Replaces the generic BuildAccountLinks with context-specific links for registration
func BuildAccountCreatedLinks(baseURL string, accountID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "accounts", accountID)

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/auth/login", baseURL),
			Rel:    "login",
			Method: "POST",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/auth/verify-email", baseURL),
			Rel:    "verify-email",
			Method: "POST",
		},
	}
}

// BuildChangePasswordLinks constructs HATEOAS links for password change response (HU57)
// Returns links to profile and login after successful password change
func BuildChangePasswordLinks(baseURL string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/me", baseURL),
			Rel:    "profile",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/auth/login", baseURL),
			Rel:    "login",
			Method: "POST",
		},
	}
}

// BuildUpdateProfileLinks constructs HATEOAS links for profile update response (HU52)
// Richardson Maturity Model Level 3: Provides discoverable next actions
func BuildUpdateProfileLinks(baseURL string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/me", baseURL),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/me", baseURL),
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/me/password", baseURL),
			Rel:    "change-password",
			Method: "PUT",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/auth/login", baseURL),
			Rel:    "login",
			Method: "POST",
		},
	}
}

// BuildPublicContactLinks constructs HATEOAS links for public contact response (HU55)
// Para motociclistas viendo info de contacto del representante
func BuildPublicContactLinks(baseURL string, personID string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/%s/contact", baseURL, personID),
			Rel:    "self",
			Method: "GET",
		},
	}
}
