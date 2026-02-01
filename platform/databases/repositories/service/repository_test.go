package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewRepository Tests
// ============================================

func TestNewRepository_NilDB(t *testing.T) {
	repo, err := NewRepository(nil)

	assert.Nil(t, repo)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestNewRepository_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect all prepared statements
	mock.ExpectPrepare("SELECT id, name, description, service_type, is_active FROM services ORDER BY name")
	mock.ExpectPrepare("SELECT id, name, description, service_type, is_active FROM services WHERE service_type")
	mock.ExpectPrepare("SELECT id, name, description, service_type, is_active FROM services WHERE id")
	mock.ExpectPrepare("SELECT s.id, s.name.*FROM branch_services bs")
	mock.ExpectPrepare("UPDATE services SET name")
	mock.ExpectPrepare("INSERT INTO branch_services")
	mock.ExpectPrepare("DELETE FROM branch_services")
	mock.ExpectPrepare("SELECT COUNT.*FROM branch_services")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT id, name, description").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestNewRepository_PrepareError_GetByType(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT id, name.*ORDER BY name")
	mock.ExpectPrepare("SELECT id, name.*WHERE service_type").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
}

func TestNewRepository_PrepareError_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT id, name.*ORDER BY name")
	mock.ExpectPrepare("SELECT id, name.*WHERE service_type")
	mock.ExpectPrepare("SELECT id, name.*WHERE id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
}

func TestNewRepository_PrepareError_GetByBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT id, name.*ORDER BY name")
	mock.ExpectPrepare("SELECT id, name.*WHERE service_type")
	mock.ExpectPrepare("SELECT id, name.*WHERE id")
	mock.ExpectPrepare("SELECT.*branch_services.*WHERE bs.branch_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
}

func TestNewRepository_PrepareError_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT id, name.*ORDER BY name")
	mock.ExpectPrepare("SELECT id, name.*WHERE service_type")
	mock.ExpectPrepare("SELECT id, name.*WHERE id")
	mock.ExpectPrepare("SELECT.*branch_services.*WHERE bs.branch_id")
	mock.ExpectPrepare("UPDATE services").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
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
// GetAllServices Tests
// ============================================

func TestGetAllServices_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "is_active"}).
		AddRow("svc-001", "Cambio de aceite", "Servicio de cambio de aceite", "MAINTENANCE", true).
		AddRow("svc-002", "Lavado completo", nil, "CLEANING", true).
		AddRow("svc-003", "Revisión de frenos", "Inspección y ajuste", "MAINTENANCE", false)

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllServices, _ = db.Prepare("SELECT id, name, description FROM services")

	services, err := repo.GetAllServices(context.Background())

	assert.NoError(t, err)
	assert.Len(t, services, 3)
	assert.Equal(t, "svc-001", services[0].ID)
	assert.Equal(t, "Cambio de aceite", services[0].Name)
	assert.True(t, services[0].IsActive)
	assert.Equal(t, "", services[1].Description) // nil description
}

func TestGetAllServices_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "is_active"})

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllServices, _ = db.Prepare("SELECT id, name, description FROM services")

	services, err := repo.GetAllServices(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, services)
}

func TestGetAllServices_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetAllServices, _ = db.Prepare("SELECT id, name, description FROM services")

	services, err := repo.GetAllServices(context.Background())

	assert.Nil(t, services)
	assert.Error(t, err)
}

// ============================================
// GetServiceByID Tests
// ============================================

func TestGetServiceByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "is_active"}).
		AddRow("svc-123", "Servicio Premium", "Descripción completa", "PREMIUM", true)

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().
		WithArgs("svc-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetServiceByID, _ = db.Prepare("SELECT id, name, description FROM services WHERE id = ?")

	service, err := repo.GetServiceByID(context.Background(), "svc-123")

	assert.NoError(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "svc-123", service.ID)
	assert.Equal(t, "Servicio Premium", service.Name)
	assert.Equal(t, "Descripción completa", service.Description)
	assert.True(t, service.IsActive)
}

func TestGetServiceByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetServiceByID, _ = db.Prepare("SELECT id, name, description FROM services WHERE id = ?")

	service, err := repo.GetServiceByID(context.Background(), "not-found")

	assert.Nil(t, service)
	assert.Equal(t, domain.ErrServiceNotFound, err)
}

func TestGetServiceByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().
		WithArgs("svc-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetServiceByID, _ = db.Prepare("SELECT id, name, description FROM services WHERE id = ?")

	service, err := repo.GetServiceByID(context.Background(), "svc-error")

	assert.Nil(t, service)
	assert.Error(t, err)
}

// ============================================
// CheckServiceAssociation Tests
// ============================================

