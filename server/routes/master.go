package routes

import (
    "github.com/gin-gonic/gin"
)

func Master() *gin.Engine {
    router := gin.Default()

    // system endpoints
    router.GET("/health", HealthCheck)

    // HTTP chat endpoint
    router.POST("/chat", HTTPChatHandler)

    return router
}
