package person

import (
	"database/sql"

	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports"
	mysql "github.com/go-sql-driver/mysql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) ports.Repository {
	return &repository{
		db: db,
	}
}

const (
	querySave       = "INSERT INTO persons (id, identity_number, first_name, last_name, second_last_name, email, phone_number, email_verified, phone_number_verified, password, role) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	queryGetByEmail = "SELECT id, identity_number, first_name, last_name, second_last_name, email, phone_number, email_verified, phone_number_verified, password, role FROM persons WHERE email = ? LIMIT 1"
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
	}

	stmt, err := r.db.Prepare(querySave)
	if err != nil {
		return domain.ErrUserCannotSave
	}
	defer stmt.Close()

	_, err = stmt.Exec(
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
	err := r.db.QueryRow(queryGetByEmail, email).Scan(
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