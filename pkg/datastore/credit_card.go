package datastore

import (
	"gorm.io/gorm"
	"time"
)

type CreditCard struct {
	ID          uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	IsDefault   bool           `json:"is_default"`
	Last4Digits string         `json:"last_4_digits"`
	Brand       string         `json:"brand"`
	PatientID   uint           `json:"patient_id" gorm:"not null"`
	CardID      string
}

type CreditCardDataStore interface {
	Create(card *CreditCard) error
	FindByPatientID(patientID uint) ([]CreditCard, error)
	Delete(id uint) error
}

type GormCreditCardDataStore struct {
	db *gorm.DB
}

func NewGormCreditCardDataStore(db *gorm.DB) (CreditCardDataStore, error) {
	return &GormCreditCardDataStore{db: db}, db.AutoMigrate(&CreditCard{})
}

func (g GormCreditCardDataStore) Create(card *CreditCard) error {
	return g.db.Create(card).Error
}

func (g GormCreditCardDataStore) FindByPatientID(patientID uint) ([]CreditCard, error) {
	var cards []CreditCard
	if err := g.db.Where(&CreditCard{PatientID: patientID}).Find(&cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func (g GormCreditCardDataStore) Delete(id uint) error {
	return g.db.Delete(&CreditCard{}, id).Error
}
