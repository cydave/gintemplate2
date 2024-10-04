package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		isAuthenticated := false
		_, hasUserID := session.Get("user_id").(string)
		if hasUserID && session.Get("authenticated") == "1" {
			isAuthenticated = true
		}
		if !isAuthenticated {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
