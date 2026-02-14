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
	// Schema validation statement (prepared and immediately closed)
	mock.ExpectPrepare("SELECT EXISTS.*completed_services.*diagnostics")
	mock.ExpectPrepare("SELECT r.category.*FROM motorcycle_references")
	mock.ExpectPrepare("SELECT r.model.*FROM motorcycle_references r.*WHERE r.category")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

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

func TestNewRepository_PrepareError_GetByOwnerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByOwnerID")
}

func TestNewRepository_PrepareError_GetByLicensePlate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.license_plate").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetByLicensePlate")
}

func TestNewRepository_PrepareError_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtUpdate")
}

func TestNewRepository_PrepareError_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles.*SET reference_id")
	mock.ExpectPrepare("UPDATE motorcycles.*SET deleted_at").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDelete")
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
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes", "profile_image_url",
		"ref_id", "brand_id", "brand_name", "model", "category", "engine_displacement",
	}).AddRow(
		"moto-123", "ABC123", "ref-456", "owner-789", 2020, 15000, "Buen estado", "https://storage.example.com/image.jpg",
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
	assert.NotNil(t, motorcycle.ProfileImageURL)
	assert.Equal(t, "https://storage.example.com/image.jpg", *motorcycle.ProfileImageURL)
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
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes", "profile_image_url",
		"ref_id", "brand_id", "brand_name", "model", "category", "engine_displacement",
	}).AddRow(
		"moto-1", "XYZ999", "ref-1", "owner-123", 2021, 5000, nil, nil,
		"ref-1", "brand-002", "Honda", "CBR-600", "SPORT", 599,
	).AddRow(
		"moto-2", "DEF456", "ref-2", "owner-123", 2019, 25000, "Revision reciente", "https://example.com/image.jpg",
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
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes", "profile_image_url",
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
		"id", "license_plate", "reference_id", "owner_id", "year", "current_mileage", "owner_notes", "profile_image_url",
		"ref_id", "brand_id", "brand_name", "model", "category", "engine_displacement",
	}).AddRow(
		"moto-by-plate", "ZZZ111", "ref-plate", "owner-plate", 2022, 1000, nil, nil,
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

// ============================================
// GetAllReferences Tests
// ============================================

func TestGetAllReferences_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "brand_id", "brand_name", "model", "category", "engine_displacement"}).
		AddRow("ref-001", "brand-001", "Yamaha", "YZF-R3", "Sport", 321).
		AddRow("ref-002", "brand-002", "Honda", "CBR500R", "Sport", 500)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllReferences, _ = db.Prepare("SELECT * FROM motorcycle_references")

	refs, err := repo.GetAllReferences(context.Background())

	assert.NoError(t, err)
	assert.Len(t, refs, 2)
	assert.Equal(t, "ref-001", refs[0].ID)
	assert.Equal(t, "Yamaha", refs[0].BrandName)
}

func TestGetAllReferences_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "brand_id", "brand_name", "model", "category", "engine_displacement"})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllReferences, _ = db.Prepare("SELECT * FROM motorcycle_references")

	refs, err := repo.GetAllReferences(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, refs)
}

func TestGetAllReferences_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetAllReferences, _ = db.Prepare("SELECT * FROM motorcycle_references")

	refs, err := repo.GetAllReferences(context.Background())

	assert.Nil(t, refs)
	assert.Error(t, err)
}

// ============================================
// ValidateReferenceExists Tests
// ============================================

func TestValidateReferenceExists_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("ref-valid").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.ValidateReferenceExists(context.Background(), "ref-valid")

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestValidateReferenceExists_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("ref-invalid").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.ValidateReferenceExists(context.Background(), "ref-invalid")

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestValidateReferenceExists_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("ref-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	exists, err := repo.ValidateReferenceExists(context.Background(), "ref-error")

	assert.False(t, exists)
	assert.Error(t, err)
}

// ============================================
// CheckLicensePlateExists Tests
// ============================================

