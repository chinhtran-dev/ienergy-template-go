package app

import (
	"ienergy-template-go/internal/http/handler"
	"ienergy-template-go/internal/http/router"
	"ienergy-template-go/internal/repository"
	"ienergy-template-go/internal/service"

	"go.uber.org/fx"
)

var Module = fx.Options(
	handler.Module,
	router.Module,
	repository.Module,
	service.Module,
)
