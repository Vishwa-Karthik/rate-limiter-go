package main

import (
	"net/http"

	ipclientratelimit "github.com/Vishwa-Karthik/rate-limiter/ip-client-rate-limit"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	port := ":8080"

	// router.Use(tokenbucket.RateLimiter())

	router.Use(ipclientratelimit.RateLimiter())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Go server 🚀"})
	})

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "uptime": "running"})
	})

	// Start the server
	if err := router.Run(port); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
