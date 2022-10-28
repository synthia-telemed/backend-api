package datastore

import (
	"errors"
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
	FindByID(id uint) (*Notification, error)
	SetAsRead(id uint) error
	SetAllAsRead(patientID uint) error
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
	tx := g.db.Model(&Notification{}).Where("patient_id = ? AND is_read = ?", patientID, false).Count(&count)
	return int(count), tx.Error
}

func (g GormNotificationDataStore) ListLatest(patientID uint) ([]Notification, error) {
	var notifications []Notification
	tx := g.db.Where(&Notification{PatientID: patientID}).Order("created_at desc").Find(&notifications)
	return notifications, tx.Error
}

func (g GormNotificationDataStore) FindByID(id uint) (*Notification, error) {
	var notification Notification
	if err := g.db.First(&notification, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

func (g GormNotificationDataStore) SetAsRead(id uint) error {
	return g.db.Model(&Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (g GormNotificationDataStore) SetAllAsRead(patientID uint) error {
	return g.db.Model(&Notification{}).Where("patient_id = ?", patientID).Update("is_read", true).Error
}
