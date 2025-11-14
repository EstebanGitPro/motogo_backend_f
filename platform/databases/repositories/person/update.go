
package person

import (
    "context"

    "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
    "github.com/EstebanGitPro/motogo-backend/core/ports/output"
)

func (r *repository) UpdatePerson(ctx context.Context, tx output.Tx, person domain.Person) error {
    personToUpdate := FromDomain(person)

    dbTx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    _, err = dbTx.ExecContext(ctx, queryUpdate,
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
        dbTx.Rollback()
        return domain.ErrUserCannotSave
    }

    err = dbTx.Commit()
    if err != nil {
        return err
    }

    return nil
}