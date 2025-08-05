package routes

import (
	"github.com/gin-gonic/gin"
)

func Master() *gin.Engine {
	router := gin.Default()
	
	router.GET("/health", HealthCheck)
	
	return router
}
