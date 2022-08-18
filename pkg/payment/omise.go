package payment

import (
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type Config struct {
	PublicKey string `env:"OMISE_PUBLIC_KEY,required"`
	SecretKey string `env:"OMISE_SECRET_KEY,required"`
}

type OmisePaymentClient struct {
	client *omise.Client
}

func NewOmisePaymentClient(c *Config) (*OmisePaymentClient, error) {
	client, err := omise.NewClient(c.PublicKey, c.SecretKey)
	if err != nil {
		return nil, err
	}
	return &OmisePaymentClient{client: client}, nil
}

func (c OmisePaymentClient) CreateCustomer(patientID uint) (string, error) {
	customer, createCustomer := &omise.Customer{}, &operations.CreateCustomer{
		Metadata: map[string]interface{}{
			"patient_id": patientID,
		},
	}
	if err := c.client.Do(customer, createCustomer); err != nil {
		return "", err
	}
	return customer.ID, nil
}

func (c OmisePaymentClient) AddCreditCard(customerID, cardToken string) error {
	addCardOps := &operations.UpdateCustomer{
		CustomerID: customerID,
		Card:       cardToken,
	}
	return c.client.Do(nil, addCardOps)
}

func (c OmisePaymentClient) ListCards(customerID string) ([]Card, error) {
	customer, retrieveCustomerOps := &omise.Customer{}, &operations.RetrieveCustomer{CustomerID: customerID}
	if err := c.client.Do(customer, retrieveCustomerOps); err != nil {
		return nil, err
	}
	cards := make([]Card, customer.Cards.Total)
	for i, c := range customer.Cards.Data {
		cards[i] = Card{
			ID:         c.ID,
			LastDigits: c.LastDigits,
			Brand:      c.Brand,
			Default:    customer.DefaultCard == c.ID,
		}
	}
	return cards, nil
}

func (c OmisePaymentClient) IsOwnCreditCard(customerID, cardID string) (bool, error) {
	card, retrieveCardOps := &omise.Card{}, &operations.RetrieveCard{
		CustomerID: customerID,
		CardID:     cardID,
	}
	if err := c.client.Do(card, retrieveCardOps); err != nil {
		e := err.(*omise.Error)
		if e.Code == "not_found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c OmisePaymentClient) PayWithCreditCard(customerID, cardID, refID string, amount int) (*Payment, error) {
	charge, createChargeOps := &omise.Charge{}, &operations.CreateCharge{
		Customer:    customerID,
		Card:        cardID,
		Amount:      int64(amount),
		Currency:    "THB",
		DontCapture: false,
		Metadata:    map[string]interface{}{"ref_id": refID},
	}
	if err := c.client.Do(charge, createChargeOps); err != nil {
		return nil, err
	}
	return &Payment{
		ID:             charge.ID,
		Amount:         int(charge.Amount),
		Currency:       charge.Currency,
		CreatedAt:      charge.Created,
		Paid:           charge.Paid,
		Success:        charge.Status == omise.ChargeSuccessful,
		FailureCode:    charge.FailureCode,
		FailureMessage: charge.FailureMessage,
	}, nil
}
