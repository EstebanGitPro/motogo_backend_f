package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// API base path — all HATEOAS links use relative paths to avoid open redirects (S5146)
const apiBase = "/motogo/api/v1"

// Path format strings (extracted to avoid duplication — SonarCloud S1192)
const (
	adminMessagesPath        = apiBase + "/admin/messages"
	adminMessageIDFmt        = adminMessagesPath + "/%s"
	personsProfilePath       = apiBase + "/persons/me"
	personsPasswordPath      = personsProfilePath + "/password"
	authLoginPath            = apiBase + "/auth/login"
	branchTypesPath          = apiBase + "/branch-types"
	motorcycleReferencesPath = apiBase + "/motorcycle-references"
	motorcycleCategoriesPath = apiBase + "/motorcycle-categories"
	motorcycleEvidenceFmt    = apiBase + "/motorcycles/%s/evidence"
	motorcycleEvidenceIDFmt  = motorcycleEvidenceFmt + "/%s"
)

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type HATEOASResource struct {
	Links []Link `json:"_links"`
}

// GetBaseURL is kept for backward compatibility but returns only the relative API base path.
// This prevents open redirects via user-controlled Host header injection (S5146).
func GetBaseURL(_ *gin.Context) string {
	return apiBase
}

// SetLocationHeader establece el Location header con la URL del recurso.
// Uses a relative path to avoid open redirects via user-controlled Host header (S5146).
func SetLocationHeader(c *gin.Context, _, resource, resourceID string) {
	path := fmt.Sprintf("%s/%s", apiBase, resource)
	if resourceID != "" {
		path = fmt.Sprintf("%s/%s", path, resourceID)
	}
	c.Header("Location", path)
}

// BuildResourceURL construye una URL relativa para un recurso
func BuildResourceURL(_, resource, resourceID string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, resource, resourceID)
}

// BuildCollectionURL construye una URL relativa para una colección
func BuildCollectionURL(_, resource string) string {
	return fmt.Sprintf("%s/%s", apiBase, resource)
}

// BuildResourceLinks construye links HATEOAS genéricos para un recurso
func BuildResourceLinks(_, resource, resourceID string) []Link {
	resourceURL := BuildResourceURL("", resource, resourceID)
	collectionURL := BuildCollectionURL("", resource)

	return []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
		{Href: resourceURL, Rel: "update", Method: "PUT"},
		{Href: resourceURL, Rel: "delete", Method: "DELETE"},
		{Href: collectionURL, Rel: "collection", Method: "GET"},
	}
}

// BuildMessageLinks construye links HATEOAS para un mensaje específico
func BuildMessageLinks(_, messageID string) []Link {
	resourceURL := fmt.Sprintf(adminMessageIDFmt, messageID)
	collectionURL := adminMessagesPath
	return []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
		{Href: resourceURL, Rel: "update", Method: "PUT"},
		{Href: resourceURL, Rel: "delete", Method: "DELETE"},
		{Href: collectionURL, Rel: "collection", Method: "GET"},
	}
}

// BuildMessageCreatedLinks construye links para un mensaje recién creado
func BuildMessageCreatedLinks(_, messageID string) []Link {
	resourceURL := fmt.Sprintf(adminMessageIDFmt, messageID)
	collectionURL := adminMessagesPath
	return []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
		{Href: resourceURL, Rel: "update", Method: "PUT"},
		{Href: resourceURL, Rel: "delete", Method: "DELETE"},
		{Href: collectionURL, Rel: "list", Method: "GET"},
	}
}

// BuildMessageUpdatedLinks construye links para un mensaje actualizado
func BuildMessageUpdatedLinks(_, messageID string) []Link {
	resourceURL := fmt.Sprintf(adminMessageIDFmt, messageID)
	collectionURL := adminMessagesPath
	return []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
		{Href: resourceURL, Rel: "delete", Method: "DELETE"},
		{Href: collectionURL, Rel: "list", Method: "GET"},
	}
}

// BuildMessageListLinks construye links para la lista de mensajes
func BuildMessageListLinks(_ string) []Link {
	collectionURL := adminMessagesPath
	return []Link{
		{Href: collectionURL, Rel: "self", Method: "GET"},
		{Href: collectionURL, Rel: "create", Method: "POST"},
	}
}

// ========================================
// Personalized HATEOAS Functions
// ========================================

