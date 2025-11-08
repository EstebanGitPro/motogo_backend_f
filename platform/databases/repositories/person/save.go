package person

import (
	"context"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) SavePerson(person domain.Person) error {

	personToSave := FromDomain(person)

	// begin transaction
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(context.Background(), querySave,
		personToSave.ID,
		personToSave.IdentityNumber,
		personToSave.FirstName,
		 personToSave.LastName,
		personToSave.SecondLastName,
		personToSave.Email,
		personToSave.PhoneNumber,
		personToSave.Role,
		personToSave.KeycloakUserID, )

	if err != nil {
		tx.Rollback()
		// TODO log error
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return domain.ErrDuplicateUser
		} else {
			return domain.ErrUserCannotSave
		}

	}

		
	return nil

}


