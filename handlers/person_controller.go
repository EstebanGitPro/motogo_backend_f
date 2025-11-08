package handlers

import (
	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

func (h handler) GetPersonByEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Param("email")

		person, err := h.PersonService.GetPersonByEmail(c,email)
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

		result, err := h.PersonService.RegisterPerson(c,personRequest.ToDomain())
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
				Role:                result.Person.Role,
			},
			Message: result.Message,
		}

		c.JSON(http.StatusCreated, response)
	}
}
