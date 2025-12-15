package app

import (
	"context"

	_ "github.com/ruslanonly/blindtyping/src/docs"
	"github.com/ruslanonly/blindtyping/src/internal/app/config"
	"github.com/ruslanonly/blindtyping/src/internal/app/di"
)

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	ctx := context.Background()
	cfg := config.MustParse()
	diContainer := di.NewContainer(cfg)
	logger := diContainer.Logger()
	logger.Info(logger.WithMsg(ctx, "di container created"))

	// Postgres
	postgres := diContainer.Postgres()
	postgres.MustConnect(ctx, cfg.Postgres.Connection)
	defer postgres.MustClose()
	postgres.MustMigrate(cfg.Postgres.Migrations, cfg.Postgres.Connection)
	logger.Info(logger.WithMsg(ctx, "connected to postgres"))

	// Redis
	redisClient := diContainer.Redis()
	redisClient.MustConnect(ctx)
	defer redisClient.MustClose()
	logger.Info(logger.WithMsg(ctx, "connected to redis"))

	// Scheduler
	scheduler := diContainer.Scheduler()
	scheduler.Start()
	defer scheduler.Stop()
	logger.Info(logger.WithMsg(ctx, "scheduler started"))

	// Antifroad
	if !cfg.Antifroad.IsDisabled {
		antifroadService := diContainer.AntifroadService()
		err := antifroadService.Init(ctx)
		if err != nil {
			panic(err)
		}
		logger.Info(logger.WithMsg(ctx, "antifroad service initialized"))
	}

	// HTTP server
	server := diContainer.Server()
	server.MustRun(ctx)
	logger.Info(logger.WithMsg(ctx, "server stopped"))
}
