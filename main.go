package main

import (
	"github.com/OmarDardery/solve-the-x-backend/database"
	"github.com/OmarDardery/solve-the-x-backend/models"
	"github.com/OmarDardery/solve-the-x-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db, closer := database.NewDatabase()
	defer closer()

	if err := db.AutoMigrate(&models.Professor{}, &models.Student{}); err != nil {
		panic("failed to migrate database")
	}

	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.POST("/sign-up/:role", routes.SignUpHandler(db))
	server.POST("/sign-in/:role", routes.SignInHandler(db))

	server.Run(":8000")
}
