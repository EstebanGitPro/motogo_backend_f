package person

import (
	"context"
	"errors"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/core/ports/output"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) SavePerson(ctx context.Context, tx output.Tx, person domain.Person) error {
	personToSave := FromDomain(person)

	dbTx, ok := tx.(*common.SQLTx)
	if !ok {
		return domain.ErrInvalidTransaction
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
		mysqlErr := &mysql.MySQLError{}
		if errors.As(err, &mysqlErr) {
			return domain.ErrDuplicateUser
		}
		return domain.ErrUserCannotSave
	}

	return nil
}
