package message

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/EstebanGitPro/motogo-backend/platform/databases/common"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageRepository_NilDB(t *testing.T) {
	// Act
	repo, err := NewMessageRepository(nil)

	// Assert
	assert.Nil(t, repo)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestNewMessageRepository_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expect all prepared statements
	mock.ExpectPrepare("SELECT .* FROM system_messages WHERE is_active = true")
	mock.ExpectPrepare("SELECT .* FROM system_messages WHERE message_code")
	mock.ExpectPrepare("SELECT .* FROM system_messages WHERE message_code .* AND is_active = true")
	mock.ExpectPrepare("SELECT .* FROM system_messages WHERE id")
	mock.ExpectPrepare("SELECT .* FROM system_messages WHERE type")
	mock.ExpectPrepare("SELECT .* FROM system_messages WHERE module")
	mock.ExpectPrepare("INSERT INTO system_messages")
	mock.ExpectPrepare("DELETE FROM system_messages WHERE id")
	mock.ExpectPrepare("SELECT .* FROM system_messages WHERE message_code")

	// Act
	repo, err := NewMessageRepository(db)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllActive_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	}).AddRow(
		"msg-1", "ERR_001", "ERROR", "USER", "MOD_U",
		"Error Title", "Error content", true, now, now,
	).AddRow(
		"msg-2", "SUC_001", "EXITO", "SYSTEM", "MOD_S",
		"Success Title", "Success content", true, now, now,
	)

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE is_active = true").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetAllActive(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.Equal(t, "msg-1", messages[0].ID)
	assert.Equal(t, "msg-2", messages[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllActive_Empty(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE is_active = true").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetAllActive(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllActive_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE is_active = true").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetAllActive(context.Background())

	// Assert
	assert.Nil(t, messages)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByCode_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	}).AddRow(
		"msg-1", "ERR_001", "ERROR", "USER", "MOD_U",
		"Error Title", "Error content", true, now, now,
	)

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE message_code").
		WithArgs("ERR_001").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	message, err := repo.GetByCode(context.Background(), "ERR_001")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "msg-1", message.ID)
	assert.Equal(t, "ERR_001", message.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByCode_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE message_code").
		WithArgs("NONEXISTENT").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}

	// Act
	message, err := repo.GetByCode(context.Background(), "NONEXISTENT")

	// Assert
	assert.Nil(t, message)
	assert.Nil(t, err) // GetByCode returns nil, nil for not found
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByCode_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE message_code").
		WithArgs("ERR_001").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	// Act
	message, err := repo.GetByCode(context.Background(), "ERR_001")

	// Assert
	assert.Nil(t, message)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	}).AddRow(
		"msg-1", "ERR_001", "ERROR", "USER", "MOD_U",
		"Error Title", "Error content", true, now, now,
	)

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE id").
		WithArgs("msg-1").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	message, err := repo.GetByID(context.Background(), "msg-1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, "msg-1", message.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE id").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	repo := &repository{db: db}

	// Act
	message, err := repo.GetByID(context.Background(), "nonexistent")

	// Assert
	assert.Nil(t, message)
	assert.Nil(t, err) // GetByID returns nil, nil for not found
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveMessage_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO system_messages").
		WithArgs("msg-123", "ERR_001", "ERROR", "USER", "MOD_U", "Error Title", "Error content", true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	message := domain.Message{
		ID:       "msg-123",
		Code:     "ERR_001",
		Type:     domain.TypeError,
		Category: "USER",
		Module:   "MOD_U",
		Title:    "Error Title",
		Content:  "Error content",
		Active:   true,
	}

	// Act
	err = repo.SaveMessage(context.Background(), sqlTx, message)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveMessage_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO system_messages").
		WithArgs("msg-123", "ERR_001", "ERROR", "USER", "MOD_U", "Error Title", "Error content", true).
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	message := domain.Message{
		ID:       "msg-123",
		Code:     "ERR_001",
		Type:     domain.TypeError,
		Category: "USER",
		Module:   "MOD_U",
		Title:    "Error Title",
		Content:  "Error content",
		Active:   true,
	}

	// Act
	err = repo.SaveMessage(context.Background(), sqlTx, message)

	// Assert
	assert.Equal(t, domain.ErrMessageCannotSave, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveMessage_InvalidTransaction(t *testing.T) {
	// Arrange
	repo := &repository{}
	message := domain.Message{ID: "msg-123"}

	// Act
	err := repo.SaveMessage(context.Background(), nil, message)

	// Assert
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

func TestDeleteMessage_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM system_messages WHERE id").
		WithArgs("msg-123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	// Act
	err = repo.DeleteMessage(context.Background(), sqlTx, "msg-123")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMessage_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM system_messages WHERE id").
		WithArgs("msg-123").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	// Act
	err = repo.DeleteMessage(context.Background(), sqlTx, "msg-123")

	// Assert
	assert.Equal(t, domain.ErrMessageCannotDelete, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteMessage_InvalidTransaction(t *testing.T) {
	// Arrange
	repo := &repository{}

	// Act
	err := repo.DeleteMessage(context.Background(), nil, "msg-123")

	// Assert
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}

// ============================================
// BeginTx Tests
// ============================================

func TestMessageBeginTx_Success(t *testing.T) {
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

func TestMessageBeginTx_Error(t *testing.T) {
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

// ============================================
// GetByType Tests
// ============================================

func TestGetByType_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	}).AddRow(
		"msg-1", "ERR_001", "ERROR", "USER", "MOD_U",
		"Error Title", "Error content", true, now, now,
	).AddRow(
		"msg-2", "ERR_002", "ERROR", "SYSTEM", "MOD_S",
		"Another Error", "Another content", true, now, now,
	)

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE type").
		WithArgs("ERROR").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetByType(context.Background(), "ERROR")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, messages, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByType_Empty(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE type").
		WithArgs("NONEXISTENT").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetByType(context.Background(), "NONEXISTENT")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByType_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE type").
		WithArgs("ERROR").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetByType(context.Background(), "ERROR")

	// Assert
	assert.Nil(t, messages)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// GetByModule Tests
// ============================================

func TestGetByModule_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	}).AddRow(
		"msg-1", "ERR_001", "ERROR", "USER", "MOD_U",
		"Error Title", "Error content", true, now, now,
	)

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE module").
		WithArgs("MOD_U").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetByModule(context.Background(), "MOD_U")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, "MOD_U", messages[0].Module)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByModule_Empty(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "message_code", "type", "category", "module",
		"message_title", "message_content", "is_active", "created_at", "updated_at",
	})

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE module").
		WithArgs("NONEXISTENT").
		WillReturnRows(rows)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetByModule(context.Background(), "NONEXISTENT")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, messages)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByModule_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT .* FROM system_messages WHERE module").
		WithArgs("MOD_U").
		WillReturnError(sql.ErrConnDone)

	repo := &repository{db: db}

	// Act
	messages, err := repo.GetByModule(context.Background(), "MOD_U")

	// Assert
	assert.Nil(t, messages)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ============================================
// UpdateMessage Tests
// ============================================

func TestUpdateMessage_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE system_messages SET").
		WithArgs("New Title", "New Content", true, "msg-123", "ERR_001").
		WillReturnResult(sqlmock.NewResult(0, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	message := domain.Message{
		ID:      "msg-123",
		Code:    "ERR_001",
		Title:   "New Title",
		Content: "New Content",
		Active:  true,
	}

	// Act
	err = repo.UpdateMessage(context.Background(), sqlTx, message)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateMessage_Error(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE system_messages SET").
		WithArgs("New Title", "New Content", true, "msg-123", "ERR_001").
		WillReturnError(sql.ErrConnDone)

	tx, err := db.Begin()
	assert.NoError(t, err)

	sqlTx := common.NewSQLTx(tx)
	repo := &repository{db: db}

	message := domain.Message{
		ID:      "msg-123",
		Code:    "ERR_001",
		Title:   "New Title",
		Content: "New Content",
		Active:  true,
	}

	// Act
	err = repo.UpdateMessage(context.Background(), sqlTx, message)

	// Assert
	assert.Equal(t, domain.ErrMessageCannotUpdate, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateMessage_InvalidTransaction(t *testing.T) {
	// Arrange
	repo := &repository{}
	message := domain.Message{ID: "msg-123", Code: "ERR_001"}

	// Act
	err := repo.UpdateMessage(context.Background(), nil, message)

	// Assert
	assert.Equal(t, domain.ErrInvalidTransaction, err)
}
