package location

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewRepository Tests
// ============================================

func TestNewRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect all prepared statements
	mock.ExpectPrepare("SELECT id, name FROM departments")
	mock.ExpectPrepare("SELECT id, name, department_id FROM cities")
	mock.ExpectPrepare("SELECT 1 FROM cities WHERE id")
	mock.ExpectPrepare("SELECT id, name FROM departments WHERE id")
	mock.ExpectPrepare("INSERT INTO locations")
	mock.ExpectPrepare("UPDATE locations")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_GetAllDepartments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT id, name FROM departments").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GetAllDepartments")
}

// ============================================
// ValidateCityInDepartment Tests
// ============================================

func TestValidateCityInDepartment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT 1 FROM cities WHERE id")

	stmt, err := db.Prepare("SELECT 1 FROM cities WHERE id = ? AND department_id = ?")
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)

	mock.ExpectQuery("SELECT 1 FROM cities WHERE id").
		WithArgs("city-1", "dept-1").
		WillReturnRows(rows)

	repo := &repository{stmtValidateCityInDepartment: stmt}

	err = repo.ValidateCityInDepartment(context.Background(), "city-1", "dept-1")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestValidateCityInDepartment_CityNotInDepartment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT 1 FROM cities WHERE id")

	stmt, err := db.Prepare("SELECT 1 FROM cities WHERE id = ? AND department_id = ?")
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT 1 FROM cities WHERE id").
		WithArgs("city-1", "dept-2").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{stmtValidateCityInDepartment: stmt}

	err = repo.ValidateCityInDepartment(context.Background(), "city-1", "dept-2")

	assert.Equal(t, domain.ErrCityNotInDepartment, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestValidateCityInDepartment_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT 1 FROM cities WHERE id")

	stmt, err := db.Prepare("SELECT 1 FROM cities WHERE id = ? AND department_id = ?")
	assert.NoError(t, err)

	mock.ExpectQuery("SELECT 1 FROM cities WHERE id").
		WithArgs("city-1", "dept-1").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{stmtValidateCityInDepartment: stmt}

	err = repo.ValidateCityInDepartment(context.Background(), "city-1", "dept-1")

	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// BeginTx Tests
// ============================================

func TestBeginTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	repo := &repository{db: db}

	tx, err := repo.BeginTx(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	tx, err := repo.BeginTx(context.Background())

	assert.Nil(t, tx)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// GetAllDepartments Tests
// ============================================

func TestGetAllDepartments_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("dept-001", "Bogotá D.C.").
		AddRow("dept-002", "Antioquia").
		AddRow("dept-003", "Valle del Cauca")

	stmt := mock.ExpectPrepare("SELECT id, name FROM departments")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllDepartments, _ = db.Prepare("SELECT id, name FROM departments")

	departments, err := repo.GetAllDepartments(context.Background())

	assert.NoError(t, err)
	assert.Len(t, departments, 3)
	assert.Equal(t, "dept-001", departments[0].ID)
	assert.Equal(t, "Bogotá D.C.", departments[0].Name)
}

func TestGetAllDepartments_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"})

	stmt := mock.ExpectPrepare("SELECT id, name FROM departments")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllDepartments, _ = db.Prepare("SELECT id, name FROM departments")

	departments, err := repo.GetAllDepartments(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, departments)
}

func TestGetAllDepartments_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name FROM departments")
	stmt.ExpectQuery().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetAllDepartments, _ = db.Prepare("SELECT id, name FROM departments")

	departments, err := repo.GetAllDepartments(context.Background())

	assert.Nil(t, departments)
	assert.Error(t, err)
}

// ============================================
// GetCitiesByDepartment Tests
// ============================================

