package datastore

import "gorm.io/gorm"

type Patient struct {
	gorm.Model
	CitizenID string `gorm:"size:13;uniqueIndex"`
	RefID     string `gorm:"unique"`
}

type PatientDataStore interface {
	Create(patient *Patient) error
}
