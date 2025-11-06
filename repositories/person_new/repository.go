package personnew

import (
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

type repository struct {
	keycloak output.AuthClient
	db *sql.DB
}

func NewClient(db *sql.DB, keycloak output.AuthClient) (*repository, error) {
	return &repository{
		keycloak: keycloak,
		db: db,
	}, nil
}