package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	utils "github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/stretchr/testify/assert"
)

func TestFindModuleRoot_Success(t *testing.T) {
	// Act
	root, err := utils.FindModuleRoot()

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, root)

	// Verify go.mod exists in the returned root
	gomodPath := filepath.Join(root, "go.mod")
	_, statErr := os.Stat(gomodPath)
	assert.NoError(t, statErr, "go.mod should exist in module root")
}

func TestFindModuleRoot_ReturnsAbsolutePath(t *testing.T) {
	// Act
	root, err := utils.FindModuleRoot()

	// Assert
	assert.NoError(t, err)
	assert.True(t, filepath.IsAbs(root), "Module root should be an absolute path")
}

func TestFindModuleRoot_ConsistentResults(t *testing.T) {
	// Act - call multiple times
	root1, err1 := utils.FindModuleRoot()
	root2, err2 := utils.FindModuleRoot()

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, root1, root2, "Should return same root on multiple calls")
}
