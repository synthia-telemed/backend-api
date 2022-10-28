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
	ListLatest(patientID uint) ([]Notification, error)
}

type GormNotificationDataStore struct {
	db *gorm.DB
}

func NewGormNotificationDataStore(db *gorm.DB) (NotificationDataStore, error) {
	return &GormNotificationDataStore{db: db}, db.AutoMigrate(&Notification{})
}

func (g GormNotificationDataStore) Create(notification *Notification) error {
	return g.db.Create(&notification).Error
}

func (g GormNotificationDataStore) CountUnRead(patientID uint) (int, error) {
	var count int64
	tx := g.db.Model(&Notification{}).Where(&Notification{PatientID: patientID, IsRead: false}).Count(&count)
	return int(count), tx.Error
}

func (g GormNotificationDataStore) ListLatest(patientID uint) ([]Notification, error) {
	var notifications []Notification
	tx := g.db.Where(&Notification{PatientID: patientID}).Find(&notifications)
	return notifications, tx.Error
}
