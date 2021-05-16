package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/handler"
	"github.com/l-orlov/matcha/internal/repository"
	"github.com/l-orlov/matcha/internal/repository/postgres"
	"github.com/l-orlov/matcha/internal/server"
	"github.com/l-orlov/matcha/internal/service"
	"github.com/l-orlov/task-tracker/pkg/logger"
	"github.com/l-orlov/task-tracker/pkg/mailer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Run initializes whole application.
func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	lg, err := logger.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	// Dependencies
	db, err := postgres.ConnectToDB(cfg.PostgresDB)
	if err != nil {
		lg.Fatalf("failed to connect to db: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			lg.Errorf("failed to close db: %v", err)
		}
	}()

	if cfg.PostgresDB.MigrationMode {
		if err = postgres.MigrateSchema(db.DB, cfg.PostgresDB); err != nil {
			log.Fatalf("failed to do migration: %v", err)
		}
	}

	m := mailer.New(
		mailer.Config{
			Host:              cfg.Mailer.ServerAddress.Host,
			Port:              cfg.Mailer.ServerAddress.Port,
			Username:          cfg.Mailer.Username,
			Password:          cfg.Mailer.Password.String(),
			Timeout:           cfg.Mailer.Timeout.Duration(),
			MsgToSendChanSize: cfg.Mailer.MsgToSendChanSize,
			WorkersNum:        cfg.Mailer.WorkersNum,
		},
		logrus.NewEntry(lg).WithFields(logrus.Fields{"source": "mailer"}),
	)
	m.Init()
	defer m.Shutdown()

	// Repo, Service & API Handlers
	repo, err := repository.NewRepository(cfg, lg, db)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	svc, err := service.NewService(cfg, lg, repo, m)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	h := handler.New(cfg, lg, svc)

	// HTTP Server
	srv := server.New(cfg.Port, h.InitRoutes())
	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Fatalf("error occurred while running http server: %v", err)
		}
	}()

	lg.Infof("service started on port %s", cfg.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	lg.Info("service shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		lg.Errorf("failed to shut down: %v", err)
	}
}
