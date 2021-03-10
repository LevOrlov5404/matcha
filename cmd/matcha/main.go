package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/LevOrlov5404/matcha/internal/config"
	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/LevOrlov5404/matcha/internal/server"
	"github.com/LevOrlov5404/matcha/internal/service"
	"github.com/LevOrlov5404/matcha/pkg/logger"
	_ "github.com/lib/pq"
)

func main() {
	cfg := &config.Config{}
	if err := config.ReadFromFileAndSetEnv(os.Getenv("CONFIG_PATH"), cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	lg, err := logger.New(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	db, err := repository.ConnectToDB(repository.Config{
		Host:     cfg.DB.Address.Host,
		Port:     cfg.DB.Address.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Database: cfg.DB.Database,
	})
	if err != nil {
		lg.Fatalf("failed to connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			lg.Errorf("failed to close db: %v", err)
		}
	}()

	repo := repository.NewRepository(db, cfg.DB.Timeout.Duration())
	services := service.NewService(repo, service.Options{
		TokenLifetime: cfg.JWT.TokenLifetime.Duration(),
		SigningKey:    cfg.JWT.SigningKey,
	})

	srv := server.NewServer(cfg, lg, services)
	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Fatalf("error occurred while running http server: %v", err)
		}
	}()

	lg.Info("service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-quit

	lg.Info("service shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		lg.Errorf("failed to shut down: %v", err)
	}
}
