package person

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// ============================================
// NewClientRepository Tests
// ============================================

func TestNewClientRepository_NilDB(t *testing.T) {
	// Act
	repo, err := NewClientRepository(nil)

	// Assert
	assert.Nil(t, repo)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestNewClientRepository_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect all prepared statements
	mock.ExpectPrepare("INSERT INTO persons")
	mock.ExpectPrepare("SELECT .* FROM persons WHERE email")
	mock.ExpectPrepare("SELECT .* FROM persons WHERE id")
	mock.ExpectPrepare("SELECT .* FROM persons WHERE keycloak_user_id")
	mock.ExpectPrepare("UPDATE persons SET")
	mock.ExpectPrepare("DELETE FROM persons WHERE id")
	mock.ExpectPrepare("UPDATE persons SET keycloak_user_id")

	// Act
	repo, err := NewClientRepository(db)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// GetPersonByEmail Tests
// ============================================

func TestGetPersonByEmail_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "identity_number", "first_name", "last_name", "second_last_name",
		"email", "phone_number", "role", "keycloak_user_id",
	}).AddRow(
		"person-123", "12345", "John", "Doe", "Smith",
		"test@example.com", "3001234567", "CLIENT", "kc-123",
	)

	mock.ExpectQuery("SELECT .* FROM persons WHERE email").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	person, err := repo.GetPersonByEmail(context.Background(), "test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, "person-123", person.ID)
	assert.Equal(t, "test@example.com", person.Email)
	assert.Equal(t, "John", person.FirstName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPersonByEmail_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM persons WHERE email").
		WithArgs("notfound@example.com").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}

	// Act
	person, err := repo.GetPersonByEmail(context.Background(), "notfound@example.com")

	// Assert
	assert.Nil(t, person)
	assert.Equal(t, domain.ErrPersonNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPersonByEmail_DBError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbError := sql.ErrConnDone
	mock.ExpectQuery("SELECT .* FROM persons WHERE email").
		WithArgs("test@example.com").
		WillReturnError(dbError)

	repo := &repository{db: db}

	// Act
	person, err := repo.GetPersonByEmail(context.Background(), "test@example.com")

	// Assert
	assert.Nil(t, person)
	assert.Equal(t, dbError, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// GetPersonByID Tests
// ============================================

func TestGetPersonByID_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "identity_number", "first_name", "last_name", "second_last_name",
		"email", "phone_number", "role", "keycloak_user_id",
	}).AddRow(
		"person-123", "12345", "John", "Doe", "Smith",
		"test@example.com", "3001234567", "CLIENT", "kc-123",
	)

	mock.ExpectQuery("SELECT .* FROM persons WHERE id").
		WithArgs("person-123").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	person, err := repo.GetPersonByID(context.Background(), "person-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, "person-123", person.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPersonByID_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM persons WHERE id").
		WithArgs("nonexistent-id").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}

	// Act
	person, err := repo.GetPersonByID(context.Background(), "nonexistent-id")

	// Assert
	assert.Nil(t, person)
	assert.Equal(t, domain.ErrPersonNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// GetPersonByKeycloakID Tests
// ============================================

func TestGetPersonByKeycloakID_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create prepared statement mock
	mock.ExpectPrepare("SELECT .* FROM persons WHERE keycloak_user_id")

	stmt, err := db.Prepare(queryGetByKeycloakID)
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{
		"id", "identity_number", "first_name", "last_name", "second_last_name",
		"email", "phone_number", "role", "keycloak_user_id",
	}).AddRow(
		"person-123", "12345", "John", "Doe", "Smith",
		"test@example.com", "3001234567", "CLIENT", "kc-user-123",
	)

	mock.ExpectQuery("SELECT .* FROM persons WHERE keycloak_user_id").
		WithArgs("kc-user-123").
		WillReturnRows(rows)

	repo := &repository{stmtGetByKeycloakID: stmt}

	// Act
	person, err := repo.GetPersonByKeycloakID(context.Background(), "kc-user-123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, person)
	assert.Equal(t, "kc-user-123", person.KeycloakUserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// SavePerson Tests
// ============================================

func TestSavePerson_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO persons").
		WithArgs("person-123", "12345", "John", "Doe", "Smith", "test@example.com", "3001234567", "CLIENT", "kc-123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	person := domain.Person{
		ID:             "person-123",
		IdentityNumber: "12345",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		Email:          "test@example.com",
		PhoneNumber:    "3001234567",
		Role:           "CLIENT",
		KeycloakUserID: "kc-123",
	}

	// Act
	err = repo.SavePerson(context.Background(), sqlTx, person)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSavePerson_DuplicateUser(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO persons").
		WithArgs("person-123", "12345", "John", "Doe", "Smith", "test@example.com", "3001234567", "CLIENT", "kc-123").
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	person := domain.Person{
		ID:             "person-123",
		IdentityNumber: "12345",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		Email:          "test@example.com",
		PhoneNumber:    "3001234567",
		Role:           "CLIENT",
		KeycloakUserID: "kc-123",
	}

	// Act
	err = repo.SavePerson(context.Background(), sqlTx, person)

	// Assert
	assert.Equal(t, domain.ErrDuplicateUser, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSavePerson_GenericError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO persons").
		WithArgs("person-123", "12345", "John", "Doe", "Smith", "test@example.com", "3001234567", "CLIENT", "kc-123").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	person := domain.Person{
		ID:             "person-123",
		IdentityNumber: "12345",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		Email:          "test@example.com",
		PhoneNumber:    "3001234567",
		Role:           "CLIENT",
		KeycloakUserID: "kc-123",
	}

	// Act
	err = repo.SavePerson(context.Background(), sqlTx, person)

	// Assert
	assert.Equal(t, domain.ErrUserCannotSave, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSavePerson_InvalidTransaction(t *testing.T) {
	// Arrange
	repo := &repository{}
	person := domain.Person{ID: "person-123"}

	// Act - pass nil as transaction
	err := repo.SavePerson(context.Background(), nil, person)

	// Assert
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// UpdatePerson Tests
// ============================================

func TestUpdatePerson_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE persons SET").
		WithArgs("12345", "John", "Doe", "Smith", "test@example.com", "3001234567", "CLIENT", "kc-123", "person-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	person := domain.Person{
		ID:             "person-123",
		IdentityNumber: "12345",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		Email:          "test@example.com",
		PhoneNumber:    "3001234567",
		Role:           "CLIENT",
		KeycloakUserID: "kc-123",
	}

	// Act
	err = repo.UpdatePerson(context.Background(), sqlTx, person)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePerson_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE persons SET").
		WithArgs("12345", "John", "Doe", "Smith", "test@example.com", "3001234567", "CLIENT", "kc-123", "person-123").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	person := domain.Person{
		ID:             "person-123",
		IdentityNumber: "12345",
		FirstName:      "John",
		LastName:       "Doe",
		SecondLastName: "Smith",
		Email:          "test@example.com",
		PhoneNumber:    "3001234567",
		Role:           "CLIENT",
		KeycloakUserID: "kc-123",
	}

	// Act
	err = repo.UpdatePerson(context.Background(), sqlTx, person)

	// Assert
	assert.Equal(t, domain.ErrUserCannotSave, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePerson_InvalidTransaction(t *testing.T) {
	// Arrange
	repo := &repository{}
	person := domain.Person{ID: "person-123"}

	// Act
	err := repo.UpdatePerson(context.Background(), nil, person)

	// Assert
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// PatchPerson Tests
// ============================================

func TestPatchPerson_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE persons SET keycloak_user_id").
		WithArgs("kc-new-123", "person-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	// Act
	err = repo.PatchPerson(context.Background(), sqlTx, "person-123", "kc-new-123")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPatchPerson_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE persons SET keycloak_user_id").
		WithArgs("kc-new-123", "nonexistent-id").
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	// Act
	err = repo.PatchPerson(context.Background(), sqlTx, "nonexistent-id", "kc-new-123")

	// Assert
	assert.Equal(t, domain.ErrPersonNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPatchPerson_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE persons SET keycloak_user_id").
		WithArgs("kc-new-123", "person-123").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	// Act
	err = repo.PatchPerson(context.Background(), sqlTx, "person-123", "kc-new-123")

	// Assert
	assert.Equal(t, domain.ErrUserCannotSave, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPatchPerson_InvalidTransaction(t *testing.T) {
	// Arrange
	repo := &repository{}

	// Act
	err := repo.PatchPerson(context.Background(), nil, "person-123", "kc-123")

	// Assert
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// DeletePerson Tests
// ============================================

func TestDeletePerson_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM persons WHERE id").
		WithArgs("person-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	// Act
	err = repo.DeletePerson(context.Background(), sqlTx, "person-123")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePerson_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM persons WHERE id").
		WithArgs("person-123").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	// Act
	err = repo.DeletePerson(context.Background(), sqlTx, "person-123")

	// Assert
	assert.Equal(t, domain.ErrUserCannotDelete, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePerson_InvalidTransaction(t *testing.T) {
	// Arrange
	repo := &repository{}

	// Act
	err := repo.DeletePerson(context.Background(), nil, "person-123")

	// Assert
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// BeginTx Tests
// ============================================

func TestBeginTx_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()

	repo := &repository{db: db}

	// Act
	tx, err := repo.BeginTx(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBeginTx_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	// Act
	tx, err := repo.BeginTx(context.Background())

	// Assert
	assert.Nil(t, tx)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
