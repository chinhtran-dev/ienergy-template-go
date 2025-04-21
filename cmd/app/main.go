package main

import (
	"context"
	"fmt"
	"ienergy-template-go/config"
	"ienergy-template-go/internal/app"
	"ienergy-template-go/pkg/database"
	"ienergy-template-go/pkg/graceful"
	"ienergy-template-go/pkg/logger"
	"ienergy-template-go/pkg/swagger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/fx"
)

func registerSwaggerHandler(g *gin.Engine) {
	swaggerAPI := g.Group("/swagger")
	swag := swagger.NewSwagger()
	swaggerAPI.Use(swag.SwaggerHandler(false))
	swag.Register(swaggerAPI)
}

func startServer(g *gin.Engine, lifecycle fx.Lifecycle, logger *logger.StandardLogger, config *config.Config) {
	gracefulService := graceful.NewService(graceful.WithStopTimeout(time.Second), graceful.WithWaitTime(time.Second))
	gracefulService.Register(g)
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port := fmt.Sprintf("%d", cast.ToInt(config.Server.Port))
				fmt.Println("run on port:", port)
				go gracefulService.StartServer(g, port)
				return nil
			},
			OnStop: func(context.Context) error {
				gracefulService.Close(logger)
				return nil
			},
		},
	)
}

func main() {
	fx.New(
		fx.Provide(config.NewConfig),
		fx.Provide(database.NewDatabase),
		fx.Provide(logger.NewLogger),
		app.Module,
		fx.Invoke(
			registerSwaggerHandler,
			startServer,
		),
	).Run()
}
