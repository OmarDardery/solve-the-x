package models

import "gorm.io/gorm"

type Coins struct {
	gorm.Model
	Amount    int      `gorm:"not null"`
	StudentID uint     `gorm:"not null;uniqueIndex"`
	Student   *Student `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:StudentID;references:ID"`
}

func CreateCoins(db *gorm.DB, studentID uint) error {
	coins := Coins{
		Amount:    0,
		StudentID: studentID,
	}
	return db.Create(&coins).Error
}
