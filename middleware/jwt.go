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
		const bearerPrefix = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing Authorization header"})
			c.Abort()
			return
		}

		tokenString := authHeader[len(bearerPrefix):]
		claims, err := jwt_service.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		idFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}
		id := uint(idFloat)

		var user interface{}
		switch role {
		case "student":
			user, err = models.GetStudentByID(db, id)
		case "professor":
			user, err = models.GetProfessorByID(db, id)
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid role: %s", role)})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("%s not found", role)})
			c.Abort()
			return
		}

		c.Set("role", role)
		c.Set("user", user)
		c.Next()
	}
}
