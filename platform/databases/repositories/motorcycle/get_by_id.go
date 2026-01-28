package motorcycle

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetByID(ctx context.Context, motorcycleID string) (*domain.Motorcycle, error) {
	row := r.stmtGetByID.QueryRowContext(ctx, motorcycleID)

	var m Motorcycle
	var ref MotorcycleReference

	err := row.Scan(
		&m.ID, &m.LicensePlate, &m.ReferenceID, &m.OwnerID, &m.Year, &m.CurrentMileage, &m.OwnerNotes,
		&ref.ID, &ref.BrandID, &ref.BrandName, &ref.Model, &ref.Category, &ref.EngineDisplacement,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrMotorcycleNotFound
	}
	if err != nil {
		log.Error(logger.LogMotorcycleRepoGetByIDError, "error", err, "motorcycle_id", motorcycleID)
		return nil, err
	}

	motorcycle := m.ToDomain(&ref)
	return &motorcycle, nil
}

func (r *repository) GetByOwnerID(ctx context.Context, ownerID string) ([]domain.Motorcycle, error) {
	rows, err := r.stmtGetByOwnerID.QueryContext(ctx, ownerID)
	if err != nil {
		log.Error(logger.LogMotorcycleRepoGetByOwnerError, "error", err, "owner_id", ownerID)
		return nil, err
	}
	defer rows.Close()

	var motorcycles []domain.Motorcycle

	for rows.Next() {
		var m Motorcycle
		var ref MotorcycleReference

		err := rows.Scan(
			&m.ID, &m.LicensePlate, &m.ReferenceID, &m.OwnerID, &m.Year, &m.CurrentMileage, &m.OwnerNotes,
			&ref.ID, &ref.BrandID, &ref.BrandName, &ref.Model, &ref.Category, &ref.EngineDisplacement,
		)
		if err != nil {
			log.Error(logger.LogMotorcycleRepoGetByOwnerScan, "error", err, "owner_id", ownerID)
			return nil, err
		}

		motorcycle := m.ToDomain(&ref)
		motorcycles = append(motorcycles, motorcycle)
	}

	if err = rows.Err(); err != nil {
		log.Error(logger.LogMotorcycleRepoGetByOwnerIter, "error", err, "owner_id", ownerID)
		return nil, err
	}

	return motorcycles, nil
}
