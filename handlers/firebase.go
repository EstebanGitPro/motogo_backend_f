package handlers

// FirebaseTokenResponse represents the response for GET /auth/firebase-token
type FirebaseTokenResponse struct {
	FirebaseToken string `json:"firebase_token"`
}
