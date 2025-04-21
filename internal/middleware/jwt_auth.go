package middleware

import (
	"net/http"

	"ienergy-template-go/config"
	"ienergy-template-go/pkg/errors"
	"ienergy-template-go/pkg/util"
	"ienergy-template-go/pkg/wrapper"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := util.TokenValid(c, config.JWT)
		if err != nil {
			c.JSON(http.StatusUnauthorized, wrapper.NewErrorResponse(
				errors.NewUnauthorizedError("Unauthorized"),
			))
			c.Abort()
			return
		}
		c.Next()
	}
}
