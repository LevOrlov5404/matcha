package server

import (
	"context"
	"net/http"
	"time"

	"github.com/LevOrlov5404/matcha/internal/config"
	"github.com/LevOrlov5404/matcha/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const timeout = 10 * time.Second

type (
	Server struct {
		cfg        *config.Config
		log        *logrus.Logger
		services   *service.Service
		httpServer *http.Server
	}
)

func NewServer(
	cfg *config.Config, log *logrus.Logger, services *service.Service,
) *Server {
	s := &Server{
		cfg:      cfg,
		log:      log,
		services: services,
	}

	s.httpServer = &http.Server{
		Addr:           cfg.Address.String(),
		Handler:        s.InitRoutes(),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
	}

	return s
}

func (s *Server) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(s.InitMiddleware)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", s.CreateUser)
		auth.POST("/sign-in", s.SignIn)
		router.POST("/reset-password", s.ResetPassword)
		// auth.POST("/refresh-session", s.RefreshSession)
		// auth.POST("/logout", s.Logout)
	}

	router.POST("/confirm-email", s.ConfirmEmail)
	router.POST("/confirm-reset-password", s.ConfirmPasswordReset)

	api := router.Group("/api/v1", s.UserIdentityMiddleware)
	{
		users := api.Group("/users")
		{
			users.POST("", s.CreateUser)
			users.GET("", s.GetAllUsers)
			users.GET("/by-id/:id", s.GetUserByID)
			users.PUT("", s.UpdateUser)
			users.PUT("/set-password", s.SetUserPassword)
			users.PUT("/change-password", s.ChangeUserPassword)
			users.DELETE("/:id", s.DeleteUser)
		}
	}

	return router
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
