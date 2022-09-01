package payment

type Client interface {
	CreateCustomer(patientID uint) (string, error)
	AddCreditCard(customerID, cardToken string) (*Card, error)
	RemoveCreditCard(customerID, cardID string) error
	PayWithCreditCard(customerID, cardID, refID string, amount int) (*Payment, error)
}

type Card struct {
	ID          string `json:"id"`
	Last4Digits string `json:"last_4_digits"`
	Brand       string `json:"brand"`
}

type Payment struct {
	ID             string  `json:"id"`
	Amount         int     `json:"amount"`
	Currency       string  `json:"currency"`
	Paid           bool    `json:"paid"`
	Success        bool    `json:"success"`
	FailureCode    *string `json:"failure_code"`
	FailureMessage *string `json:"failure_message"`
}
