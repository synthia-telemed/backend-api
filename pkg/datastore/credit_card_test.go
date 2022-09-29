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

var _ = Describe("Credit Card Datastore", Ordered, func() {
	var (
		db                  *gorm.DB
		creditCardDataStore datastore.CreditCardDataStore
		patient             *datastore.Patient
		card                *datastore.CreditCard
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
		var err error
		creditCardDataStore, err = datastore.NewGormCreditCardDataStore(db)
		Expect(err).To(BeNil())

		patient = generatePatient()
		Expect(db.AutoMigrate(&datastore.Patient{})).To(Succeed())
		Expect(db.Create(patient).Error).To(Succeed())
		card = generateCreditCard(patient.ID, true)
		Expect(db.Create(card).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.CreditCard{}, &datastore.Patient{})).To(Succeed())
	})

	Context("Create credit card", func() {
		It("should create new credit card", func() {
			newCard := generateCreditCard(patient.ID, false)
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

	Context("FindByID", func() {
		When("card is not existed", func() {
			It("should return nil with no error", func() {
				c, err := creditCardDataStore.FindByID(uint(rand.Uint32()))
				Expect(err).To(BeNil())
				Expect(c).To(BeNil())
			})
		})
		When("card is existed", func() {
			It("should return the card with no error", func() {
				c, err := creditCardDataStore.FindByID(card.ID)
				Expect(err).To(BeNil())
				Expect(c.ID).To(Equal(card.ID))
			})
		})
	})

	Context("Count credit card", func() {
		When("patient has no credit card", func() {
			It("should return 0", func() {
				isFirst, err := creditCardDataStore.Count(uint(rand.Uint32()))
				Expect(err).To(BeNil())
				Expect(isFirst).To(BeZero())
			})
		})
		When("patient has a credit card", func() {
			It("should return more than 0", func() {
				isFirst, err := creditCardDataStore.Count(patient.ID)
				Expect(err).To(BeNil())
				Expect(isFirst).ToNot(BeZero())
			})
		})
	})

	Context("SetAllToNonDefault", func() {
		BeforeEach(func() {
			c := generateCreditCard(patient.ID, true)
			Expect(db.Create(c).Error).To(BeNil())
		})
		It("should set isDefault of all cards to be false", func() {
			Expect(creditCardDataStore.SetAllToNonDefault(patient.ID)).To(Succeed())
			var cards []datastore.CreditCard
			Expect(db.Where(&datastore.CreditCard{PatientID: patient.ID}).Find(&cards).Error).To(BeNil())
			for _, c := range cards {
				Expect(c.IsDefault).To(BeFalse())
			}
		})
	})

	DescribeTable("SetIsDefault", func(initialIsDefault bool) {
		card := generateCreditCard(patient.ID, initialIsDefault)
		Expect(db.Create(card).Error).To(BeNil())
		Expect(creditCardDataStore.SetIsDefault(card.ID, !initialIsDefault)).To(Succeed())
		var retrievedCard datastore.CreditCard
		Expect(db.Where(card.ID).First(&retrievedCard).Error).To(BeNil())
		Expect(retrievedCard.IsDefault).To(Equal(!initialIsDefault))
	},
		Entry("Set to false from true", true),
		Entry("Set to true from false", false))
})
