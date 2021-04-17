package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/repository"
	userpostgres "github.com/l-orlov/matcha/internal/repository/user-postgres"
	"github.com/l-orlov/matcha/internal/server"
	"github.com/l-orlov/matcha/internal/service"
	"github.com/l-orlov/matcha/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"github.com/sirupsen/logrus"
)

const (
	passwordAllowedLowerLetters = "abcdefghijklmnopqrstuvwxyz"
	passwordAllowedUpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	passwordAllowedDigits       = "0123456789"
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

	db, err := userpostgres.ConnectToDB(cfg)
	if err != nil {
		lg.Fatalf("failed to connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			lg.Errorf("failed to close db: %v", err)
		}
	}()

	repo, err := repository.NewRepository(cfg, lg, db)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	randomSymbolsGenerator, err := password.NewGenerator(&password.GeneratorInput{
		LowerLetters: passwordAllowedLowerLetters,
		UpperLetters: passwordAllowedUpperLetters,
		Digits:       passwordAllowedDigits,
	})
	if err != nil {
		log.Fatalf("failed to create random symbols generator: %v", err)
	}

	mailerLogEntry := logrus.NewEntry(lg).WithFields(logrus.Fields{"source": "mailerService"})
	mailer := service.NewMailerService(cfg.Mailer, mailerLogEntry)
	defer mailer.Close()

	svc := service.NewService(cfg, lg, repo, randomSymbolsGenerator, mailer)

	srv := server.NewServer(cfg, lg, svc)
	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Fatalf("error occurred while running http server: %v", err)
		}
	}()

	lg.Info("service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	lg.Info("service shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		lg.Errorf("failed to shut down: %v", err)
	}
}
