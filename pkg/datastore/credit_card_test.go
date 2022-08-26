package datastore_test

import (
	"fmt"
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

var _ = Describe("Credit Card Datastore", Ordered, func() {
	var (
		db                  *gorm.DB
		creditCardDataStore datastore.CreditCardDataStore
		patient             *datastore.Patient
		card                *datastore.CreditCard
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

		patient = generatePatient()
		Expect(db.AutoMigrate(&datastore.Patient{})).To(Succeed())
		Expect(db.Create(patient).Error).To(Succeed())
		card = generateCreditCard(patient.ID)
		Expect(db.Create(card).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.CreditCard{}, &datastore.Patient{})).To(Succeed())
	})

	Context("Create credit card", func() {
		It("should create new credit card", func() {
			newCard := generateCreditCard(patient.ID)
			Expect(creditCardDataStore.Create(newCard)).To(Succeed())
			var retrievedCard datastore.CreditCard
			Expect(db.First(&retrievedCard, newCard.ID).Error).To(Succeed())
			Expect(newCard.Last4Digits).To(Equal(retrievedCard.Last4Digits))
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

	Context("Find by patientID", func() {
		When("patient has no card", func() {
			It("should return empty slice", func() {
				cards, err := creditCardDataStore.FindByPatientID(uint(rand.Uint32()))
				Expect(err).To(BeNil())
				Expect(cards).To(HaveLen(0))
			})
		})
		When("patient has cards", func() {
			It("should return slice of credit card", func() {
				cards, err := creditCardDataStore.FindByPatientID(patient.ID)
				Expect(err).To(BeNil())
				Expect(cards).To(HaveLen(1))
			})
		})
	})

	Context("Delete credit card", func() {
		It("should soft delete credit card", func() {
			Expect(creditCardDataStore.Delete(card.ID)).To(Succeed())
			var deletedCard datastore.CreditCard
			Expect(db.First(&deletedCard, card.ID).Error).To(Equal(gorm.ErrRecordNotFound))
			Expect(db.Unscoped().First(&deletedCard, card.ID).Error).To(Succeed())
			Expect(deletedCard.CardID).To(Equal(card.CardID))
		})
	})

	Context("IsOwnCreditCard", func() {
		When("patient doesn't own the card", func() {
			It("should return false with no error", func() {
				isOwn, err := creditCardDataStore.IsOwnCreditCard(uint(rand.Uint32()), card.ID)
				Expect(err).To(BeNil())
				Expect(isOwn).To(BeFalse())
			})
		})
		When("patient own the card", func() {
			It("should return true with no error", func() {
				isOwn, err := creditCardDataStore.IsOwnCreditCard(patient.ID, card.ID)
				Expect(err).To(BeNil())
				Expect(isOwn).To(BeTrue())
			})
		})
	})
})

func generateCreditCard(patientID uint) *datastore.CreditCard {
	return &datastore.CreditCard{
		IsDefault:   true,
		Last4Digits: fmt.Sprintf("%d", rand.Intn(10000)),
		Brand:       "Visa",
		PatientID:   patientID,
		CardID:      uuid.New().String(),
		Name:        "test_card",
	}
}
