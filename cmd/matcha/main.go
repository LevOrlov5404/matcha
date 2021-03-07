package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/LevOrlov5404/matcha/internal/handler"
	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/LevOrlov5404/matcha/internal/server"
	"github.com/LevOrlov5404/matcha/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("failed to initialize config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("failed to load env variables: %s", err.Error())
	}

	db, err := repository.ConnectToDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to connect to db: %s", err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Errorf("error occurred on db connection closing: %s", err.Error())
		}
	}()

	repo := repository.NewRepository(db)
	services := service.NewService(repo, os.Getenv("SALT"), os.Getenv("SIGNING_KEY"))
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TaskTracker started")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TaskTracker shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