// BuildLoginLinks constructs HATEOAS links for successful login response
func BuildLoginLinks(_ string) []Link {
	return []Link{
		{Href: personsProfilePath, Rel: "profile", Method: "GET"},
		{Href: fmt.Sprintf("%s/auth/password-reset", apiBase), Rel: "password-reset", Method: "POST"},
	}
}

// BuildAuthMeLinks constructs HATEOAS links for authenticated user profile (/auth/me)
func BuildAuthMeLinks(_ string) []Link {
	return []Link{
		{Href: personsProfilePath, Rel: "self", Method: "GET"},
		{Href: personsPasswordPath, Rel: "change-password", Method: "PUT"},
		{Href: authLoginPath, Rel: "login", Method: "POST"},
	}
}

// BuildAccountCreatedLinks constructs HATEOAS links for newly created account
func BuildAccountCreatedLinks(_ string) []Link {
	return []Link{
		{Href: authLoginPath, Rel: "login", Method: "POST"},
		{Href: fmt.Sprintf("%s/auth/verify-email", apiBase), Rel: "verify-email", Method: "POST"},
	}
}

// BuildChangePasswordLinks constructs HATEOAS links for password change response (HU57)
func BuildChangePasswordLinks(_ string) []Link {
	return []Link{
		{Href: personsProfilePath, Rel: "profile", Method: "GET"},
		{Href: authLoginPath, Rel: "login", Method: "POST"},
	}
}

// BuildUpdateProfileLinks constructs HATEOAS links for profile update response (HU52)
func BuildUpdateProfileLinks(_ string) []Link {
	return []Link{
		{Href: personsProfilePath, Rel: "self", Method: "GET"},
		{Href: personsProfilePath, Rel: "update", Method: "PUT"},
		{Href: personsPasswordPath, Rel: "change-password", Method: "PUT"},
		{Href: authLoginPath, Rel: "login", Method: "POST"},
	}
}

// BuildPublicContactLinks constructs HATEOAS links for public contact response (HU55)
func BuildPublicContactLinks(_, personID string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/persons/%s/contact", apiBase, personID), Rel: "self", Method: "GET"},
	}
}

// BuildBranchCreatedLinks constructs HATEOAS links for newly created branch (HU59)
func BuildBranchCreatedLinks(_, branchID string) []Link {
	resourceURL := BuildResourceURL("", "branches", branchID)
	collectionURL := BuildCollectionURL("", "branches")
	return []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
		{Href: resourceURL, Rel: "update", Method: "PUT"},
		{Href: resourceURL, Rel: "delete", Method: "DELETE"},
		{Href: collectionURL, Rel: "list", Method: "GET"},
		{Href: BuildCollectionURL("", "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildBrandListLinks constructs HATEOAS links for the brands catalog
func BuildBrandListLinks(_ string) []Link {
	return []Link{
		{Href: BuildCollectionURL("", "brands"), Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL("", "branches"), Rel: "create-branch", Method: "POST"},
	}
}

// BuildDepartmentListLinks constructs HATEOAS links for the departments catalog
func BuildDepartmentListLinks() []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/departments", apiBase), Rel: "self", Method: "GET"},
	}
}

// BuildCityListLinks constructs HATEOAS links for cities of a department
func BuildCityListLinks(departmentID string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/departments/%s/cities", apiBase, departmentID), Rel: "self", Method: "GET"},
		{Href: fmt.Sprintf("%s/departments", apiBase), Rel: "departments", Method: "GET"},
	}
}

// BuildBranchDetailLinks constructs HATEOAS links for branch query response (HU62)
func BuildBranchDetailLinks(_ string, branchID string, isOwner bool) []Link {
	resourceURL := BuildResourceURL("", "branches", branchID)
	collectionURL := BuildCollectionURL("", "branches")

	links := []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
	}

	if isOwner {
		links = append(links,
			Link{Href: resourceURL, Rel: "update", Method: "PUT"},
			Link{Href: resourceURL, Rel: "delete", Method: "DELETE"},
			Link{Href: collectionURL, Rel: "list", Method: "GET"},
		)
	}
	return links
}

// BuildBranchTypesLinks constructs HATEOAS links for branch types catalog (HU76)
func BuildBranchTypesLinks(_ string) []Link {
	return []Link{
		{Href: branchTypesPath, Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL("", "branches"), Rel: "create-branch", Method: "POST"},
	}
}

