package person

import (
    "database/sql"
    "context"

    "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
    "github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

func (r *repository) GetPersonByEmail(ctx context.Context, tx output.Tx, email string) (*domain.Person, error) {
    var p Person
    err := r.db.QueryRowContext(ctx, queryGetByEmail, email).Scan(
        &p.ID,
        &p.IdentityNumber,
        &p.FirstName,
        &p.LastName,
        &p.SecondLastName,
        &p.Email,
        &p.PhoneNumber,
        &p.Role,
        &p.KeycloakUserID,
    )

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPersonNotFound
		}
		return nil, err
	}

	domainPerson := p.ToDomain()
	return &domainPerson, nil
}