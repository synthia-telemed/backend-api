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

var _ = Describe("Doctor Datastore", Ordered, func() {
	var (
		db              *gorm.DB
		doctorDataStore datastore.DoctorDataStore
		doctors         []*datastore.Doctor
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
})

func generateDoctors(num int) []*datastore.Doctor {
	doctors := make([]*datastore.Doctor, num)
	for i := 0; i < num; i++ {
		doctors[i] = &datastore.Doctor{RefID: uuid.New().String()}
	}
	return doctors
}

func getRandomDoctor(docs []*datastore.Doctor) *datastore.Doctor {
	return docs[rand.Int()%len(docs)]
}
