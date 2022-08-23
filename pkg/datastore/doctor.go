package datastore

import (
	"gorm.io/gorm"
	"time"
)

type Doctor struct {
	ID        uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RefID     string         `json:"refID" gorm:"unique"`
}

type DoctorDataStore interface {
	FindOrCreate(doctor *Doctor) error
}

type GormDoctorDataStore struct {
	db *gorm.DB
}

func NewGormDoctorDataStore(db *gorm.DB) (DoctorDataStore, error) {
	return &GormDoctorDataStore{db}, db.AutoMigrate(&Doctor{})
}

func (g GormDoctorDataStore) FindOrCreate(doctor *Doctor) error {
	return g.db.FirstOrCreate(doctor, doctor).Error
}
