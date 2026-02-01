package handlers_test

import (
	"testing"

	"github.com/EstebanGitPro/motogo-backend/handlers"
	"github.com/stretchr/testify/assert"
)

func TestTrimString_RemovesWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"leading spaces", "  hello", "hello"},
		{"trailing spaces", "hello  ", "hello"},
		{"both sides", "  hello  ", "hello"},
		{"tabs", "\thello\t", "hello"},
		{"newlines", "\nhello\n", "hello"},
		{"mixed whitespace", "  \t\n hello \n\t  ", "hello"},
		{"no whitespace", "hello", "hello"},
		{"empty string", "", ""},
		{"only whitespace", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.TrimString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTrimStringPtr_NilInput(t *testing.T) {
	result := handlers.TrimStringPtr(nil)
	assert.Nil(t, result)
}

func TestTrimStringPtr_ValidInput(t *testing.T) {
	input := "  hello  "
	result := handlers.TrimStringPtr(&input)

	assert.NotNil(t, result)
	assert.Equal(t, "hello", *result)
}

func TestTrimStringPtr_EmptyString(t *testing.T) {
	input := "   "
	result := handlers.TrimStringPtr(&input)

	assert.NotNil(t, result)
	assert.Equal(t, "", *result)
}
