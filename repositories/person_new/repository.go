package personnew

import (
	"database/sql"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/ports"
)

// TODO: Terminar de pasar los metodos restantes y revisar conecion a base de datos.
// TODO: Manejo de conexiones: Apertura y cierre
// TODO: Manejo de transacciones: Commit y Rollback
// TODO: Manejo de recursos de bases de datos: Apertura y cierre
// TODO: Manejo de patch para actualizar solo el ID de keycloak en db de negocio
type repository struct {
	db             *sql.DB
	StmtSave       *sql.Stmt
	StmtGetByEmail *sql.Stmt
	StmtGetByID    *sql.Stmt
	StmtUpdate     *sql.Stmt
	StmtDelete     *sql.Stmt
	keycloak ports.AuthClient
}


const (
	querySave       = "INSERT INTO persons (id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	queryGetByEmail = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE email = ? LIMIT 1"
	queryGetByID    = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, role, keycloak_user_id FROM persons WHERE id = ? LIMIT 1"
	queryUpdate     = "UPDATE persons SET identity_number = ?, first_name = ?, last_name = ?, second_last_name = ?, email = ?, phone_number = ?, role = ?, keycloak_user_id = ? WHERE id = ?"
	queryDelete     = "DELETE FROM persons WHERE id = ?"
)


func NewClient(db *sql.DB, keycloak ports.AuthClient) (*repository, error) {
	repo := &repository{
		db:       db,
		keycloak: keycloak,
	}
	
	var err error
	
	repo.StmtSave, err = db.Prepare(querySave)
	if err != nil {
		repo.Close() 
		return nil, fmt.Errorf("error preparing stmtSave: %w", err)
	}

	repo.StmtGetByEmail, err = db.Prepare(queryGetByEmail)
	if err != nil {
		repo.Close() 
		return nil, fmt.Errorf("error preparing stmtGetByEmail: %w", err)
	}

	repo.StmtGetByID, err = db.Prepare(queryGetByID)
	if err != nil {
		repo.Close() 
		return nil, fmt.Errorf("error preparing stmtGetByID: %w", err)
	}

	repo.StmtUpdate, err = db.Prepare(queryUpdate)
	if err != nil {
		repo.Close()
		return nil, fmt.Errorf("error preparing stmtUpdate: %w", err)
	}
	
	repo.StmtDelete, err = db.Prepare(queryDelete)
	if err != nil {
		repo.Close()
		return nil, fmt.Errorf("error preparing stmtDelete: %w", err)
	}
	
	return repo, nil
}



func (r *repository) Close() error {
	var firstErr error

	if r.StmtSave != nil {
		if err := r.StmtSave.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if r.StmtGetByEmail != nil {
		if err := r.StmtGetByEmail.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if r.StmtGetByID != nil {
		if err := r.StmtGetByID.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if r.StmtUpdate != nil {
		if err := r.StmtUpdate.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if r.StmtDelete != nil {
		if err := r.StmtDelete.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}



