package service

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
	uuid "github.com/EstebanGitPro/motogo-backend/tools/utils"
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
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	db                      *sql.DB
	stmtGetAllServices      *sql.Stmt
	stmtGetServicesByType   *sql.Stmt
	stmtGetServiceByID      *sql.Stmt
	stmtGetServicesByBranch *sql.Stmt
}

// NewRepository creates a new ServiceRepository with prepared statements (fail-fast pattern)
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

	return &repository{
		db:                      db,
		stmtGetAllServices:      stmtGetAllServices,
		stmtGetServicesByType:   stmtGetServicesByType,
		stmtGetServiceByID:      stmtGetServiceByID,
		stmtGetServicesByBranch: stmtGetServicesByBranch,
	}, nil
}

// GetAllServices retrieves all services from the catalog ordered by name
func (r *repository) GetAllServices(ctx context.Context) ([]domain.Service, error) {
	log.Info(logger.LogServiceRepoGetAll)

	rows, err := r.stmtGetAllServices.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogServiceRepoGetAllError, "error", err)
		return nil, err
	}
	defer rows.Close()

	var services []domain.Service
	for rows.Next() {
		var svc domain.Service
		var description sql.NullString
		if err := rows.Scan(&svc.ID, &svc.Name, &description, &svc.ServiceType, &svc.IsActive); err != nil {
			log.Error(logger.LogServiceRepoScanError, "error scanning service", err)
			continue
		}
		if description.Valid {
			svc.Description = description.String
		}
		services = append(services, svc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GetServicesByType retrieves services filtered by type
func (r *repository) GetServicesByType(ctx context.Context, serviceType string) ([]domain.Service, error) {
	log.Info(logger.LogServiceRepoGetByType, "type", serviceType)

	rows, err := r.stmtGetServicesByType.QueryContext(ctx, serviceType)
	if err != nil {
		log.Error(logger.LogServiceRepoGetByTypeError, "error", err)
		return nil, err
	}
	defer rows.Close()

	var services []domain.Service
	for rows.Next() {
		var svc domain.Service
		var description sql.NullString
		if err := rows.Scan(&svc.ID, &svc.Name, &description, &svc.ServiceType, &svc.IsActive); err != nil {
			log.Error(logger.LogServiceRepoScanError, "error scanning service", err)
			continue
		}
		if description.Valid {
			svc.Description = description.String
		}
		services = append(services, svc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// GetServicesByBranch retrieves services associated with a specific branch
// Returns the services with their added_at timestamp from branch_services table
func (r *repository) GetServicesByBranch(ctx context.Context, branchID string) ([]domain.BranchServiceInfo, error) {
	log.Info(logger.LogBranchServicesRepoGetByBranch, "branch_id", branchID)

	rows, err := r.stmtGetServicesByBranch.QueryContext(ctx, branchID)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoGetByBranchErr, "error", err)
		return nil, err
	}
	defer rows.Close()

	var services []domain.BranchServiceInfo
	for rows.Next() {
		var info domain.BranchServiceInfo
		var description sql.NullString
		var createdAt sql.NullTime
		if err := rows.Scan(
			&info.Service.ID,
			&info.Service.Name,
			&description,
			&info.Service.ServiceType,
			&createdAt,
			&info.Active,
		); err != nil {
			log.Error(logger.LogServiceRepoScanError, "error scanning branch service", err)
			continue
		}
		if description.Valid {
			info.Service.Description = description.String
		}
		if createdAt.Valid {
			info.AddedAt = createdAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		services = append(services, info)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

// BeginTx starts a new database transaction
func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// AssociateBranchServices associates multiple services to a branch
func (r *repository) AssociateBranchServices(ctx context.Context, tx output.Tx, branchID string, serviceIDs []string) error {
	log.Info(logger.LogBranchServicesRepoAssociate, "branch_id", branchID, "service_count", len(serviceIDs))

	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrInternalServer
	}

	query := "INSERT INTO branch_services (id, branch_id, service_id, active) VALUES (?, ?, ?, TRUE)"
	stmt, err := sqlTx.PrepareContext(ctx, query)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoPrepareErr, "error", err)
		return err
	}
	defer stmt.Close()

	for _, serviceID := range serviceIDs {
		id := generateUUID()
		_, err := stmt.ExecContext(ctx, id, branchID, serviceID)
		if err != nil {
			log.Error(logger.LogBranchServicesRepoAssociateErr, "branch_id", branchID, "service_id", serviceID, "error", err)
			return err
		}
	}

	log.Success(logger.LogBranchServicesRepoAssociateOK, "branch_id", branchID, "service_count", len(serviceIDs))
	return nil
}

// DissociateBranchService removes a service from a branch
func (r *repository) DissociateBranchService(ctx context.Context, tx output.Tx, branchID, serviceID string) error {
	log.Info(logger.LogBranchServicesRepoDissociate, "branch_id", branchID, "service_id", serviceID)

	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return domain.ErrInternalServer
	}

	query := "DELETE FROM branch_services WHERE branch_id = ? AND service_id = ?"
	result, err := sqlTx.ExecContext(ctx, query, branchID, serviceID)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoDissociateErr, "error", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Warn(logger.LogBranchServicesRepoNotFound, "branch_id", branchID, "service_id", serviceID)
		return domain.ErrServiceNotFound
	}

	log.Success(logger.LogBranchServicesRepoDissociateOK, "branch_id", branchID, "service_id", serviceID)
	return nil
}

// ValidateServiceIDs checks if all provided service IDs exist in the services table
func (r *repository) ValidateServiceIDs(ctx context.Context, serviceIDs []string) error {
	if len(serviceIDs) == 0 {
		return nil
	}

	log.Info(logger.LogBranchServicesRepoValidateIDs, "count", len(serviceIDs))

	// Build placeholders for IN clause
	placeholders := ""
	args := make([]interface{}, len(serviceIDs))
	for i, id := range serviceIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	query := "SELECT COUNT(*) FROM services WHERE id IN (" + placeholders + ")"
	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoValidateErr, "error", err)
		return err
	}

	if count != len(serviceIDs) {
		log.Warn(logger.LogBranchServicesRepoValidateMiss, "expected", len(serviceIDs), "found", count)
		return domain.ErrServiceNotFound
	}

	return nil
}

