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

var _ = Describe("Patient Datastore", Ordered, func() {

	var (
		db               *gorm.DB
		patientDataStore datastore.PatientDataStore
		patients         []*datastore.Patient
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
		patientDataStore, err = datastore.NewGormPatientDataStore(db)
		Expect(err).To(BeNil())

		patients = generatePatients(10)
		Expect(db.Create(&patients).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{})).To(Succeed())
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
			Expect(patientDataStore.Create(patient)).To(Succeed())
			Expect(patient.ID).ToNot(BeZero())

			var foundPatient datastore.Patient
			Expect(db.First(&foundPatient, patient.ID).Error).To(Succeed())
			Expect(foundPatient.ID).To(Equal(patient.ID))
			Expect(foundPatient.RefID).To(Equal(patient.RefID))
		})
	})

	When("FindOrCreate", func() {
		It("should create patient", func() {
			patient := generatePatient()
			Expect(patientDataStore.FindOrCreate(patient)).To(Succeed())
			Expect(patient.ID).ToNot(BeZero())

			var foundPatient datastore.Patient
			Expect(db.First(&foundPatient, patient.ID).Error).To(Succeed())
			Expect(foundPatient.ID).To(Equal(patient.ID))
			Expect(foundPatient.RefID).To(Equal(patient.RefID))
		})

		It("should found patient", func() {
			patient := getRandomPatient(patients)
			Expect(patientDataStore.FindOrCreate(patient)).To(Succeed())
			Expect(patient.ID).ToNot(BeZero())
		})
	})

	When("Updating", func() {
		It("should update patient", func() {
			patient := getRandomPatient(patients)
			patient.RefID = "updated-ref-id"
			Expect(patientDataStore.Save(patient)).To(Succeed())
			Expect(patient.ID).ToNot(BeZero())

			var foundPatient datastore.Patient
			Expect(db.First(&foundPatient, patient.ID).Error).To(Succeed())
			Expect(foundPatient.ID).To(Equal(patient.ID))
			Expect(foundPatient.RefID).To(Equal(patient.RefID))
		})
	})
})

func generatePatient() *datastore.Patient {
	return &datastore.Patient{RefID: uuid.New().String()}
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
