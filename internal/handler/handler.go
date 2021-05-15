package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/service"
	"github.com/sirupsen/logrus"
)

type (
	Options struct {
		AccessTokenCookieMaxAge  int
		RefreshTokenCookieMaxAge int
		SecureCookie             *securecookie.SecureCookie
	}
	Handler struct {
		cfg     *config.Config
		log     *logrus.Logger
		options Options
		svc     *service.Service
	}
)

func New(
	cfg *config.Config, log *logrus.Logger, svc *service.Service,
) *Handler {
	c := &Handler{
		cfg: cfg,
		log: log,
		options: Options{
			AccessTokenCookieMaxAge:  int(cfg.JWT.AccessTokenLifetime.Duration().Seconds()),
			RefreshTokenCookieMaxAge: int(cfg.JWT.RefreshTokenLifetime.Duration().Seconds()),
			SecureCookie:             securecookie.New(cfg.Cookie.HashKey, cfg.Cookie.BlockKey),
		},
		svc: svc,
	}

	return c
}

func (h *Handler) InitRoutes() http.Handler {
	router := gin.New()

	router.Use(h.InitMiddleware)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.CreateUser)
		auth.POST("/sign-in", h.SignIn)
		router.POST("/reset-password", h.ResetPassword)
		auth.POST("/validate-access-token", h.ValidateAccessToken)
		auth.POST("/refresh-session", h.RefreshSession)
		auth.POST("/logout", h.Logout)
	}

	router.POST("/confirm-email", h.ConfirmEmail)
	router.POST("/confirm-reset-password", h.ConfirmPasswordReset)

	api := router.Group("/api/v1", h.UserAuthorizationMiddleware)
	{
		users := api.Group("/users")
		{
			users.POST("/", h.CreateUser)
			users.GET("/", h.GetAllUsers)
			users.GET("/by-id/:id", h.GetUserByID)
			users.PUT("/", h.UpdateUser)
			users.PUT("/set-password", h.SetUserPassword)
			users.PUT("/change-password", h.ChangeUserPassword)
			users.DELETE("/by-id/:id", h.DeleteUser)
			users.GET("/profile/by-id/:id", h.GetUserProfileByID)
			users.PUT("/profile", h.UpdateUserProfile)

			usersPictures := users.Group("/pictures")
			{
				usersPictures.POST("/avatar", h.UploadUserAvatar)
				usersPictures.DELETE("/avatar", h.DeleteUserAvatar)
				usersPictures.POST("/", h.UploadUserPicture)
				usersPictures.GET("/", h.GetUserPictures)
				usersPictures.DELETE("/", h.DeleteUserPicture)
			}
		}
	}

	return CORS(router)
}
