package domain

import "time"

// Diagnostic represents a motorcycle diagnostic request (HU11-14)
type Diagnostic struct {
	ID                 string               // UUID primary key
	MotorcycleID       string               // FK to motorcycle
	BranchID           string               // FK to branch issuing diagnostic
	Date               time.Time            // Diagnostic date
	ProblemDescription *string              // User description of the problem
	PossibleSolution   *string              // Response from the branch
	SentViaWhatsApp    bool                 // Whether sent via WhatsApp
	Evidence           []DiagnosticEvidence // Associated photos (populated on read)
}

// DiagnosticEvidence represents a photo attached to a diagnostic request (HU11)
type DiagnosticEvidence struct {
	ID           string    // UUID primary key
	DiagnosticID string    // FK to diagnostic
	ImageURL     string    // Firebase Storage URL
	Description  *string   // Optional description
	CreatedAt    time.Time // Upload timestamp
}