func TestCheckLicensePlateExists_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("ABC123").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.CheckLicensePlateExists(context.Background(), "ABC123")

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCheckLicensePlateExists_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("NEW123").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.CheckLicensePlateExists(context.Background(), "NEW123")

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCheckLicensePlateExists_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("ERROR123").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	exists, err := repo.CheckLicensePlateExists(context.Background(), "ERROR123")

	assert.False(t, exists)
	assert.Error(t, err)
}

// ============================================
// Motorcycle Model Tests
// ============================================

func TestMotorcycle_ToDomain(t *testing.T) {
	year := 2021
	mileage := 15000
	notes := "Buen estado"

	ref := &MotorcycleReference{
		ID:                 sql.NullString{String: "ref-001", Valid: true},
		BrandID:            sql.NullString{String: "brand-001", Valid: true},
		BrandName:          sql.NullString{String: "Yamaha", Valid: true},
		Model:              sql.NullString{String: "MT-07", Valid: true},
		Category:           sql.NullString{String: "Naked", Valid: true},
		EngineDisplacement: sql.NullInt64{Int64: 689, Valid: true},
	}

	m := &Motorcycle{
		ID:             "moto-001",
		LicensePlate:   "ABC123",
		ReferenceID:    sql.NullString{String: "ref-001", Valid: true},
		OwnerID:        "owner-001",
		Year:           sql.NullInt64{Int64: int64(year), Valid: true},
		CurrentMileage: sql.NullInt64{Int64: int64(mileage), Valid: true},
		OwnerNotes:     sql.NullString{String: notes, Valid: true},
	}

	dm := m.ToDomain(ref)

	assert.Equal(t, "moto-001", dm.ID)
	assert.Equal(t, "ABC123", dm.LicensePlate)
	assert.Equal(t, "ref-001", dm.ReferenceID)
	assert.Equal(t, "owner-001", dm.OwnerID)
	assert.NotNil(t, dm.Year)
	assert.Equal(t, year, *dm.Year)
	assert.NotNil(t, dm.CurrentMileage)
	assert.Equal(t, mileage, *dm.CurrentMileage)
	assert.NotNil(t, dm.OwnerNotes)
	assert.Equal(t, notes, *dm.OwnerNotes)
	assert.NotNil(t, dm.Reference)
	assert.Equal(t, "Yamaha", dm.Reference.BrandName)
}

func TestMotorcycle_ToDomain_NullFields(t *testing.T) {
	m := &Motorcycle{
		ID:             "moto-002",
		LicensePlate:   "XYZ789",
		ReferenceID:    sql.NullString{Valid: false},
		OwnerID:        "owner-002",
		Year:           sql.NullInt64{Valid: false},
		CurrentMileage: sql.NullInt64{Valid: false},
		OwnerNotes:     sql.NullString{Valid: false},
	}

	dm := m.ToDomain(nil)

	assert.Equal(t, "moto-002", dm.ID)
	assert.Empty(t, dm.ReferenceID)
	assert.Nil(t, dm.Year)
	assert.Nil(t, dm.CurrentMileage)
	assert.Nil(t, dm.OwnerNotes)
	assert.Nil(t, dm.Reference)
}

func TestMotorcycle_ToDomain_WithReference(t *testing.T) {
	ref := &MotorcycleReference{
		ID:                 sql.NullString{String: "ref-001", Valid: true},
		BrandID:            sql.NullString{String: "brand-001", Valid: true},
		BrandName:          sql.NullString{String: "Yamaha", Valid: true},
		Model:              sql.NullString{String: "MT-07", Valid: true},
		Category:           sql.NullString{String: "Naked", Valid: true},
		EngineDisplacement: sql.NullInt64{Int64: 689, Valid: true},
	}

	m := &Motorcycle{
		ID:           "moto-003",
		LicensePlate: "DEF456",
		OwnerID:      "owner-003",
	}

	dm := m.ToDomain(ref)

	assert.NotNil(t, dm.Reference)
	assert.Equal(t, "ref-001", dm.Reference.ID)
	assert.Equal(t, "brand-001", dm.Reference.BrandID)
	assert.Equal(t, "Yamaha", dm.Reference.BrandName)
	assert.Equal(t, "MT-07", dm.Reference.Model)
	assert.Equal(t, "Naked", dm.Reference.Category)
	assert.Equal(t, 689, dm.Reference.EngineDisplacement)
}

