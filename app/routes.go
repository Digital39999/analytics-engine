package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	router.GET("/", infoHandler)
	router.GET("/stats", statsHandler)
	router.POST("/event", eventHandler)

	router.GET("/analytics", analyticsHandler)
	router.DELETE("/analytics", flushHandler)

	router.NoRoute(notFoundHandler)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/" {
			c.Next()
			return
		}

		apiAuth := os.Getenv("API_AUTH")
		authHeader := c.GetHeader("Authorization")

		if authHeader != apiAuth {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "data": "Unauthorized."})
			c.Abort()
			return
		}

		c.Next()
	}
}