// BuildBranchListLinks constructs HATEOAS links for branch list response (HU62)
func BuildBranchListLinks(_ string) []Link {
	collectionURL := BuildCollectionURL("", "branches")
	return []Link{
		{Href: collectionURL, Rel: "self", Method: "GET"},
		{Href: collectionURL, Rel: "create", Method: "POST"},
		{Href: BuildCollectionURL("", "brands"), Rel: "brands", Method: "GET"},
		{Href: branchTypesPath, Rel: "branch-types", Method: "GET"},
	}
}

// BuildBranchDeletedLinks constructs HATEOAS links after branch deletion (HU61)
func BuildBranchDeletedLinks(_ string) []Link {
	collectionURL := BuildCollectionURL("", "branches")
	return []Link{
		{Href: collectionURL, Rel: "list", Method: "GET"},
		{Href: collectionURL, Rel: "create", Method: "POST"},
	}
}

// BuildServiceTypeListLinks constructs HATEOAS links for service types catalog (HU75)
func BuildServiceTypeListLinks(_ string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/service-types", apiBase), Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL("", "services"), Rel: "services", Method: "GET"},
	}
}

// BuildServiceListLinks constructs HATEOAS links for services catalog (HU63)
func BuildServiceListLinks(_, filterType string) []Link {
	servicesURL := BuildCollectionURL("", "services")
	links := []Link{
		{Href: servicesURL, Rel: "self", Method: "GET"},
		{Href: fmt.Sprintf("%s/service-types", apiBase), Rel: "service-types", Method: "GET"},
	}
	if filterType != "" {
		links = append(links, Link{Href: servicesURL, Rel: "all-services", Method: "GET"})
	}
	return links
}

// ========================================
// Motorcycle HATEOAS Functions (HU43-47)
// ========================================

// BuildMotorcycleDetailLinks constructs HATEOAS links for motorcycle response (HU43-47)
func BuildMotorcycleDetailLinks(_ string, motorcycleID string, isOwner bool) []Link {
	resourceURL := BuildResourceURL("", "motorcycles", motorcycleID)
	collectionURL := BuildCollectionURL("", "motorcycles")

	links := []Link{
		{Href: resourceURL, Rel: "self", Method: "GET"},
	}
	if isOwner {
		links = append(links,
			Link{Href: resourceURL, Rel: "update", Method: "PUT"},
			Link{Href: resourceURL, Rel: "delete", Method: "DELETE"},
			Link{Href: collectionURL, Rel: "list", Method: "GET"},
		)
	}
	return links
}

// BuildMotorcycleDeletedLinks constructs HATEOAS links after motorcycle deletion (HU45)
func BuildMotorcycleDeletedLinks(_ string) []Link {
	collectionURL := BuildCollectionURL("", "motorcycles")
	return []Link{
		{Href: collectionURL, Rel: "list", Method: "GET"},
		{Href: collectionURL, Rel: "create", Method: "POST"},
	}
}

// BuildNearbyBranchesLinks constructs HATEOAS links for nearby branches search (HU89)
func BuildNearbyBranchesLinks(_ string, lat, lng, radiusKm float64) []Link {
	nearbyURL := fmt.Sprintf("%s/branches/nearby?lat=%.8f&lng=%.8f&radius=%.1f", apiBase, lat, lng, radiusKm)
	return []Link{
		{Href: nearbyURL, Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL("", "branches"), Rel: "branches", Method: "GET"},
		{Href: branchTypesPath, Rel: "branch-types", Method: "GET"},
	}
}

