package models

import (
	"errors"

	"gorm.io/gorm"
)

type Coins struct {
	gorm.Model
	Amount    int      `gorm:"not null"`
	StudentID uint     `gorm:"not null;uniqueIndex"`
	Student   *Student `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:StudentID;references:ID"`
}

// CreateCoins creates a new Coins record for a student
func CreateCoins(db *gorm.DB, studentID uint) error {
	coins := Coins{
		Amount:    0,
		StudentID: studentID,
	}
	return db.Create(&coins).Error
}

// GetCoinsByStudentID retrieves the coin record for a student
func GetCoinsByStudentID(db *gorm.DB, studentID uint) (*Coins, error) {
	var coins Coins
	if err := db.Where("student_id = ?", studentID).First(&coins).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("coins record not found")
		}
		return nil, err
	}
	return &coins, nil
}

// IncrementCoins increases a student's coins by a specified amount
func IncrementCoins(db *gorm.DB, studentID uint, amount int) error {
	coins, err := GetCoinsByStudentID(db, studentID)
	if err != nil {
		return err
	}

	coins.Amount += amount
	if err := db.Save(coins).Error; err != nil {
		return err
	}

	return nil
}

// DecrementCoins decreases a student's coins by a specified amount (with validation)
func DecrementCoins(db *gorm.DB, studentID uint, amount int) error {
	coins, err := GetCoinsByStudentID(db, studentID)
	if err != nil {
		return err
	}

	if coins.Amount < amount {
		return errors.New("insufficient coins")
	}

	coins.Amount -= amount
	if err := db.Save(coins).Error; err != nil {
		return err
	}

	return nil
}
