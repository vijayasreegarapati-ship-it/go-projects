package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("DEBUG: Admin Middleware Intercepted Request!")

		role := c.GetHeader("X-User-Role")

		if role != "admin" {
			log.Println("DEBUG: Blocked! Role was:", role)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: admin access required"})
			return
		}

		log.Println("DEBUG: Passed! Role was admin.")
		c.Next()
	}
}
