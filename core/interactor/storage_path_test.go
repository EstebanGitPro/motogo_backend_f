package interactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// extractStoragePath Tests
// ============================================

func TestExtractStoragePath_ValidFirebaseURL(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/branches%2Fid123%2Fprofile.jpg?alt=media&token=abc-123"
	path := extractStoragePath(url)
	assert.Equal(t, "branches/id123/profile.jpg", path)
}

func TestExtractStoragePath_NestedPath(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/branches%2Fid123%2Fimages%2Fprofile.jpg?alt=media"
	path := extractStoragePath(url)
	assert.Equal(t, "branches/id123/images/profile.jpg", path)
}

func TestExtractStoragePath_InvalidURL(t *testing.T) {
	path := extractStoragePath("://invalid-url")
	assert.Equal(t, "", path)
}

func TestExtractStoragePath_NonFirebaseURL(t *testing.T) {
	path := extractStoragePath("https://example.com/images/photo.jpg")
	assert.Equal(t, "", path)
}

func TestExtractStoragePath_EmptyString(t *testing.T) {
	path := extractStoragePath("")
	assert.Equal(t, "", path)
}

// ============================================
// sameStoragePath Tests
// ============================================

func TestSameStoragePath_SamePathDifferentToken(t *testing.T) {
	url1 := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/branches%2Fid123%2Fprofile.jpg?alt=media&token=old-token"
	url2 := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/branches%2Fid123%2Fprofile.jpg?alt=media&token=new-token"
	assert.True(t, sameStoragePath(url1, url2), "should be true: same path, different token")
}

func TestSameStoragePath_DifferentPaths(t *testing.T) {
	url1 := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/branches%2Fid123%2Fprofile.jpg?alt=media&token=token1"
	url2 := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/branches%2Fid456%2Fprofile.jpg?alt=media&token=token2"
	assert.False(t, sameStoragePath(url1, url2), "should be false: different branch IDs")
}

func TestSameStoragePath_IdenticalURLs(t *testing.T) {
	url := "https://firebasestorage.googleapis.com/v0/b/motogo-xxx.appspot.com/o/branches%2Fid123%2Fprofile.jpg?alt=media&token=same"
	assert.True(t, sameStoragePath(url, url), "identical URLs should match")
}

func TestSameStoragePath_FallbackWhenBothInvalid(t *testing.T) {
	// Both non-Firebase URLs, falls back to string comparison
	assert.True(t, sameStoragePath("http://example.com/a.jpg", "http://example.com/a.jpg"))
	assert.False(t, sameStoragePath("http://example.com/a.jpg", "http://example.com/b.jpg"))
}

func TestSameStoragePath_FallbackWhenOneInvalid(t *testing.T) {
	valid := "https://firebasestorage.googleapis.com/v0/b/bucket/o/file.jpg?alt=media&token=x"
	invalid := "http://example.com/file.jpg"
	// One path is empty → falls back to full URL comparison → different → false
	assert.False(t, sameStoragePath(valid, invalid))
}
