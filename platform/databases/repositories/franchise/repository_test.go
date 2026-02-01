package franchise

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
	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises")
	mock.ExpectPrepare("DELETE FROM franchises")
	mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE id")
	mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE name")
	mock.ExpectPrepare("SELECT DISTINCT f.id, f.name, f.description")
	mock.ExpectPrepare("SELECT COUNT")
	mock.ExpectPrepare("UPDATE branches SET franchise_id")
	mock.ExpectPrepare("UPDATE branches SET franchise_id = NULL")
	mock.ExpectPrepare("UPDATE branches SET franchise_id = NULL")

	repo, err := NewRepository(db)

	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtSaveFranchise")
}

func TestNewRepository_PrepareError_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtUpdateFranchise")
}

func TestNewRepository_PrepareError_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises")
	mock.ExpectPrepare("DELETE FROM franchises").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtDeleteFranchise")
}

func TestNewRepository_PrepareError_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises")
	mock.ExpectPrepare("DELETE FROM franchises")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetFranchiseByID")
}

func TestNewRepository_PrepareError_GetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises")
	mock.ExpectPrepare("DELETE FROM franchises")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE id")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE name").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetFranchiseByName")
}

func TestNewRepository_PrepareError_GetByRep(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises")
	mock.ExpectPrepare("DELETE FROM franchises")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE id")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE name")
	mock.ExpectPrepare("SELECT.*franchises.*INNER JOIN branches").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtGetFranchisesByRepresentative")
}

func TestNewRepository_PrepareError_CountBranches(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises")
	mock.ExpectPrepare("DELETE FROM franchises")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE id")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE name")
	mock.ExpectPrepare("SELECT.*franchises.*INNER JOIN branches")
	mock.ExpectPrepare("SELECT COUNT.*FROM branches WHERE franchise_id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtCountBranchesByFranchise")
}

func TestNewRepository_PrepareError_AssociateBranch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO franchises")
	mock.ExpectPrepare("UPDATE franchises")
	mock.ExpectPrepare("DELETE FROM franchises")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE id")
	mock.ExpectPrepare("SELECT.*FROM franchises.*WHERE name")
	mock.ExpectPrepare("SELECT.*franchises.*INNER JOIN branches")
	mock.ExpectPrepare("SELECT COUNT.*FROM branches WHERE franchise_id")
	mock.ExpectPrepare("UPDATE branches SET franchise_id.*WHERE id").
		WillReturnError(sql.ErrConnDone)

	repo, err := NewRepository(db)

	assert.Nil(t, repo)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error preparing stmtAssociateBranchToFranchise")
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
// GetFranchiseByID Tests
// ============================================

func TestGetFranchiseByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	description := "Test Description"
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow("franchise-123", "Franquicia Test", description)

	stmt := mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE id")
	stmt.ExpectQuery().
		WithArgs("franchise-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetFranchiseByID, _ = db.Prepare("SELECT id, name, description FROM franchises WHERE id = ?")

	franchise, err := repo.GetFranchiseByID(context.Background(), "franchise-123")

	assert.NoError(t, err)
	assert.NotNil(t, franchise)
	assert.Equal(t, "franchise-123", franchise.ID)
	assert.Equal(t, "Franquicia Test", franchise.Name)
	assert.Equal(t, description, *franchise.Description)
}

func TestGetFranchiseByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE id")
	stmt.ExpectQuery().
		WithArgs("not-found").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetFranchiseByID, _ = db.Prepare("SELECT id, name, description FROM franchises WHERE id = ?")

	franchise, err := repo.GetFranchiseByID(context.Background(), "not-found")

	assert.Nil(t, franchise)
	assert.Equal(t, domain.ErrFranchiseNotFound, err)
}

func TestGetFranchiseByID_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE id")
	stmt.ExpectQuery().
		WithArgs("franchise-123").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetFranchiseByID, _ = db.Prepare("SELECT id, name, description FROM franchises WHERE id = ?")

	franchise, err := repo.GetFranchiseByID(context.Background(), "franchise-123")

	assert.Nil(t, franchise)
	assert.Equal(t, sql.ErrConnDone, err)
}

// ============================================
// GetFranchiseByName Tests
// ============================================

