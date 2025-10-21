package middleware

import (
	"fmt"
	"net/http"

	"github.com/OmarDardery/solve-the-x-backend/jwt_service"
	"github.com/OmarDardery/solve-the-x-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JWTMiddleware verifies JWT tokens and attaches user info (role + struct) to the Gin context
func JWTMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing Authorization header"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		claims, err := jwt_service.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Extract role from claims
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
			c.Abort()
			return
		}

		// Extract user ID
		idFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}
		id := uint(idFloat)

		// Fetch user based on role
		var user interface{}
		switch role {
		case "student":
			var student models.Student
			if err := db.First(&student, id).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Student not found"})
				c.Abort()
				return
			}
			user = &student

		case "professor":
			var professor models.Professor
			if err := db.First(&professor, id).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Professor not found"})
				c.Abort()
				return
			}
			user = &professor

		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid role: %s", role)})
			c.Abort()
			return
		}

		// Store both role and user in context
		c.Set("role", role)
		c.Set("user", user)

		c.Next()
	}
}
