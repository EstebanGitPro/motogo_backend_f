package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// URL path format strings (extracted to avoid duplication — SonarCloud S1192)
const (
	motorcycleReferencesPath = "%s/motogo/api/v1/motorcycle-references"
	motorcycleCategoriesPath = "%s/motogo/api/v1/motorcycle-categories"
	branchTypesPath          = "%s/motogo/api/v1/branch-types"
	motorcycleEvidencePath   = "%s/motogo/api/v1/motorcycles/%s/evidence"
	motorcycleEvidenceIDPath = "%s/motogo/api/v1/motorcycles/%s/evidence/%s"
	personsProfilePath       = "%s/motogo/api/v1/persons/me"
	personsPasswordPath      = "%s/motogo/api/v1/persons/me/password"
	authLoginPath            = "%s/motogo/api/v1/auth/login"
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
// The host is validated to prevent open redirects via Host header injection (S5146)
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := sanitizeHost(c.Request.Host)
	return scheme + "://" + host
}

// sanitizeHost validates that the host contains only valid hostname characters.
// Returns "localhost" if the host contains suspicious characters that could enable open redirects.
func sanitizeHost(host string) string {
	for _, ch := range host {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && ch != '.' && ch != ':' && ch != '-' {
			return "localhost"
		}
	}
	if host == "" {
		return "localhost"
	}
	return host
}

// SetLocationHeader establece el Location header con la URL del recurso.
// Uses a relative path to avoid open redirects via user-controlled Host header (S5146).
func SetLocationHeader(c *gin.Context, _, resource, resourceID string) {
	path := fmt.Sprintf("/motogo/api/v1/%s", resource)
	if resourceID != "" {
		path = fmt.Sprintf("%s/%s", path, resourceID)
	}
	c.Header("Location", path)
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
func BuildAccountLinks(baseURL, accountID string) []Link {
	return BuildResourceLinks(baseURL, "accounts", accountID)
}

// BuildMessageLinks construye links HATEOAS para un mensaje específico
func BuildMessageLinks(baseURL, messageID string) []Link {
	return BuildResourceLinks(baseURL, "messages", messageID)
}

// BuildMessageCreatedLinks construye links para un mensaje recién creado
func BuildMessageCreatedLinks(baseURL, messageID string) []Link {
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
func BuildMessageUpdatedLinks(baseURL, messageID string) []Link {
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
			Href:   fmt.Sprintf(personsProfilePath, baseURL),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf(personsPasswordPath, baseURL),
			Rel:    "change-password",
			Method: "PUT",
		},
		{
			Href:   fmt.Sprintf(authLoginPath, baseURL),
			Rel:    "login",
			Method: "POST",
		},
	}
}

// BuildAccountCreatedLinks constructs HATEOAS links for newly created account
// Replaces the generic BuildAccountLinks with context-specific links for registration
func BuildAccountCreatedLinks(baseURL, accountID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "accounts", accountID)

	return []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf(authLoginPath, baseURL),
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
			Href:   fmt.Sprintf(personsProfilePath, baseURL),
			Rel:    "profile",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf(authLoginPath, baseURL),
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
			Href:   fmt.Sprintf(personsProfilePath, baseURL),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf(personsProfilePath, baseURL),
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   fmt.Sprintf(personsPasswordPath, baseURL),
			Rel:    "change-password",
			Method: "PUT",
		},
		{
			Href:   fmt.Sprintf(authLoginPath, baseURL),
			Rel:    "login",
			Method: "POST",
		},
	}
}

// BuildPublicContactLinks constructs HATEOAS links for public contact response (HU55)
// Para motociclistas viendo info de contacto del representante
func BuildPublicContactLinks(baseURL, personID string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/persons/%s/contact", baseURL, personID),
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildBranchCreatedLinks constructs HATEOAS links for newly created branch (HU59)
// Richardson Maturity Model Level 3: Provides discoverable next actions
func BuildBranchCreatedLinks(baseURL, branchID string) []Link {
	resourceURL := BuildResourceURL(baseURL, "branches", branchID)
	collectionURL := BuildCollectionURL(baseURL, "branches")

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
		{
			Href:   BuildCollectionURL(baseURL, "brands"),
			Rel:    "brands",
			Method: "GET",
		},
	}
}