// ============================================
// GetByLicensePlate Additional Tests
// ============================================

func TestGetByLicensePlate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT m.id, m.license_plate")
	stmt.ExpectQuery().
		WithArgs("ERROR-PLATE").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetByLicensePlate, _ = db.Prepare("SELECT m.id, m.license_plate FROM motorcycles m WHERE m.license_plate = ?")

	motorcycle, err := repo.GetByLicensePlate(context.Background(), "ERROR-PLATE")

	assert.Nil(t, motorcycle)
	assert.Error(t, err)
}

// ============================================
// ToDomain Additional Tests
// ============================================

func TestMotorcycle_ToDomain_WithPartialReference(t *testing.T) {
	// Reference with only some fields valid
	ref := &MotorcycleReference{
		ID:                 sql.NullString{String: "ref-002", Valid: true},
		BrandID:            sql.NullString{String: "brand-002", Valid: true},
		BrandName:          sql.NullString{Valid: false}, // NULL
		Model:              sql.NullString{Valid: false}, // NULL
		Category:           sql.NullString{Valid: false}, // NULL
		EngineDisplacement: sql.NullInt64{Valid: false},  // NULL
	}

	m := &Motorcycle{
		ID:           "moto-partial",
		LicensePlate: "PAR123",
		OwnerID:      "owner-partial",
		ReferenceID:  sql.NullString{String: "ref-002", Valid: true},
	}

	dm := m.ToDomain(ref)

	assert.NotNil(t, dm.Reference)
	assert.Equal(t, "ref-002", dm.Reference.ID)
	assert.Equal(t, "brand-002", dm.Reference.BrandID)
	assert.Equal(t, "", dm.Reference.BrandName) // Empty because NULL
	assert.Equal(t, "", dm.Reference.Model)
	assert.Equal(t, "", dm.Reference.Category)
	assert.Equal(t, 0, dm.Reference.EngineDisplacement)
}

func TestMotorcycle_ToDomain_PreservesReferenceID(t *testing.T) {
	m := &Motorcycle{
		ID:          "moto-refid",
		ReferenceID: sql.NullString{String: "some-reference-uuid", Valid: true},
		OwnerID:     "owner-test",
	}

	dm := m.ToDomain(nil)

	assert.Equal(t, "some-reference-uuid", dm.ReferenceID)
	assert.Nil(t, dm.Reference)
}

func TestMotorcycle_ToDomain_AllOptionalFieldsSet(t *testing.T) {
	m := &Motorcycle{
		ID:             "moto-full",
		LicensePlate:   "FUL999",
		ReferenceID:    sql.NullString{String: "ref-full", Valid: true},
		OwnerID:        "owner-full",
		Year:           sql.NullInt64{Int64: 2023, Valid: true},
		CurrentMileage: sql.NullInt64{Int64: 15000, Valid: true},
		OwnerNotes:     sql.NullString{String: "Excellent condition", Valid: true},
	}

	dm := m.ToDomain(nil)

	assert.Equal(t, "moto-full", dm.ID)
	assert.Equal(t, "FUL999", dm.LicensePlate)
	assert.Equal(t, "ref-full", dm.ReferenceID)
	assert.NotNil(t, dm.Year)
	assert.Equal(t, 2023, *dm.Year)
	assert.NotNil(t, dm.CurrentMileage)
	assert.Equal(t, 15000, *dm.CurrentMileage)
	assert.NotNil(t, dm.OwnerNotes)
	assert.Equal(t, "Excellent condition", *dm.OwnerNotes)
}

func TestMotorcycle_ToDomain_ReferenceWithNullID(t *testing.T) {
	// Reference exists but ID is not valid - should not create domain reference
	ref := &MotorcycleReference{
		ID:        sql.NullString{Valid: false}, // NULL ID
		BrandID:   sql.NullString{String: "brand-x", Valid: true},
		BrandName: sql.NullString{String: "Honda", Valid: true},
	}

	m := &Motorcycle{
		ID:           "moto-nullref",
		LicensePlate: "NUL000",
		OwnerID:      "owner-nullref",
	}

	dm := m.ToDomain(ref)

	assert.Nil(t, dm.Reference) // Should be nil because ref.ID is not valid
}

