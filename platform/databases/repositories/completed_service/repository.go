package completed_service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	queryInsert = `
		INSERT INTO completed_services (
			id, branch_id, motorcycle_id, diagnostic_id,
			request_date, status,
			quoted_price, final_price, representative_notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	queryInsertItem = `
		INSERT INTO completed_service_items (id, completed_service_id, service_id)
		VALUES (?, ?, ?)
	`

	queryInsertStatusHistory = `
		INSERT INTO service_status_transitions (id, completed_service_id, previous_status, new_status, created_by)
		VALUES (?, ?, ?, ?, ?)
	`

	queryGetByID = `
		SELECT id, branch_id, motorcycle_id, diagnostic_id,
			request_date, completion_date, status,
			quoted_price, final_price, representative_notes,
			created_at, updated_at
		FROM completed_services
		WHERE id = ?
	`

	queryGetByMotorcycleID = `
		SELECT cs.id, cs.branch_id, b.name AS branch_name,
			cs.motorcycle_id, cs.diagnostic_id,
			cs.request_date, cs.completion_date, cs.status,
			cs.quoted_price, cs.final_price, cs.representative_notes,
			cs.created_at, cs.updated_at
		FROM completed_services cs
		LEFT JOIN branches b ON b.id = cs.branch_id
		WHERE cs.motorcycle_id = ? AND cs.deleted_at IS NULL
		ORDER BY cs.request_date DESC
	`

	queryGetByBranchID = `
		SELECT cs.id, cs.branch_id, b.name AS branch_name,
			cs.motorcycle_id, cs.diagnostic_id,
			cs.request_date, cs.completion_date, cs.status,
			cs.quoted_price, cs.final_price, cs.representative_notes,
			cs.created_at, cs.updated_at
		FROM completed_services cs
		LEFT JOIN branches b ON b.id = cs.branch_id
		WHERE cs.branch_id = ? AND cs.deleted_at IS NULL
		ORDER BY cs.request_date DESC
	`

	queryGetItemsByCompletedServiceID = `
		SELECT csi.id, csi.completed_service_id, csi.service_id,
			s.name AS service_name,
			csi.rating, csi.comment, csi.rated_at, csi.is_offensive_comment
		FROM completed_service_items csi
		LEFT JOIN services s ON s.id = csi.service_id
		WHERE csi.completed_service_id = ?
	`

	queryValidateBranchServices = `
		SELECT COUNT(DISTINCT s.id)
		FROM services s
		INNER JOIN branch_services bs ON bs.service_id = s.id
		WHERE bs.branch_id = ? AND s.id IN (%s) AND s.is_active = 1
	`

	queryHasActiveService = `
		SELECT COUNT(*) FROM completed_services
		WHERE branch_id = ? AND motorcycle_id = ? AND status IN ('PENDIENTE', 'EN_PROCESO')
	`

	queryDelete = `
		DELETE FROM completed_services WHERE id = ?
	`

	querySoftDelete = `
		UPDATE completed_services SET deleted_at = NOW() WHERE id = ?
	`

	queryUpdateStatus = `
		UPDATE completed_services
		SET status = ?, completion_date = ?
		WHERE id = ?
	`

	queryUpdateStatusWithPrice = `
		UPDATE completed_services
		SET status = ?, completion_date = ?, final_price = ?
		WHERE id = ?
	`

	queryUpdateDetails = `
		UPDATE completed_services
		SET quoted_price = ?, final_price = ?, representative_notes = ?
		WHERE id = ?
	`

	queryGetStatusHistory = `
		SELECT id, completed_service_id, previous_status, new_status, created_by, created_at
		FROM service_status_transitions
		WHERE completed_service_id = ?
		ORDER BY created_at ASC
	`
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                      *sql.DB
	stmtInsert              *sql.Stmt
	stmtInsertItem          *sql.Stmt
	stmtInsertStatusHistory *sql.Stmt
	stmtGetByID             *sql.Stmt
	stmtGetByMotorcycleID   *sql.Stmt
	stmtGetByBranchID       *sql.Stmt
	stmtGetItemsByCSID      *sql.Stmt
	stmtHasActiveService    *sql.Stmt
	stmtDelete              *sql.Stmt
	stmtSoftDelete          *sql.Stmt
	stmtUpdateStatus        *sql.Stmt
	stmtGetStatusHistory    *sql.Stmt
}

// NewRepository creates a new completed service repository with prepared statements
func NewRepository(db *sql.DB) (output.CompletedServiceRepository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtInsert, err := db.Prepare(queryInsert)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareInsert, err)
		return nil, fmt.Errorf("error preparing stmtInsert: %w", err)
	}

	stmtInsertItem, err := db.Prepare(queryInsertItem)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareInsertItem, err)
		return nil, fmt.Errorf("error preparing stmtInsertItem: %w", err)
	}

	stmtInsertStatusHistory, err := db.Prepare(queryInsertStatusHistory)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareHistory, err)
		return nil, fmt.Errorf("error preparing stmtInsertStatusHistory: %w", err)
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareGetByID, err)
		return nil, fmt.Errorf("error preparing stmtGetByID: %w", err)
	}

	stmtGetByMotorcycleID, err := db.Prepare(queryGetByMotorcycleID)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareGetByMoto, err)
		return nil, fmt.Errorf("error preparing stmtGetByMotorcycleID: %w", err)
	}

	stmtGetByBranchID, err := db.Prepare(queryGetByBranchID)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareGetByBranch, err)
		return nil, fmt.Errorf("error preparing stmtGetByBranchID: %w", err)
	}

	stmtGetItemsByCSID, err := db.Prepare(queryGetItemsByCompletedServiceID)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareGetItems, err)
		return nil, fmt.Errorf("error preparing stmtGetItemsByCSID: %w", err)
	}

	stmtHasActiveService, err := db.Prepare(queryHasActiveService)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareHasActive, err)
		return nil, fmt.Errorf("error preparing stmtHasActiveService: %w", err)
	}

	stmtDelete, err := db.Prepare(queryDelete)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareDelete, err)
		return nil, fmt.Errorf("error preparing stmtDelete: %w", err)
	}

	stmtSoftDelete, err := db.Prepare(querySoftDelete)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareDelete, err)
		return nil, fmt.Errorf("error preparing stmtSoftDelete: %w", err)
	}

	stmtUpdateStatus, err := db.Prepare(queryUpdateStatus)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareUpdateStatus, err)
		return nil, fmt.Errorf("error preparing stmtUpdateStatus: %w", err)
	}

	stmtGetStatusHistory, err := db.Prepare(queryGetStatusHistory)
	if err != nil {
		log.Error(logger.LogCSRepoPrepareGetHistory, err)
		return nil, fmt.Errorf("error preparing stmtGetStatusHistory: %w", err)
	}

	return &repository{
		db:                      db,
		stmtInsert:              stmtInsert,
		stmtInsertItem:          stmtInsertItem,
		stmtInsertStatusHistory: stmtInsertStatusHistory,
		stmtGetByID:             stmtGetByID,
		stmtGetByMotorcycleID:   stmtGetByMotorcycleID,
		stmtGetByBranchID:       stmtGetByBranchID,
		stmtGetItemsByCSID:      stmtGetItemsByCSID,
		stmtHasActiveService:    stmtHasActiveService,
		stmtDelete:              stmtDelete,
		stmtSoftDelete:          stmtSoftDelete,
		stmtUpdateStatus:        stmtUpdateStatus,
		stmtGetStatusHistory:    stmtGetStatusHistory,
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
