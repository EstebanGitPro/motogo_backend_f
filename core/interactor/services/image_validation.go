package services

import (
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// ValidFirebaseStorageDomains lists valid Firebase Storage URL domains
var ValidFirebaseStorageDomains = []string{
	"firebasestorage.googleapis.com",
}

// ValidateImageURL checks if a single URL is a valid Firebase Storage URL
// Returns nil if empty (optional field) or valid
func ValidateImageURL(url string) error {
	if url == "" {
		return nil
	}

	for _, validDomain := range ValidFirebaseStorageDomains {
		if strings.Contains(url, validDomain) {
			return nil
		}
	}

	return domain.ErrInvalidImageURL
}

// ValidateImageURLs checks if multiple URLs are valid Firebase Storage URLs
// Returns nil if all URLs are empty or valid
func ValidateImageURLs(urls []string) error {
	for _, url := range urls {
		if err := ValidateImageURL(url); err != nil {
			return err
		}
	}
	return nil
}
