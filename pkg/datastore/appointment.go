package datastore

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Appointment struct {
	ID          uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	RefID       int            `json:"ref_id" gorm:"unique"`
	Duration    float64        `json:"duration"`
	StartedTime time.Time      `json:"started_time"`
}

type AppointmentDataStore interface {
	Create(appointment *Appointment) error
	FindByRefID(refID int) (*Appointment, error)
}

type GormAppointmentDataStore struct {
	db *gorm.DB
}

func NewGormAppointmentDataStore(db *gorm.DB) (AppointmentDataStore, error) {

	return &GormAppointmentDataStore{db: db}, db.AutoMigrate(&Appointment{})
}

func (g GormAppointmentDataStore) Create(appointment *Appointment) error {
	return g.db.Create(appointment).Error
}

func (g GormAppointmentDataStore) FindByRefID(refID int) (*Appointment, error) {
	var appointment Appointment
	if err := g.db.Where("ref_id = ?", refID).First(&appointment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &appointment, nil
}
