package payment

type Client interface {
	CreateCustomer(patientID uint) (string, error)
	AddCreditCard(customerID, cardToken string) error
	ListCards(customerID string) ([]Card, error)
}

type Card struct {
	ID         string `json:"id"`
	LastDigits string `json:"last_digits"`
	Brand      string `json:"brand"`
}
