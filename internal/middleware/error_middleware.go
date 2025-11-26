package middleware

import (
	"log"
	"net/http"
	"renew-guard/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ErrorMiddleware handles panics and errors globally
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors from handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Request error: %v", err.Error())
			
			// Check if response was already written
			if !c.Writer.Written() {
				utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
		}
	}
}
