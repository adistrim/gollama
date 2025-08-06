package routes

import (
	"github.com/gin-gonic/gin"
)

func Master() *gin.Engine {
	router := gin.Default()
	
	// system endpoints
	router.GET("/health", HealthCheck)

	// websocket endpoint
	router.GET("/chat", WebSocketHandler)

	return router
}
