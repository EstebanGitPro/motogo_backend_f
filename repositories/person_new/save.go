package personnew

import (
	"context"

	"github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/Nerzal/gocloak/v13"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) SavePerson(ctx context.Context, person *domain.Person) error {

	personToSave := FromDomain(person)

	// begin transaction
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, querySave,
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

	// TODO logica de keycloak
	user := gocloak.User{
		Email:         &person.Email,
		FirstName:     &person.FirstName,
		LastName:      &person.LastName,
		Enabled:       gocloak.BoolP(true),
		Username:      &person.Email,
	}
	_, err = r.keycloak.CreateUser(
		ctx,
		user,
		r.config.Keycloak.Realm,
		
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	person.KeycloakUserID = user.ID
	


	// commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}