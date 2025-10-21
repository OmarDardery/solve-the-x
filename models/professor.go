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

// GetProfessorByID retrieves a professor by ID
func GetProfessorByID(db *gorm.DB, id uint) (*Professor, error) {
	var professor Professor
	if err := db.First(&professor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("professor not found")
		}
		return nil, err
	}
	return &professor, nil
}

// UpdateProfessor updates a professor’s fields and rehashes password if changed
func UpdateProfessor(db *gorm.DB, id uint, updates map[string]interface{}) (*Professor, error) {
	professor, err := GetProfessorByID(db, id)
	if err != nil {
		return nil, err
	}

	// Handle password update
	if newPasswordRaw, ok := updates["Password"]; ok {
		newPassword, ok := newPasswordRaw.(string)
		if !ok {
			return nil, errors.New("password must be a string")
		}

		hashed, err := HashPassword(newPassword)
		if err != nil {
			return nil, err
		}

		updates["Password"] = hashed
		updates["LastChangedPassword"] = time.Now()
	}

	if err := db.Model(professor).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Refresh record
	if err := db.First(professor, id).Error; err != nil {
		return nil, err
	}

	return professor, nil
}

// DeleteProfessor deletes a professor by ID
func DeleteProfessor(db *gorm.DB, id uint) error {
	professor, err := GetProfessorByID(db, id)
	if err != nil {
		return err
	}

	if err := db.Delete(professor).Error; err != nil {
		return err
	}

	return nil
}
