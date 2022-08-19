package datastore

import (
	"gorm.io/gorm"
	"time"
)

type MeasurementDataStore interface {
	CreateBloodPressure(bp *BloodPressure) error
	CreateGlucose(g *Glucose) error
}

type BloodPressure struct {
	ID        uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PatientID uint           `json:"patient_id"`
	DateTime  time.Time      `json:"date_time"`
	Systolic  uint           `json:"systolic" gorm:"not null"`
	Diastolic uint           `json:"diastolic" gorm:"not null"`
	Pulse     uint           `json:"pulse" gorm:"not null"`
}

type Glucose struct {
	ID           uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	PatientID    uint           `json:"patient_id"`
	DateTime     time.Time      `json:"date_time"`
	Value        uint           `json:"value" gorm:"not null"`
	IsBeforeMeal bool           `json:"is_before_meal" gorm:"not null"`
}

type GormMeasurementDataStore struct {
	db *gorm.DB
}

func NewGormMeasurementDataStore(db *gorm.DB) (MeasurementDataStore, error) {
	return &GormMeasurementDataStore{db: db}, db.AutoMigrate(&BloodPressure{}, &Glucose{})
}

func (d GormMeasurementDataStore) CreateBloodPressure(bp *BloodPressure) error {
	return d.db.Create(bp).Error
}

func (d GormMeasurementDataStore) CreateGlucose(g *Glucose) error {
	return d.db.Create(g).Error
}
