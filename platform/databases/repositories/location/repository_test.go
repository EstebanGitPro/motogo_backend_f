package location

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
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
