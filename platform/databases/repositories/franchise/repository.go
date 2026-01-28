package franchise

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
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

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return common.NewSQLTx(tx), nil
}








