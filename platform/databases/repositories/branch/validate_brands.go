package branch

import (
	"context"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

// ValidateBrands checks if all provided brands exist in motorcycle_references table
func (r *repository) ValidateBrands(ctx context.Context, brands []string) error {
	if len(brands) == 0 {
		return nil
	}

	// Build dynamic query with placeholders
	placeholders := ""
	args := make([]interface{}, len(brands))
	for i, brand := range brands {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		args[i] = brand
	}

	query := fmt.Sprintf(queryValidateBrands, placeholders)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("error validating brands", "error", err)
		return err
	}
	defer rows.Close()

	// Collect found brands
	foundBrands := make(map[string]bool)
	for rows.Next() {
		var brand string
		if err := rows.Scan(&brand); err != nil {
			continue
		}
		foundBrands[brand] = true
	}

	// Check if all brands were found
	for _, brand := range brands {
		if !foundBrands[brand] {
			log.Warn("brand not found in brands catalog", "brand_id", brand)
			return domain.ErrBrandNotFound
		}
	}

	return nil
}