func TestGetCitiesByDepartment_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "department_id"}).
		AddRow("city-001", "Medellín", "dept-002").
		AddRow("city-002", "Envigado", "dept-002").
		AddRow("city-003", "Bello", "dept-002")

	stmt := mock.ExpectPrepare("SELECT id, name, department_id FROM cities")
	stmt.ExpectQuery().
		WithArgs("dept-002").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetCitiesByDepartment, _ = db.Prepare("SELECT id, name, department_id FROM cities WHERE department_id = ?")

	cities, err := repo.GetCitiesByDepartment(context.Background(), "dept-002")

	assert.NoError(t, err)
	assert.Len(t, cities, 3)
	assert.Equal(t, "city-001", cities[0].ID)
	assert.Equal(t, "Medellín", cities[0].Name)
	assert.Equal(t, "dept-002", cities[0].DepartmentID)
}

func TestGetCitiesByDepartment_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "department_id"})

	stmt := mock.ExpectPrepare("SELECT id, name, department_id FROM cities")
	stmt.ExpectQuery().
		WithArgs("dept-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetCitiesByDepartment, _ = db.Prepare("SELECT id, name, department_id FROM cities WHERE department_id = ?")

	cities, err := repo.GetCitiesByDepartment(context.Background(), "dept-empty")

	assert.NoError(t, err)
	assert.Empty(t, cities)
}

func TestGetCitiesByDepartment_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, department_id FROM cities")
	stmt.ExpectQuery().
		WithArgs("dept-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetCitiesByDepartment, _ = db.Prepare("SELECT id, name, department_id FROM cities WHERE department_id = ?")

	cities, err := repo.GetCitiesByDepartment(context.Background(), "dept-error")

	assert.Nil(t, cities)
	assert.Error(t, err)
}

// ============================================
// GetDepartmentByID Tests
// ============================================

func TestGetDepartmentByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("dept-001", "Antioquia")

	stmt := mock.ExpectPrepare("SELECT id, name FROM departments")
	stmt.ExpectQuery().
		WithArgs("dept-001").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDepartmentByID, _ = db.Prepare("SELECT id, name FROM departments WHERE id = ?")

	dept, err := repo.GetDepartmentByID(context.Background(), "dept-001")

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, "dept-001", dept.ID)
	assert.Equal(t, "Antioquia", dept.Name)
}

func TestGetDepartmentByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name FROM departments")
	stmt.ExpectQuery().
		WithArgs("dept-not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetDepartmentByID, _ = db.Prepare("SELECT id, name FROM departments WHERE id = ?")

	dept, err := repo.GetDepartmentByID(context.Background(), "dept-not-found")

	assert.Nil(t, dept)
	assert.Equal(t, domain.ErrDepartmentNotFound, err)
}

func TestGetDepartmentByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name FROM departments")
	stmt.ExpectQuery().
		WithArgs("dept-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetDepartmentByID, _ = db.Prepare("SELECT id, name FROM departments WHERE id = ?")

	dept, err := repo.GetDepartmentByID(context.Background(), "dept-error")

	assert.Nil(t, dept)
	assert.Error(t, err)
}

// ============================================
// Location ToDomain Tests
// ============================================

func TestLocation_ToDomain_WithCoordinates(t *testing.T) {
	loc := Location{
		ID:        "loc-001",
		BranchID:  "branch-001",
		CityID:    "city-001",
		Address:   "Calle 123",
		Latitude:  sql.NullFloat64{Float64: 4.6097, Valid: true},
		Longitude: sql.NullFloat64{Float64: -74.0817, Valid: true},
	}

	result := loc.ToDomain()

	assert.Equal(t, "loc-001", result.ID)
	assert.Equal(t, "branch-001", result.BranchID)
	assert.Equal(t, "city-001", result.CityID)
	assert.Equal(t, "Calle 123", result.Address)
	assert.NotNil(t, result.Latitude)
	assert.Equal(t, 4.6097, *result.Latitude)
	assert.NotNil(t, result.Longitude)
	assert.Equal(t, -74.0817, *result.Longitude)
}

func TestLocation_ToDomain_WithoutCoordinates(t *testing.T) {
	loc := Location{
		ID:       "loc-002",
		BranchID: "branch-002",
		CityID:   "city-002",
		Address:  "Carrera 45",
	}

	result := loc.ToDomain()

	assert.Equal(t, "loc-002", result.ID)
	assert.Nil(t, result.Latitude)
	assert.Nil(t, result.Longitude)
}

