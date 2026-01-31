package motorcycle

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
	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles.*SET deleted_at")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*FROM motorcycle_references r.*ORDER BY b.name")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*FROM motorcycle_references r.*WHERE r.brand_id")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByID")
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
// GetByID Tests
// ============================================

func TestGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes",
		"ref_id", "brand_id", "brand_name", "model", "category", "engine_displacement",
	}).AddRow(
		"moto-123", "ABC123", "ref-456", "owner-789", 2020, 15000, "Buen estado",
		"ref-456", "brand-001", "Yamaha", "MT-07", "SPORT", 689,
	)

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("moto-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.id = ?")

	motorcycle, err := repo.GetByID(context.Background(), "moto-123")

	assert.NoError(t, err)
	assert.NotNil(t, motorcycle)
	assert.Equal(t, "moto-123", motorcycle.ID)
	assert.Equal(t, "ABC123", motorcycle.LicensePlate)
	assert.NotNil(t, motorcycle.Year)
	assert.Equal(t, 2020, *motorcycle.Year)
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.id = ?")

	motorcycle, err := repo.GetByID(context.Background(), "not-found")

	assert.Nil(t, motorcycle)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestGetByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("moto-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByID, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.id = ?")

	motorcycle, err := repo.GetByID(context.Background(), "moto-error")

	assert.Nil(t, motorcycle)
	assert.Error(t, err)
}

// ============================================
// GetByOwnerID Tests
// ============================================

func TestGetByOwnerID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes",
		"ref_id", "brand_id", "brand_name", "model", "category", "engine_displacement",
	}).AddRow(
		"moto-1", "XYZ999", "ref-1", "owner-123", 2021, 5000, nil,
		"ref-1", "brand-002", "Honda", "CBR-600", "SPORT", 599,
	).AddRow(
		"moto-2", "DEF456", "ref-2", "owner-123", 2019, 25000, "Revision reciente",
		"ref-2", "brand-001", "Yamaha", "R6", "SUPERSPORT", 599,
	)

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("owner-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByOwnerID, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.owner_id = ?")

	motorcycles, err := repo.GetByOwnerID(context.Background(), "owner-123")

	assert.NoError(t, err)
	assert.Len(t, motorcycles, 2)
	assert.Equal(t, "moto-1", motorcycles[0].ID)
	assert.Equal(t, "moto-2", motorcycles[1].ID)
}

func TestGetByOwnerID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes",
		"ref_id", "brand_id", "brand_name", "model", "category", "engine_displacement",
	})

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("no-motos").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByOwnerID, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.owner_id = ?")

	motorcycles, err := repo.GetByOwnerID(context.Background(), "no-motos")

	assert.NoError(t, err)
	assert.Empty(t, motorcycles)
}

func TestGetByOwnerID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("owner-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByOwnerID, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.owner_id = ?")

	motorcycles, err := repo.GetByOwnerID(context.Background(), "owner-error")

	assert.Nil(t, motorcycles)
	assert.Error(t, err)
}

// ============================================
// GetByLicensePlate Tests
// ============================================

func TestGetByLicensePlate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes",
		"ref_id", "brand_id", "brand_name", "model", "category", "engine_displacement",
	}).AddRow(
		"moto-by-plate", "ZZZ111", "ref-plate", "owner-plate", 2022, 1000, nil,
		"ref-plate", "brand-003", "Kawasaki", "Ninja-400", "SPORT", 399,
	)

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("ZZZ111").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetByLicensePlate, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.license_plate = ?")

	motorcycle, err := repo.GetByLicensePlate(context.Background(), "ZZZ111")

	assert.NoError(t, err)
	assert.NotNil(t, motorcycle)
	assert.Equal(t, "moto-by-plate", motorcycle.ID)
	assert.Equal(t, "ZZZ111", motorcycle.LicensePlate)
}

func TestGetByLicensePlate_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("NOT-EXIST").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetByLicensePlate, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.license_plate = ?")

	motorcycle, err := repo.GetByLicensePlate(context.Background(), "NOT-EXIST")

	assert.Nil(t, motorcycle)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

// ============================================
// GetReferencesByBrandID Tests
// ============================================

func TestGetReferencesByBrandID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "brand_id", "brand_name", "model", "category", "engine_displacement"}).
		AddRow("ref-001", "brand-001", "Yamaha", "YZF-R3", "Sport", 321).
		AddRow("ref-002", "brand-001", "Yamaha", "MT-07", "Naked", 689)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("brand-001").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetReferencesByBrandID, _ = db.Prepare("SELECT * FROM motorcycle_references WHERE brand_id = ?")

	refs, err := repo.GetReferencesByBrandID(context.Background(), "brand-001")

	assert.NoError(t, err)
	assert.Len(t, refs, 2)
	assert.Equal(t, "ref-001", refs[0].ID)
	assert.Equal(t, "Yamaha", refs[0].BrandName)
	assert.Equal(t, "YZF-R3", refs[0].Model)
}

func TestGetReferencesByBrandID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "brand_id", "brand_name", "model", "category", "engine_displacement"})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("brand-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetReferencesByBrandID, _ = db.Prepare("SELECT * FROM motorcycle_references WHERE brand_id = ?")

	refs, err := repo.GetReferencesByBrandID(context.Background(), "brand-empty")

	assert.NoError(t, err)
	assert.Empty(t, refs)
}

func TestGetReferencesByBrandID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("brand-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetReferencesByBrandID, _ = db.Prepare("SELECT * FROM motorcycle_references WHERE brand_id = ?")

	refs, err := repo.GetReferencesByBrandID(context.Background(), "brand-error")

	assert.Nil(t, refs)
	assert.Error(t, err)
}
