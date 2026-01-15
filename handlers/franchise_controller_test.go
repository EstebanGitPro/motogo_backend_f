package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// BuildFranchiseLinks Tests
// ============================================

func TestBuildFranchiseLinks_AllLinksPresent(t *testing.T) {
	links := BuildFranchiseLinks("http://localhost:8080", "encoded-franchise-123")

	assert.NotEmpty(t, links)
	assert.Len(t, links, 6, "Should have 6 HATEOAS links")

	// Check for expected link relations
	relFound := make(map[string]bool)
	for _, link := range links {
		relFound[link.Rel] = true
	}

	assert.True(t, relFound["self"], "Should have 'self' link")
	assert.True(t, relFound["update"], "Should have 'update' link")
	assert.True(t, relFound["delete"], "Should have 'delete' link")
	assert.True(t, relFound["list"], "Should have 'list' link")
	assert.True(t, relFound["branches"], "Should have 'branches' link")
	assert.True(t, relFound["add-branch"], "Should have 'add-branch' link")
}

func TestBuildFranchiseLinks_CorrectURLs(t *testing.T) {
	baseURL := "http://localhost:8080"
	franchiseID := "abc123"

	links := BuildFranchiseLinks(baseURL, franchiseID)

	// Verify URLs contain expected patterns
	for _, link := range links {
		assert.Contains(t, link.Href, baseURL, "Link should contain base URL")

		switch link.Rel {
		case "self":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "GET", link.Method)
		case "update":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "PUT", link.Method)
		case "delete":
			assert.Contains(t, link.Href, franchiseID)
			assert.Equal(t, "DELETE", link.Method)
		case "list":
			assert.Equal(t, "GET", link.Method)
		case "branches":
			assert.Contains(t, link.Href, "franchise_id="+franchiseID)
			assert.Equal(t, "GET", link.Method)
		case "add-branch":
			assert.Contains(t, link.Href, franchiseID+"/branches")
			assert.Equal(t, "POST", link.Method)
		}
	}
}

func TestBuildFranchiseLinks_EmptyID(t *testing.T) {
	links := BuildFranchiseLinks("http://localhost:8080", "")

	assert.NotEmpty(t, links)
	// Should still return links, just with empty ID in URL
}

func TestBuildFranchiseLinks_DifferentBaseURLs(t *testing.T) {
	testCases := []string{
		"http://localhost:8080",
		"https://api.example.com",
		"http://192.168.1.1:3000",
	}

	for _, baseURL := range testCases {
		links := BuildFranchiseLinks(baseURL, "test-id")
		assert.NotEmpty(t, links)
		for _, link := range links {
			assert.Contains(t, link.Href, baseURL)
		}
	}
}
