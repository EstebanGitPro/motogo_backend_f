package diagnostic

import (
	"database/sql"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Diagnostic represents the database model for diagnosticos table
type Diagnostic struct {
	ID                 string          `db:"id"`
	MotorcycleID       string          `db:"motocicleta_id"`
	BranchID           string          `db:"sede_id"`
	Date               time.Time       `db:"fecha"`
	ProblemDescription sql.NullString  `db:"descripcion_problema"`
	PossibleSolution   sql.NullString  `db:"posible_solucion"`
	LaborQuote         sql.NullFloat64 `db:"cotizacion_mano_obra"`
	PartsQuote         sql.NullFloat64 `db:"cotizacion_repuestos"`
	EstimatedTime      sql.NullString  `db:"tiempo_estimado"`
	SentViaWhatsApp    bool            `db:"enviado_whatsapp"`
}

// DiagnosticEvidence represents the database model for evidencias_diagnostico table
type DiagnosticEvidence struct {
	ID           string         `db:"id"`
	DiagnosticID string         `db:"diagnostico_id"`
	ImageURL     string         `db:"url_imagen"`
	Description  sql.NullString `db:"descripcion"`
	CreatedAt    time.Time      `db:"created_at"`
}

// ToDomain converts the database Diagnostic model to domain entity
func (d *Diagnostic) ToDomain() domain.Diagnostic {
	diagnostic := domain.Diagnostic{
		ID:              d.ID,
		MotorcycleID:    d.MotorcycleID,
		BranchID:        d.BranchID,
		Date:            d.Date,
		SentViaWhatsApp: d.SentViaWhatsApp,
	}

	if d.ProblemDescription.Valid {
		desc := d.ProblemDescription.String
		diagnostic.ProblemDescription = &desc
	}

	if d.PossibleSolution.Valid {
		sol := d.PossibleSolution.String
		diagnostic.PossibleSolution = &sol
	}

	if d.LaborQuote.Valid {
		lq := d.LaborQuote.Float64
		diagnostic.LaborQuote = &lq
	}

	if d.PartsQuote.Valid {
		pq := d.PartsQuote.Float64
		diagnostic.PartsQuote = &pq
	}

	if d.EstimatedTime.Valid {
		et := d.EstimatedTime.String
		diagnostic.EstimatedTime = &et
	}

	return diagnostic
}

// FromDomain converts a domain Diagnostic entity to database model
func FromDomain(diagnostic *domain.Diagnostic) *Diagnostic {
	d := &Diagnostic{
		ID:              diagnostic.ID,
		MotorcycleID:    diagnostic.MotorcycleID,
		BranchID:        diagnostic.BranchID,
		Date:            diagnostic.Date,
		SentViaWhatsApp: diagnostic.SentViaWhatsApp,
	}

	if diagnostic.ProblemDescription != nil {
		d.ProblemDescription = sql.NullString{String: *diagnostic.ProblemDescription, Valid: true}
	}

	if diagnostic.PossibleSolution != nil {
		d.PossibleSolution = sql.NullString{String: *diagnostic.PossibleSolution, Valid: true}
	}

	if diagnostic.LaborQuote != nil {
		d.LaborQuote = sql.NullFloat64{Float64: *diagnostic.LaborQuote, Valid: true}
	}

	if diagnostic.PartsQuote != nil {
		d.PartsQuote = sql.NullFloat64{Float64: *diagnostic.PartsQuote, Valid: true}
	}

	if diagnostic.EstimatedTime != nil {
		d.EstimatedTime = sql.NullString{String: *diagnostic.EstimatedTime, Valid: true}
	}

	return d
}

// EvidenceToDomain converts the database DiagnosticEvidence model to domain entity
func (e *DiagnosticEvidence) EvidenceToDomain() domain.DiagnosticEvidence {
	evidence := domain.DiagnosticEvidence{
		ID:           e.ID,
		DiagnosticID: e.DiagnosticID,
		ImageURL:     e.ImageURL,
		CreatedAt:    e.CreatedAt,
	}

	if e.Description.Valid {
		desc := e.Description.String
		evidence.Description = &desc
	}

	return evidence
}

// EvidenceFromDomain converts a domain DiagnosticEvidence entity to database model
func EvidenceFromDomain(evidence *domain.DiagnosticEvidence) *DiagnosticEvidence {
	e := &DiagnosticEvidence{
		ID:           evidence.ID,
		DiagnosticID: evidence.DiagnosticID,
		ImageURL:     evidence.ImageURL,
		CreatedAt:    evidence.CreatedAt,
	}

	if evidence.Description != nil {
		e.Description = sql.NullString{String: *evidence.Description, Valid: true}
	}

	return e
}
