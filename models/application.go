package models

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	StudentID     uint         `json:"student_id" gorm:"not null"`
	Student       *Student     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:StudentID;references:ID"`
	OpportunityID uint         `json:"opportunity_id" gorm:"not null"`
	Opportunity   *Opportunity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:OpportunityID;references:ID"`
	Status        string       `json:"status" gorm:"type:TEXT CHECK(status IN ('pending','accepted','rejected'));not null;default:'pending'"`
}

const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
	StatusRejected = "rejected"
)

func (s Student) CreateApplication(db *gorm.DB, opportunityID uint) error {
	application := Application{
		StudentID:     s.ID,
		OpportunityID: opportunityID,
		Status:        StatusPending,
	}
	return db.Create(&application).Error
}
