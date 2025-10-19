package models

import (
	"errors"
	"time"

	"github.com/OmarDardery/solve-the-x-backend/middleware"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Email               string    `json:"email" gorm:"unique"`
	Password            string    `json:"password"`
	LastChangedPassword time.Time `json:"last_changed_password"`
}

// Generate JWT for the student
func (s *Student) GetJWT() (string, error) {
	return middleware.GenerateJWT(s.ID, s.Email, "student")
}

// AuthenticateStudent checks credentials and returns the student if valid
func AuthenticateStudent(db *gorm.DB, email, password string) (*Student, error) {
	var student Student

	// Find by email
	if err := db.Where("email = ?", email).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}

	// Check password hash
	if !CheckPasswordHash(password, student.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &student, nil
}

// CreateStudent registers a new student with hashed password
func CreateStudent(db *gorm.DB, firstName, lastName, email, password string) error {
	var existing Student
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	student := &Student{
		FirstName:           firstName,
		LastName:            lastName,
		Email:               email,
		Password:            passwordHash,
		LastChangedPassword: time.Now(),
	}

	return db.Create(student).Error
}

func (s Student) Notify(subject, content string) error {
	return middleware.SendNotification(s.Email, subject, content)
}
