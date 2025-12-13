package domain_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

func TestMessageSetID_GeneratesValidID(t *testing.T) {
	// Arrange
	message := domain.Message{
		Code:    "TEST_001",
		Title:   "Test",
		Content: "Test content",
	}

	// Act
	message.SetID()

	// Assert
	assert.NotEmpty(t, message.ID)
	assert.Greater(t, len(message.ID), 0)
}

func TestMessageSetID_GeneratesUniqueIDs(t *testing.T) {
	// Arrange
	message1 := domain.Message{Code: "TEST_001"}
	message2 := domain.Message{Code: "TEST_002"}

	// Act
	message1.SetID()
	message2.SetID()

	// Assert
	assert.NotEqual(t, message1.ID, message2.ID)
}

func TestMessageToLogger_ReturnsCorrectFormat(t *testing.T) {
	// Arrange
	message := domain.Message{
		ID:     "msg-123",
		Code:   "TEST_CODE_001",
		Type:   domain.TypeError,
		Module: "TEST_MODULE",
	}

	// Act
	logFields := message.ToLogger()

	// Assert
	assert.Equal(t, 4, len(logFields))
	assert.Contains(t, logFields, "id:msg-123")
	assert.Contains(t, logFields, "code:TEST_CODE_001")
	assert.Contains(t, logFields, "type:ERROR")
	assert.Contains(t, logFields, "module:TEST_MODULE")
}

func TestMessageValidate_Success(t *testing.T) {
	// Arrange
	message := domain.Message{
		Code:    "TEST_CODE_001",
		Title:   "Test",
		Content: "Content",
	}

	// Act
	err := message.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestMessageValidate_CodeRequired(t *testing.T) {
	// Arrange
	message := domain.Message{
		Title:   "Test",
		Content: "Content",
		// Code missing
	}

	// Act
	err := message.Validate()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMessageCodeRequired, err)
}

func TestMessageTypes_Constants(t *testing.T) {
	// Test that message type constants are defined correctly
	assert.Equal(t, domain.MessageType("ERROR"), domain.TypeError)
	assert.Equal(t, domain.MessageType("EXITO"), domain.TypeSuccess)
	assert.Equal(t, domain.MessageType("WARNING"), domain.TypeWarning)
	assert.Equal(t, domain.MessageType("INFO"), domain.TypeInfo)
	assert.Equal(t, domain.MessageType("DEBUG"), domain.TypeDebug)
}
