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
	"time"
)

var _ = Describe("Measurement Datastore", Ordered, func() {
	var (
		db                   *gorm.DB
		measurementDataStore datastore.MeasurementDataStore
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
		measurementDataStore, err = datastore.NewGormMeasurementDataStore(db)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.BloodPressure{}, &datastore.Glucose{})).To(Succeed())
	})

	Context("CreateBloodPressure", func() {
		var bp *datastore.BloodPressure
		BeforeEach(func() {
			bp = &datastore.BloodPressure{
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

func assertRecord(db *gorm.DB, t interface{}) {
	Expect(db.Where(t).First(t).Error).To(Succeed())
}