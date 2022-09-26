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
	Expiry      string `json:"expiry"`
}

type Payment struct {
	FailureMessage *string `json:"failure_message"`
	FailureCode    *string `json:"failure_code"`
	Currency       string  `json:"currency"`
	ID             string  `json:"id"`
	Amount         int     `json:"amount"`
	Paid           bool    `json:"paid"`
	Success        bool    `json:"success"`
}
