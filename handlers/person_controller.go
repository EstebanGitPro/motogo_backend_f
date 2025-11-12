package handlers

import (
	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

func (h handler) RegisterPerson() func(c *gin.Context) {
	return func(c *gin.Context) {

		var personRequest PersonRequest
		if err := c.ShouldBindJSON(&personRequest); err != nil {
			c.Error(domain.ErrInvalidJSONFormat)
			return
		}

		result, err := h.Interactor.RegisterPerson(c, personRequest.ToDomain())
		if err != nil {
			c.Error(err)
			return
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		baseURL := scheme + "://" + c.Request.Host
		links := BuildAccountLinks(baseURL, result.Person.ID)
		
		locationURL := baseURL + "/motogo/api/v1/accounts/" + result.Person.ID
		c.Header("Location", locationURL)

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
			Links:   links,
		}

	
		c.JSON(http.StatusCreated, response)
	}
}
