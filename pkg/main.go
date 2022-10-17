package main

import (
	"Skipper_cms_auth/pkg/handlers"
	"Skipper_cms_auth/pkg/models"
	"Skipper_cms_auth/pkg/repository"
	"Skipper_cms_auth/pkg/servises"
)

func main() {
	db := models.GetDB()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlerses := handlers.NewHandler(services)
	handlerses.InitRoutes()
}
