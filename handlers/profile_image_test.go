package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

// ============================================
// ProfileImageRequest.Sanitize Tests
// ============================================

func TestProfileImageRequest_Sanitize(t *testing.T) {
	// Arrange
	req := &handlers.ProfileImageRequest{
		ImageURL: "  https://storage.googleapis.com/image.jpg  ",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "https://storage.googleapis.com/image.jpg", req.ImageURL)
}

func TestProfileImageRequest_Sanitize_WithTabs(t *testing.T) {
	// Arrange
	req := &handlers.ProfileImageRequest{
		ImageURL: "\thttps://storage.googleapis.com/test.jpg\n",
	}

	// Act
	req.Sanitize()

	// Assert
	assert.Equal(t, "https://storage.googleapis.com/test.jpg", req.ImageURL)
}

// ============================================
// ProfileImageRequest.ToDomain Tests
// ============================================

func TestProfileImageRequest_ToDomain(t *testing.T) {
	// Arrange
	imageURL := "https://storage.googleapis.com/image.jpg"
	req := &handlers.ProfileImageRequest{
		ImageURL: imageURL,
	}

	// Act
	domain := req.ToDomain()

	// Assert
	assert.NotNil(t, domain.ProfileImageURL)
	assert.Equal(t, imageURL, *domain.ProfileImageURL)
}

// ============================================
// ToProfileImageResponse Tests
// ============================================

func TestToProfileImageResponse_WithImage(t *testing.T) {
	// Arrange
	motorcycleID := "moto-123"
	imageURL := "https://storage.googleapis.com/profile.jpg"

	// Act
	response := handlers.ToProfileImageResponse(motorcycleID, &imageURL)

	// Assert
	assert.Equal(t, motorcycleID, response.MotorcycleID)
	assert.Equal(t, imageURL, response.ProfileImageURL)
}

func TestToProfileImageResponse_WithoutImage(t *testing.T) {
	// Arrange
	motorcycleID := "moto-456"

	// Act
	response := handlers.ToProfileImageResponse(motorcycleID, nil)

	// Assert
	assert.Equal(t, motorcycleID, response.MotorcycleID)
	assert.Empty(t, response.ProfileImageURL)
}

// ============================================
// BuildProfileImageLinks Tests
// ============================================

func TestBuildProfileImageLinks_WithImage(t *testing.T) {
	// Act
	links := handlers.BuildProfileImageLinks("http://api.test.com", "moto-enc-123", true)

	// Assert
	assert.Len(t, links, 4) // self, update, motorcycle, delete

	expectedRels := []string{"self", "update", "motorcycle", "delete"}
	foundRels := make(map[string]bool)
	for _, l := range links {
		foundRels[l.Rel] = true
	}

	for _, rel := range expectedRels {
		assert.True(t, foundRels[rel], "Missing expected rel: "+rel)
	}
}

func TestBuildProfileImageLinks_WithoutImage(t *testing.T) {
	// Act
	links := handlers.BuildProfileImageLinks("http://api.test.com", "moto-enc-456", false)

	// Assert
	assert.Len(t, links, 3) // self, update, motorcycle (no delete)

	// Verify no delete link
	for _, l := range links {
		assert.NotEqual(t, "delete", l.Rel, "Should not have delete link when no image")
	}
}

func TestBuildProfileImageLinks_URLs(t *testing.T) {
	// Act
	links := handlers.BuildProfileImageLinks("http://api", "moto-1", true)

	// Assert
	var selfLink, updateLink, motorcycleLink, deleteLink *handlers.Link
	for i, l := range links {
		switch l.Rel {
		case "self":
			selfLink = &links[i]
		case "update":
			updateLink = &links[i]
		case "motorcycle":
			motorcycleLink = &links[i]
		case "delete":
			deleteLink = &links[i]
		}
	}

	assert.NotNil(t, selfLink)
	assert.Equal(t, "http://api/motorcycles/moto-1/profile-image", selfLink.Href)
	assert.Equal(t, "GET", selfLink.Method)

	assert.NotNil(t, updateLink)
	assert.Equal(t, "http://api/motorcycles/moto-1/profile-image", updateLink.Href)
	assert.Equal(t, "PUT", updateLink.Method)

	assert.NotNil(t, motorcycleLink)
	assert.Equal(t, "http://api/motorcycles/moto-1", motorcycleLink.Href)
	assert.Equal(t, "GET", motorcycleLink.Method)

	assert.NotNil(t, deleteLink)
	assert.Equal(t, "http://api/motorcycles/moto-1/profile-image", deleteLink.Href)
	assert.Equal(t, "DELETE", deleteLink.Method)
}
