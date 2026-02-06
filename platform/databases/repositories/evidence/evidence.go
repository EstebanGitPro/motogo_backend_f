package evidence

import (
	"database/sql"
	"time"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// Evidence represents the database model for motorcycle_evidence table
type Evidence struct {
	ID           string         `db:"id"`
	MotorcycleID string         `db:"motorcycle_id"`
	Angle        sql.NullString `db:"angle"`
	ImageURL     string         `db:"image_url"`
	UploadDate   time.Time      `db:"upload_date"`
}

// ToDomain converts the database model to domain entity
func (e *Evidence) ToDomain() domain.MotorcycleEvidence {
	evidence := domain.MotorcycleEvidence{
		ID:           e.ID,
		MotorcycleID: e.MotorcycleID,
		ImageURL:     e.ImageURL,
		UploadDate:   e.UploadDate,
	}

	if e.Angle.Valid {
		angle := e.Angle.String
		evidence.Angle = &angle
	}

	return evidence
}

// FromDomain converts a domain entity to database model
func FromDomain(evidence *domain.MotorcycleEvidence) *Evidence {
	e := &Evidence{
		ID:           evidence.ID,
		MotorcycleID: evidence.MotorcycleID,
		ImageURL:     evidence.ImageURL,
		UploadDate:   evidence.UploadDate,
	}

	if evidence.Angle != nil {
		e.Angle = sql.NullString{String: *evidence.Angle, Valid: true}
	}

	return e
}
