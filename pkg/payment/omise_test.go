package payment_test

import (
	"github.com/caarlos0/env/v6"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"time"
)

var _ = Describe("Omise Payment Client", func() {
	var (
		client        *omise.Client
		paymentClient payment.Client
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
	})

	Context("Create customer", func() {
		var (
			token          *omise.Token
			createTokenOps *operations.CreateToken
			customerID     string
		)

		BeforeEach(func() {
			token, createTokenOps = &omise.Token{}, &operations.CreateToken{
				Name:            "John Doe (Testing)",
				Number:          "4242424242424242",
				ExpirationMonth: 12,
				ExpirationYear:  time.Now().Year(),
				SecurityCode:    "123",
				City:            "Bangkok",
				PostalCode:      "10500",
			}
			err := client.Do(token, createTokenOps)
			Expect(err).To(BeNil())
		})
		It("should create Omise's customer", func() {
			cusID, err := paymentClient.CreateCustomer(1)
			Expect(err).To(BeNil())
			Expect(cusID).ToNot(BeEmpty())
			Expect(cusID).To(MatchRegexp("cust(_test)?_[0-9a-z]+"))
			customerID = cusID
		})

		AfterEach(func() {
			cus, deleteCustomerOps := &omise.Customer{}, &operations.DestroyCustomer{CustomerID: customerID}
			err := client.Do(cus, deleteCustomerOps)
			Expect(err).To(BeNil())
		})
	})
})
