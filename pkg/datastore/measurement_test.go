package datastore_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
	"time"
)

var _ = Describe("Measurement Datastore", Ordered, func() {
	var (
		db                   *gorm.DB
		measurementDataStore datastore.MeasurementDataStore
		patient              *datastore.Patient
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
		measurementDataStore, err = datastore.NewGormMeasurementDataStore(db)
		Expect(err).To(BeNil())
		Expect(db.AutoMigrate(&datastore.Patient{})).To(Succeed())
		patient = &datastore.Patient{RefID: fmt.Sprintf("ref-id-%d", rand.Int())}
		Expect(db.Create(patient).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.BloodPressure{}, &datastore.Glucose{}, &datastore.Patient{})).To(Succeed())
	})

	Context("CreateBloodPressure", func() {
		var bp *datastore.BloodPressure
		BeforeEach(func() {
			bp = &datastore.BloodPressure{
				PatientID: patient.ID,
				DateTime:  time.Now().UTC(),
				Systolic:  uint(rand.Uint32()),
				Diastolic: uint(rand.Uint32()),
				Pulse:     uint(rand.Uint32()),
			}
		})
		It("should save blood pressure", func() {
			Expect(measurementDataStore.CreateBloodPressure(bp)).To(Succeed())
			assertRecord(db, bp)
		})
	})

	Context("CreateGlucose", func() {
		var g *datastore.Glucose
		BeforeEach(func() {
			g = &datastore.Glucose{
				PatientID:    patient.ID,
				DateTime:     time.Now().UTC(),
				IsBeforeMeal: true,
				Value:        uint(rand.Uint32()),
			}
		})
		It("should save blood pressure", func() {
			Expect(measurementDataStore.CreateGlucose(g)).To(Succeed())
			assertRecord(db, g)
		})
	})
})
