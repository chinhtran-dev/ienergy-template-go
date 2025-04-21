package router

import (
	"ienergy-template-go/internal/middleware"
	"ienergy-template-go/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouterParams struct {
	fx.In
	AuthRoutes   AuthRoutes
	UserRoutes   UserRoutes
	Logger       *logger.StandardLogger
	ErrorHandler *middleware.ErrorHandler
}

func NewRouter(params RouterParams) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.LoggingMiddleware(params.Logger))
	router.Use(params.ErrorHandler.Handle())

	api := router.Group("/api/v1")
	params.AuthRoutes.Setup(api)
	params.UserRoutes.Setup(api)
	return router
}

var Module = fx.Options(
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(middleware.NewErrorHandler),
	fx.Provide(NewRouter),
)
