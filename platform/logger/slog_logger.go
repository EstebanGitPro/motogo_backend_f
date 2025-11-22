package logger

import "log/slog"


type SlogLogger struct{}


func NewSlogLogger() Logger {
	return &SlogLogger{}
}

func (s *SlogLogger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func (s *SlogLogger) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func (s *SlogLogger) Success(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (s *SlogLogger) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func (s *SlogLogger) Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
}

func (s *SlogLogger) Panic(msg string, args ...any) {
	slog.Error(msg, args...)
}



