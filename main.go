package main

import (
	"github.com/OmarDardery/solve-the-x-backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	server.POST("/test-email", func(c *gin.Context) {
		err := middleware.SendVerificationEmail("omar71epic@gmail.com", "123456")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Email sent!"})
	})

	server.Run(":8000")
}
