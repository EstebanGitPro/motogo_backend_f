package branch

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetBranchesNearby retrieves branches within radius of given coordinates (HU89)
// Uses Haversine formula for accurate distance calculation with bounding box optimization
func (r *repository) GetBranchesNearby(
	ctx context.Context,
	lat, lng, radiusKm float64,
	establishmentType string,
	latMin, latMax, lngMin, lngMax float64,
) ([]domain.NearbyBranch, error) {
	// Query parameters order:
	// 1. lat (for COS(RADIANS(?)))
	// 2. lng (for COS(RADIANS(l.longitude) - RADIANS(?)))
	// 3. lat (for SIN(RADIANS(?)) * SIN(RADIANS(l.latitude)))
	// 4. establishmentType (for type = '' check)
	// 5. establishmentType (for comparison)
	// 6. latMin (for latitude BETWEEN)
	// 7. latMax (for latitude BETWEEN)
	// 8. lngMin (for longitude BETWEEN)
	// 9. lngMax (for longitude BETWEEN)
	// 10. radiusKm (for HAVING distance_km <= ?)
	rows, err := r.stmtGetBranchesNearby.QueryContext(ctx,
		lat, lng, lat,
		establishmentType, establishmentType,
		latMin, latMax, lngMin, lngMax,
		radiusKm,
	)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error querying nearby branches", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	var branches []domain.NearbyBranch
	for rows.Next() {
		var branch domain.NearbyBranch
		var location domain.NearbyLocation
		var profileImageURL *string
		var status string

		err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.EstablishmentType,
			&profileImageURL,
			&status,
			&location.Address,
			&location.Latitude,
			&location.Longitude,
			&location.CityName,
			&location.DepartmentName,
			&branch.ContactPhone,
			&branch.DistanceKm,
		)
		if err != nil {
			log.Error(logger.LogDatabaseUnavailable, "error scanning nearby branch row", err)
			return nil, err
		}

		branch.ProfileImageURL = profileImageURL
		branch.Location = &location
		branches = append(branches, branch)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error iterating nearby branches", err)
		return nil, err
	}

	return branches, nil
}
