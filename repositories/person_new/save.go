package personnew

import (
	"context"
	"fmt"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/go-sql-driver/mysql"
)

func (r *repository) SavePerson(ctx context.Context, person *domain.Person) error {

	personToSave := FromDomain(*person)

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
		keycloakUserID, err := r.keycloak.CreateUser(ctx, person)

		if err != nil {

		existingUser, getErr := r.keycloak.GetUserByEmail(ctx, person.Email)
		if getErr != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create or get user in keycloak: %w", err)
		}
		keycloakUserID = *existingUser.ID
		}
	
		

		person.KeycloakUserID = keycloakUserID
		err = r.UpdatePerson(ctx, person)
		if err != nil {
			_ = r.keycloak.DeleteUser(ctx, keycloakUserID)
			tx.Rollback()
			return fmt.Errorf("failed to update person with keycloak user id: %w", err)
		}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}


