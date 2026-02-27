package branch

import (
	"context"
	"fmt"
	"strings"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

// GetBranchesNearby retrieves branches within radius of given coordinates (HU89)
// Uses Haversine formula with bounding box optimization and optional brand/displacement filters
func (r *repository) GetBranchesNearby(
	ctx context.Context,
	params domain.NearbySearchParams,
) ([]domain.NearbyBranch, error) {
	// Build dynamic query with optional filters
	var query strings.Builder
	args := make([]interface{}, 0, 12)

	query.WriteString(`
		SELECT 
			b.id, b.name, b.establishment_type, b.profile_image_url, b.status,
			l.address, l.latitude, l.longitude,
			ci.name AS city_name, d.name AS department_name,
			p.phone_number,
			(6371 * ACOS(
				COS(RADIANS(?)) * COS(RADIANS(l.latitude)) *
				COS(RADIANS(l.longitude) - RADIANS(?)) +
				SIN(RADIANS(?)) * SIN(RADIANS(l.latitude))
			)) AS distance_km
		FROM branches b
		INNER JOIN locations l ON b.id = l.branch_id
		INNER JOIN cities ci ON l.city_id = ci.id
		INNER JOIN departments d ON ci.department_id = d.id
		LEFT JOIN persons p ON b.representative_id = p.id
		WHERE b.status = 'ACTIVE'
		  AND l.latitude IS NOT NULL
		  AND l.longitude IS NOT NULL
	`)
	args = append(args, params.Lat, params.Lng, params.Lat)

	// Optional: establishment type filter
	if params.EstablishmentType != "" {
		query.WriteString(fmt.Sprintf("  AND b.establishment_type LIKE CONCAT('%%', ?, '%%')\n"))
		args = append(args, params.EstablishmentType)
	}

	// Bounding box filter (always applied)
	query.WriteString("  AND l.latitude BETWEEN ? AND ?\n")
	query.WriteString("  AND l.longitude BETWEEN ? AND ?\n")
	args = append(args, params.LatMin, params.LatMax, params.LngMin, params.LngMax)
	// Optional filters: brand and/or displacement range
	// When both are provided, use OR (match either); when only one, use AND
	hasBrand := params.BrandID != ""
	hasDisplacement := params.DisplacementRange != ""

	switch {
	case hasBrand && hasDisplacement:
		query.WriteString("  AND (EXISTS (SELECT 1 FROM branch_brands bb WHERE bb.branch_id = b.id AND bb.brand_id = ?)")
		query.WriteString("    OR EXISTS (SELECT 1 FROM branch_displacement_ranges bdr WHERE bdr.branch_id = b.id AND bdr.displacement_range = ? AND bdr.active = TRUE))\n")
		args = append(args, params.BrandID, params.DisplacementRange)
	case hasBrand:
		query.WriteString("  AND EXISTS (SELECT 1 FROM branch_brands bb WHERE bb.branch_id = b.id AND bb.brand_id = ?)\n")
		args = append(args, params.BrandID)
	case hasDisplacement:
		query.WriteString("  AND EXISTS (SELECT 1 FROM branch_displacement_ranges bdr WHERE bdr.branch_id = b.id AND bdr.displacement_range = ? AND bdr.active = TRUE)\n")
		args = append(args, params.DisplacementRange)
	}

	// HAVING, ORDER, LIMIT
	query.WriteString("HAVING distance_km <= ?\n")
	query.WriteString("ORDER BY distance_km ASC\n")
	query.WriteString("LIMIT 50\n")
	args = append(args, params.RadiusKm)

	// Execute dynamic query
	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error querying nearby branches", err)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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

	// Hydrate each branch with brands and displacement ranges
	for i := range branches {
		brands, err := r.getBranchBrands(ctx, branches[i].ID)
		if err != nil {
			log.Warn("Error obteniendo marcas de sede cercana", "branch_id", branches[i].ID, "error", err)
		}
		branches[i].Brands = brands

		ranges, err := r.getBranchDisplacementRanges(ctx, branches[i].ID)
		if err != nil {
			log.Warn("Error obteniendo rangos de cilindraje de sede cercana", "branch_id", branches[i].ID, "error", err)
		}
		branches[i].DisplacementRanges = ranges

		serviceNames, err := r.getBranchServiceNames(ctx, branches[i].ID)
		if err != nil {
			log.Warn("Error obteniendo servicios de sede cercana", "branch_id", branches[i].ID, "error", err)
		}
		branches[i].ServiceNames = serviceNames
	}

	return branches, nil
}
