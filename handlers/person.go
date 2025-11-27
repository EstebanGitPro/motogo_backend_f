package handlers

import (
	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

type PersonRequest struct {
	IdentityNumber string `json:"identity_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Password       string `json:"password"`
	Role           string `json:"role"`
}

type PersonResponse struct {
	ID             string `json:"id"`
	IdentityNumber string `json:"identity_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	SecondLastName string `json:"second_last_name"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	Role           string `json:"role"`
}


type RegistrationResponse struct {
	Message string         `json:"message"`
	Links   []Link         `json:"_links"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (p PersonRequest) ToDomain() domain.Person {
	return domain.Person{
		IdentityNumber: p.IdentityNumber,
		FirstName:      p.FirstName,
		LastName:       p.LastName,
		SecondLastName: p.SecondLastName,
		Email:          p.Email,
		PhoneNumber:    p.PhoneNumber,
		Password:       p.Password,
		Role:           p.Role,
	}
}


// GetBaseURL extrae la URL base de la petición (scheme + host)
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// EncodeID ofusca un UUID usando el encoder del handler
// Retorna el ID ofuscado o un error si falla
func (h *handler) EncodeID(uuid string) (string, error) {
	encodedID, err := h.IDEncoder.Encode(uuid)
	if err != nil {
		h.Logger.Error("Error ofuscando ID",
			"uuid", uuid,
			"error", err)
		return "", err
	}
	return encodedID, nil
}

// DecodeID desofusca un ID ofuscado a UUID usando el encoder del handler
// Retorna el UUID o un error si falla
func (h *handler) DecodeID(encodedID string) (string, error) {
	uuid, err := h.IDEncoder.Decode(encodedID)
	if err != nil {
		h.Logger.Error("Error decodificando ID",
			"encoded_id", encodedID,
			"error", err)
		return "", err
	}
	return uuid, nil
}

// HandleIDEncodingError maneja errores de ofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDEncodingError(c *gin.Context, uuid string, err error) {
	h.Logger.Error("Error ofuscando ID",
		"uuid", uuid,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInternalServer)
}

// HandleIDDecodingError maneja errores de desofuscamiento y envía respuesta apropiada
func (h *handler) HandleIDDecodingError(c *gin.Context, encodedID string, err error) {
	h.Logger.Error("Error decodificando ID",
		"encoded_id", encodedID,
		"error", err,
		"client_ip", c.ClientIP())
	c.Error(domain.ErrInvalidID)
}

// SetLocationHeader establece el Location header con la URL del recurso
func SetLocationHeader(c *gin.Context, baseURL, resource, encodedID string) {
	locationURL := BuildResourceURL(baseURL, resource, encodedID)
	c.Header("Location", locationURL)
}
