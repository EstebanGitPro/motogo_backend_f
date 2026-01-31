package branch

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
	mock.ExpectPrepare("INSERT INTO branches")
	mock.ExpectPrepare("UPDATE branches")
	mock.ExpectPrepare("DELETE FROM branches")
	mock.ExpectPrepare("SELECT b.id, b.representative_id.*WHERE b.id")
	mock.ExpectPrepare("SELECT id, representative_id.*WHERE franchise_id")
	mock.ExpectPrepare("SELECT b.id, b.representative_id.*WHERE b.representative_id")
	mock.ExpectPrepare("INSERT INTO branch_brands")
	mock.ExpectPrepare("DELETE FROM branch_brands")
	mock.ExpectPrepare("SELECT brand_id FROM branch_brands")
	mock.ExpectPrepare("SELECT.*FROM branches b.*INNER JOIN locations l")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError_SaveBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO branches").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtSaveBranch")
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
// GetBranchByID Tests
// ============================================

func TestGetBranchByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetBranchByID, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches b WHERE b.id = ?")

	branch, err := repo.GetBranchByID(context.Background(), "not-found")

	assert.Nil(t, branch)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestGetBranchByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT b.id, b.representative_id")
	stmt.ExpectQuery().
		WithArgs("branch-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetBranchByID, _ = db.Prepare("SELECT b.id, b.representative_id FROM branches b WHERE b.id = ?")

	branch, err := repo.GetBranchByID(context.Background(), "branch-error")

	assert.Nil(t, branch)
	assert.Error(t, err)
}

// ============================================
// HasBranchesByRepresentative Tests
// ============================================

func TestHasBranchesByRepresentative_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("rep-123").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.HasBranchesByRepresentative(context.Background(), "rep-123")

	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestHasBranchesByRepresentative_False(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("rep-no-branches").
		WillReturnRows(rows)

	repo := &repository{db: db}

	exists, err := repo.HasBranchesByRepresentative(context.Background(), "rep-no-branches")

	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestHasBranchesByRepresentative_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("rep-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	exists, err := repo.HasBranchesByRepresentative(context.Background(), "rep-error")

	assert.False(t, exists)
	assert.Error(t, err)
}

// ============================================
// ValidateBrands Tests
// ============================================

func TestValidateBrands_EmptySlice(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{})

	assert.NoError(t, err)
}

func TestValidateBrands_AllFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("brand-001").
		AddRow("brand-002")

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-001", "brand-002").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{"brand-001", "brand-002"})

	assert.NoError(t, err)
}

func TestValidateBrands_SomeNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("brand-001")

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-001", "brand-notfound").
		WillReturnRows(rows)

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{"brand-001", "brand-notfound"})

	assert.Equal(t, domain.ErrBrandNotFound, err)
}

func TestValidateBrands_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT id FROM brands WHERE id IN").
		WithArgs("brand-001").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	err = repo.ValidateBrands(context.Background(), []string{"brand-001"})

	assert.Error(t, err)
}

// ============================================
// GetBranchByFranchiseAndName Tests
// ============================================

func TestGetBranchByFranchiseAndName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "representative_id", "franchise_id", "name", "establishment_type", "profile_image_url", "status",
	}).AddRow("branch-001", "rep-123", "franchise-001", "Sucursal Norte", "WORKSHOP", "http://example.com/img.jpg", "ACTIVE")

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("franchise-001", "Sucursal Norte").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetBranchByFranchiseAndName, _ = db.Prepare("SELECT * FROM branches WHERE franchise_id = ? AND name = ?")

	branch, err := repo.GetBranchByFranchiseAndName(context.Background(), "franchise-001", "Sucursal Norte")

	assert.NoError(t, err)
	assert.NotNil(t, branch)
	assert.Equal(t, "branch-001", branch.ID)
	assert.Equal(t, "Sucursal Norte", branch.Name)
}

func TestGetBranchByFranchiseAndName_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("franchise-001", "No Existe").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetBranchByFranchiseAndName, _ = db.Prepare("SELECT * FROM branches WHERE franchise_id = ? AND name = ?")

	branch, err := repo.GetBranchByFranchiseAndName(context.Background(), "franchise-001", "No Existe")

	assert.Nil(t, branch)
	assert.Equal(t, domain.ErrBranchNotFound, err)
}

func TestGetBranchByFranchiseAndName_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT")
	stmt.ExpectQuery().
		WithArgs("franchise-error", "Test").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetBranchByFranchiseAndName, _ = db.Prepare("SELECT * FROM branches WHERE franchise_id = ? AND name = ?")

	branch, err := repo.GetBranchByFranchiseAndName(context.Background(), "franchise-error", "Test")

	assert.Nil(t, branch)
	assert.Error(t, err)
}
