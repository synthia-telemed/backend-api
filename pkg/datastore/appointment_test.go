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
	"time"
)

var _ = Describe("Appointment Datastore", Ordered, func() {
	var (
		db                   *gorm.DB
		appointmentDataStore datastore.AppointmentDataStore
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
		appointmentDataStore, err = datastore.NewGormAppointmentDataStore(db)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Appointment{})).To(Succeed())
	})

	Context("Create", func() {
		It("should save new appointment to db", func() {
			appointment := datastore.Appointment{
				RefID:       uuid.NewString(),
				Duration:    (time.Minute * 10).Seconds(),
				StartedTime: time.Now(),
			}
			Expect(appointmentDataStore.Create(&appointment)).To(Succeed())
			var retriedAppointment datastore.Appointment
			Expect(db.First(&retriedAppointment, appointment.ID).Error).To(Succeed())
			Expect(retriedAppointment).ToNot(BeNil())
			Expect(retriedAppointment.RefID).To(Equal(appointment.RefID))
		})
	})

	Context("FindByRefID", func() {
		When("appointment is not found", func() {
			It("should return nil with no error", func() {
				app, err := appointmentDataStore.FindByRefID(uuid.NewString())
				Expect(err).To(BeNil())
				Expect(app).To(BeNil())
			})
		})
		When("appointment is found", func() {
			var appointment datastore.Appointment
			BeforeEach(func() {
				appointment = datastore.Appointment{
					RefID:       uuid.NewString(),
					Duration:    (time.Minute * 10).Seconds(),
					StartedTime: time.Now(),
				}
				Expect(db.Create(&appointment).Error).To(Succeed())
			})
			It("should return appointment with no error", func() {
				app, err := appointmentDataStore.FindByRefID(appointment.RefID)
				Expect(err).To(BeNil())
				Expect(app).ToNot(BeNil())
				Expect(app.ID).To(Equal(app.ID))
				Expect(app.RefID).To(Equal(app.RefID))
			})
		})
	})
})
