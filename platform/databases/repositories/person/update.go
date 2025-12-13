package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
)

func (r *repository) UpdatePerson(ctx context.Context, tx output.Tx, person domain.Person) error {
	personToUpdate := FromDomain(person)

	// Type assertion segura
	dbTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
	}

	// Solo ejecutar el query - NO manejar commit/rollback
	_, err := dbTx.ExecContext(ctx, queryUpdate,
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
		return domain.ErrUserCannotSave
	}

	return nil
}
