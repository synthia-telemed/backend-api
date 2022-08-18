package payment

import "time"

type Client interface {
	CreateCustomer(patientID uint) (string, error)
	AddCreditCard(customerID, cardToken string) error
	ListCards(customerID string) ([]Card, error)
	PayWithCreditCard(customerID, cardID, refID string, amount int) (*Payment, error)
	IsOwnCreditCard(customerID, cardID string) (bool, error)
}

type Card struct {
	ID         string `json:"id"`
	LastDigits string `json:"last_digits"`
	Brand      string `json:"brand"`
	Default    bool   `json:"default"`
}

type Payment struct {
	ID             string    `json:"id"`
	Amount         int       `json:"amount"`
	Currency       string    `json:"currency"`
	CreatedAt      time.Time `json:"created_at"`
	Paid           bool      `json:"paid"`
	Success        bool      `json:"success"`
	FailureCode    *string   `json:"failure_code"`
	FailureMessage *string   `json:"failure_message"`
}