// ============================================
// Save Tests
// ============================================

func TestSave_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("INSERT INTO motorcycles").
		WithArgs("moto-123", "ABC123", sqlmock.AnyArg(), "owner-456", 2020, 10000, "some notes", "https://image.url").
		WillReturnResult(sqlmock.NewResult(1, 1))

	year := 2020
	mileage := 10000
	notes := "some notes"
	imageURL := "https://image.url"
	motorcycle := &domain.Motorcycle{
		ID:              "moto-123",
		LicensePlate:    "ABC123",
		ReferenceID:     "ref-789",
		OwnerID:         "owner-456",
		Year:            &year,
		CurrentMileage:  &mileage,
		OwnerNotes:      &notes,
		ProfileImageURL: &imageURL,
	}

	err = repo.Save(context.Background(), tx, motorcycle)
	assert.NoError(t, err)
}

func TestSave_InvalidTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	// Pass nil transaction
	err = repo.Save(context.Background(), nil, &domain.Motorcycle{})
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleCannotSave, err)
}

func TestSave_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("INSERT INTO motorcycles").
		WillReturnError(sql.ErrConnDone)

	motorcycle := &domain.Motorcycle{
		ID:           "moto-123",
		LicensePlate: "ABC123",
	}

	err = repo.Save(context.Background(), tx, motorcycle)
	assert.Error(t, err)
}

// ============================================
// Update Tests
// ============================================

func TestUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles").
		WithArgs("ref-789", 2021, 15000, "updated notes", "https://new-image.url", "moto-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	year := 2021
	mileage := 15000
	notes := "updated notes"
	imageURL := "https://new-image.url"
	motorcycle := &domain.Motorcycle{
		ID:              "moto-123",
		ReferenceID:     "ref-789",
		Year:            &year,
		CurrentMileage:  &mileage,
		OwnerNotes:      &notes,
		ProfileImageURL: &imageURL,
	}

	err = repo.Update(context.Background(), tx, motorcycle)
	assert.NoError(t, err)
}

func TestUpdate_InvalidTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	err = repo.Update(context.Background(), nil, &domain.Motorcycle{})
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestUpdate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles").
		WillReturnError(sql.ErrConnDone)

	err = repo.Update(context.Background(), tx, &domain.Motorcycle{ID: "moto-123"})
	assert.Error(t, err)
}

// ============================================
// Delete Tests
// ============================================

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles").
		WithArgs("moto-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), tx, "moto-123")
	assert.NoError(t, err)
}

func TestDelete_InvalidTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	err = repo.Delete(context.Background(), nil, "moto-123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles").
		WithArgs("moto-not-found").
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.Delete(context.Background(), tx, "moto-not-found")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles").
		WillReturnError(sql.ErrConnDone)

	err = repo.Delete(context.Background(), tx, "moto-123")
	assert.Error(t, err)
}

// ============================================
// GetDistinctCategories Tests (HU41)
// ============================================

func TestGetDistinctCategories_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"category", "line_count"}).
		AddRow("Sport", 15).
		AddRow("Naked", 8).
		AddRow("Touring", 3)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDistinctCategories, _ = db.Prepare("SELECT r.category, COUNT(*) FROM motorcycle_references r GROUP BY r.category")

	cats, err := repo.GetDistinctCategories(context.Background())

	assert.NoError(t, err)
	assert.Len(t, cats, 3)
	assert.Equal(t, "Sport", cats[0].Name)
	assert.Equal(t, 15, cats[0].LineCount)
	assert.Equal(t, "Naked", cats[1].Name)
}

func TestGetDistinctCategories_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"category", "line_count"})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDistinctCategories, _ = db.Prepare("SELECT r.category, COUNT(*) FROM motorcycle_references r GROUP BY r.category")

	cats, err := repo.GetDistinctCategories(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, cats)
}

