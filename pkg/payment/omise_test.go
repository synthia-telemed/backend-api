package payment_test

import (
	"github.com/caarlos0/env/v6"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"math/rand"
	"time"
)

var _ = Describe("Omise Payment Client", func() {
	var (
		client         *omise.Client
		paymentClient  payment.Client
		testCustomerID string
	)

	BeforeEach(func() {
		var (
			err error
			c   payment.Config
		)
		Expect(env.Parse(&c)).To(BeNil())
		client, err = omise.NewClient(c.PublicKey, c.SecretKey)
		Expect(err).To(BeNil())
		paymentClient, err = payment.NewOmisePaymentClient(&payment.Config{PublicKey: c.PublicKey, SecretKey: c.SecretKey})
		Expect(err).To(BeNil())

		testCustomerID = createCustomer(client)
	})

	AfterEach(func() {
		deleteCustomer(client, testCustomerID)
	})

	Context("Create customer", func() {
		var customerID string
		It("should create Omise's customer", func() {
			cusID, err := paymentClient.CreateCustomer(1)
			Expect(err).To(BeNil())
			Expect(cusID).To(MatchRegexp("cust(_test)?_[0-9a-z]+"))
			customerID = cusID
		})

		AfterEach(func() {
			deleteCustomer(client, customerID)
		})
	})

	Context("Add credit card", func() {
		var cardToken string
		BeforeEach(func() {
			cardToken = createCardToken(client)
		})

		It("should add credit card to Omise's customer", func() {
			err := paymentClient.AddCreditCard(testCustomerID, cardToken)
			Expect(err).To(BeNil())

			customer, getCustomer := &omise.Customer{}, &operations.RetrieveCustomer{CustomerID: testCustomerID}
			Expect(client.Do(customer, getCustomer)).To(BeNil())
			Expect(customer.Cards).ToNot(BeNil())
			Expect(customer.Cards.Total).To(Equal(1))
		})
	})

	Context("List cards", func() {
		var n int
		BeforeEach(func() {
			n = 3
			for i := 0; i < n; i++ {
				t := createCardToken(client)
				attachCardToCustomer(client, testCustomerID, t)
			}
		})
		It("should list cards of Omise's customer", func() {
			cards, err := paymentClient.ListCards(testCustomerID)
			Expect(err).To(BeNil())
			Expect(len(cards)).To(Equal(n))
			for _, c := range cards {
				Expect(c).To(HaveExistingField("ID"))
				Expect(c).To(HaveExistingField("LastDigits"))
				Expect(c).To(HaveExistingField("Brand"))
			}
		})
	})
})

func createCardToken(client *omise.Client) string {
	token, createTokenOps := &omise.Token{}, &operations.CreateToken{
		Name:            "John Doe (Testing)",
		Number:          "4242424242424242",
		ExpirationMonth: 12,
		ExpirationYear:  time.Now().Year(),
		SecurityCode:    "123",
		City:            "Bangkok",
		PostalCode:      "10500",
	}
	Expect(client.Do(token, createTokenOps)).To(Succeed())
	return token.ID
}

func createCustomer(client *omise.Client) string {
	customer, createCustomerOps := &omise.Customer{}, &operations.CreateCustomer{
		Email:       "test@synthia.tech",
		Description: "test customer for unit testing",
		Metadata: map[string]interface{}{
			"patient_id":       rand.Uint32(),
			"testing_customer": true,
		},
	}
	Expect(client.Do(customer, createCustomerOps)).To(Succeed())
	return customer.ID
}

func deleteCustomer(client *omise.Client, customerID string) {
	deleteCustomerOps := &operations.DestroyCustomer{CustomerID: customerID}
	Expect(client.Do(nil, deleteCustomerOps)).To(Succeed())
}

func attachCardToCustomer(client *omise.Client, customerID, cardToken string) {
	addCardOps := &operations.UpdateCustomer{
		CustomerID: customerID,
		Card:       cardToken,
	}
	Expect(client.Do(nil, addCardOps)).To(Succeed())
}
