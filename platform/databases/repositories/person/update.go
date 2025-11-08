
package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
)

func (r *repository) UpdatePerson(person domain.Person) error {
	personToUpdate := FromDomain(person)

	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(context.Background(), queryUpdate,
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