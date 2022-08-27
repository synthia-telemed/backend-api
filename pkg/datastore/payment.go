package datastore

import (
	"gorm.io/gorm"
	"time"
)

type PaymentMethod string
type PaymentStatus string

const (
	CreditCardPaymentMethod PaymentMethod = "credit_card"
	SuccessPaymentStatus    PaymentStatus = "success"
	FailedPaymentStatus     PaymentStatus = "failed"
	PendingPaymentStatus    PaymentStatus = "pending"
)

type Payment struct {
	ID           uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Method       PaymentMethod  `json:"method" gorm:"not null"`
	Amount       float64        `json:"amount" gorm:"not null"`
	PaidAt       *time.Time     `json:"timestamp"`
	ChargeID     string         `json:"charge_id" gorm:"not null"`
	InvoiceID    string         `json:"invoice_id" gorm:"not null"`
	Status       PaymentStatus  `json:"status" gorm:"not null"`
	CreditCard   *CreditCard    `json:"credit_card"`
	CreditCardID *uint
}

type PaymentDataStore interface {
	Create(payment *Payment) error
}

type GormPaymentDataStore struct {
	db *gorm.DB
}

func NewGormPaymentDataStore(db *gorm.DB) (PaymentDataStore, error) {
	return &GormPaymentDataStore{db: db}, db.AutoMigrate(&Payment{})
}

func (g GormPaymentDataStore) Create(payment *Payment) error {
	return g.db.Create(payment).Error
}
