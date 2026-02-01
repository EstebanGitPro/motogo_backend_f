package brand

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

	mock.ExpectPrepare("SELECT id, name FROM brands ORDER BY name")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("SELECT id, name FROM brands ORDER BY name").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetAllBrands")
}

// ============================================
// ValidateBrandIDs Tests
// ============================================

func TestValidateBrandIDs_EmptyList(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.ValidateBrandIDs(context.Background(), []string{})

	assert.NoError(t, err)
}

func TestValidateBrandIDs_AllFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("brand-1").
		AddRow("brand-2")

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-1", "brand-2").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBrandIDs(context.Background(), []string{"brand-1", "brand-2"})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestValidateBrandIDs_SomeNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Only return brand-1, brand-2 is not found
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("brand-1")

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-1", "brand-2").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBrandIDs(context.Background(), []string{"brand-1", "brand-2"})

	assert.Equal(t, domain.ErrBrandNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestValidateBrandIDs_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-1").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	err = repo.ValidateBrandIDs(context.Background(), []string{"brand-1"})

	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestValidateBrandIDs_SingleBrand(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("brand-1")

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-1").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBrandIDs(context.Background(), []string{"brand-1"})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// GetAllBrands Tests
// ============================================

func TestGetAllBrands_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("brand-001", "Yamaha").
		AddRow("brand-002", "Honda").
		AddRow("brand-003", "Suzuki")

	stmt := mock.ExpectPrepare("SELECT id, name FROM brands")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllBrands, _ = db.Prepare("SELECT id, name FROM brands ORDER BY name")

	brands, err := repo.GetAllBrands(context.Background())

	assert.NoError(t, err)
	assert.Len(t, brands, 3)
	assert.Equal(t, "brand-001", brands[0].ID)
	assert.Equal(t, "Yamaha", brands[0].Name)
}

func TestGetAllBrands_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"})

	stmt := mock.ExpectPrepare("SELECT id, name FROM brands")
	stmt.ExpectQuery().WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetAllBrands, _ = db.Prepare("SELECT id, name FROM brands ORDER BY name")

	brands, err := repo.GetAllBrands(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, brands)
}

func TestGetAllBrands_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name FROM brands")
	stmt.ExpectQuery().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetAllBrands, _ = db.Prepare("SELECT id, name FROM brands ORDER BY name")

	brands, err := repo.GetAllBrands(context.Background())

	assert.Nil(t, brands)
	assert.Error(t, err)
}
