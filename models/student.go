package models

import (
	"errors"
	"time"

	jwt_service "github.com/OmarDardery/solve-the-x-backend/jwt_service"
	mail_service "github.com/OmarDardery/solve-the-x-backend/mail_service"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	FirstName           string
	LastName            string
	Email               string `gorm:"unique"`
	Password            string
	LastChangedPassword time.Time
	Tags                []Tag `gorm:"many2many:student_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Coins               Coins `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Generate JWT for the student
func (s *Student) GetJWT() (string, error) {
	return jwt_service.GenerateJWT(s.ID, s.Email, "student")
}

// AuthenticateStudent checks credentials and returns the student if valid
func AuthenticateStudent(db *gorm.DB, email, password string) (*Student, error) {
	var student Student

	if err := db.Where("email = ?", email).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}

	if !CheckPasswordHash(password, student.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &student, nil
}

// CreateStudent registers a new student and automatically creates a Coins record
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

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(student).Error; err != nil {
			return err
		}

		if err := CreateCoins(tx, student.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s Student) Notify(subject, content string) error {
	return mail_service.SendNotification(s.Email, subject, content)
}