// CheckServiceAssociation checks if a service is already associated with a branch
func (r *repository) CheckServiceAssociation(ctx context.Context, branchID, serviceID string) (bool, error) {
	log.Info(logger.LogBranchServicesRepoCheckAssoc, "branch_id", branchID, "service_id", serviceID)

	query := "SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?"
	var count int
	err := r.db.QueryRowContext(ctx, query, branchID, serviceID).Scan(&count)
	if err != nil {
		log.Error(logger.LogBranchServicesRepoCheckAssocErr, "error", err)
		return false, err
	}

	return count > 0, nil
}

// GetServiceByID retrieves a service by its UUID (HU68)
func (r *repository) GetServiceByID(ctx context.Context, serviceID string) (*domain.Service, error) {
	log.Info(logger.LogServiceRepoGetByID, "service_id", serviceID)

	var svc domain.Service
	var description sql.NullString
	err := r.stmtGetServiceByID.QueryRowContext(ctx, serviceID).Scan(
		&svc.ID,
		&svc.Name,
		&description,
		&svc.ServiceType,
		&svc.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn(logger.LogServiceRepoNotFound, "service_id", serviceID)
			return nil, domain.ErrServiceNotFound
		}
		log.Error(logger.LogServiceRepoGetByIDError, "error", err, "service_id", serviceID)
		return nil, err
	}

	if description.Valid {
		svc.Description = description.String
	}

	log.Success(logger.LogServiceRepoGetByIDOK, "service_id", serviceID)
	return &svc, nil
}

// UpdateService updates an existing service in the catalog (HU68 - Admin only)
func (r *repository) UpdateService(ctx context.Context, tx output.Tx, service domain.Service) error {
	log.Info(logger.LogServiceRepoUpdate, "service_id", service.ID)

	// Use direct DB connection for now (transactional support can be added later if needed)
	result, err := r.db.ExecContext(ctx, queryUpdateService,
		service.Name,
		service.Description,
		string(service.ServiceType),
		service.IsActive,
		service.ID,
	)
	if err != nil {
		log.Error(logger.LogServiceRepoUpdateError, "error", err, "service_id", service.ID)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	// Note: RowsAffected can be 0 if the values are the same as current.
	// The service existence is already verified in the controller before calling this.
	// So we don't fail on RowsAffected=0 - it just means no changes were needed.
	log.Success(logger.LogServiceRepoUpdateOK, "service_id", service.ID, "rows_affected", rowsAffected)
	return nil
}

// generateUUID generates a new UUID
func generateUUID() string {
	return uuid.Generate()
}
