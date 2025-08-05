package routes

import (
	"github.com/gin-gonic/gin"
)

func Master() *gin.Engine {
	router := gin.Default()
	
	// system endpoints
	router.GET("/health", HealthCheck)

	// application endpoints
	router.POST("/chat", ChatHandler)

	return router
}
