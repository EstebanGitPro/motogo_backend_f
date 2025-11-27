package handlers

import "fmt"

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type HATEOASResource struct {
	Links []Link `json:"_links"`
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
