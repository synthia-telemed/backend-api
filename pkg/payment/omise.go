package payment

import (
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type Client interface {
	CreateCustomerWithCard(patientID uint, cardToken string) (string, error)
}

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

func (c OmisePaymentClient) CreateCustomerWithCard(patientID uint, cardToken string) (string, error) {
	customer, createCustomer := &omise.Customer{}, &operations.CreateCustomer{
		Card: cardToken,
		Metadata: map[string]interface{}{
			"patient_id": patientID,
		},
	}
	if err := c.client.Do(customer, createCustomer); err != nil {
		return "", err
	}
	return customer.ID, nil
}