func TestGetDistinctCategories_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetDistinctCategories, _ = db.Prepare("SELECT r.category, COUNT(*) FROM motorcycle_references r GROUP BY r.category")

	cats, err := repo.GetDistinctCategories(context.Background())

	assert.Nil(t, cats)
	assert.Error(t, err)
}

// ============================================
// GetLinesByCategory Tests (HU41)
// ============================================

func TestGetLinesByCategory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"model", "brand_name", "engine_displacement"}).
		AddRow("MT-07", "Yamaha", 689).
		AddRow("CBR-600", "Honda", 599)

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("Sport").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetLinesByCategory, _ = db.Prepare("SELECT r.model, b.name, r.engine_displacement FROM motorcycle_references r WHERE r.category = ?")

	lines, err := repo.GetLinesByCategory(context.Background(), "Sport")

	assert.NoError(t, err)
	assert.Len(t, lines, 2)
	assert.Equal(t, "MT-07", lines[0].Model)
	assert.Equal(t, "Yamaha", lines[0].BrandName)
	assert.Equal(t, 689, lines[0].EngineDisplacement)
}

func TestGetLinesByCategory_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"model", "brand_name", "engine_displacement"})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("NonExistent").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetLinesByCategory, _ = db.Prepare("SELECT r.model, b.name, r.engine_displacement FROM motorcycle_references r WHERE r.category = ?")

	lines, err := repo.GetLinesByCategory(context.Background(), "NonExistent")

	assert.NoError(t, err)
	assert.Empty(t, lines)
}

func TestGetLinesByCategory_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("Error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetLinesByCategory, _ = db.Prepare("SELECT r.model, b.name, r.engine_displacement FROM motorcycle_references r WHERE r.category = ?")

	lines, err := repo.GetLinesByCategory(context.Background(), "Error")

	assert.Nil(t, lines)
	assert.Error(t, err)
}

// ============================================
// GetDistinctDisplacements Tests (HU49)
// ============================================

func TestGetDistinctDisplacements_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"displacement_range"}).
		AddRow("BAJO").
		AddRow("MEDIO").
		AddRow("ALTO")

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDistinctDisplacements, _ = db.Prepare("SELECT DISTINCT r.displacement_range FROM motorcycle_references r")

	disps, err := repo.GetDistinctDisplacements(context.Background())

	assert.NoError(t, err)
	assert.Len(t, disps, 3)
	assert.Equal(t, domain.DisplacementRange("BAJO"), disps[0].Range)
	assert.Equal(t, domain.DisplacementRange("MEDIO"), disps[1].Range)
	assert.Equal(t, domain.DisplacementRange("ALTO"), disps[2].Range)
}

func TestGetDistinctDisplacements_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"displacement_range"})

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetDistinctDisplacements, _ = db.Prepare("SELECT DISTINCT r.displacement_range FROM motorcycle_references r")

	disps, err := repo.GetDistinctDisplacements(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, disps)
}

func TestGetDistinctDisplacements_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetDistinctDisplacements, _ = db.Prepare("SELECT DISTINCT r.displacement_range FROM motorcycle_references r")

	disps, err := repo.GetDistinctDisplacements(context.Background())

	assert.Nil(t, disps)
	assert.Error(t, err)
}

// ============================================
// HardDelete Tests (HU45)
// ============================================

func TestHardDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM motorcycles").
		WithArgs("moto-hard-del").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.HardDelete(context.Background(), tx, "moto-hard-del")
	assert.NoError(t, err)
}

func TestHardDelete_InvalidTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	err = repo.HardDelete(context.Background(), nil, "moto-123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestHardDelete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM motorcycles").
		WithArgs("not-found").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.HardDelete(context.Background(), tx, "not-found")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestHardDelete_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("DELETE FROM motorcycles").
		WillReturnError(sql.ErrConnDone)

	err = repo.HardDelete(context.Background(), tx, "moto-123")
	assert.Error(t, err)
}

// ============================================
// ClearProfileImageURL Tests (HU39)
// ============================================

