package location

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
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
)

type repository struct {
	db                           *sql.DB
	stmtGetAllDepartments        *sql.Stmt
	stmtGetCitiesByDepartment    *sql.Stmt
	stmtValidateCityInDepartment *sql.Stmt
	stmtGetDepartmentByID        *sql.Stmt
}

// NewRepository creates a new LocationRepository implementation
func NewRepository(db *sql.DB) (output.LocationRepository, error) {
	stmtGetAllDepartments, err := db.Prepare(queryGetAllDepartments)
	if err != nil {
		log.Error("error preparing GetAllDepartments statement", "error", err)
		return nil, fmt.Errorf("error preparing GetAllDepartments: %w", err)
	}

	stmtGetCitiesByDepartment, err := db.Prepare(queryGetCitiesByDepartment)
	if err != nil {
		log.Error("error preparing GetCitiesByDepartment statement", "error", err)
		return nil, fmt.Errorf("error preparing GetCitiesByDepartment: %w", err)
	}

	stmtValidateCityInDepartment, err := db.Prepare(queryValidateCityInDepartment)
	if err != nil {
		log.Error("error preparing ValidateCityInDepartment statement", "error", err)
		return nil, fmt.Errorf("error preparing ValidateCityInDepartment: %w", err)
	}

	stmtGetDepartmentByID, err := db.Prepare(queryGetDepartmentByID)
	if err != nil {
		log.Error("error preparing GetDepartmentByID statement", "error", err)
		return nil, fmt.Errorf("error preparing GetDepartmentByID: %w", err)
	}

	return &repository{
		db:                           db,
		stmtGetAllDepartments:        stmtGetAllDepartments,
		stmtGetCitiesByDepartment:    stmtGetCitiesByDepartment,
		stmtValidateCityInDepartment: stmtValidateCityInDepartment,
		stmtGetDepartmentByID:        stmtGetDepartmentByID,
	}, nil
}

// GetAllDepartments retrieves all departments ordered by name
func (r *repository) GetAllDepartments(ctx context.Context) ([]domain.Department, error) {
	rows, err := r.stmtGetAllDepartments.QueryContext(ctx)
	if err != nil {
		log.Error("error querying departments", "error", err)
		return nil, err
	}
	defer rows.Close()

	var departments []domain.Department
	for rows.Next() {
		var dept domain.Department
		if err := rows.Scan(&dept.ID, &dept.Name); err != nil {
			log.Error("error scanning department", "error", err)
			continue
		}
		departments = append(departments, dept)
	}

	if err := rows.Err(); err != nil {
		log.Error("error iterating departments", "error", err)
		return nil, err
	}

	return departments, nil
}

// GetCitiesByDepartment retrieves all cities for a specific department
func (r *repository) GetCitiesByDepartment(ctx context.Context, departmentID string) ([]domain.City, error) {
	rows, err := r.stmtGetCitiesByDepartment.QueryContext(ctx, departmentID)
	if err != nil {
		log.Error("error querying cities", "error", err, "department_id", departmentID)
		return nil, err
	}
	defer rows.Close()

	var cities []domain.City
	for rows.Next() {
		var city domain.City
		if err := rows.Scan(&city.ID, &city.Name, &city.DepartmentID); err != nil {
			log.Error("error scanning city", "error", err)
			continue
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		log.Error("error iterating cities", "error", err)
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
		log.Error("error validating city in department", "error", err, "city_id", cityID, "department_id", departmentID)
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
		log.Error("error getting department by ID", "error", err, "department_id", departmentID)
		return nil, err
	}
	return &dept, nil
}
