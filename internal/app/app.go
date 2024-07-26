package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/canyouhearthemusic/project-management/config"
	"github.com/canyouhearthemusic/project-management/internal/handler"
	"github.com/canyouhearthemusic/project-management/internal/repository"
	"github.com/canyouhearthemusic/project-management/internal/repository/postgres"
	"github.com/canyouhearthemusic/project-management/internal/service/management"
	"github.com/canyouhearthemusic/project-management/pkg/server"
	"github.com/sirupsen/logrus"
)

func Run() {
	logger := logrus.New().WithContext(context.Background())

	configs, err := config.New()
	if err != nil {
		logger.Errorln("failed to load configurations")
		return
	}

	db, err := postgres.New(configs.DB)
	if err != nil {
		logger.Errorln("failed to connect to database")
		return
	}
	defer db.Close()

	repositories, err := repository.New(repository.WithPostgresStore(configs.DB))
	if err != nil {
		logger.Errorln("failed to create repositories")
		return
	}

	managementService := management.New(
		management.WithProjectRepository(repositories.Project),
		management.WithTaskRepository(repositories.Task),
		management.WithUserRepository(repositories.User),
	)

	handler := handler.New(
		handler.Dependencies{
			ManagementService: managementService,
		},
		handler.WithHTTPHandler())

	server, err := server.New(server.WithHTTPServer(handler.Mux, configs.APP.Port))
	if err != nil {
		logger.Errorln("failed to create server")
		return
	}

	if err := server.Start(); err != nil {
		logger.Errorln("failed to start server")
		return
	}

	logger.Infof("server is running on port %s, swagger is at /swagger/index.html\n", configs.APP.Port)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown
	logger.Infoln("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logger.Errorln("failed to stop server")
		return
	}

	logger.Infoln("server stopped")
}
