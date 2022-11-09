package handlers

import (
	docs "Skipper_cms_auth/pkg/docs"
	service "Skipper_cms_auth/pkg/servises"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() {
	router := gin.Default()
	router.Use(corsMiddleware())
	api_v1 := router.Group("/api/v1")
	{
		auth := api_v1.Group("/auth")
		{
			auth.POST("/sign-in", h.signIn)
			auth.POST("/refresh-token", h.refreshToken)
			auth.POST("/reset-password", h.sendEmailToResetPassword)
		}
	}
	router.GET("/reset-password")

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run()
}
