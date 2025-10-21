package models

import (
	"errors"
	"time"

	jwt_service "github.com/OmarDardery/solve-the-x-backend/jwt_service"
	mail_service "github.com/OmarDardery/solve-the-x-backend/mail_service"
	"gorm.io/gorm"
)

type Professor struct {
	gorm.Model
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email" gorm:"unique"`
	Password            string    `json:"password"`
	LastChangedPassword time.Time `json:"last_changed_password"`
}

// Generate JWT for the professor
func (p *Professor) GetJWT() (string, error) {
	return jwt_service.GenerateJWT(p.ID, p.Email, "professor")
}

// AuthenticateProfessor checks credentials and returns the professor if valid
func AuthenticateProfessor(db *gorm.DB, email, password string) (*Professor, error) {
	var professor Professor

	// Find by email
	if err := db.Where("email = ?", email).First(&professor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("professor not found")
		}
		return nil, err
	}

	// Check password hash
	if !CheckPasswordHash(password, professor.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &professor, nil
}

// CreateProfessor registers a new professor with hashed password
func CreateProfessor(db *gorm.DB, firstName, lastName, email, password string) error {
	var existing Professor
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	professor := &Professor{
		FirstName:           firstName,
		LastName:            lastName,
		Email:               email,
		Password:            passwordHash,
		LastChangedPassword: time.Now(),
	}

	return db.Create(professor).Error
}

func (p Professor) Notify(subject, content string) error {
	return mail_service.SendNotification(p.Email, subject, content)
}
