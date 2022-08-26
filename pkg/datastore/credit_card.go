package datastore

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type CreditCard struct {
	ID          uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Last4Digits string         `json:"last_4_digits"`
	Brand       string         `json:"brand"`
	PatientID   uint           `json:"patient_id" gorm:"not null"`
	Name        string         `json:"name"`
	CardID      string         `json:"-"`
}

type CreditCardDataStore interface {
	Create(card *CreditCard) error
	FindByID(id uint) (*CreditCard, error)
	FindByPatientID(patientID uint) ([]CreditCard, error)
	IsOwnCreditCard(patientID, cardID uint) (bool, error)
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

func (g GormCreditCardDataStore) IsOwnCreditCard(patientID, id uint) (bool, error) {
	var c CreditCard
	if err := g.db.Where(&CreditCard{PatientID: patientID, ID: id}).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (g GormCreditCardDataStore) Delete(id uint) error {
	return g.db.Delete(&CreditCard{}, id).Error
}

func (g GormCreditCardDataStore) FindByID(id uint) (*CreditCard, error) {
	var c CreditCard
	if err := g.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