// BuildBrandListLinks constructs HATEOAS links for the brands catalog
func BuildBrandListLinks(baseURL string) []Link {
	return []Link{
		{
			Href:   BuildCollectionURL(baseURL, "brands"),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   BuildCollectionURL(baseURL, "branches"),
			Rel:    "create-branch",
			Method: "POST",
		},
	}
}

// BuildDepartmentListLinks constructs HATEOAS links for the departments catalog
func BuildDepartmentListLinks() []Link {
	return []Link{
		{
			Href:   "/motogo/api/v1/departments",
			Rel:    "self",
			Method: "GET",
		},
	}
}

// BuildCityListLinks constructs HATEOAS links for cities of a department
func BuildCityListLinks(departmentID string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("/motogo/api/v1/departments/%s/cities", departmentID),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   "/motogo/api/v1/departments",
			Rel:    "departments",
			Method: "GET",
		},
	}
}

// BuildBranchDetailLinks constructs HATEOAS links for branch query response (HU62)
// Richardson Maturity Model Level 3: Provides discoverable next actions
// isOwner: true if the authenticated user is the branch owner (shows edit/delete links)
func BuildBranchDetailLinks(baseURL string, branchID string, isOwner bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "branches", branchID)
	collectionURL := BuildCollectionURL(baseURL, "branches")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Only show edit/delete links if user is the owner
	if isOwner {
		links = append(links,
			Link{
				Href:   resourceURL,
				Rel:    "update",
				Method: "PUT",
			},
			Link{
				Href:   resourceURL,
				Rel:    "delete",
				Method: "DELETE",
			},
			Link{
				Href:   collectionURL,
				Rel:    "list",
				Method: "GET",
			},
		)
	}

	return links
}

// BuildBranchTypesLinks constructs HATEOAS links for branch types catalog (HU76)
func BuildBranchTypesLinks(baseURL string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf(branchTypesPath, baseURL),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   BuildCollectionURL(baseURL, "branches"),
			Rel:    "create-branch",
			Method: "POST",
		},
	}
}

// BuildBranchListLinks constructs HATEOAS links for branch list response (HU62)
func BuildBranchListLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "branches")
	return []Link{
		{Href: collectionURL, Rel: "self", Method: "GET"},
		{Href: collectionURL, Rel: "create", Method: "POST"},
		{Href: BuildCollectionURL(baseURL, "brands"), Rel: "brands", Method: "GET"},
		{Href: fmt.Sprintf(branchTypesPath, baseURL), Rel: "branch-types", Method: "GET"},
	}
}

// BuildBranchDeletedLinks constructs HATEOAS links after branch deletion (HU61)
// Shows next possible actions: list branches or create new branch
func BuildBranchDeletedLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "branches")
	return []Link{
		{Href: collectionURL, Rel: "list", Method: "GET"},
		{Href: collectionURL, Rel: "create", Method: "POST"},
	}
}

// BuildServiceTypeListLinks constructs HATEOAS links for service types catalog (HU75)
func BuildServiceTypeListLinks(baseURL string) []Link {
	return []Link{
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/service-types", baseURL),
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   BuildCollectionURL(baseURL, "services"),
			Rel:    "services",
			Method: "GET",
		},
	}
}

// BuildServiceListLinks constructs HATEOAS links for services catalog (HU63)
func BuildServiceListLinks(baseURL, filterType string) []Link {
	servicesURL := BuildCollectionURL(baseURL, "services")
	links := []Link{
		{
			Href:   servicesURL,
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   fmt.Sprintf("%s/motogo/api/v1/service-types", baseURL),
			Rel:    "service-types",
			Method: "GET",
		},
	}

	// If filtered, add link to unfiltered list
	if filterType != "" {
		links = append(links, Link{
			Href:   servicesURL,
			Rel:    "all-services",
			Method: "GET",
		})
	}

	return links
}

// ========================================
// Motorcycle HATEOAS Functions (HU43-47)
// ========================================

