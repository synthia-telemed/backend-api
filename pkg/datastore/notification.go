package datastore

import (
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	IsRead    bool           `json:"is_read"`
	PatientID uint           `json:"patient_id"`
}

type NotificationDataStore interface {
	Create(notification *Notification) error
	CountUnRead(patientID uint) (int, error)
	ListLatest(patientID uint) ([]*Notification, error)
}
