package handlers

import "strings"

// RateServiceItemRequest represents the POST /completed-services/:id/items/:itemId/rating request body (HU48).
type RateServiceItemRequest struct {
	Rating  int     `json:"rating" binding:"required"`
	Comment *string `json:"comment,omitempty"`
}

// Sanitize trims whitespace from string fields.
func (r *RateServiceItemRequest) Sanitize() {
	if r.Comment != nil {
		trimmed := strings.TrimSpace(*r.Comment)
		r.Comment = &trimmed
	}
}

// ──────────────────────────────────────────
// GET /services/:id/reviews — Response DTOs
// ──────────────────────────────────────────

// ServiceReviewResponse is the top-level JSON response for GET /services/:id/reviews
type ServiceReviewResponse struct {
	ServiceID     string                 `json:"service_id"`
	ServiceName   string                 `json:"service_name"`
	AverageRating float64                `json:"average_rating"`
	TotalReviews  int                    `json:"total_reviews"`
	Breakdown     ReviewBreakdownDTO     `json:"breakdown"`
	Reviews       []ServiceReviewItemDTO `json:"reviews"`
}

// ReviewBreakdownDTO represents the star-count breakdown
type ReviewBreakdownDTO struct {
	Star5 int `json:"5"`
	Star4 int `json:"4"`
	Star3 int `json:"3"`
	Star2 int `json:"2"`
	Star1 int `json:"1"`
}

// ServiceReviewItemDTO represents a single review in the response
type ServiceReviewItemDTO struct {
	ReviewerName    string  `json:"reviewer_name"`
	Rating          int     `json:"rating"`
	Comment         *string `json:"comment,omitempty"`
	RatedAt         string  `json:"rated_at"`
	MotorcycleModel *string `json:"motorcycle_model,omitempty"`
}