// BuildMotorcycleDetailLinks constructs HATEOAS links for motorcycle response (HU43-47)
// Implements Richardson Maturity Level 3 with hypermedia controls
// isOwner: true if the authenticated user is the motorcycle owner (shows edit/delete links)
func BuildMotorcycleDetailLinks(baseURL string, motorcycleID string, isOwner bool) []Link {
	resourceURL := BuildResourceURL(baseURL, "motorcycles", motorcycleID)
	collectionURL := BuildCollectionURL(baseURL, "motorcycles")

	links := []Link{
		{
			Href:   resourceURL,
			Rel:    "self",
			Method: "GET",
		},
	}

	// Only show edit/delete links if user is the owner
	if isOwner {
		links = append(links,
			Link{
				Href:   resourceURL,
				Rel:    "update",
				Method: "PUT",
			},
			Link{
				Href:   resourceURL,
				Rel:    "delete",
				Method: "DELETE",
			},
			Link{
				Href:   collectionURL,
				Rel:    "list",
				Method: "GET",
			},
		)
	}

	return links
}

// BuildMotorcycleDeletedLinks constructs HATEOAS links after motorcycle deletion (HU45)
// Shows next possible actions: list motorcycles or create new motorcycle
func BuildMotorcycleDeletedLinks(baseURL string) []Link {
	collectionURL := BuildCollectionURL(baseURL, "motorcycles")
	return []Link{
		{Href: collectionURL, Rel: "list", Method: "GET"},
		{Href: collectionURL, Rel: "create", Method: "POST"},
	}
}

// BuildNearbyBranchesLinks constructs HATEOAS links for nearby branches search (HU89)
func BuildNearbyBranchesLinks(baseURL string, lat, lng, radiusKm float64) []Link {
	nearbyURL := fmt.Sprintf("%s/motogo/api/v1/branches/nearby?lat=%.8f&lng=%.8f&radius=%.1f", baseURL, lat, lng, radiusKm)
	return []Link{
		{Href: nearbyURL, Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL(baseURL, "branches"), Rel: "branches", Method: "GET"},
		{Href: fmt.Sprintf(branchTypesPath, baseURL), Rel: "branch-types", Method: "GET"},
	}
}

