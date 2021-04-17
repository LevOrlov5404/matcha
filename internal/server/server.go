package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/service"
	"github.com/sirupsen/logrus"
)

const timeout = 10 * time.Second

type (
	Options struct {
		AccessTokenCookieMaxAge  int
		RefreshTokenCookieMaxAge int
		SecureCookie             *securecookie.SecureCookie
	}
	Server struct {
		cfg        *config.Config
		log        *logrus.Logger
		options    Options
		svc        *service.Service
		httpServer *http.Server
	}
)

func NewServer(
	cfg *config.Config, log *logrus.Logger, svc *service.Service,
) *Server {
	s := &Server{
		cfg: cfg,
		log: log,
		options: Options{
			AccessTokenCookieMaxAge:  int(cfg.JWT.AccessTokenLifetime.Duration().Seconds()),
			RefreshTokenCookieMaxAge: int(cfg.JWT.RefreshTokenLifetime.Duration().Seconds()),
			SecureCookie:             securecookie.New(cfg.Cookie.HashKey, cfg.Cookie.BlockKey),
		},
		svc: svc,
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
		auth.POST("/validate-access-token", s.ValidateAccessToken)
		auth.POST("/refresh-session", s.RefreshSession)
		auth.POST("/logout", s.Logout)
	}

	router.POST("/confirm-email", s.ConfirmEmail)
	router.POST("/confirm-reset-password", s.ConfirmPasswordReset)

	api := router.Group("/api/v1", s.UserAuthorizationMiddleware)
	{
		users := api.Group("/users")
		{
			users.POST("/", s.CreateUser)
			users.GET("/", s.GetAllUsers)
			users.GET("/by-id/:id", s.GetUserByID)
			users.PUT("/", s.UpdateUser)
			users.PUT("/set-password", s.SetUserPassword)
			users.PUT("/change-password", s.ChangeUserPassword)
			users.DELETE("/by-id/:id", s.DeleteUser)
			users.GET("/profile/by-id/:id", s.GetUserProfileByID)
			users.PUT("/profile", s.UpdateUserProfile)

			usersPictures := users.Group("/pictures")
			{
				usersPictures.POST("/avatar", s.UploadUserAvatar)
				usersPictures.DELETE("/avatar", s.DeleteUserAvatar)
				usersPictures.POST("/", s.UploadUserPicture)
				usersPictures.GET("/", s.GetUserPictures)
				usersPictures.DELETE("/", s.GetUserPictures)
			}
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