func TestFromDomainLocation_WithCoordinates(t *testing.T) {
	lat := 4.6097
	lon := -74.0817
	domainLoc := domain.Location{
		ID:        "loc-001",
		BranchID:  "branch-001",
		CityID:    "city-001",
		Address:   "Calle 123",
		Latitude:  &lat,
		Longitude: &lon,
	}

	result := FromDomainLocation(domainLoc)

	assert.Equal(t, "loc-001", result.ID)
	assert.Equal(t, "branch-001", result.BranchID)
	assert.True(t, result.Latitude.Valid)
	assert.Equal(t, 4.6097, result.Latitude.Float64)
	assert.True(t, result.Longitude.Valid)
	assert.Equal(t, -74.0817, result.Longitude.Float64)
}

func TestFromDomainLocation_WithoutCoordinates(t *testing.T) {
	domainLoc := domain.Location{
		ID:       "loc-002",
		BranchID: "branch-002",
		CityID:   "city-002",
		Address:  "Carrera 45",
	}

	result := FromDomainLocation(domainLoc)

	assert.Equal(t, "loc-002", result.ID)
	assert.False(t, result.Latitude.Valid)
	assert.False(t, result.Longitude.Valid)
}

// ============================================
// Department ToDomain Tests
// ============================================

func TestDepartment_ToDomain(t *testing.T) {
	dept := Department{
		ID:   "dept-001",
		Name: "Cundinamarca",
	}

	result := dept.ToDomain()

	assert.Equal(t, "dept-001", result.ID)
	assert.Equal(t, "Cundinamarca", result.Name)
}

// ============================================
// City ToDomain Tests
// ============================================

func TestCity_ToDomain(t *testing.T) {
	city := City{
		ID:           "city-001",
		Name:         "Bogotá",
		DepartmentID: "dept-001",
	}

	result := city.ToDomain()

	assert.Equal(t, "city-001", result.ID)
	assert.Equal(t, "Bogotá", result.Name)
	assert.Equal(t, "dept-001", result.DepartmentID)
}

// ============================================
// CheckAddressExists Tests
// ============================================

func TestCheckAddressExists_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"1"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM locations").
		WithArgs("Calle 123").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.CheckAddressExists(context.Background(), "Calle 123")

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCheckAddressExists_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT 1 FROM locations").
		WithArgs("Nonexistent Address").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}

	exists, err := repo.CheckAddressExists(context.Background(), "Nonexistent Address")

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCheckAddressExists_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT 1 FROM locations").
		WithArgs("Error Address").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	exists, err := repo.CheckAddressExists(context.Background(), "Error Address")

	assert.Error(t, err)
	assert.False(t, exists)
}

// ============================================
// SaveLocation Tests
// ============================================

func TestSaveLocation_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO locations").
		WithArgs(sqlmock.AnyArg(), "branch-123", "city-123", "Test Address", nil, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	location := domain.Location{BranchID: "branch-123", CityID: "city-123", Address: "Test Address"}
	err = repo.SaveLocation(context.Background(), sqlTx, location)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveLocation_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO locations").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveLocation(context.Background(), sqlTx, domain.Location{})
	assert.Equal(t, domain.ErrLocationCannotSave, err)
}

func TestSaveLocation_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.SaveLocation(context.Background(), nil, domain.Location{})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// UpdateLocation Tests
// ============================================

func TestUpdateLocation_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE locations").
		WithArgs("city-456", "Updated Address", nil, nil, "branch-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	location := domain.Location{BranchID: "branch-123", CityID: "city-456", Address: "Updated Address"}
	err = repo.UpdateLocation(context.Background(), sqlTx, location)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateLocation_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE locations").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.UpdateLocation(context.Background(), sqlTx, domain.Location{})
	assert.Equal(t, domain.ErrLocationCannotSave, err)
}

func TestUpdateLocation_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.UpdateLocation(context.Background(), nil, domain.Location{})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}
