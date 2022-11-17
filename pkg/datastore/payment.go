package datastore

import (
	"errors"
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
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	PaidAt       *time.Time     `json:"paid_at"`
	CreditCard   *CreditCard    `json:"credit_card"`
	CreditCardID *uint          `json:"credit_card_id"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	ChargeID     string         `json:"-" gorm:"not null"`
	Status       PaymentStatus  `json:"status" gorm:"not null"`
	Method       PaymentMethod  `json:"method" gorm:"not null"`
	ID           uint           `json:"id" gorm:"autoIncrement,primaryKey"`
	InvoiceID    int            `json:"invoice_id" gorm:"not null,unique"`
	Amount       float64        `json:"amount" gorm:"not null"`
}

type PaymentDataStore interface {
	Create(payment *Payment) error
	FindLatestByInvoiceIDAndStatus(invoiceID int, status PaymentStatus) (*Payment, error)
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

func (g GormPaymentDataStore) FindLatestByInvoiceIDAndStatus(invoiceID int, status PaymentStatus) (*Payment, error) {
	var payment Payment
	tx := g.db.Preload("CreditCard", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Where(&Payment{InvoiceID: invoiceID, Status: status}).Order("created_at").First(&payment)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &payment, nil
}
