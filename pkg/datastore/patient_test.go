package datastore_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
)

var _ = Describe("Patient Datastore", Ordered, func() {

	var (
		db               *gorm.DB
		patientDataStore datastore.PatientDataStore
		users            []datastore.Patient
	)

	BeforeAll(func() {
		var err error
		config := datastore.Config{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "synthia",
			SSLMode:  "disable",
			TimeZone: "Asia/Bangkok",
		}
		db, err = gorm.Open(postgres.Open(config.DSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		rand.Seed(GinkgoRandomSeed())
		var err error
		patientDataStore, err = datastore.NewGormPatientDataStore(db)
		Expect(err).To(BeNil())

		users = generatePatients(10)
		db.Create(&users)
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{})).To(BeNil())
	})

	When("FindByID", func() {
		It("should found the patient", func() {
			user := getRandomPatient(users)
			foundUser, err := patientDataStore.FindByID(user.ID)
			Expect(err).To(BeNil())
			Expect(foundUser.ID).To(Equal(user.ID))
			Expect(foundUser.RefID).To(Equal(user.RefID))
		})

		It("should return nil if patient not found", func() {
			foundUser, err := patientDataStore.FindByID(uint(rand.Uint32()))
			Expect(err).To(BeNil())
			Expect(foundUser).To(BeNil())
		})
	})

	When("FindByRefID", func() {
		It("should found the patient", func() {
			user := getRandomPatient(users)
			foundUser, err := patientDataStore.FindByRefID(user.RefID)
			Expect(err).To(BeNil())
			Expect(foundUser.ID).To(Equal(user.ID))
			Expect(foundUser.RefID).To(Equal(user.RefID))
		})

		It("should return nil if patient not found", func() {
			foundUser, err := patientDataStore.FindByRefID("not-exist")
			Expect(err).To(BeNil())
			Expect(foundUser).To(BeNil())
		})
	})

	When("Creating", func() {
		It("should create user", func() {
			user := generatePatient()
			err := patientDataStore.Create(&user)
			Expect(err).To(BeNil())
			Expect(user.ID).ToNot(BeZero())

			var foundUser datastore.Patient
			Expect(db.First(&foundUser, user.ID).Error).To(BeNil())
			Expect(foundUser.ID).To(Equal(user.ID))
			Expect(foundUser.RefID).To(Equal(user.RefID))
		})
	})

	When("FindOrCreate", func() {
		It("should create user", func() {
			patient := generatePatient()
			err := patientDataStore.FindOrCreate(&patient)
			Expect(err).To(BeNil())
			Expect(patient.ID).ToNot(BeZero())

			var foundUser datastore.Patient
			Expect(db.First(&foundUser, patient.ID).Error).To(BeNil())
			Expect(foundUser.ID).To(Equal(patient.ID))
			Expect(foundUser.RefID).To(Equal(patient.RefID))
		})

		It("should found user", func() {
			patient := getRandomPatient(users)
			err := patientDataStore.FindOrCreate(&patient)
			Expect(err).To(BeNil())
			Expect(patient.ID).ToNot(BeZero())
		})
	})
})

func generatePatient() datastore.Patient {
	return datastore.Patient{RefID: fmt.Sprintf("HN-%d", rand.Uint32())}
}

func generatePatients(num int) []datastore.Patient {
	users := make([]datastore.Patient, num)
	for i := 0; i < num; i++ {
		users[i] = generatePatient()
	}
	return users
}

func getRandomPatient(users []datastore.Patient) datastore.Patient {
	return users[rand.Int()%len(users)]
}
