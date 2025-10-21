package main

import (
	"log"

	"github.com/OmarDardery/solve-the-x-backend/database"
	"github.com/OmarDardery/solve-the-x-backend/middleware"
	"github.com/OmarDardery/solve-the-x-backend/models"
	"github.com/OmarDardery/solve-the-x-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, continuing...")
	}

	// Initialize database
	db, closer := database.NewDatabase()
	defer closer()

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Professor{}, &models.Tag{}, &models.Student{}, &models.Coins{}, &models.Opportunity{}, &models.Application{}); err != nil {
		panic("failed to migrate database")
	}

	// Initialize Gin router
	server := gin.Default()

	// verification code map
	verificationCodes := make(map[string]int)
	// Public routes
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	server.POST("/sign-up/:role", routes.SignUpHandler(db, &verificationCodes))
	server.POST("/sign-in/:role", routes.SignInHandler(db))
	server.POST("/send-code", routes.SendCodeHandler(db, &verificationCodes))

	// Protected routes
	protected := server.Group("/api")
	protected.Use(middleware.JWTMiddleware(db))

	// Example protected route
	protected.GET("/profile", func(c *gin.Context) {
		role, _ := c.Get("role")
		user, _ := c.Get("user")

		c.JSON(200, gin.H{
			"role": role,
			"user": user,
		})
	})

	// Run server
	if err := server.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
