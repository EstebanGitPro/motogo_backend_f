package service

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryGetAllServices      = "SELECT id, name, description, service_type, is_active FROM services ORDER BY name"
	queryGetServicesByType   = "SELECT id, name, description, service_type, is_active FROM services WHERE service_type = ? ORDER BY name"
	queryGetServiceByID      = "SELECT id, name, description, service_type, is_active FROM services WHERE id = ?"
	queryUpdateService       = "UPDATE services SET name = ?, description = ?, service_type = ?, is_active = ? WHERE id = ?"
	queryGetServicesByBranch = `
		SELECT s.id, s.name, s.description, s.service_type, bs.created_at, bs.active
		FROM branch_services bs
		JOIN services s ON s.id = bs.service_id
		WHERE bs.branch_id = ?
		ORDER BY bs.created_at DESC
	`
	queryInsertBranchService     = "INSERT INTO branch_services (id, branch_id, service_id, active) VALUES (?, ?, ?, TRUE)"
	queryDeleteBranchService     = "DELETE FROM branch_services WHERE branch_id = ? AND service_id = ?"
	queryCheckServiceAssociation = "SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                          *sql.DB
	stmtGetAllServices          *sql.Stmt
	stmtGetServicesByType       *sql.Stmt
	stmtGetServiceByID          *sql.Stmt
	stmtGetServicesByBranch     *sql.Stmt
	stmtUpdateService           *sql.Stmt
	stmtInsertBranchService     *sql.Stmt
	stmtDeleteBranchService     *sql.Stmt
	stmtCheckServiceAssociation *sql.Stmt
}

func NewRepository(db *sql.DB) (output.ServiceRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtGetAllServices, err := db.Prepare(queryGetAllServices)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtGetAllServices", err)
		return nil, err
	}

	stmtGetServicesByType, err := db.Prepare(queryGetServicesByType)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtGetServicesByType", err)
		return nil, err
	}

	stmtGetServiceByID, err := db.Prepare(queryGetServiceByID)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtGetServiceByID", err)
		return nil, err
	}

	stmtGetServicesByBranch, err := db.Prepare(queryGetServicesByBranch)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtGetServicesByBranch", err)
		return nil, err
	}

	stmtUpdateService, err := db.Prepare(queryUpdateService)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtUpdateService", err)
		return nil, err
	}

	stmtInsertBranchService, err := db.Prepare(queryInsertBranchService)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtInsertBranchService", err)
		return nil, err
	}

	stmtDeleteBranchService, err := db.Prepare(queryDeleteBranchService)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtDeleteBranchService", err)
		return nil, err
	}

	stmtCheckServiceAssociation, err := db.Prepare(queryCheckServiceAssociation)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error preparing stmtCheckServiceAssociation", err)
		return nil, err
	}

	return &repository{
		db:                          db,
		stmtGetAllServices:          stmtGetAllServices,
		stmtGetServicesByType:       stmtGetServicesByType,
		stmtGetServiceByID:          stmtGetServiceByID,
		stmtGetServicesByBranch:     stmtGetServicesByBranch,
		stmtUpdateService:           stmtUpdateService,
		stmtInsertBranchService:     stmtInsertBranchService,
		stmtDeleteBranchService:     stmtDeleteBranchService,
		stmtCheckServiceAssociation: stmtCheckServiceAssociation,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return common.BeginSQLTx(ctx, r.db)
}
