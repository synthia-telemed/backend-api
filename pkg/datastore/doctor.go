package datastore

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Doctor struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RefID     string         `json:"refID" gorm:"unique"`
	ID        uint           `json:"id" gorm:"autoIncrement,primaryKey"`
}

type DoctorDataStore interface {
	FindOrCreate(doctor *Doctor) error
	FindByID(id uint) (*Doctor, error)
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

func (g GormDoctorDataStore) FindByID(id uint) (*Doctor, error) {
	var doc Doctor
	if err := g.db.First(&doc, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}
