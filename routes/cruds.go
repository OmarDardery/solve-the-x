package routes

import (
	"fmt"
	"net/http"

	"github.com/OmarDardery/solve-the-x-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ------------------ HELPERS ------------------

// Get user from context and type-assert
func getUser[T any](c *gin.Context) (*T, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	u, ok := user.(*T)
	return u, ok
}

// Require role middleware helper
func requireRole(c *gin.Context, role string) bool {
	r, _ := c.Get("role")
	if r != role {
		c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("only %s can perform this action", role)})
		return false
	}
	return true
}

// Bind JSON and handle errors
func bindJSON[T any](c *gin.Context) (*T, bool) {
	var input T
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, false
	}
	return &input, true
}

// Convert param to uint
func uintFromParam(param string) uint {
	var id uint
	fmt.Sscanf(param, "%d", &id)
	return id
}

// ------------------ CRUD ROUTES ------------------

func RegisterCRUDRoutes(api *gin.RouterGroup, db *gorm.DB) {
	registerStudentRoutes(api.Group("/students"), db)
	registerProfessorRoutes(api.Group("/professors"), db)
	registerOpportunityRoutes(api.Group("/opportunities"), db)
	registerApplicationRoutes(api.Group("/applications"), db)
	registerCoinRoutes(api.Group("/coins"), db)
}

// ------------------ STUDENTS ------------------

func registerStudentRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/me", func(c *gin.Context) {
		if student, ok := getUser[models.Student](c); ok {
			c.JSON(http.StatusOK, student)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get student"})
		}
	})

	rg.PUT("/me", func(c *gin.Context) {
		student, ok := getUser[models.Student](c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get student"})
			return
		}

		updates, ok := bindJSON[map[string]interface{}](c)
		if !ok {
			return
		}

		updated, err := models.UpdateStudent(db, student.ID, *updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	rg.DELETE("/me", func(c *gin.Context) {
		student, ok := getUser[models.Student](c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get student"})
			return
		}
		if err := models.DeleteStudent(db, student.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "student deleted"})
	})
}

// ------------------ PROFESSORS ------------------

func registerProfessorRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/me", func(c *gin.Context) {
		if prof, ok := getUser[models.Professor](c); ok {
			c.JSON(http.StatusOK, prof)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get professor"})
		}
	})

	rg.PUT("/me", func(c *gin.Context) {
		prof, ok := getUser[models.Professor](c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get professor"})
			return
		}

		updates, ok := bindJSON[map[string]interface{}](c)
		if !ok {
			return
		}

		updated, err := models.UpdateProfessor(db, prof.ID, *updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updated)
	})

	rg.DELETE("/me", func(c *gin.Context) {
		prof, ok := getUser[models.Professor](c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get professor"})
			return
		}
		if err := models.DeleteProfessor(db, prof.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "professor deleted"})
	})
}

// ------------------ OPPORTUNITIES ------------------

func registerOpportunityRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/:id", func(c *gin.Context) {
		op, err := models.GetOpportunityByID(db, uintFromParam(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, op)
	})

	rg.GET("/me", func(c *gin.Context) {
		if !requireRole(c, "professor") {
			return
		}
		prof, _ := getUser[models.Professor](c)
		ops, err := models.GetOpportunitiesByProfessorID(db, prof.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ops)
	})

	// PUT / DELETE
	rg.PUT("/:id", updateOrDeleteOpportunity(db, "update"))
	rg.DELETE("/:id", updateOrDeleteOpportunity(db, "delete"))
}

func updateOrDeleteOpportunity(db *gorm.DB, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !requireRole(c, "professor") {
			return
		}
		prof, _ := getUser[models.Professor](c)
		id := uintFromParam(c.Param("id"))
		op, err := models.GetOpportunityByID(db, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if op.ProfessorID != prof.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot modify opportunities you don't own"})
			return
		}

		switch action {
		case "update":
			updates, ok := bindJSON[map[string]interface{}](c)
			if !ok {
				return
			}
			updated, err := models.UpdateOpportunity(db, id, *updates)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, updated)
		case "delete":
			if err := models.DeleteOpportunity(db, id); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "opportunity deleted"})
		}
	}
}

// ------------------ APPLICATIONS ------------------

func registerApplicationRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/opportunity/:id", func(c *gin.Context) {
		if !requireRole(c, "professor") {
			return
		}
		prof, _ := getUser[models.Professor](c)
		op, err := models.GetOpportunityByID(db, uintFromParam(c.Param("id")))
		if err != nil || op.ProfessorID != prof.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "cannot view applications"})
			return
		}
		apps, err := models.GetApplicationsByOpportunityID(db, op.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, apps)
	})

	rg.DELETE("/", func(c *gin.Context) {
		if !requireRole(c, "student") {
			return
		}
		student, _ := getUser[models.Student](c)
		req, ok := bindJSON[struct {
			OpportunityID uint `json:"opportunity_id"`
		}](c)
		if !ok {
			return
		}

		if err := student.DeleteApplication(db, req.OpportunityID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "application deleted"})
	})
}

// ------------------ COINS ------------------

func registerCoinRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	rg.GET("/me", func(c *gin.Context) {
		if !requireRole(c, "student") {
			return
		}
		student, _ := getUser[models.Student](c)
		coins, err := models.GetCoinsByStudentID(db, student.ID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, coins)
	})

	rg.PUT("/increment", func(c *gin.Context) {
		if !requireRole(c, "student") {
			return
		}
		student, _ := getUser[models.Student](c)
		req, ok := bindJSON[struct {
			Amount int `json:"amount"`
		}](c)
		if !ok {
			return
		}

		if err := models.IncrementCoins(db, student.ID, req.Amount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "coins incremented"})
	})
}
