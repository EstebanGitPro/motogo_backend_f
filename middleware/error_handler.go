package middleware

import (
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/core/interactor/services/domain"
	"github.com/gin-gonic/gin"
)

var mapError = map[error]ErrorResponse{

	domain.ErrDuplicateUser: {
		Code:    "MOD_U_DUP_ERR_00001",
		Message: "Usuario duplicado",
		Status:  http.StatusConflict,
	},
	
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}


func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() 

        // Step2: Check if any errors were added to the context
        if len(c.Errors) > 0 {
            // Step3: Use the last error
            err := c.Errors.Last().Err
			
			if response, ok := mapError[err]; ok {
				c.JSON(response.Status, response)
				return
			}

            // Step4: Respond with a generic error message
            c.JSON(http.StatusInternalServerError, map[string]any{
                "success": false,
                "message": err.Error(),
            })
        }

        // Any other steps if no errors are found
    }
}