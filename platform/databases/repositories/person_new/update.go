
package personnew

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
)

func (r *repository) UpdatePerson(ctx context.Context, person *domain.Person) error {
	personToUpdate := FromDomain(*person)

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, queryUpdate,
		personToUpdate.IdentityNumber,
		personToUpdate.FirstName,
		personToUpdate.LastName,
		personToUpdate.SecondLastName,
		personToUpdate.Email,
		personToUpdate.PhoneNumber,
		personToUpdate.Role,
		personToUpdate.KeycloakUserID,
		personToUpdate.ID, // WHERE clause
	)

	if err != nil {
		tx.Rollback()
		return domain.ErrUserCannotSave
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}