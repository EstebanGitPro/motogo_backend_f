package brand

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryGetAllBrands = "SELECT id, name FROM brands ORDER BY name"

	queryValidateBrandIDs = "SELECT id FROM brands WHERE id IN (%s)"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db               *sql.DB
	stmtGetAllBrands *sql.Stmt
}

func NewRepository(db *sql.DB) (output.BrandRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetAllBrands, err := db.Prepare(queryGetAllBrands)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetAllBrands", err)
		return nil, fmt.Errorf("error preparing stmtGetAllBrands: %w", err)
	}

	return &repository{
		db:               db,
		stmtGetAllBrands: stmtGetAllBrands,
	}, nil
}

func (r *repository) ValidateBrandIDs(ctx context.Context, brandIDs []string) error {
	if len(brandIDs) == 0 {
		return nil
	}
	placeholders := ""
	args := make([]interface{}, len(brandIDs))
	for i, id := range brandIDs {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		args[i] = id
	}

	query := fmt.Sprintf(queryValidateBrandIDs, placeholders)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error(logger.LogBranchRepoBrandValidateErr, "error", err)
		return err
	}
	defer func() { _ = rows.Close() }() // Rows close error intentionally ignored

	foundBrands := make(map[string]bool)
	for rows.Next() {
		var id string
		if rows.Scan(&id) != nil {
			continue
		}
		foundBrands[id] = true
	}

	for _, id := range brandIDs {
		if !foundBrands[id] {
			log.Warn(logger.LogBranchRepoBrandValidateErr, "brand_id", id)
			return domain.ErrBrandNotFound
		}
	}

	return nil
}
