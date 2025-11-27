package message

import (
	"context"
	"database/sql"
)

const (
	queryGetAllActive = `
		SELECT id, codigo_mensaje, tipo, categoria, modulo, titulo_mensaje, contenido_mensaje, activo, created_at, updated_at 
		FROM mensajes_sistema 
		WHERE activo = true`

	queryGetByCode = `
		SELECT id, codigo_mensaje, tipo, categoria, modulo, titulo_mensaje, contenido_mensaje, activo, created_at, updated_at 
		FROM mensajes_sistema 
		WHERE codigo_mensaje = ? AND activo = true 
		LIMIT 1`

	queryGetByType = `
		SELECT id, codigo_mensaje, tipo, categoria, modulo, titulo_mensaje, contenido_mensaje, activo, created_at, updated_at 
		FROM mensajes_sistema 
		WHERE tipo = ? AND activo = true`

	queryGetByModule = `
		SELECT id, codigo_mensaje, tipo, categoria, modulo, titulo_mensaje, contenido_mensaje, activo, created_at, updated_at 
		FROM mensajes_sistema 
		WHERE modulo = ? AND activo = true`
)

// Repository defines the interface for message data access
type Repository interface {
	GetAllActive(ctx context.Context) ([]SystemMessage, error)
	GetByCode(ctx context.Context, code string) (*SystemMessage, error)
	GetByType(ctx context.Context, msgType MessageType) ([]SystemMessage, error)
	GetByModule(ctx context.Context, module string) ([]SystemMessage, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository creates a new message repository
func NewRepository(db *sql.DB) (Repository, error) {
	if db == nil {
		return nil, sql.ErrConnDone
	}
	return &repository{db: db}, nil
}

func (r *repository) GetAllActive(ctx context.Context) ([]SystemMessage, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAllActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []SystemMessage
	for rows.Next() {
		var m SystemMessage
		err := rows.Scan(
			&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
			&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}

func (r *repository) GetByCode(ctx context.Context, code string) (*SystemMessage, error) {
	var m SystemMessage
	err := r.db.QueryRowContext(ctx, queryGetByCode, code).Scan(
		&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
		&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *repository) GetByType(ctx context.Context, msgType MessageType) ([]SystemMessage, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByType, msgType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []SystemMessage
	for rows.Next() {
		var m SystemMessage
		err := rows.Scan(
			&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
			&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}

func (r *repository) GetByModule(ctx context.Context, module string) ([]SystemMessage, error) {
	rows, err := r.db.QueryContext(ctx, queryGetByModule, module)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []SystemMessage
	for rows.Next() {
		var m SystemMessage
		err := rows.Scan(
			&m.ID, &m.Code, &m.Type, &m.Category, &m.Module,
			&m.Title, &m.Content, &m.Active, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, rows.Err()
}
