package main

import (
	"log"
	"os"

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

	server.GET("/", func(c *gin.Context) {
		// Path to your PDF file
		filePath := "./doc.pdf"

		// Set headers
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "inline; filename=example.pdf")
		c.File(filePath)
	})
	auth := server.Group("/auth")
	auth.POST("/sign-up/:role", routes.SignUpHandler(db, &verificationCodes))
	auth.POST("/sign-in/:role", routes.SignInHandler(db))
	auth.POST("/send-code", routes.SendCodeHandler(db, &verificationCodes))

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
	routes.RegisterCRUDRoutes(protected, db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Run server
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