func TestGetFranchiseByName_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow("franchise-456", "Super Franquicia", nil)

	stmt := mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE name")
	stmt.ExpectQuery().
		WithArgs("Super Franquicia").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetFranchiseByName, _ = db.Prepare("SELECT id, name, description FROM franchises WHERE name = ?")

	franchise, err := repo.GetFranchiseByName(context.Background(), "Super Franquicia")

	assert.NoError(t, err)
	assert.NotNil(t, franchise)
	assert.Equal(t, "franchise-456", franchise.ID)
	assert.Equal(t, "Super Franquicia", franchise.Name)
	assert.Nil(t, franchise.Description)
}

func TestGetFranchiseByName_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE name")
	stmt.ExpectQuery().
		WithArgs("No Existe").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}
	repo.stmtGetFranchiseByName, _ = db.Prepare("SELECT id, name, description FROM franchises WHERE name = ?")

	franchise, err := repo.GetFranchiseByName(context.Background(), "No Existe")

	// GetFranchiseByName returns nil, nil for not found (for duplicate checking)
	assert.Nil(t, franchise)
	assert.Nil(t, err)
}

func TestGetFranchiseByName_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT id, name, description FROM franchises WHERE name")
	stmt.ExpectQuery().
		WithArgs("Error Test").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetFranchiseByName, _ = db.Prepare("SELECT id, name, description FROM franchises WHERE name = ?")

	franchise, err := repo.GetFranchiseByName(context.Background(), "Error Test")

	assert.Nil(t, franchise)
	assert.Error(t, err)
}

// ============================================
// CountBranchesByFranchise Tests
// ============================================

func TestCountBranchesByFranchise_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("franchise-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCountBranchesByFranchise, _ = db.Prepare("SELECT COUNT(*) FROM branches WHERE franchise_id = ?")

	count, err := repo.CountBranchesByFranchise(context.Background(), "franchise-123")

	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestCountBranchesByFranchise_ZeroCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("franchise-789").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtCountBranchesByFranchise, _ = db.Prepare("SELECT COUNT(*) FROM branches WHERE franchise_id = ?")

	count, err := repo.CountBranchesByFranchise(context.Background(), "franchise-789")

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestCountBranchesByFranchise_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT COUNT")
	stmt.ExpectQuery().
		WithArgs("franchise-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtCountBranchesByFranchise, _ = db.Prepare("SELECT COUNT(*) FROM branches WHERE franchise_id = ?")

	count, err := repo.CountBranchesByFranchise(context.Background(), "franchise-error")

	assert.Equal(t, 0, count)
	assert.Equal(t, sql.ErrConnDone, err)
}

// ============================================
// GetFranchisesByRepresentative Tests
// ============================================

func TestGetFranchisesByRepresentative_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	desc := "Descripción de prueba"
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow("franchise-001", "Franquicia 1", desc).
		AddRow("franchise-002", "Franquicia 2", nil)

	stmt := mock.ExpectPrepare("SELECT DISTINCT f.id, f.name")
	stmt.ExpectQuery().
		WithArgs("rep-123").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetFranchisesByRepresentative, _ = db.Prepare("SELECT DISTINCT f.id, f.name, f.description FROM franchises f WHERE rep_id = ?")

	franchises, err := repo.GetFranchisesByRepresentative(context.Background(), "rep-123")

	assert.NoError(t, err)
	assert.Len(t, franchises, 2)
	assert.Equal(t, "franchise-001", franchises[0].ID)
	assert.Equal(t, "Franquicia 1", franchises[0].Name)
	assert.NotNil(t, franchises[0].Description)
	assert.Equal(t, desc, *franchises[0].Description)
	assert.Nil(t, franchises[1].Description)
}

func TestGetFranchisesByRepresentative_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "description"})

	stmt := mock.ExpectPrepare("SELECT DISTINCT f.id, f.name")
	stmt.ExpectQuery().
		WithArgs("rep-empty").
		WillReturnRows(rows)

	repo := &repository{db: db}
	repo.stmtGetFranchisesByRepresentative, _ = db.Prepare("SELECT DISTINCT f.id, f.name, f.description FROM franchises f WHERE rep_id = ?")

	franchises, err := repo.GetFranchisesByRepresentative(context.Background(), "rep-empty")

	assert.NoError(t, err)
	assert.Empty(t, franchises)
}

