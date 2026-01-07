package location

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

var log logger.Logger = logger.NewSlogLogger()

const (
	queryGetAllDepartments = `
		SELECT id, name FROM departments ORDER BY name
	`

	queryGetCitiesByDepartment = `
		SELECT id, name, department_id FROM cities WHERE department_id = ? ORDER BY name
	`

	queryValidateCityInDepartment = `
		SELECT 1 FROM cities WHERE id = ? AND department_id = ?
	`

	queryGetDepartmentByID = `
		SELECT id, name FROM departments WHERE id = ?
	`

	querySaveLocation = `
		INSERT INTO locations (id, branch_id, city_id, address, latitude, longitude, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	queryUpdateLocation = `
		UPDATE locations 
		SET city_id = ?, address = ?, latitude = ?, longitude = ?, updated_at = NOW()
		WHERE branch_id = ? `
)

type repository struct {
	db                           *sql.DB
	stmtGetAllDepartments        *sql.Stmt
	stmtGetCitiesByDepartment    *sql.Stmt
	stmtValidateCityInDepartment *sql.Stmt
	stmtGetDepartmentByID        *sql.Stmt
	stmtSaveLocation             *sql.Stmt
	stmtUpdateLocation           *sql.Stmt
}

// NewRepository creates a new LocationRepository implementation
func NewRepository(db *sql.DB) (output.LocationRepository, error) {
	stmtGetAllDepartments, err := db.Prepare(queryGetAllDepartments)
	if err != nil {
		log.Error(logger.LogLocationRepoPrepareError, "statement", "GetAllDepartments", "error", err)
		return nil, fmt.Errorf("error preparing GetAllDepartments: %w", err)
	}

	stmtGetCitiesByDepartment, err := db.Prepare(queryGetCitiesByDepartment)
	if err != nil {
		log.Error(logger.LogLocationRepoPrepareError, "statement", "GetCitiesByDepartment", "error", err)
		return nil, fmt.Errorf("error preparing GetCitiesByDepartment: %w", err)
	}

	stmtValidateCityInDepartment, err := db.Prepare(queryValidateCityInDepartment)
	if err != nil {
		log.Error(logger.LogLocationRepoPrepareError, "statement", "ValidateCityInDepartment", "error", err)
		return nil, fmt.Errorf("error preparing ValidateCityInDepartment: %w", err)
	}

	stmtGetDepartmentByID, err := db.Prepare(queryGetDepartmentByID)
	if err != nil {
		log.Error(logger.LogLocationRepoPrepareError, "statement", "GetDepartmentByID", "error", err)
		return nil, fmt.Errorf("error preparing GetDepartmentByID: %w", err)
	}

	stmtSaveLocation, err := db.Prepare(querySaveLocation)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtSaveLocation", err)
		return nil, fmt.Errorf("error preparing stmtSaveLocation: %w", err)
	}

	stmtUpdateLocation, err := db.Prepare(queryUpdateLocation)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtUpdateLocation", err)
		return nil, fmt.Errorf("error preparing stmtUpdateLocation: %w", err)
	}

	return &repository{
		db:                           db,
		stmtGetAllDepartments:        stmtGetAllDepartments,
		stmtGetCitiesByDepartment:    stmtGetCitiesByDepartment,
		stmtValidateCityInDepartment: stmtValidateCityInDepartment,
		stmtGetDepartmentByID:        stmtGetDepartmentByID,
		stmtSaveLocation:             stmtSaveLocation,
		stmtUpdateLocation:           stmtUpdateLocation,
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

// GetAllDepartments retrieves all departments ordered by name
func (r *repository) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	rows, err := r.stmtGetAllDepartments.QueryContext(ctx)
	if err != nil {
		log.Error(logger.LogLocationRepoGetDepartmentsError, "error", err)
		return nil, err
	}
	defer rows.Close()

	var departments []domain.Department
	for rows.Next() {
		var dept domain.Department
		if err := rows.Scan(&dept.ID, &dept.Name); err != nil {
			log.Error(logger.LogLocationRepoGetDepartmentsScanError, "error", err)
			continue
		}
		departments = append(departments, dept)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogLocationRepoGetDepartmentsIterError, "error", err)
		return nil, err
	}

	return departments, nil
}

// GetCitiesByDepartment retrieves all cities for a specific department
func (r *repository) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	rows, err := r.stmtGetCitiesByDepartment.QueryContext(ctx, departmentID)
	if err != nil {
		log.Error(logger.LogLocationRepoGetCitiesError, "error", err, "department_id", departmentID)
		return nil, err
	}
	defer rows.Close()

	var cities []domain.City
	for rows.Next() {
		var city domain.City
		if err := rows.Scan(&city.ID, &city.Name, &city.DepartmentID); err != nil {
			log.Error(logger.LogLocationRepoGetCitiesScanError, "error", err)
			continue
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		log.Error(logger.LogLocationRepoGetCitiesIterError, "error", err)
		return nil, err
	}

	return cities, nil
}

// ValidateCityInDepartment checks if the city belongs to the specified department
func (r *repository) ValidateCityInDepartment(ctx context.Context, cityID, departmentID string) error {
	var exists int
	err := r.stmtValidateCityInDepartment.QueryRowContext(ctx, cityID, departmentID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrCityNotInDepartment
		}
		log.Error(logger.LogLocationRepoValidateCityError, "error", err, "city_id", cityID, "department_id", departmentID)
		return err
	}
	return nil
}

// GetDepartmentByID retrieves a department by its ID
func (r *repository) GetDepartmentByID(ctx context.Context, departmentID string) (*domain.Department, error) {
	var dept domain.Department
	err := r.stmtGetDepartmentByID.QueryRowContext(ctx, departmentID).Scan(&dept.ID, &dept.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDepartmentNotFound
		}
		log.Error(logger.LogLocationRepoGetDeptByIDError, "error", err, "department_id", departmentID)
		return nil, err
	}
	return &dept, nil
}
