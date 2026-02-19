package domain

import (
	"time"

	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
)

// Diagnostic represents a motorcycle diagnostic request (HU11-14)
type Diagnostic struct {
	ID                 string    // UUID primary key
	MotorcycleID       string    // FK to motorcycle
	BranchID           string    // FK to branch issuing diagnostic
	Date               time.Time // Diagnostic date
	ProblemDescription *string   // User description of the problem
	PossibleSolution   *string   // Response from the branch
}

func (d *Diagnostic) SetID() {
	d.ID = uuid.Generate()
}
