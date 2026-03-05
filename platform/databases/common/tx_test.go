package common

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewSQLTx Tests
// ============================================

func TestNewSQLTx_CreatesWrapper(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)

	assert.NotNil(t, sqlTx)
	assert.False(t, sqlTx.closed)
}

// ============================================
// Commit Tests
// ============================================

func TestSQLTx_Commit_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)

	err = sqlTx.Commit()

	assert.NoError(t, err)
	assert.True(t, sqlTx.closed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLTx_Commit_ReturnsTxDone_WhenClosed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)
	_ = sqlTx.Commit() // First commit

	// Second commit should return sql.ErrTxDone (idempotent)
	err = sqlTx.Commit()
	assert.ErrorIs(t, err, sql.ErrTxDone)
}

// ============================================
// Rollback Tests
// ============================================

func TestSQLTx_Rollback_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)

	err = sqlTx.Rollback()

	assert.NoError(t, err)
	assert.True(t, sqlTx.closed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSQLTx_Rollback_ReturnsNil_WhenClosed(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)
	_ = sqlTx.Commit() // Commit first

	// Rollback after commit should return nil (idempotent — safe for defer)
	err = sqlTx.Rollback()
	assert.NoError(t, err)
}

// ============================================
// BeginSQLTx Tests
// ============================================

func TestBeginSQLTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	tx, err := BeginSQLTx(context.Background(), db)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBeginSQLTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(fmt.Errorf("connection lost"))

	tx, err := BeginSQLTx(context.Background(), db)

	assert.Error(t, err)
	assert.Nil(t, tx)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// ExecContext Tests
// ============================================

func TestSQLTx_ExecContext_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test").
		WithArgs("value1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)

	result, err := sqlTx.ExecContext(context.Background(), "INSERT INTO test VALUES (?)", "value1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	affected, _ := result.RowsAffected()
	assert.Equal(t, int64(1), affected)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// QueryRowContext Tests
// ============================================

func TestSQLTx_QueryRowContext_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("1", "test")
	mock.ExpectQuery("SELECT id, name FROM test").
		WillReturnRows(rows)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)

	row := sqlTx.QueryRowContext(context.Background(), "SELECT id, name FROM test")

	var id, name string
	err = row.Scan(&id, &name)
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
	assert.Equal(t, "test", name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// QueryContext Tests
// ============================================

func TestSQLTx_QueryContext_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow("1").
		AddRow("2")
	mock.ExpectQuery("SELECT id FROM test").
		WillReturnRows(rows)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := NewSQLTx(tx)

	resultRows, err := sqlTx.QueryContext(context.Background(), "SELECT id FROM test")

	assert.NoError(t, err)
	assert.NotNil(t, resultRows)
	defer resultRows.Close()

	count := 0
	for resultRows.Next() {
		count++
	}
	assert.Equal(t, 2, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
