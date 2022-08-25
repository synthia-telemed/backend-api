package datastore_test

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
)

var _ = Describe("Credit Card Datastore", Ordered, func() {
	var (
		db                  *gorm.DB
		creditCardDataStore datastore.CreditCardDataStore
		patient             *datastore.Patient
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
		var err error
		creditCardDataStore, err = datastore.NewGormCreditCardDataStore(db)
		Expect(err).To(BeNil())
		Expect(db.AutoMigrate(&datastore.Patient{})).To(Succeed())
		patient = &datastore.Patient{RefID: fmt.Sprintf("ref-id-%d", rand.Int())}
		Expect(db.Create(patient).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.CreditCard{}, &datastore.Patient{})).To(Succeed())
	})

	Context("Create credit card", func() {
		var card *datastore.CreditCard
		BeforeEach(func() {
			card = generateCreditCard(patient.ID)
		})

		It("should create new credit card", func() {
			Expect(creditCardDataStore.Create(card)).To(Succeed())
			var retrievedCard datastore.CreditCard
			Expect(db.First(&retrievedCard, card.ID).Error).To(Succeed())
			Expect(card.Last4Digits).To(Equal(retrievedCard.Last4Digits))
		})

		When("patient ID is not valid", func() {
			BeforeEach(func() {
				card.PatientID = 0
			})
			It("should return error", func() {
				Expect(creditCardDataStore.Create(card)).ToNot(Succeed())
			})
		})
	})
})

func generateCreditCard(patientID uint) *datastore.CreditCard {
	return &datastore.CreditCard{
		Fingerprint: "fp",
		IsDefault:   true,
		Last4Digits: fmt.Sprintf("%d", rand.Intn(10000)),
		Brand:       "Visa",
		PatientID:   patientID,
		CardID:      fmt.Sprintf("test_%d", rand.Int()),
	}
}
