package domain_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

func TestSetID_GeneratesValidID(t *testing.T) {
	// Arrange
	person := domain.Person{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Act
	person.SetID()

	// Assert
	assert.NotEmpty(t, person.ID, "ID should not be empty after SetID()")
	assert.Greater(t, len(person.ID), 0, "ID should have a length greater than 0")
}

func TestSetID_GeneratesUniqueIDs(t *testing.T) {
	// Arrange
	person1 := domain.Person{Email: "person1@example.com"}
	person2 := domain.Person{Email: "person2@example.com"}

	// Act
	person1.SetID()
	person2.SetID()

	// Assert
	assert.NotEqual(t, person1.ID, person2.ID, "SetID() should generate unique IDs")
}

func TestToLogger_ReturnsCorrectFormat(t *testing.T) {
	// Arrange
	person := domain.Person{
		ID:    "test-id-123",
		Email: "test@example.com",
		Role:  "customer",
	}

	// Act
	logFields := person.ToLogger()

	// Assert
	assert.Equal(t, 3, len(logFields), "ToLogger() should return 3 fields")
	assert.Contains(t, logFields, "id:test-id-123")
	assert.Contains(t, logFields, "email:test@example.com")
	assert.Contains(t, logFields, "role:customer")
}

func TestToLogger_HandlesEmptyFields(t *testing.T) {
	// Arrange
	person := domain.Person{}

	// Act
	logFields := person.ToLogger()

	// Assert
	assert.Equal(t, 3, len(logFields), "ToLogger() should return 3 fields even with empty values")
	assert.Contains(t, logFields, "id:")
	assert.Contains(t, logFields, "email:")
	assert.Contains(t, logFields, "role:")
}
