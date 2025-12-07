package logger

// Logger defines the interface for centralized logging in the application.
// All components should use this interface for structured logging.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Success(msg string, args ...any)
	Warn(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)

	// WithTraceID returns a new Logger instance with the traceID pre-configured.
	// All logs from this logger will include the traceID for correlation in Loki.
	WithTraceID(traceID string) Logger
}
