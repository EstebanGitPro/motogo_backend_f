package handlers

// StandardResponse is the unified API response wrapper used in Swagger annotations.
// This type mirrors middleware.APIResponse but is defined here so swag can find it
// when parsing handler annotations.
//
// @Description Standard API response wrapper with success flag, message and optional data
type StandardResponse struct {
	// Whether the operation was successful
	Success bool `json:"success" example:"true"`
	// Business message code (e.g., "MSG_BRANCH_REGISTERED")
	Code string `json:"code,omitempty" example:"MSG_BRANCH_REGISTERED"`
	// Human-readable message content
	Message string `json:"message,omitempty" example:"Sede registrada exitosamente"`
	// Response data payload (type varies by endpoint)
	Data interface{} `json:"data,omitempty"`
}
