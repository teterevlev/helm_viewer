package router

import (
	"helm-viewer/handlers"
	"helm-viewer/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	helmService := services.NewHELMService()

	helmHandler := handlers.NewHELMHandler(helmService)

	api := r.Group("/api")
	{
		api.POST("/helm/load", helmHandler.LoadHELM)
	}

	return r
}