func TestClearProfileImageURL_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles.*SET profile_image_url").
		WithArgs("moto-clear").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.ClearProfileImageURL(context.Background(), tx, "moto-clear")
	assert.NoError(t, err)
}

func TestClearProfileImageURL_InvalidTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	err = repo.ClearProfileImageURL(context.Background(), nil, "moto-123")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestClearProfileImageURL_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles.*SET profile_image_url").
		WithArgs("not-found").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.ClearProfileImageURL(context.Background(), tx, "not-found")
	assert.Error(t, err)
	assert.Equal(t, domain.ErrMotorcycleNotFound, err)
}

func TestClearProfileImageURL_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("SELECT m.id, m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("UPDATE motorcycles")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT r.id, r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS")
	mock.ExpectPrepare("SELECT r.category")
	mock.ExpectPrepare("SELECT r.model")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range")

	repo, err := NewRepository(db)
	assert.NoError(t, err)

	mock.ExpectBegin()
	tx, err := repo.BeginTx(context.Background())
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE motorcycles.*SET profile_image_url").
		WillReturnError(sql.ErrConnDone)

	err = repo.ClearProfileImageURL(context.Background(), tx, "moto-123")
	assert.Error(t, err)
}

// ============================================
// HasServiceHistory Tests (HU45)
// ============================================

func TestHasServiceHistory_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("moto-with-history", "moto-with-history").
		WillReturnRows(rows)

	repo := &repository{db: db}

	has, err := repo.HasServiceHistory(context.Background(), "moto-with-history")

	assert.NoError(t, err)
	assert.True(t, has)
}

func TestHasServiceHistory_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("moto-no-history", "moto-no-history").
		WillReturnRows(rows)

	repo := &repository{db: db}

	has, err := repo.HasServiceHistory(context.Background(), "moto-no-history")

	assert.NoError(t, err)
	assert.False(t, has)
}

func TestHasServiceHistory_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("moto-error", "moto-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	has, err := repo.HasServiceHistory(context.Background(), "moto-error")

	assert.False(t, has)
	assert.Error(t, err)
}

// ============================================
// NewRepository Prepare Error Tests (new statements)
// ============================================

func TestNewRepository_PrepareError_HasServiceHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles.*SET reference_id")
	mock.ExpectPrepare("UPDATE motorcycles.*SET deleted_at")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*ORDER BY b.name")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*WHERE r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS.*completed_services.*diagnostics").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtHasServiceHistory")
}

func TestNewRepository_PrepareError_GetDistinctCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles.*SET reference_id")
	mock.ExpectPrepare("UPDATE motorcycles.*SET deleted_at")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*ORDER BY b.name")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*WHERE r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS.*completed_services.*diagnostics")
	mock.ExpectPrepare("SELECT r.category.*FROM motorcycle_references").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetDistinctCategories")
}

func TestNewRepository_PrepareError_GetLinesByCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles.*SET reference_id")
	mock.ExpectPrepare("UPDATE motorcycles.*SET deleted_at")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*ORDER BY b.name")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*WHERE r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS.*completed_services.*diagnostics")
	mock.ExpectPrepare("SELECT r.category.*FROM motorcycle_references")
	mock.ExpectPrepare("SELECT r.model.*FROM motorcycle_references").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetLinesByCategory")
}

func TestNewRepository_PrepareError_GetDistinctDisplacements(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT m.id, m.license_plate.*WHERE m.id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.owner_id")
	mock.ExpectPrepare("SELECT m.id.*WHERE m.license_plate")
	mock.ExpectPrepare("UPDATE motorcycles.*SET reference_id")
	mock.ExpectPrepare("UPDATE motorcycles.*SET deleted_at")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*ORDER BY b.name")
	mock.ExpectPrepare("SELECT r.id, r.brand_id.*WHERE r.brand_id")
	mock.ExpectPrepare("SELECT EXISTS.*completed_services.*diagnostics")
	mock.ExpectPrepare("SELECT r.category.*FROM motorcycle_references")
	mock.ExpectPrepare("SELECT r.model.*FROM motorcycle_references")
	mock.ExpectPrepare("SELECT DISTINCT r.displacement_range").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetDistinctDisplacements")
}
