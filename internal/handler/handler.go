package handler

import (
	"github.com/LevOrlov5404/matcha/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sigh-up", h.CreateUser)
		auth.POST("/sigh-in", h.SignIn)
	}

	api := router.Group("/api/v1", h.UserIdentity)
	{
		users := api.Group("/users")
		{
			users.POST("/", h.CreateUser)
			users.GET("/", h.GetAllUsers)
			users.GET("/by-id/:id", h.GetUserByID)
			users.GET("/by-email-password", h.GetUserByEmailPassword)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}
	}

	return router
}
