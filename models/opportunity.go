package models

import "gorm.io/gorm"

// Opportunity represents a research/project/internship post.
type Opportunity struct {
	gorm.Model
	ProfessorID  uint       `json:"professor_id" gorm:"not null"`
	Professor    *Professor `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ProfessorID;references:ID"`
	Name         string     `json:"name" gorm:"not null"`
	Details      string     `json:"details" gorm:"type:text"`
	Requirements string     `json:"requirements" gorm:"type:text"`

	// Relationship: each opportunity can have multiple tags
	RequirementTags []Tag `gorm:"many2many:opportunity_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Enforce type constraint
	Type string `json:"type" gorm:"type:TEXT CHECK(type IN ('research','project','internship'));not null"`
}
