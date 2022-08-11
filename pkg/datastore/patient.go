package datastore

import (
	"gorm.io/gorm"
	"time"
)

const (
	AB BloodType = "AB"
	A  BloodType = "A"
	B  BloodType = "B"
	O  BloodType = "O"
)

type BloodType string

type Patient struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	RefID       string         `json:"refID" gorm:"unique"`
	BirthDate   time.Time      `json:"birthDate"`
	BloodType   BloodType      `json:"bloodType"`
	FirstnameEn string         `json:"firstname_en"`
	FirstnameTh string         `json:"firstname_th"`
	InitialEn   string         `json:"initial_en"`
	InitialTh   string         `json:"initial_th"`
	LastnameEn  string         `json:"lastname_en"`
	LastnameTh  string         `json:"lastname_th"`
	NationalID  string         `json:"nationalId" gorm:"unique"`
	Nationality string         `json:"nationality"`
	PhoneNumber string         `json:"phoneNumber"`
	Weight      float32        `json:"weight"`
	Height      float32        `json:"height"`
}

type PatientDataStore interface {
	New(patient *Patient) error
	FindByID(id uint) (*Patient, error)
	FindByNationalID(nationalID string) (*Patient, error)
}

type GormPatientDataStore struct {
	db *gorm.DB
}

func NewGormPatientDataStore(db *gorm.DB) *GormPatientDataStore {
	return &GormPatientDataStore{db}
}

func (g GormPatientDataStore) New(patient *Patient) error {
	return g.db.Create(patient).Error
}

func (g GormPatientDataStore) FindByID(id uint) (*Patient, error) {
	var patient *Patient
	err := g.db.Limit(1).Find(patient, id).Error
	return patient, err
}

func (g GormPatientDataStore) FindByNationalID(nationalID string) (*Patient, error) {
	var patient *Patient
	err := g.db.Limit(1).Where("national_id = ?", nationalID).Find(&patient).Error
	return patient, err
}
