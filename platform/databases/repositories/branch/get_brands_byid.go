package branch

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetBranchByID retrieves a branch by its ID including location and brands
func (r *repository) GetBranchByID(ctx context.Context, branchID string) (*domain.Branch, error) {
	var branch domain.Branch
	var franchiseID, profileImageURL sql.NullString
	var locationID, cityID, address sql.NullString
	var latitude, longitude sql.NullFloat64
	var departmentID sql.NullString

	err := r.stmtGetBranchByID.QueryRowContext(ctx, branchID).Scan(
		&branch.ID,
		&branch.RepresentativeID,
		&franchiseID,
		&branch.Name,
		&branch.EstablishmentType,
		&profileImageURL,
		&branch.Status,
		&locationID,
		&cityID,
		&address,
		&latitude,
		&longitude,
		&departmentID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrBranchNotFound
		}
		log.Error(logger.LogBranchRepoGetByIDError, "error", err, "branch_id", branchID)
		return nil, err
	}

	if franchiseID.Valid {
		branch.FranchiseID = &franchiseID.String
	}
	if profileImageURL.Valid {
		branch.ProfileImageURL = &profileImageURL.String
	}

	// Parse location if exists
	if locationID.Valid {
		branch.Location = &domain.Location{
			ID:           locationID.String,
			BranchID:     branchID,
			CityID:       cityID.String,
			DepartmentID: departmentID.String,
			Address:      address.String,
		}
		if latitude.Valid {
			branch.Location.Latitude = &latitude.Float64
		}
		if longitude.Valid {
			branch.Location.Longitude = &longitude.Float64
		}
	}

	// Get brands
	brands, err := r.getBranchBrands(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchRepoBrandGetError, "error", err, "branch_id", branchID)
	}
	branch.Brands = brands

	return &branch, nil
}
