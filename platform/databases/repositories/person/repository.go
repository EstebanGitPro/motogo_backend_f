package person

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

const (
	querySave       = "INSERT INTO persons (id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	queryGetByEmail = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE email = ? LIMIT 1"
	queryGetByID    = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE id = ? LIMIT 1"
	queryUpdate     = "UPDATE persons SET identity_number = ?, first_name = ?, last_name = ?, second_last_name = ?, email = ?, phone_number = ?, role = ?, keycloak_user_id = ? WHERE id = ?"
	queryDelete     = "DELETE FROM persons WHERE id = ?"
)


type sqlTx struct {
	tx     *sql.Tx
	closed bool 
}
func (t *sqlTx) Commit() error {
	if t.closed {
		panic("sqlTx: commit on closed transaction")
	}
	t.closed = true
	return t.tx.Commit()
}

func (t *sqlTx) Rollback() error {
	if t.closed {
		panic("sqlTx: rollback on closed transaction")
	}
	t.closed = true
	return t.tx.Rollback()
}


func (t *sqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *sqlTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

type repository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}
	return &repository{
		db: db,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTx{tx: tx}, nil
}
