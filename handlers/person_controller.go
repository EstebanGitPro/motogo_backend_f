package handlers

import (
	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/core/domain"
	"github.com/gin-gonic/gin"
)

func (h handler) GetPersonByEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Param("email")

		person, err := h.PersonService.GetPersonByEmail(email)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(http.StatusOK, person)
	}
}

func (h handler) RegisterPerson() func(c *gin.Context) {
	return func(c *gin.Context) {

		var personRequest PersonRequest
		if err := c.ShouldBindJSON(&personRequest); err != nil {
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		result, err := h.PersonService.RegisterPerson(personRequest.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		response := RegistrationResponse{
			User: PersonResponse{
				ID:                  result.Person.ID,
				IdentityNumber:      result.Person.IdentityNumber,
				FirstName:           result.Person.FirstName,
				LastName:            result.Person.LastName,
				SecondLastName:      result.Person.SecondLastName,
				Email:               result.Person.Email,
				PhoneNumber:         result.Person.PhoneNumber,
				EmailVerified:       result.Person.EmailVerified,
				PhoneNumberVerified: result.Person.PhoneNumberVerified,
				Role:                result.Person.Role,
				KeycloakUserID:      result.Person.KeycloakUserID,
			},
			Message: result.Message,
		}

		c.JSON(http.StatusCreated, response)
	}
}

func (h handler) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var loginRequest LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		token, err := h.PersonService.LoginPerson(loginRequest.Email, loginRequest.Password)
		if err != nil {
			c.Error(err)
			return
		}

		response := LoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresIn,
			TokenType:    token.TokenType,
		}

		c.JSON(http.StatusOK, response)
	}
}