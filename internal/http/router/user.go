package router

import (
	"ienergy-template-go/config"
	"ienergy-template-go/internal/http/handler"
	"ienergy-template-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoutes interface {
	Setup(r *gin.RouterGroup)
}

type userRoutes struct {
	userHandler handler.UserHandler
	config      *config.Config
}

func (sr *userRoutes) Setup(r *gin.RouterGroup) {
	userInfo := r.Group("/user/info")
	userInfo.Use(middleware.JwtAuthMiddleware(sr.config))
	{
		userInfo.GET("", sr.userHandler.Info())
	}
}

func NewUserRoutes(userHandler handler.UserHandler, config *config.Config) UserRoutes {
	return &userRoutes{
		userHandler: userHandler,
		config:      config,
	}
}
