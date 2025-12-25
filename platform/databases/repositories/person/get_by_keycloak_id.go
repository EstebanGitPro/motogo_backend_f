package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) GetPersonByKeycloakID(ctx context.Context, keycloakUserID string) (*domain.Person, error) {
	var p Person
	err := r.stmtGetByKeycloakID.QueryRowContext(ctx, keycloakUserID).Scan(
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
		return nil, err
	}

	d := p.ToDomain()
	return &d, nil
}
