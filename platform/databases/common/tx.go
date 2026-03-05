package common

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
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

// Commit commits the transaction. It is idempotent: calling Commit on a
// closed transaction returns sql.ErrTxDone instead of panicking.
func (t *SQLTx) Commit() error {
	if t.closed {
		return sql.ErrTxDone
	}
	t.closed = true
	return t.tx.Commit()
}

// Rollback rolls back the transaction. It is idempotent: calling Rollback on a
// closed transaction (e.g. after a successful Commit) returns nil instead of
// panicking — this is safe for the defer tx.Rollback() pattern.
func (t *SQLTx) Rollback() error {
	if t.closed {
		return nil
	}
	t.closed = true
	return t.tx.Rollback()
}

// BeginSQLTx starts a new database transaction and wraps it in an SQLTx.
// Repositories should call this instead of db.BeginTx + NewSQLTx directly.
func BeginSQLTx(ctx context.Context, db *sql.DB) (_ output.Tx, err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	return NewSQLTx(tx), nil
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
