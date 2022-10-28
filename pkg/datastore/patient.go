package datastore

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type BloodType string

type Patient struct {
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	PaymentCustomerID *string        `gorm:"unique"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	RefID             string         `json:"refID" gorm:"unique"`
	CreditCards       []CreditCard   `gorm:"foreignKey:PatientID"`
	ID                uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	Notification      []Notification `gorm:"foreignKey:PatientID"`
	NotificationToken string         `json:"-"`
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

func NewGormPatientDataStore(db *gorm.DB) (PatientDataStore, error) {
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
