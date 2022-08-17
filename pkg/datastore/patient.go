package datastore

import (
	"errors"
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
	ID                uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	RefID             string         `json:"refID" gorm:"unique"`
	PaymentCustomerID *string        `gorm:"unique"`
	//BirthDate   time.Time      `json:"birthDate"`
	//BloodType   BloodType      `json:"bloodType"`
	//FirstnameEn string         `json:"firstname_en"`
	//FirstnameTh string         `json:"firstname_th"`
	//InitialEn   string         `json:"initial_en"`
	//InitialTh   string         `json:"initial_th"`
	//LastnameEn  string         `json:"lastname_en"`
	//LastnameTh  string         `json:"lastname_th"`
	//NationalID  *string        `json:"nationalId" gorm:"unique"`
	//PassportID  *string        `json:"passportId" gorm:"unique"`
	//Nationality string         `json:"nationality"`
	//PhoneNumber string         `json:"phoneNumber"`
	//Weight      float32        `json:"weight"`
	//Height      float32        `json:"height"`
}

type PatientDataStore interface {
	Create(patient *Patient) error
	FindByID(id uint) (*Patient, error)
	FindByRefID(refID string) (*Patient, error)
	FindOrCreate(patient *Patient) error
	Save(patient *Patient) error
	//FindByGovCredential(nationalID string) (*Patient, error)
}

type GormPatientDataStore struct {
	db *gorm.DB
}

func NewGormPatientDataStore(db *gorm.DB) (*GormPatientDataStore, error) {
	return &GormPatientDataStore{db}, db.AutoMigrate(&Patient{})
}

func (g GormPatientDataStore) Create(patient *Patient) error {
	return g.db.Create(patient).Error
}

func (g GormPatientDataStore) FindByID(id uint) (*Patient, error) {
	var patient Patient
	if err := g.db.First(&patient, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

func (g GormPatientDataStore) FindByRefID(refID string) (*Patient, error) {
	var patient Patient
	if err := g.db.First(&patient, "ref_id = ?", refID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &patient, nil
}

func (g GormPatientDataStore) FindOrCreate(patient *Patient) error {
	return g.db.FirstOrCreate(patient, patient).Error
}

func (g GormPatientDataStore) Save(patient *Patient) error {
	return g.db.Save(patient).Error
}

//func (g GormPatientDataStore) FindByGovCredential(cred string) (*Patient, error) {
//	var patient *Patient
//	err := g.db.Limit(1).Where("national_id = ?", cred).Or("passport_id = ?", cred).Find(&patient).Error
//	return patient, err
//}
