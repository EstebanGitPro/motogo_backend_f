package motorcycle

import (
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

type Motorcycle struct {
	ID             string         `db:"id"`
	LicensePlate   string         `db:"license_plate"`
	ReferenceID    sql.NullString `db:"reference_id"`
	OwnerID        string         `db:"owner_id"`
	Year           sql.NullInt64  `db:"year"`
	CurrentMileage sql.NullInt64  `db:"current_mileage"`
	OwnerNotes     sql.NullString `db:"owner_notes"`
}

type MotorcycleReference struct {
	ID                 sql.NullString `db:"ref_id"`
	BrandID            sql.NullString `db:"brand_id"`
	BrandName          sql.NullString `db:"brand_name"`
	Model              sql.NullString `db:"model"`
	Category           sql.NullString `db:"category"`
	EngineDisplacement sql.NullInt64  `db:"engine_displacement"`
}

func (m *Motorcycle) ToDomain(ref *MotorcycleReference) domain.Motorcycle {
	motorcycle := domain.Motorcycle{
		ID:           m.ID,
		LicensePlate: m.LicensePlate,
		OwnerID:      m.OwnerID,
	}

	if m.Year.Valid {
		year := int(m.Year.Int64)
		motorcycle.Year = &year
	}
	if m.CurrentMileage.Valid {
		mileage := int(m.CurrentMileage.Int64)
		motorcycle.CurrentMileage = &mileage
	}
	if m.OwnerNotes.Valid {
		notes := m.OwnerNotes.String
		motorcycle.OwnerNotes = &notes
	}

	if ref != nil && ref.ID.Valid {
		motorcycle.Reference = &domain.MotorcycleReference{
			ID:      ref.ID.String,
			BrandID: ref.BrandID.String,
		}
		if ref.BrandName.Valid {
			motorcycle.Reference.BrandName = ref.BrandName.String
		}
		if ref.Model.Valid {
			motorcycle.Reference.Model = ref.Model.String
		}
		if ref.Category.Valid {
			motorcycle.Reference.Category = ref.Category.String
		}
		if ref.EngineDisplacement.Valid {
			motorcycle.Reference.EngineDisplacement = int(ref.EngineDisplacement.Int64)
		}
	}

	return motorcycle
}
