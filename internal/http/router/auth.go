package router

import (
	"ienergy-template-go/internal/http/handler"

	"github.com/gin-gonic/gin"
)

type AuthRoutes interface {
	Setup(r *gin.RouterGroup)
}

type authRoutes struct {
	authHandler handler.AuthHandler
}

func (sr *authRoutes) Setup(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", sr.authHandler.Register())
		auth.POST("/login", sr.authHandler.Login())
	}
}

func NewAuthRoutes(authHandler handler.AuthHandler) AuthRoutes {
	return &authRoutes{
		authHandler: authHandler,
	}
}
