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
		patients         []*datastore.Patient
	)

	BeforeAll(func() {
		config := datastore.Config{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "synthia",
			SSLMode:  "disable",
			TimeZone: "Asia/Bangkok",
		}
		var err error
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

		patients = generatePatients(10)
		err = db.Create(&patients).Error
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{})).To(BeNil())
	})

	When("FindByID", func() {
		It("should found the patient", func() {
			patient := getRandomPatient(patients)
			foundPatient, err := patientDataStore.FindByID(patient.ID)
			Expect(err).To(BeNil())
			Expect(foundPatient.ID).To(Equal(patient.ID))
			Expect(foundPatient.RefID).To(Equal(patient.RefID))
		})

		It("should return nil if patient not found", func() {
			foundPatient, err := patientDataStore.FindByID(getRandomID())
			Expect(err).To(BeNil())
			Expect(foundPatient).To(BeNil())
		})
	})

	When("FindByRefID", func() {
		It("should found the patient", func() {
			patient := getRandomPatient(patients)
			foundPatient, err := patientDataStore.FindByRefID(patient.RefID)
			Expect(err).To(BeNil())
			Expect(foundPatient.ID).To(Equal(patient.ID))
			Expect(foundPatient.RefID).To(Equal(patient.RefID))
		})

		It("should return nil if patient not found", func() {
			foundPatient, err := patientDataStore.FindByRefID("not-exist")
			Expect(err).To(BeNil())
			Expect(foundPatient).To(BeNil())
		})
	})

	When("Creating", func() {
		It("should create patient", func() {
			patient := generatePatient()
			err := patientDataStore.Create(patient)
			Expect(err).To(BeNil())
			Expect(patient.ID).ToNot(BeZero())

			var foundPatient datastore.Patient
			err = db.First(&foundPatient, patient.ID).Error
			Expect(err).To(BeNil())
			Expect(foundPatient.ID).To(Equal(patient.ID))
			Expect(foundPatient.RefID).To(Equal(patient.RefID))
		})
	})

	When("FindOrCreate", func() {
		It("should create patient", func() {
			patient := generatePatient()
			err := patientDataStore.FindOrCreate(patient)
			Expect(err).To(BeNil())
			Expect(patient.ID).ToNot(BeZero())

			var foundPatient datastore.Patient
			err = db.First(&foundPatient, patient.ID).Error
			Expect(err).To(BeNil())
			Expect(foundPatient.ID).To(Equal(patient.ID))
			Expect(foundPatient.RefID).To(Equal(patient.RefID))
		})

		It("should found patient", func() {
			patient := getRandomPatient(patients)
			err := patientDataStore.FindOrCreate(patient)
			Expect(err).To(BeNil())
			Expect(patient.ID).ToNot(BeZero())
		})
	})
})

func generatePatient() *datastore.Patient {
	return &datastore.Patient{RefID: fmt.Sprintf("HN-%d", rand.Uint32())}
}

func generatePatients(num int) []*datastore.Patient {
	users := make([]*datastore.Patient, num)
	for i := 0; i < num; i++ {
		users[i] = generatePatient()
	}
	return users
}

func getRandomPatient(users []*datastore.Patient) *datastore.Patient {
	return users[rand.Int()%len(users)]
}

func getRandomID() uint {
	return uint(rand.Uint32())
}
