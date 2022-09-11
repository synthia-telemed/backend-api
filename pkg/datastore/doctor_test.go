package datastore_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
)

var _ = Describe("Doctor Datastore", Ordered, func() {
	var (
		db              *gorm.DB
		doctorDataStore datastore.DoctorDataStore
		doctors         []*datastore.Doctor
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
		doctorDataStore, err = datastore.NewGormDoctorDataStore(db)
		Expect(err).To(BeNil())

		doctors = generateDoctors(10)
		Expect(db.Create(&doctors).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Doctor{})).To(Succeed())
	})

	Context("FindOrCreate", func() {
		var (
			refID  string
			doctor *datastore.Doctor
		)
		BeforeEach(func() {
			refID = uuid.New().String()
			doctor = &datastore.Doctor{RefID: refID}
		})
		JustBeforeEach(func() {
			Expect(doctorDataStore.FindOrCreate(doctor)).To(Succeed())
		})

		When("doctor is not in the db", func() {
			It("should create", func() {
				Expect(doctor.ID).ToNot(BeZero())
				Expect(doctor.RefID).To(Equal(refID))
				Expect(db.First(&datastore.Doctor{}, doctor.ID).Error).To(Succeed())
			})
		})
		When("doctor is existed", func() {
			BeforeEach(func() {
				Expect(db.Create(doctor).Error).To(Succeed())
			})
			It("should found doctor", func() {
				Expect(doctor.ID).ToNot(BeZero())
				Expect(doctor.RefID).To(Equal(refID))
			})
		})
	})

	Context("FindByID", func() {
		When("doctor is found", func() {
			It("should return doctor with no error", func() {
				d := doctors[3]
				foundDoc, err := doctorDataStore.FindByID(d.ID)
				Expect(err).To(BeNil())
				Expect(foundDoc.ID).To(Equal(d.ID))
				Expect(foundDoc.RefID).To(Equal(d.RefID))
			})
		})
		When("doctor is not found", func() {
			It("should return nil with no error", func() {
				foundDoc, err := doctorDataStore.FindByID(uint(rand.Uint32()))
				Expect(err).To(BeNil())
				Expect(foundDoc).To(BeNil())
			})
		})
	})
})
