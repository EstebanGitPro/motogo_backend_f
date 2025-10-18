package person

import (
	"database/sql"
	"fmt"

	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	"github.com/go-sql-driver/mysql"
	//"github.com/go-sql-driver/mysql"
)

type repository struct {
	db             *sql.DB
	stmtSave       *sql.Stmt
	stmtGetByEmail *sql.Stmt
	stmtGetByID    *sql.Stmt
	stmtUpdate     *sql.Stmt
}

func NewRepository(db *sql.DB) (ports.Repository, error) {
	stmtSave, err := db.Prepare(querySave)
	if err != nil {
		return nil, fmt.Errorf("error preparing stmtSave: %w", err)
	}

	stmtGetByEmail, err := db.Prepare(queryGetByEmail)
	if err != nil {
		return nil, fmt.Errorf("error preparing stmtGetByEmail: %w", err)
	}

	stmtGetByID, err := db.Prepare(queryGetByID)
	if err != nil {
		return nil, fmt.Errorf("error preparing stmtGetByID: %w", err)
	}

	stmtUpdate, err := db.Prepare(queryUpdate)
	if err != nil {
		return nil, fmt.Errorf("error preparing stmtUpdate: %w", err)
	}

	return &repository{
		db:             db,
		stmtSave:       stmtSave,
		stmtGetByEmail: stmtGetByEmail,
		stmtGetByID:    stmtGetByID,
		stmtUpdate:     stmtUpdate,
	}, nil
}

const (
	querySave       = "INSERT INTO persons (id, identity_number, first_name, last_name, second_last_name, email, phone_number, email_verified, phone_number_verified, password, role, keycloak_user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	queryGetByEmail = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, email_verified, phone_number_verified, password, role, keycloak_user_id FROM persons WHERE email = ? LIMIT 1"
	queryGetByID    = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, email_verified, phone_number_verified, password, role, keycloak_user_id FROM persons WHERE id = ? LIMIT 1"
	queryUpdate     = "UPDATE persons SET identity_number = ?, first_name = ?, last_name = ?, second_last_name = ?, email = ?, phone_number = ?, email_verified = ?, phone_number_verified = ?, password = ?, role = ?, keycloak_user_id = ? WHERE id = ?"
	queryDelete     = "DELETE FROM persons WHERE id = ?"
)

func (r *repository) Save(person domain.Person) error {

	personToSave := Person{
		ID:                  person.ID,
		IdentityNumber:      person.IdentityNumber,
		FirstName:           person.FirstName,
		LastName:            person.LastName,
		SecondLastName:      person.SecondLastName,
		Email:               person.Email,
		PhoneNumber:         person.PhoneNumber,
		EmailVerified:       person.EmailVerified,
		PhoneNumberVerified: person.PhoneNumberVerified,
		Password:            person.Password,
		Role:                person.Role,
		KeycloakUserID:      person.KeycloakUserID,
	}

	_, err := r.stmtSave.Exec(
		personToSave.ID,
		personToSave.IdentityNumber,
		personToSave.FirstName,
		personToSave.LastName,
		personToSave.SecondLastName,
		personToSave.Email,
		personToSave.PhoneNumber,
		personToSave.EmailVerified,
		personToSave.PhoneNumberVerified,
		personToSave.Password,
		personToSave.Role,
		personToSave.KeycloakUserID,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return domain.ErrDuplicateUser
		} else {
			return domain.ErrUserCannotSave
		}

	}

	return nil

}

func (r *repository) GetPersonByEmail(email string) (*domain.Person, error) {
	var p Person
	err := r.stmtGetByEmail.QueryRow(email).Scan(
		&p.ID,
		&p.IdentityNumber,
		&p.FirstName,
		&p.LastName,
		&p.SecondLastName,
		&p.Email,
		&p.PhoneNumber,
		&p.EmailVerified,
		&p.PhoneNumberVerified,
		&p.Password,
		&p.Role,
		&p.KeycloakUserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}
	d := p.ToDomain()
	return &d, nil
}

func (r *repository) GetPersonByID(id string) (*domain.Person, error) {
	var p Person
	err := r.stmtGetByID.QueryRow(id).Scan(
		&p.ID,
		&p.IdentityNumber,
		&p.FirstName,
		&p.LastName,
		&p.SecondLastName,
		&p.Email,
		&p.PhoneNumber,
		&p.EmailVerified,
		&p.PhoneNumberVerified,
		&p.Password,
		&p.Role,
		&p.KeycloakUserID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}
	d := p.ToDomain()
	return &d, nil
}

func (r *repository) Update(person domain.Person) error {
	personToUpdate := Person{
		ID:                  person.ID,
		IdentityNumber:      person.IdentityNumber,
		FirstName:           person.FirstName,
		LastName:            person.LastName,
		SecondLastName:      person.SecondLastName,
		Email:               person.Email,
		PhoneNumber:         person.PhoneNumber,
		EmailVerified:       person.EmailVerified,
		PhoneNumberVerified: person.PhoneNumberVerified,
		Password:            person.Password,
		Role:                person.Role,
		KeycloakUserID:      person.KeycloakUserID,
	}

	_, err := r.stmtUpdate.Exec(
		personToUpdate.IdentityNumber,
		personToUpdate.FirstName,
		personToUpdate.LastName,
		personToUpdate.SecondLastName,
		personToUpdate.Email,
		personToUpdate.PhoneNumber,
		personToUpdate.EmailVerified,
		personToUpdate.PhoneNumberVerified,
		personToUpdate.Password,
		personToUpdate.Role,
		personToUpdate.KeycloakUserID,
		personToUpdate.ID, // WHERE clause
	)
	if err != nil {
		return domain.ErrUserCannotSave
	}

	return nil
}

func (r *repository) Delete(id string) error {
	result, err := r.db.Exec(queryDelete, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPersonNotFound
	}

	return nil
}
