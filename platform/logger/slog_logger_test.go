package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================
// NewSlogLogger Tests
// ============================================

func TestNewSlogLogger_ReturnsInstance(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotNil(t, logger)
	assert.Implements(t, (*Logger)(nil), logger)
}

// ============================================
// WithTraceID Tests
// ============================================

func TestWithTraceID_ReturnsNewLogger(t *testing.T) {
	logger := NewSlogLogger()
	traceID := "trace-123-abc"

	loggerWithTrace := logger.WithTraceID(traceID)

	assert.NotNil(t, loggerWithTrace)
	assert.Implements(t, (*Logger)(nil), loggerWithTrace)
	// The new logger should be different from the original
	assert.NotSame(t, logger, loggerWithTrace)
}

func TestWithTraceID_ChainingWorks(t *testing.T) {
	logger := NewSlogLogger()

	logger1 := logger.WithTraceID("trace-1")
	logger2 := logger1.WithTraceID("trace-2")

	assert.NotNil(t, logger1)
	assert.NotNil(t, logger2)
	assert.NotSame(t, logger1, logger2)
}

// ============================================
// Log Method Tests (basic invocation - no panics)
// ============================================

func TestInfo_DoesNotPanic(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Info("test message", "key", "value")
	})
}

func TestError_DoesNotPanic(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Error("error message", "error", "some error")
	})
}

func TestDebug_DoesNotPanic(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Debug("debug message", "detail", "value")
	})
}

func TestSuccess_DoesNotPanic(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Success("success message")
	})
}

func TestWarn_DoesNotPanic(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Warn("warning message", "warning", "details")
	})
}

func TestFatal_DoesNotPanic(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Fatal("fatal message")
	})
}

func TestPanic_DoesNotPanic(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Panic("panic message")
	})
}

// ============================================
// WithTraceID enrichWithContext Tests
// ============================================

func TestEnrichWithContext_WithTraceID(t *testing.T) {
	logger := NewSlogLogger().(*SlogLogger)
	loggerWithTrace := logger.WithTraceID("trace-abc-123").(*SlogLogger)

	// Test that logging with trace ID doesn't panic
	assert.NotPanics(t, func() {
		loggerWithTrace.Info("message with trace", "key", "value")
	})
}

func TestEnrichWithContext_WithoutTraceID(t *testing.T) {
	logger := NewSlogLogger().(*SlogLogger)

	args := logger.enrichWithContext("key", "value")

	// Without traceID, args should be unchanged
	assert.Equal(t, []any{"key", "value"}, args)
}

func TestEnrichWithContext_WithTraceIDSet(t *testing.T) {
	logger := NewSlogLogger().(*SlogLogger)
	loggerWithTrace := logger.WithTraceID("test-trace").(*SlogLogger)

	args := loggerWithTrace.enrichWithContext("key", "value")

	// With traceID, it should be prepended
	assert.Len(t, args, 4)
	assert.Equal(t, "traceID", args[0])
	assert.Equal(t, "test-trace", args[1])
	assert.Equal(t, "key", args[2])
	assert.Equal(t, "value", args[3])
}

// ============================================
// Log Method with Multiple Args Tests
// ============================================

func TestInfo_WithMultipleArgs(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Info("message", "key1", "value1", "key2", 123, "key3", true)
	})
}

func TestError_WithNoArgs(t *testing.T) {
	logger := NewSlogLogger()

	assert.NotPanics(t, func() {
		logger.Error("simple error")
	})
}

// ============================================
// Interface Compliance Tests
// ============================================

func TestSlogLogger_ImplementsLoggerInterface(t *testing.T) {
	var _ Logger = &SlogLogger{}
	var _ Logger = NewSlogLogger()
}
