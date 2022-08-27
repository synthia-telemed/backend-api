package datastore_test

import (
	"github.com/caarlos0/env/v6"
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
		It("should create succeed payment", func() {
			p := generateCreditCardPayment(datastore.SuccessPaymentStatus, creditCard.ID)
			Expect(paymentDataStore.Create(p)).To(Succeed())
			Expect(p.ID).ToNot(BeZero())
			Expect(p.CreatedAt).ToNot(BeZero())
			Expect(p.PaidAt).ToNot(BeNil())
		})
		It("should create failed payment", func() {
			p := generateCreditCardPayment(datastore.FailedPaymentStatus, creditCard.ID)
			Expect(paymentDataStore.Create(p)).To(Succeed())
			Expect(p.ID).ToNot(BeZero())
			Expect(p.CreatedAt).ToNot(BeZero())
			Expect(p.PaidAt).To(BeNil())
		})
		It("should create pending payment", func() {
			p := generateCreditCardPayment(datastore.PendingPaymentStatus, creditCard.ID)
			Expect(paymentDataStore.Create(p)).To(Succeed())
			Expect(p.ID).ToNot(BeZero())
			Expect(p.CreatedAt).ToNot(BeZero())
			Expect(p.PaidAt).To(BeNil())
		})
	})
})
