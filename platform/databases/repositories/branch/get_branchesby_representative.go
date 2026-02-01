package branch

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

func (r *repository) GetBranchesByRepresentative(ctx context.Context, representativeID string) ([]domain.Branch, error) {
	rows, err := r.stmtGetBranchesByRepresentative.QueryContext(ctx, representativeID)
	if err != nil {
		log.Error(logger.LogBranchRepoGetByRepError, "error", err, "representative_id", representativeID)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var branches []domain.Branch
	for rows.Next() {
		var branch domain.Branch
		var franchiseID, profileImageURL sql.NullString
		var locationID, cityID, address sql.NullString
		var departmentID sql.NullString
		var latitude, longitude sql.NullFloat64

		err := rows.Scan(
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
			log.Error(logger.LogBranchRepoScanError, "error", err)
			continue
		}

		if franchiseID.Valid {
			branch.FranchiseID = &franchiseID.String
		}
		if profileImageURL.Valid {
			branch.ProfileImageURL = &profileImageURL.String
		}

		if locationID.Valid {
			branch.Location = &domain.Location{
				ID:       locationID.String,
				BranchID: branch.ID,
				CityID:   cityID.String,
				Address:  address.String,
			}
			if latitude.Valid {
				branch.Location.Latitude = &latitude.Float64
			}
			if longitude.Valid {
				branch.Location.Longitude = &longitude.Float64
			}
			if departmentID.Valid {
				branch.Location.DepartmentID = departmentID.String
			}
		}

		brands, err := r.getBranchBrands(ctx, branch.ID)
		if err != nil {
			log.Error(logger.LogBranchRepoBrandGetError, "error", err, "branch_id", branch.ID)
		}
		branch.Brands = brands

		branches = append(branches, branch)
	}

	return branches, nil
}
