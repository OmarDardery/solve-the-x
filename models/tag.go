package models

import (
	"errors"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
}

func CreateTag(db *gorm.DB, name, description string) error {
	var existing Tag
	if err := db.Where("name = ?", name).First(&existing).Error; err == nil {
		return nil // Tag already exists, do nothing
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	tag := &Tag{
		Name:        name,
		Description: description,
	}
	return db.Create(tag).Error
}
