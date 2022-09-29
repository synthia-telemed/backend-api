package datastore

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type CreditCard struct {
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Last4Digits string         `json:"last_4_digits"`
	Brand       string         `json:"brand"`
	Name        string         `json:"name"`
	CardID      string         `json:"-"`
	Payments    []Payment      `json:"-" gorm:"foreignKey:CreditCardID"`
	ID          uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	PatientID   uint           `json:"patient_id" gorm:"not null"`
	IsDefault   bool           `json:"is_default"`
	Expiry      string         `json:"expiry"`
}

type CreditCardDataStore interface {
	Create(card *CreditCard) error
	FindByID(id uint) (*CreditCard, error)
	FindByPatientID(patientID uint) ([]CreditCard, error)
	IsOwnCreditCard(patientID, cardID uint) (bool, error)
	Count(patientID uint) (int, error)
	SetAllToNonDefault(patientID uint) error
	SetIsDefault(cardID uint, isDefault bool) error
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

func (g GormCreditCardDataStore) Count(patientID uint) (int, error) {
	var count int64
	if err := g.db.Model(&CreditCard{}).Where(&CreditCard{PatientID: patientID}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (g GormCreditCardDataStore) SetAllToNonDefault(patientID uint) error {
	return g.db.Model(&CreditCard{}).Where(&CreditCard{PatientID: patientID}).Update("is_default", false).Error
}

func (g GormCreditCardDataStore) SetIsDefault(cardID uint, isDefault bool) error {
	return g.db.Model(&CreditCard{}).Where(&CreditCard{ID: cardID}).Update("is_default", isDefault).Error
}