func TestCheckServiceAssociation_Associated(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("branch-123", "svc-456").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckServiceAssociation, _ = db.Prepare("SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?")

	associated, err := repo.CheckServiceAssociation(context.Background(), "branch-123", "svc-456")

	assert.NoError(t, err)
	assert.True(t, associated)
}

func TestCheckServiceAssociation_NotAssociated(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("branch-123", "svc-789").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCheckServiceAssociation, _ = db.Prepare("SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?")

	associated, err := repo.CheckServiceAssociation(context.Background(), "branch-123", "svc-789")

	assert.NoError(t, err)
	assert.False(t, associated)
}

func TestCheckServiceAssociation_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("branch-error", "svc-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCheckServiceAssociation, _ = db.Prepare("SELECT COUNT(*) FROM branch_services WHERE branch_id = ? AND service_id = ?")

	associated, err := repo.CheckServiceAssociation(context.Background(), "branch-error", "svc-error")

	assert.False(t, associated)
	assert.Error(t, err)
}

// ============================================
// GetServicesByType Tests
// ============================================

func TestGetServicesByType_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "is_active"}).
		AddRow("svc-001", "Cambio de aceite", "Servicio básico", "MAINTENANCE", true).
		AddRow("svc-002", "Revisión de frenos", nil, "MAINTENANCE", true)

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().
		WithArgs("MAINTENANCE").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetServicesByType, _ = db.Prepare("SELECT id, name, description, service_type, is_active FROM services WHERE service_type = ?")

	services, err := repo.GetServicesByType(context.Background(), "MAINTENANCE")

	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "svc-001", services[0].ID)
	assert.Equal(t, domain.ServiceType("MAINTENANCE"), services[0].ServiceType)
}

func TestGetServicesByType_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "is_active"})

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().
		WithArgs("NONEXISTENT").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetServicesByType, _ = db.Prepare("SELECT id, name, description, service_type, is_active FROM services WHERE service_type = ?")

	services, err := repo.GetServicesByType(context.Background(), "NONEXISTENT")

	assert.NoError(t, err)
	assert.Empty(t, services)
}

func TestGetServicesByType_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description")
	stmt.ExpectQuery().
		WithArgs("ERROR").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetServicesByType, _ = db.Prepare("SELECT id, name, description, service_type, is_active FROM services WHERE service_type = ?")

	services, err := repo.GetServicesByType(context.Background(), "ERROR")

	assert.Nil(t, services)
	assert.Error(t, err)
}

// ============================================
// ValidateServiceIDs Tests
// ============================================

func TestValidateServiceIDs_Empty(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.ValidateServiceIDs(context.Background(), []string{})

	assert.NoError(t, err)
}

func TestValidateServiceIDs_AllFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("svc-001", "svc-002").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateServiceIDs(context.Background(), []string{"svc-001", "svc-002"})

	assert.NoError(t, err)
}

func TestValidateServiceIDs_SomeNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1) // Only 1 found out of 2

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("svc-001", "svc-notfound").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateServiceIDs(context.Background(), []string{"svc-001", "svc-notfound"})

	assert.Equal(t, domain.ErrServiceNotFound, err)
}

func TestValidateServiceIDs_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT COUNT").
		WithArgs("svc-001").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	err = repo.ValidateServiceIDs(context.Background(), []string{"svc-001"})

	assert.Error(t, err)
}

// ============================================
// GetServicesByBranch Tests
// ============================================

func TestGetServicesByBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "created_at", "active"}).
		AddRow("svc-001", "Cambio de aceite", "Servicio básico", "MAINTENANCE", now, true).
		AddRow("svc-002", "Revisión general", nil, "INSPECTION", now, false)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("branch-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetServicesByBranch, _ = db.Prepare("SELECT * FROM branch_services WHERE branch_id = ?")

	services, err := repo.GetServicesByBranch(context.Background(), "branch-123")

	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "svc-001", services[0].Service.ID)
	assert.True(t, services[0].Active)
}

func TestGetServicesByBranch_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "service_type", "created_at", "active"})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("branch-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetServicesByBranch, _ = db.Prepare("SELECT * FROM branch_services WHERE branch_id = ?")

	services, err := repo.GetServicesByBranch(context.Background(), "branch-empty")

	assert.NoError(t, err)
	assert.Empty(t, services)
}

func TestGetServicesByBranch_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("branch-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetServicesByBranch, _ = db.Prepare("SELECT * FROM branch_services WHERE branch_id = ?")

	services, err := repo.GetServicesByBranch(context.Background(), "branch-error")

	assert.Nil(t, services)
	assert.Error(t, err)
}
