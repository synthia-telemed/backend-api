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
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DateTime  time.Time      `json:"date_time"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	PatientID uint           `json:"patient_id" gorm:"not null"`
	Systolic  uint           `json:"systolic" gorm:"not null"`
	Diastolic uint           `json:"diastolic" gorm:"not null"`
	Pulse     uint           `json:"pulse" gorm:"not null"`
}

type Glucose struct {
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DateTime     time.Time      `json:"date_time"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	ID           uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	PatientID    uint           `json:"patient_id" gorm:"not null"`
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
