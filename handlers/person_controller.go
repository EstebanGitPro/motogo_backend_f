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
			h.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, person)
	}
}

func (h handler) RegisterPerson() func(c *gin.Context) {
	return func(c *gin.Context) {

		var personRequest PersonRequest
		if err := c.ShouldBindJSON(&personRequest); err != nil {
			h.HandleError(c, domain.ErrInvalidJSONFormat)
			return
		}

		person, err := h.PersonService.RegisterPerson(personRequest.ToDomain())
		if err != nil {

			switch err {
			case domain.ErrDuplicateUser:
				h.HandleError(c, domain.ErrDuplicateUser)
			case domain.ErrUserCannotSave:
				h.HandleError(c, domain.ErrUserCannotSave)
			default:
				h.HandleError(c, domain.ErrUserCannotSave)
			}
			return
		}

		response := PersonResponse{
			ID:                  person.ID,
			IdentityNumber:      person.IdentityNumber,
			FirstName:           person.FirstName,
			LastName:            person.LastName,
			SecondLastName:      person.SecondLastName,
			Email:               person.Email,
			PhoneNumber:         person.PhoneNumber,
			EmailVerified:       person.EmailVerified,
			PhoneNumberVerified: person.PhoneNumberVerified,
			Role:                person.Role,
		}

		c.JSON(http.StatusCreated, response)
	}
}