// BuildMotorcycleReferencesLinks constructs HATEOAS links for motorcycle references catalog (HU50)
func BuildMotorcycleReferencesLinks(baseURL string) []Link {
	return []Link{
		{Href: fmt.Sprintf(motorcycleReferencesPath, baseURL), Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL(baseURL, "motorcycles"), Rel: "motorcycles", Method: "GET"},
		{Href: BuildCollectionURL(baseURL, "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildBrandLinesLinks constructs HATEOAS links for brand lines catalog (HU40 - Admin)
func BuildBrandLinesLinks(baseURL, brandID string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/motogo/api/v1/admin/brands/%s/lines", baseURL, brandID), Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL(baseURL, "brands"), Rel: "brands", Method: "GET"},
		{Href: fmt.Sprintf(motorcycleReferencesPath, baseURL), Rel: "references", Method: "GET"},
	}
}

// ============================================
// MOTORCYCLE CATEGORY HATEOAS (HU41)
// ============================================

// BuildMotorcycleCategoryListLinks constructs HATEOAS links for the motorcycle categories list
func BuildMotorcycleCategoryListLinks(baseURL string) []Link {
	return []Link{
		{Href: fmt.Sprintf(motorcycleCategoriesPath, baseURL), Rel: "self", Method: "GET"},
		{Href: fmt.Sprintf(motorcycleReferencesPath, baseURL), Rel: "references", Method: "GET"},
		{Href: BuildCollectionURL(baseURL, "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildMotorcycleCategoryItemLinks constructs drill-down HATEOAS links for a single category item
func BuildMotorcycleCategoryItemLinks(baseURL, categoryName string) []Link {
	return []Link{
		{Href: fmt.Sprintf(motorcycleCategoriesPath+"/%s/lines", baseURL, categoryName), Rel: "lines", Method: "GET"},
	}
}

// BuildCategoryLinesLinks constructs HATEOAS links for lines within a specific category
func BuildCategoryLinesLinks(baseURL, categoryName string) []Link {
	return []Link{
		{Href: fmt.Sprintf(motorcycleCategoriesPath+"/%s/lines", baseURL, categoryName), Rel: "self", Method: "GET"},
		{Href: fmt.Sprintf(motorcycleCategoriesPath, baseURL), Rel: "categories", Method: "GET"},
		{Href: BuildCollectionURL(baseURL, "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildEngineDisplacementLinks constructs HATEOAS links for the engine displacement list (HU49)
func BuildEngineDisplacementLinks(baseURL string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/motogo/api/v1/engine-displacements", baseURL), Rel: "self", Method: "GET"},
		{Href: fmt.Sprintf(motorcycleCategoriesPath, baseURL), Rel: "categories", Method: "GET"},
		{Href: BuildCollectionURL(baseURL, "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildRatingRangeLinks constructs HATEOAS links for the rating range list (HU48)
func BuildRatingRangeLinks(baseURL string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/motogo/api/v1/rating-ranges", baseURL), Rel: "self", Method: "GET"},
		{Href: fmt.Sprintf("%s/motogo/api/v1/engine-displacements", baseURL), Rel: "displacements", Method: "GET"},
		{Href: fmt.Sprintf(motorcycleCategoriesPath, baseURL), Rel: "categories", Method: "GET"},
	}
}

// ============================================
// MOTORCYCLE EVIDENCE HATEOAS (HU16-19)
// ============================================

// BuildEvidenceDetailLinks constructs HATEOAS links for evidence detail (HU16-19)
// isOwner is used to show/hide edit and delete actions
func BuildEvidenceDetailLinks(baseURL, motorcycleID, evidenceID string, isOwner bool) []Link {
	evidenceURL := fmt.Sprintf(motorcycleEvidenceIDPath, baseURL, motorcycleID, evidenceID)
	listURL := fmt.Sprintf(motorcycleEvidencePath, baseURL, motorcycleID)
	motorcycleURL := BuildResourceURL(baseURL, "motorcycles", motorcycleID)

	links := []Link{
		{Href: evidenceURL, Rel: "self", Method: "GET"},
		{Href: listURL, Rel: "list", Method: "GET"},
		{Href: motorcycleURL, Rel: "motorcycle", Method: "GET"},
	}

	if isOwner {
		links = append(links,
			Link{Href: evidenceURL, Rel: "delete", Method: "DELETE"},
			Link{Href: listURL, Rel: "create", Method: "POST"},
		)
	}

	return links
}

// BuildEvidenceListLinks constructs HATEOAS links for evidence list (HU18)
func BuildEvidenceListLinks(baseURL, motorcycleID string) []Link {
	listURL := fmt.Sprintf(motorcycleEvidencePath, baseURL, motorcycleID)
	motorcycleURL := BuildResourceURL(baseURL, "motorcycles", motorcycleID)

	return []Link{
		{Href: listURL, Rel: "self", Method: "GET"},
		{Href: listURL, Rel: "create", Method: "POST"},
		{Href: motorcycleURL, Rel: "motorcycle", Method: "GET"},
	}
}

// BuildEvidenceDeletedLinks constructs HATEOAS links after evidence deletion (HU19)
// Shows next possible actions: list evidence or view motorcycle
func BuildEvidenceDeletedLinks(baseURL, motorcycleID string) []Link {
	listURL := fmt.Sprintf(motorcycleEvidencePath, baseURL, motorcycleID)
	motorcycleURL := BuildResourceURL(baseURL, "motorcycles", motorcycleID)

	return []Link{
		{Href: listURL, Rel: "list", Method: "GET"},
		{Href: listURL, Rel: "create", Method: "POST"},
		{Href: motorcycleURL, Rel: "motorcycle", Method: "GET"},
	}
}

// ============================================
// DIAGNOSTIC PERMISSION HATEOAS
// ============================================

// BuildPermissionLinks constructs HATEOAS links for diagnostic permission operations
func BuildPermissionLinks(baseURL, motorcycleID string) []Link {
	permissionsURL := fmt.Sprintf("%s/motogo/api/v1/motorcycles/%s/permissions", baseURL, motorcycleID)
	motorcycleURL := BuildResourceURL(baseURL, "motorcycles", motorcycleID)

	return []Link{
		{Href: permissionsURL, Rel: "list-permissions", Method: "GET"},
		{Href: permissionsURL, Rel: "grant-permission", Method: "POST"},
		{Href: motorcycleURL, Rel: "motorcycle", Method: "GET"},
	}
}
