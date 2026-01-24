package common

import (
	"context"
	"database/sql"
)

// SQLTx wraps sql.Tx to implement output.Tx interface
type SQLTx struct {
	tx     *sql.Tx
	closed bool
}

// NewSQLTx creates a new SQLTx wrapper
func NewSQLTx(tx *sql.Tx) *SQLTx {
	return &SQLTx{
		tx:     tx,
		closed: false,
	}
}

// Commit commits the transaction
func (t *SQLTx) Commit() error {
	if t.closed {
		panic("sqlTx: commit on closed transaction")
	}
	t.closed = true
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *SQLTx) Rollback() error {
	if t.closed {
		panic("sqlTx: rollback on closed transaction")
	}
	t.closed = true
	return t.tx.Rollback()
}

// ExecContext executes a query within the transaction
func (t *SQLTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// QueryRowContext queries a single row within the transaction
func (t *SQLTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

// QueryContext executes a query within the transaction (for SELECT ... FOR UPDATE)
func (t *SQLTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

// StmtContext wraps a prepared statement for use within a transaction (for fail-fast pattern)
func (t *SQLTx) StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt {
	return t.tx.StmtContext(ctx, stmt)
}
