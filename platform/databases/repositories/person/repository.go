package person

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
)

const (
	querySave       = "INSERT INTO persons (id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	queryGetByEmail = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE email = ? LIMIT 1"
	queryGetByID    = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE id = ? LIMIT 1"
	queryUpdate     = "UPDATE persons SET identity_number = ?, first_name = ?, last_name = ?, second_last_name = ?, email = ?, phone_number = ?, role = ?, keycloak_user_id = ? WHERE id = ?"
	queryDelete     = "DELETE FROM persons WHERE id = ?"
)

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
	return common.NewSQLTx(tx), nil
}