// BuildMotorcycleReferencesLinks constructs HATEOAS links for motorcycle references catalog (HU50)
func BuildMotorcycleReferencesLinks(_ string) []Link {
	return []Link{
		{Href: motorcycleReferencesPath, Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL("", "motorcycles"), Rel: "motorcycles", Method: "GET"},
		{Href: BuildCollectionURL("", "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildBrandLinesLinks constructs HATEOAS links for brand lines catalog (HU40 - Admin)
func BuildBrandLinesLinks(_, brandID string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/admin/brands/%s/lines", apiBase, brandID), Rel: "self", Method: "GET"},
		{Href: BuildCollectionURL("", "brands"), Rel: "brands", Method: "GET"},
		{Href: motorcycleReferencesPath, Rel: "references", Method: "GET"},
	}
}

// ============================================
// MOTORCYCLE CATEGORY HATEOAS (HU41)
// ============================================

// BuildMotorcycleCategoryListLinks constructs HATEOAS links for the motorcycle categories list
func BuildMotorcycleCategoryListLinks(_ string) []Link {
	return []Link{
		{Href: motorcycleCategoriesPath, Rel: "self", Method: "GET"},
		{Href: motorcycleReferencesPath, Rel: "references", Method: "GET"},
		{Href: BuildCollectionURL("", "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildMotorcycleCategoryItemLinks constructs drill-down HATEOAS links for a single category item
func BuildMotorcycleCategoryItemLinks(_, categoryName string) []Link {
	return []Link{
		{Href: fmt.Sprintf(motorcycleCategoriesPath+"/%s/lines", categoryName), Rel: "lines", Method: "GET"},
	}
}

// BuildCategoryLinesLinks constructs HATEOAS links for lines within a specific category
func BuildCategoryLinesLinks(_, categoryName string) []Link {
	return []Link{
		{Href: fmt.Sprintf(motorcycleCategoriesPath+"/%s/lines", categoryName), Rel: "self", Method: "GET"},
		{Href: motorcycleCategoriesPath, Rel: "categories", Method: "GET"},
		{Href: BuildCollectionURL("", "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildEngineDisplacementLinks constructs HATEOAS links for the engine displacement list (HU49)
func BuildEngineDisplacementLinks(_ string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/engine-displacements", apiBase), Rel: "self", Method: "GET"},
		{Href: motorcycleCategoriesPath, Rel: "categories", Method: "GET"},
		{Href: BuildCollectionURL("", "brands"), Rel: "brands", Method: "GET"},
	}
}

// BuildRatingRangeLinks constructs HATEOAS links for the rating range list (HU48)
func BuildRatingRangeLinks(_ string) []Link {
	return []Link{
		{Href: fmt.Sprintf("%s/rating-ranges", apiBase), Rel: "self", Method: "GET"},
		{Href: fmt.Sprintf("%s/engine-displacements", apiBase), Rel: "displacements", Method: "GET"},
		{Href: motorcycleCategoriesPath, Rel: "categories", Method: "GET"},
	}
}

// ============================================
// MOTORCYCLE EVIDENCE HATEOAS (HU16-19)
// ============================================

// BuildEvidenceDetailLinks constructs HATEOAS links for evidence detail (HU16-19)
func BuildEvidenceDetailLinks(_, motorcycleID, evidenceID string, isOwner bool) []Link {
	evidenceURL := fmt.Sprintf(motorcycleEvidenceIDFmt, motorcycleID, evidenceID)
	listURL := fmt.Sprintf(motorcycleEvidenceFmt, motorcycleID)
	motorcycleURL := BuildResourceURL("", "motorcycles", motorcycleID)

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
func BuildEvidenceListLinks(_, motorcycleID string) []Link {
	listURL := fmt.Sprintf(motorcycleEvidenceFmt, motorcycleID)
	motorcycleURL := BuildResourceURL("", "motorcycles", motorcycleID)
	return []Link{
		{Href: listURL, Rel: "self", Method: "GET"},
		{Href: listURL, Rel: "create", Method: "POST"},
		{Href: motorcycleURL, Rel: "motorcycle", Method: "GET"},
	}
}

// BuildEvidenceDeletedLinks constructs HATEOAS links after evidence deletion (HU19)
func BuildEvidenceDeletedLinks(_, motorcycleID string) []Link {
	listURL := fmt.Sprintf(motorcycleEvidenceFmt, motorcycleID)
	motorcycleURL := BuildResourceURL("", "motorcycles", motorcycleID)
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
func BuildPermissionLinks(_, motorcycleID string) []Link {
	permissionsURL := fmt.Sprintf("%s/motorcycles/%s/permissions", apiBase, motorcycleID)
	motorcycleURL := BuildResourceURL("", "motorcycles", motorcycleID)
	return []Link{
		{Href: permissionsURL, Rel: "list-permissions", Method: "GET"},
		{Href: permissionsURL, Rel: "grant-permission", Method: "POST"},
		{Href: motorcycleURL, Rel: "motorcycle", Method: "GET"},
	}
}
