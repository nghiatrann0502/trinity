package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/nghiatrann0502/trinity/internal/video/app"
	"github.com/nghiatrann0502/trinity/internal/video/config"
	"github.com/nghiatrann0502/trinity/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig("./config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load configuration")
	}

	logger := logger.NewLogger(cfg.App.Name, cfg.App.Version, !cfg.App.Production)
	logger.Info(fmt.Sprintf("Starting %s service with version %s ", cfg.App.Name, cfg.App.Version), nil)
	app, cleanup, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal("cannot start the service", err, nil)
	}

	if err := app.Run(); err != nil {
		logger.Fatal("cannot start the service", err, nil)
	}

	<-ctx.Done()
	logger.Warn("Shutting down servers...", nil)

	cleanup()

	logger.Info("Servers shutdown completed", nil)
}
