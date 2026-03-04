package person

import (
	"context"
	"database/sql"

	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/EstebanGitPro/motogo-backend/platform/logger"
)

const (
	querySave            = "INSERT INTO persons (id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	queryGetByEmail      = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE email = ? LIMIT 1"
	queryGetByID         = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE id = ? LIMIT 1"
	queryGetByKeycloakID = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE keycloak_user_id = ? LIMIT 1"
	queryUpdate          = "UPDATE persons SET identity_number = ?, first_name = ?, last_name = ?, second_last_name = ?, email = ?, phone_number = ?, role = ?, keycloak_user_id = ? WHERE id = ?"
	queryDelete          = "DELETE FROM persons WHERE id = ?"
	queryPatch           = "UPDATE persons SET keycloak_user_id = ? WHERE id = ?"
)

var log logger.Logger = logger.NewSlogLogger()

type repository struct {
	stmtSave            *sql.Stmt
	stmtGetByEmail      *sql.Stmt
	stmtGetByID         *sql.Stmt
	stmtGetByKeycloakID *sql.Stmt
	stmtUpdate          *sql.Stmt
	stmtDelete          *sql.Stmt
	stmtPatch           *sql.Stmt
	db                  *sql.DB
}

func NewClientRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}

	stmtSave, err := db.Prepare(querySave)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtSave", err)
		return nil, err
	}

	stmtGetByEmail, err := db.Prepare(queryGetByEmail)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByEmail", err)
		return nil, err
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByID", err)
		return nil, err
	}

	stmtGetByKeycloakID, err := db.Prepare(queryGetByKeycloakID)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtGetByKeycloakID", err)
		return nil, err
	}

	stmtUpdate, err := db.Prepare(queryUpdate)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtUpdate", err)
		return nil, err
	}

	stmtDelete, err := db.Prepare(queryDelete)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtDelete", err)
		return nil, err
	}

	stmtPatch, err := db.Prepare(queryPatch)
	if err != nil {
		log.Error(logger.LogDatabaseUnavailable, "error preparing stmtPatch", err)
		return nil, err
	}

	return &repository{
		db:                  db,
		stmtSave:            stmtSave,
		stmtGetByEmail:      stmtGetByEmail,
		stmtGetByID:         stmtGetByID,
		stmtGetByKeycloakID: stmtGetByKeycloakID,
		stmtUpdate:          stmtUpdate,
		stmtDelete:          stmtDelete,
		stmtPatch:           stmtPatch,
	}, nil
}

func (r *repository) BeginTx(ctx context.Context) (output.Tx, error) {
	return common.BeginSQLTx(ctx, r.db)
}
