package personnew

import (
	"context"
	"database/sql"

	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
)

func (r *repository) GetPersonByEmail(ctx context.Context, email string) (*domain.Person, error) {

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var person Person
	err = tx.QueryRowContext(ctx, queryGetByEmail,
		email,
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
		if err == sql.ErrNoRows {
			tx.Rollback()
			return nil, domain.ErrPersonNotFound
		}
		tx.Rollback()
		return nil, err
	}
	domainPerson := person.ToDomain()
	return &domainPerson, nil
}



