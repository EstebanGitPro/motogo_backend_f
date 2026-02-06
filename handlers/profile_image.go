package handlers

import "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"

// === Profile Image DTOs (HU36-39) ===

// ProfileImageRequest represents the request to update profile image (HU36/37)
type ProfileImageRequest struct {
	ImageURL string `json:"image_url" binding:"required"`
}

// ProfileImageResponse represents the profile image response (HU38)
type ProfileImageResponse struct {
	MotorcycleID    string `json:"motorcycle_id"`
	ProfileImageURL string `json:"profile_image_url,omitempty"`
	Links           []Link `json:"_links,omitempty"`
}

// ToDomain converts ProfileImageRequest to domain updates
func (r *ProfileImageRequest) ToDomain() domain.Motorcycle {
	return domain.Motorcycle{
		ProfileImageURL: &r.ImageURL,
	}
}

// ToProfileImageResponse converts domain to response
func ToProfileImageResponse(motorcycleID string, imageURL *string) ProfileImageResponse {
	response := ProfileImageResponse{
		MotorcycleID: motorcycleID,
	}
	if imageURL != nil {
		response.ProfileImageURL = *imageURL
	}
	return response
}

// BuildProfileImageLinks builds HATEOAS links for profile image operations
func BuildProfileImageLinks(baseURL, motorcycleID string, hasImage bool) []Link {
	links := []Link{
		{
			Href:   baseURL + "/motorcycles/" + motorcycleID + "/profile-image",
			Rel:    "self",
			Method: "GET",
		},
		{
			Href:   baseURL + "/motorcycles/" + motorcycleID + "/profile-image",
			Rel:    "update",
			Method: "PUT",
		},
		{
			Href:   baseURL + "/motorcycles/" + motorcycleID,
			Rel:    "motorcycle",
			Method: "GET",
		},
	}

	if hasImage {
		links = append(links, Link{
			Href:   baseURL + "/motorcycles/" + motorcycleID + "/profile-image",
			Rel:    "delete",
			Method: "DELETE",
		})
	}

	return links
}
