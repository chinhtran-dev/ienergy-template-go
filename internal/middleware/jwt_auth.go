package middleware

import (
	"net/http"

	"ienergy-template-go/config"
	"ienergy-template-go/pkg/util"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := util.TokenValid(c, config.JWT)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
