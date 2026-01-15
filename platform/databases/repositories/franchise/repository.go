package franchise

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	// Franchise queries
	querySaveFranchise = `
		INSERT INTO franchises (id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`
	queryUpdateFranchise = `
		UPDATE franchises 
		SET name = ?, description = ?, updated_at = NOW()
		WHERE id = ?
	`
	queryDeleteFranchise = "DELETE FROM franchises WHERE id = ?"

	queryGetFranchiseByID = `
		SELECT id, name, description
		FROM franchises
		WHERE id = ?
	`
	queryGetFranchiseByName = `
		SELECT id, name, description
		FROM franchises
		WHERE name = ?
		LIMIT 1
	`
	queryGetFranchisesByRepresentative = `
		SELECT DISTINCT f.id, f.name, f.description
		FROM franchises f
		INNER JOIN branches b ON f.id = b.franchise_id
		WHERE b.representative_id = ?
	`
	queryCountBranchesByFranchise = `
		SELECT COUNT(*) FROM branches WHERE franchise_id = ?
	`

	// Branch association queries
	queryAssociateBranchToFranchise = `
		UPDATE branches SET franchise_id = ?, updated_at = NOW() WHERE id = ?
	`
	queryDissociateBranchesFromFranchise = `
		UPDATE branches SET franchise_id = NULL, updated_at = NOW() WHERE franchise_id = ?
	`
	queryDissociateSingleBranch = `
		UPDATE branches SET franchise_id = NULL, updated_at = NOW() WHERE id = ?
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                                  *sql.DB
	stmtSaveFranchise                   *sql.Stmt
	stmtUpdateFranchise                 *sql.Stmt
	stmtDeleteFranchise                 *sql.Stmt
	stmtGetFranchiseByID                *sql.Stmt
	stmtGetFranchiseByName              *sql.Stmt
	stmtGetFranchisesByRepresentative   *sql.Stmt
	stmtCountBranchesByFranchise        *sql.Stmt
	stmtAssociateBranchToFranchise      *sql.Stmt
	stmtDissociateBranchesFromFranchise *sql.Stmt
	stmtDissociateSingleBranch          *sql.Stmt
}

// NewRepository creates a new FranchiseRepository with prepared statements (fail-fast pattern)
func NewRepository(db *sql.DB) (output.FranchiseRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtSaveFranchise, err := db.Prepare(querySaveFranchise)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "SaveFranchise", "error", err)
		return nil, fmt.Errorf("error preparing stmtSaveFranchise: %w", err)
	}

	stmtUpdateFranchise, err := db.Prepare(queryUpdateFranchise)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "UpdateFranchise", "error", err)
		return nil, fmt.Errorf("error preparing stmtUpdateFranchise: %w", err)
	}

	stmtDeleteFranchise, err := db.Prepare(queryDeleteFranchise)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "DeleteFranchise", "error", err)
		return nil, fmt.Errorf("error preparing stmtDeleteFranchise: %w", err)
	}

	stmtGetFranchiseByID, err := db.Prepare(queryGetFranchiseByID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "GetFranchiseByID", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetFranchiseByID: %w", err)
	}

	stmtGetFranchiseByName, err := db.Prepare(queryGetFranchiseByName)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "GetFranchiseByName", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetFranchiseByName: %w", err)
	}

	stmtGetFranchisesByRepresentative, err := db.Prepare(queryGetFranchisesByRepresentative)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "GetFranchisesByRepresentative", "error", err)
		return nil, fmt.Errorf("error preparing stmtGetFranchisesByRepresentative: %w", err)
	}

	stmtCountBranchesByFranchise, err := db.Prepare(queryCountBranchesByFranchise)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "CountBranchesByFranchise", "error", err)
		return nil, fmt.Errorf("error preparing stmtCountBranchesByFranchise: %w", err)
	}

	stmtAssociateBranchToFranchise, err := db.Prepare(queryAssociateBranchToFranchise)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "AssociateBranchToFranchise", "error", err)
		return nil, fmt.Errorf("error preparing stmtAssociateBranchToFranchise: %w", err)
	}

	stmtDissociateBranchesFromFranchise, err := db.Prepare(queryDissociateBranchesFromFranchise)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "DissociateBranchesFromFranchise", "error", err)
		return nil, fmt.Errorf("error preparing stmtDissociateBranchesFromFranchise: %w", err)
	}

	stmtDissociateSingleBranch, err := db.Prepare(queryDissociateSingleBranch)
	if err != nil {
		log.Error(logger.LogFranchiseRepoPrepareError, "statement", "DissociateSingleBranch", "error", err)
		return nil, fmt.Errorf("error preparing stmtDissociateSingleBranch: %w", err)
	}

	return &repository{
		db:                                  db,
		stmtSaveFranchise:                   stmtSaveFranchise,
		stmtUpdateFranchise:                 stmtUpdateFranchise,
		stmtDeleteFranchise:                 stmtDeleteFranchise,
		stmtGetFranchiseByID:                stmtGetFranchiseByID,
		stmtGetFranchiseByName:              stmtGetFranchiseByName,
		stmtGetFranchisesByRepresentative:   stmtGetFranchisesByRepresentative,
		stmtCountBranchesByFranchise:        stmtCountBranchesByFranchise,
		stmtAssociateBranchToFranchise:      stmtAssociateBranchToFranchise,
		stmtDissociateBranchesFromFranchise: stmtDissociateBranchesFromFranchise,
		stmtDissociateSingleBranch:          stmtDissociateSingleBranch,
	}, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}

// SaveFranchise inserts a new franchise
func (r *repository) SaveFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtSaveFranchise)

	_, err := stmt.ExecContext(ctx, franchise.ID, franchise.Name, franchise.Description)
	if err != nil {
		log.Error(logger.LogFranchiseRepoSaveError, "error", err, franchise.ToLogger())
		return err
	}
	return nil
}

// UpdateFranchise updates an existing franchise
func (r *repository) UpdateFranchise(ctx context.Context, tx output.Tx, franchise domain.Franchise) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtUpdateFranchise)

	_, err := stmt.ExecContext(ctx, franchise.Name, franchise.Description, franchise.ID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoUpdateError, "error", err, franchise.ToLogger())
		return err
	}
	return nil
}

// DeleteFranchise removes a franchise by ID
func (r *repository) DeleteFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtDeleteFranchise)

	_, err := stmt.ExecContext(ctx, franchiseID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDeleteError, "franchise_id", franchiseID, "error", err)
		return err
	}
	return nil
}

// GetFranchiseByID retrieves a franchise by its ID
func (r *repository) GetFranchiseByID(ctx context.Context, franchiseID string) (*domain.Franchise, error) {
	var franchise domain.Franchise
	var description sql.NullString

	err := r.stmtGetFranchiseByID.QueryRowContext(ctx, franchiseID).Scan(
		&franchise.ID,
		&franchise.Name,
		&description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrFranchiseNotFound
		}
		log.Error(logger.LogFranchiseRepoGetByIDError, "franchise_id", franchiseID, "error", err)
		return nil, err
	}

	if description.Valid {
		franchise.Description = &description.String
	}

	return &franchise, nil
}

// GetFranchiseByName retrieves a franchise by its name (for duplicate validation)
func (r *repository) GetFranchiseByName(ctx context.Context, name string) (*domain.Franchise, error) {
	var franchise domain.Franchise
	var description sql.NullString

	err := r.stmtGetFranchiseByName.QueryRowContext(ctx, name).Scan(
		&franchise.ID,
		&franchise.Name,
		&description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found is not an error for duplicate check
		}
		log.Error(logger.LogFranchiseRepoGetByNameError, "name", name, "error", err)
		return nil, err
	}

	if description.Valid {
		franchise.Description = &description.String
	}

	return &franchise, nil
}

// GetFranchisesByRepresentative lists all franchises owned by a representative
func (r *repository) GetFranchisesByRepresentative(ctx context.Context, representativeID string) ([]domain.Franchise, error) {
	rows, err := r.stmtGetFranchisesByRepresentative.QueryContext(ctx, representativeID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoGetByRepError, "representative_id", representativeID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var franchises []domain.Franchise
	for rows.Next() {
		var franchise domain.Franchise
		var description sql.NullString

		if err := rows.Scan(&franchise.ID, &franchise.Name, &description); err != nil {
			log.Error(logger.LogFranchiseRepoScanError, "error", err)
			return nil, err
		}

		if description.Valid {
			franchise.Description = &description.String
		}

		franchises = append(franchises, franchise)
	}

	return franchises, nil
}

// CountBranchesByFranchise returns the number of branches associated with a franchise
func (r *repository) CountBranchesByFranchise(ctx context.Context, franchiseID string) (int, error) {
	var count int
	err := r.stmtCountBranchesByFranchise.QueryRowContext(ctx, franchiseID).Scan(&count)
	if err != nil {
		log.Error(logger.LogFranchiseRepoCountBranches, "franchise_id", franchiseID, "error", err)
		return 0, err
	}
	return count, nil
}

// AssociateBranchesToFranchise updates branches to link them to a franchise
func (r *repository) AssociateBranchesToFranchise(ctx context.Context, tx output.Tx, franchiseID string, branchIDs []string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtAssociateBranchToFranchise)

	for _, branchID := range branchIDs {
		_, err := stmt.ExecContext(ctx, franchiseID, branchID)
		if err != nil {
			log.Error(logger.LogFranchiseRepoAssociateError,
				"franchise_id", franchiseID, "branch_id", branchID, "error", err)
			return err
		}
	}
	return nil
}

// DissociateBranchesFromFranchise removes franchise association from all branches
func (r *repository) DissociateBranchesFromFranchise(ctx context.Context, tx output.Tx, franchiseID string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtDissociateBranchesFromFranchise)

	_, err := stmt.ExecContext(ctx, franchiseID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDissociateError, "franchise_id", franchiseID, "error", err)
		return err
	}
	return nil
}

// DissociateSingleBranch removes franchise association from a single branch
func (r *repository) DissociateSingleBranch(ctx context.Context, tx output.Tx, branchID string) error {
	sqlTx := tx.(*common.SQLTx)
	stmt := sqlTx.StmtContext(ctx, r.stmtDissociateSingleBranch)

	_, err := stmt.ExecContext(ctx, branchID)
	if err != nil {
		log.Error(logger.LogFranchiseRepoDissociateError, "branch_id", branchID, "error", err)
		return err
	}
	return nil
}
