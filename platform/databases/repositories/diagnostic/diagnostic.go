package diagnostic

import (
	"database/sql"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Diagnostic represents the database model for diagnosticos table
type Diagnostic struct {
	ID                 string         `db:"id"`
	MotorcycleID       string         `db:"motocicleta_id"`
	BranchID           string         `db:"sede_id"`
	BranchName         sql.NullString `db:"branch_name"`
	Date               time.Time      `db:"fecha"`
	ProblemDescription sql.NullString `db:"descripcion_problema"` //nolint:misspell // Spanish DB column
	PossibleSolution   sql.NullString `db:"posible_solucion"`
}

// ToDomain converts the database Diagnostic model to domain entity
func (d *Diagnostic) ToDomain() domain.Diagnostic {
	diagnostic := domain.Diagnostic{
		ID:           d.ID,
		MotorcycleID: d.MotorcycleID,
		BranchID:     d.BranchID,
		Date:         d.Date,
	}

	if d.BranchName.Valid {
		diagnostic.BranchName = d.BranchName.String
	}

	if d.ProblemDescription.Valid {
		desc := d.ProblemDescription.String
		diagnostic.ProblemDescription = &desc
	}

	if d.PossibleSolution.Valid {
		sol := d.PossibleSolution.String
		diagnostic.PossibleSolution = &sol
	}

	return diagnostic
}

// FromDomain converts a domain Diagnostic entity to database model
func FromDomain(diagnostic *domain.Diagnostic) *Diagnostic {
	d := &Diagnostic{
		ID:           diagnostic.ID,
		MotorcycleID: diagnostic.MotorcycleID,
		BranchID:     diagnostic.BranchID,
		Date:         diagnostic.Date,
	}

	if diagnostic.ProblemDescription != nil {
		d.ProblemDescription = sql.NullString{String: *diagnostic.ProblemDescription, Valid: true}
	}

	if diagnostic.PossibleSolution != nil {
		d.PossibleSolution = sql.NullString{String: *diagnostic.PossibleSolution, Valid: true}
	}

	return d
}
