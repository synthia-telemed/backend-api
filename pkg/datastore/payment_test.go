package datastore_test

import (
	"github.com/caarlos0/env/v6"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
)

var _ = Describe("Payment Datastore", Ordered, func() {
	var (
		paymentDataStore datastore.PaymentDataStore
		db               *gorm.DB

		patient    *datastore.Patient
		creditCard *datastore.CreditCard
	)

	BeforeAll(func() {
		var config datastore.Config
		Expect(env.Parse(&config)).To(Succeed())
		var err error
		db, err = gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		rand.Seed(GinkgoRandomSeed())
		Expect(db.AutoMigrate(&datastore.Patient{}, &datastore.CreditCard{})).To(Succeed())
		var err error
		paymentDataStore, err = datastore.NewGormPaymentDataStore(db)
		Expect(err).To(BeNil())

		patient = generatePatient()
		Expect(db.Create(patient).Error).To(Succeed())
		creditCard = generateCreditCard(patient.ID)
		Expect(db.Create(creditCard).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{}, &datastore.CreditCard{}, &datastore.Payment{})).To(Succeed())
	})

	Context("Create payment", func() {
		DescribeTable("create credit card payment",
			func(status datastore.PaymentStatus) {
				p := generateCreditCardPayment(status, creditCard.ID)
				Expect(paymentDataStore.Create(p)).To(Succeed())
				Expect(p.ID).ToNot(BeZero())
				Expect(p.CreatedAt).ToNot(BeZero())
				if status == datastore.SuccessPaymentStatus {
					Expect(p.PaidAt).ToNot(BeZero())
				} else {
					Expect(p.PaidAt).To(BeZero())
				}
			},
			Entry("success payment", datastore.SuccessPaymentStatus),
			Entry("failed payment", datastore.FailedPaymentStatus),
			Entry("pending payment", datastore.PendingPaymentStatus),
		)
	})

	Context("FindByInvoiceID", func() {
		Context("Credit Card Payment", func() {
			var payment *datastore.Payment
			BeforeEach(func() {
				payment = generateCreditCardPayment(datastore.SuccessPaymentStatus, creditCard.ID)
				Expect(db.Create(payment).Error).To(Succeed())
			})
			When("payment is not found", func() {
				It("should return nil with no error", func() {
					p, err := paymentDataStore.FindByInvoiceID(uuid.New().String())
					Expect(err).To(BeNil())
					Expect(p).To(BeNil())
				})
			})
			When("payment is found", func() {
				It("should return payment with credit card preloaded", func() {
					p, err := paymentDataStore.FindByInvoiceID(payment.InvoiceID)
					Expect(err).To(BeNil())
					Expect(p).ToNot(BeNil())
					Expect(*p.CreditCardID).To(Equal(creditCard.ID))
					Expect(p.CreditCard.CardID).To(Equal(creditCard.CardID))
				})
			})
		})
	})
})
