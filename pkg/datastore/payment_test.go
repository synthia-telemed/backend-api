package datastore_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	pg "gorm.io/driver/postgres"
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
		var err error
		db, err = gorm.Open(pg.Open(postgres.Config.DSN()), &gorm.Config{
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
		creditCard = generateCreditCard(patient.ID, false)
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

	Context("FindLatestByInvoiceIDAndStatus", func() {
		Context("Credit Card Payment", func() {
			var payment *datastore.Payment
			BeforeEach(func() {
				payment = generateCreditCardPayment(datastore.SuccessPaymentStatus, creditCard.ID)
				Expect(db.Create(payment).Error).To(Succeed())
			})
			When("payment is not found", func() {
				It("should return nil with no error", func() {
					p, err := paymentDataStore.FindLatestByInvoiceIDAndStatus(int(rand.Int31()), datastore.SuccessPaymentStatus)
					Expect(err).To(BeNil())
					Expect(p).To(BeNil())
				})
			})
			When("payment is found", func() {
				It("should return payment with credit card preloaded", func() {
					p, err := paymentDataStore.FindLatestByInvoiceIDAndStatus(payment.InvoiceID, datastore.SuccessPaymentStatus)
					Expect(err).To(BeNil())
					Expect(p).ToNot(BeNil())
					Expect(*p.CreditCardID).To(Equal(creditCard.ID))
					Expect(p.CreditCard.CardID).To(Equal(creditCard.CardID))
				})
			})
			When("payment is found but credit card is deleted", func() {
				BeforeEach(func() {
					tx := db.Delete(creditCard)
					Expect(tx.Error).To(BeNil())
					Expect(tx.RowsAffected).To(BeEquivalentTo(1))
				})
				It("should return payment with credit card preloaded", func() {
					p, err := paymentDataStore.FindLatestByInvoiceIDAndStatus(payment.InvoiceID, datastore.SuccessPaymentStatus)
					Expect(err).To(BeNil())
					Expect(p).ToNot(BeNil())
					Expect(*p.CreditCardID).To(Equal(creditCard.ID))
					Expect(p.CreditCard.CardID).To(Equal(creditCard.CardID))
				})
			})
		})
	})
})
