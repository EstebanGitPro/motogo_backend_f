package handlers

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

// CreateFranchiseRequest represents the request body for creating a franchise (HU26)
type CreateFranchiseRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	BranchIDs   []string `json:"branch_ids" binding:"required,min=1"`
}

// UpdateFranchiseRequest represents the request body for updating a franchise (HU27)
type UpdateFranchiseRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

// AddBranchToFranchiseRequest represents the request to add a branch to a franchise
type AddBranchToFranchiseRequest struct {
	BranchID string `json:"branch_id" binding:"required"`
}

// FranchiseResponse represents a franchise in API responses
type FranchiseResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	BranchCount int     `json:"branch_count,omitempty"`
	Links       []Link  `json:"_links,omitempty"`
}

// FranchiseListResponse represents a list of franchises
type FranchiseListResponse struct {
	Franchises []FranchiseResponse `json:"franchises"`
	Total      int                 `json:"total"`
	Links      []Link              `json:"_links,omitempty"`
}

// ToFranchiseDomain converts CreateFranchiseRequest to domain.Franchise
func (r *CreateFranchiseRequest) ToFranchiseDomain() domain.Franchise {
	return domain.Franchise{
		Name:        r.Name,
		Description: r.Description,
	}
}

// ToFranchiseDomain converts UpdateFranchiseRequest to domain.Franchise
func (r *UpdateFranchiseRequest) ToFranchiseDomain(id string) domain.Franchise {
	return domain.Franchise{
		ID:          id,
		Name:        r.Name,
		Description: r.Description,
	}
}

// ToFranchiseResponse converts domain.Franchise to FranchiseResponse
func ToFranchiseResponse(franchise *domain.Franchise, branchCount int, franchiseEncodedID string, links []Link) FranchiseResponse {
	return FranchiseResponse{
		ID:          franchiseEncodedID,
		Name:        franchise.Name,
		Description: franchise.Description,
		BranchCount: branchCount,
		Links:       links,
	}
}

// ToFranchiseListResponse converts a slice of franchises to FranchiseListResponse
func ToFranchiseListResponse(franchises []domain.Franchise, encodedIDs []string, baseURL string) FranchiseListResponse {
	responses := make([]FranchiseResponse, len(franchises))
	for i, f := range franchises {
		links := BuildFranchiseLinks(baseURL, encodedIDs[i])
		responses[i] = FranchiseResponse{
			ID:          encodedIDs[i],
			Name:        f.Name,
			Description: f.Description,
			Links:       links,
		}
	}
	return FranchiseListResponse{
		Franchises: responses,
		Total:      len(franchises),
		Links: []Link{
			{Rel: "self", Href: BuildCollectionURL(baseURL, "franchises"), Method: "GET"},
			{Rel: "create", Href: BuildCollectionURL(baseURL, "franchises"), Method: "POST"},
		},
	}
}
