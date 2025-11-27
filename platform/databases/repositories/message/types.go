package message

import "time"

// MessageType represents system message types
type MessageType string

const (
	TypeError   MessageType = "ERROR"
	TypeSuccess MessageType = "EXITO"
	TypeWarning MessageType = "WARNING"
	TypeInfo    MessageType = "INFO"
	TypeDebug   MessageType = "DEBUG"
)

// SystemMessage represents a message stored in DB/Cache
type SystemMessage struct {
	ID        int64       `json:"id" db:"id"`
	Code      string      `json:"codigo_mensaje" db:"codigo_mensaje"`
	Type      MessageType `json:"tipo" db:"tipo"`
	Category  string      `json:"categoria" db:"categoria"`
	Module    string      `json:"modulo" db:"modulo"`
	Title     string      `json:"titulo_mensaje" db:"titulo_mensaje"`
	Content   string      `json:"contenido_mensaje" db:"contenido_mensaje"`
	Active    bool        `json:"activo" db:"activo"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// MessageResponse is the response structure for the frontend
type MessageResponse struct {
	Code    string      `json:"code"`
	Type    MessageType `json:"type"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
}

// ToResponse converts to response with placeholders replaced
func (m *SystemMessage) ToResponse(params ...string) MessageResponse {
	content := m.Content
	for i, param := range params {
		placeholder := "${" + string(rune('0'+i)) + "}"
		content = replaceAll(content, placeholder, param)
	}
	return MessageResponse{
		Code:    m.Code,
		Type:    m.Type,
		Title:   m.Title,
		Content: content,
	}
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old)
		} else {
			result += string(s[i])
			i++
		}
	}
	return result
}
