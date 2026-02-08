package domain

import "errors"

// DiagnosticPermission represents a per-branch diagnostic viewing permission
type DiagnosticPermission struct {
	ID           string // UUID primary key
	MotorcycleID string // FK to motorcycle
	BranchID     string // FK to authorized branch
	Active       bool   // Whether the permission is currently active
}

// Diagnostic Permission business errors
var (
	ErrPermissionNotFound     = errors.New("DIAGNOSTIC_PERMISSION_NOT_FOUND")
	ErrPermissionCannotSave   = errors.New("DIAGNOSTIC_PERMISSION_CANNOT_SAVE")
	ErrPermissionCannotDelete = errors.New("DIAGNOSTIC_PERMISSION_CANNOT_DELETE")
)