func TestGetFranchisesByRepresentative_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	stmt := mock.ExpectPrepare("SELECT DISTINCT f.id, f.name")
	stmt.ExpectQuery().
		WithArgs("rep-error").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}
	repo.stmtGetFranchisesByRepresentative, _ = db.Prepare("SELECT DISTINCT f.id, f.name, f.description FROM franchises f WHERE rep_id = ?")

	franchises, err := repo.GetFranchisesByRepresentative(context.Background(), "rep-error")

	assert.Nil(t, franchises)
	assert.Error(t, err)
}

// ============================================
// SaveFranchise Tests
// ============================================

func TestSaveFranchise_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO franchises").
		WithArgs("franchise-123", "Test Franchise", nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	franchise := domain.Franchise{ID: "franchise-123", Name: "Test Franchise"}
	err = repo.SaveFranchise(context.Background(), sqlTx, franchise)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveFranchise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO franchises").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.SaveFranchise(context.Background(), sqlTx, domain.Franchise{})
	assert.Error(t, err)
}

func TestSaveFranchise_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.SaveFranchise(context.Background(), nil, domain.Franchise{})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// UpdateFranchise Tests
// ============================================

func TestUpdateFranchise_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE franchises").
		WithArgs("Updated Name", nil, "franchise-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	franchise := domain.Franchise{ID: "franchise-123", Name: "Updated Name"}
	err = repo.UpdateFranchise(context.Background(), sqlTx, franchise)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateFranchise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE franchises").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.UpdateFranchise(context.Background(), sqlTx, domain.Franchise{})
	assert.Error(t, err)
}

func TestUpdateFranchise_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.UpdateFranchise(context.Background(), nil, domain.Franchise{})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DeleteFranchise Tests
// ============================================

func TestDeleteFranchise_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM franchises").
		WithArgs("franchise-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteFranchise(context.Background(), sqlTx, "franchise-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteFranchise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM franchises").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DeleteFranchise(context.Background(), sqlTx, "franchise-err")
	assert.Error(t, err)
}

func TestDeleteFranchise_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.DeleteFranchise(context.Background(), nil, "franchise-123")
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// AssociateBranchesToFranchise Tests
// ============================================

func TestAssociateBranchesToFranchise_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches SET franchise_id").
		WithArgs("franchise-123", "branch-1").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE branches SET franchise_id").
		WithArgs("franchise-123", "branch-2").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.AssociateBranchesToFranchise(context.Background(), sqlTx, "franchise-123", []string{"branch-1", "branch-2"})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAssociateBranchesToFranchise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches SET franchise_id").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.AssociateBranchesToFranchise(context.Background(), sqlTx, "franchise-123", []string{"branch-1"})
	assert.Error(t, err)
}

func TestAssociateBranchesToFranchise_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.AssociateBranchesToFranchise(context.Background(), nil, "franchise-123", []string{"branch-1"})
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DissociateBranchesFromFranchise Tests
// ============================================

func TestDissociateBranchesFromFranchise_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches SET franchise_id = NULL").
		WithArgs("franchise-123").
		WillReturnResult(sqlmock.NewResult(0, 3))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DissociateBranchesFromFranchise(context.Background(), sqlTx, "franchise-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDissociateBranchesFromFranchise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches SET franchise_id = NULL").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DissociateBranchesFromFranchise(context.Background(), sqlTx, "franchise-err")
	assert.Error(t, err)
}

func TestDissociateBranchesFromFranchise_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.DissociateBranchesFromFranchise(context.Background(), nil, "franchise-123")
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DissociateSingleBranch Tests
// ============================================

func TestDissociateSingleBranch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches SET franchise_id = NULL").
		WithArgs("branch-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DissociateSingleBranch(context.Background(), sqlTx, "branch-123")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDissociateSingleBranch_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE branches SET franchise_id = NULL").WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	err = repo.DissociateSingleBranch(context.Background(), sqlTx, "branch-err")
	assert.Error(t, err)
}

func TestDissociateSingleBranch_InvalidTransaction(t *testing.T) {
	repo := &repository{}
	err := repo.DissociateSingleBranch(context.Background(), nil, "branch-123")
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}
