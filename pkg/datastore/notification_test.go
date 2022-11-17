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

var _ = Describe("Patient Datastore", Ordered, func() {

	var (
		db                    *gorm.DB
		notificationDataStore datastore.NotificationDataStore
		patients              []*datastore.Patient
		patientsReadCount     []int
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
		Expect(db.AutoMigrate(&datastore.Patient{})).To(Succeed())
		var err error
		notificationDataStore, err = datastore.NewGormNotificationDataStore(db)
		Expect(err).To(BeNil())

		patients, patientsReadCount = generatePatientWithNotifications(3)
		Expect(db.Create(&patients).Error).To(Succeed())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{}, &datastore.Notification{})).To(Succeed())
	})

	Context("Create notification", func() {
		It("should crate notification", func() {
			noti := generateNotification(patients[1].ID)
			Expect(notificationDataStore.Create(&noti)).To(Succeed())
			var foundNoti datastore.Notification
			Expect(db.Where(&noti).First(&foundNoti).Error).To(BeNil())
		})
	})

	Context("Count Unread Notification", func() {
		It("should return number of unread notification", func() {
			p := patients[0]
			expected := len(p.Notification) - patientsReadCount[0]
			count, err := notificationDataStore.CountUnRead(p.ID)
			Expect(err).To(BeNil())
			Expect(count).To(Equal(expected))
		})
	})

	Context("List latest", func() {
		It("List notification from the latest to oldest", func() {
			p := patients[0]
			notifications, err := notificationDataStore.ListLatest(p.ID)
			Expect(err).To(BeNil())
			Expect(notifications).To(HaveLen(len(p.Notification)))
			// Cannot test the order because of timestamp is all the same
		})
	})

	Context("FindByID", func() {
		When("notification is not found", func() {
			It("should return nil with error of nil", func() {
				notification, err := notificationDataStore.FindByID(473842738)
				Expect(err).To(BeNil())
				Expect(notification).To(BeNil())
			})
		})

		When("notification is found", func() {
			It("should return the notification", func() {
				expected := patients[1].Notification[0]
				notification, err := notificationDataStore.FindByID(expected.ID)
				Expect(err).To(BeNil())
				Expect(notification.ID).To(Equal(expected.ID))
			})
		})
	})

	Context("SetAsRead", func() {
		When("notification is not exist", func() {
			It("should return no error and do no action", func() {
				Expect(notificationDataStore.SetAsRead(2837428)).To(Succeed())
			})
		})

		When("notification is exist", func() {
			var targetNotification *datastore.Notification
			BeforeEach(func() {
				for _, noti := range patients[0].Notification {
					if !noti.IsRead {
						targetNotification = &noti
					}
				}
			})

			It("should return error of nil", func() {
				Expect(notificationDataStore.SetAsRead(targetNotification.ID)).To(Succeed())
				var retrievedNoti datastore.Notification
				db.First(&retrievedNoti, targetNotification.ID)
				Expect(retrievedNoti.ID).To(Equal(targetNotification.ID))
				Expect(retrievedNoti.IsRead).To(BeTrue())
			})
		})
	})

	Context("SetAllAsRead", func() {
		It("should set all notifications to read", func() {
			patient := patients[0]
			Expect(notificationDataStore.SetAllAsRead(patient.ID)).To(Succeed())
			var readCount int64
			db.Model(datastore.Notification{}).Where("patient_id = ? AND is_read = ?", patient.ID, false).Count(&readCount)
			Expect(readCount).To(BeZero())
		})
	})
})
