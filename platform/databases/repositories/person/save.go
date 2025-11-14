package person

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) SavePerson(ctx context.Context, tx output.Tx, person domain.Person) error {
	personToSave := FromDomain(person)

	var dbTx *sqlTx
	var shouldCommit bool

	if tx != nil {
		// Usar la transacción existente
		dbTx = tx.(*sqlTx)
		shouldCommit = false
	} else {
		// Crear nueva transacción
		newTx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		dbTx = &sqlTx{Tx: newTx}
		shouldCommit = true
	}

	_, err := dbTx.ExecContext(ctx, querySave,
		personToSave.ID,
		personToSave.IdentityNumber,
		personToSave.FirstName,
		personToSave.LastName,
		personToSave.SecondLastName,
		personToSave.Email,
		personToSave.PhoneNumber,
		personToSave.Role,
		personToSave.KeycloakUserID)

	if err != nil {
		if shouldCommit {
			dbTx.Rollback()
		}

		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return domain.ErrDuplicateUser
		}
		return domain.ErrUserCannotSave
	}

	if shouldCommit {
		if err = dbTx.Commit(); err != nil {
			return domain.ErrUserCannotSave
		}
	}

	return nil
}
