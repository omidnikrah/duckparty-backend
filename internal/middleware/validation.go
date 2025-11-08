package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/omidnikrah/duckparty-backend/internal/utils"
)

func ValidationErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if _, ok := err.(validator.ValidationErrors); ok {
				errorMessages := utils.FormatValidationError(err)
				c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
		}
	}
}
