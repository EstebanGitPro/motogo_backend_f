package person

import (
	"context"
	"database/sql"

	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
)

func (r *repository) GetPersonByID(ctx context.Context, id string) (*domain.Person, error) {

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var person Person
	err = tx.QueryRowContext(ctx, queryGetByID,
		id,
	).Scan(
		&person.ID,
		&person.IdentityNumber,
		&person.FirstName,
		&person.LastName,
		&person.SecondLastName,
		&person.Email,
		&person.PhoneNumber,
		&person.Role,
		&person.KeycloakUserID,
	)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			tx.Rollback()
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}
	domainPerson := person.ToDomain()
	
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	
	return &domainPerson, nil
}