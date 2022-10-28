package datastore_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

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

func assertRecord(db *gorm.DB, t interface{}) {
	Expect(db.Where(t).First(t).Error).To(Succeed())
}

func generateCreditCard(patientID uint, isDefault bool) *datastore.CreditCard {
	return &datastore.CreditCard{
		Last4Digits: fmt.Sprintf("%d", rand.Intn(10000)),
		Brand:       "Visa",
		PatientID:   patientID,
		CardID:      uuid.New().String(),
		Name:        "test_card",
		IsDefault:   isDefault,
		Expiry:      fmt.Sprintf("12/%d", time.Now().Year()),
	}
}

func generateCreditCardPayment(status datastore.PaymentStatus, creditCardID uint) *datastore.Payment {
	var paidAt *time.Time
	if status == datastore.SuccessPaymentStatus {
		now := time.Now()
		paidAt = &now
	}
	return &datastore.Payment{
		Method:       datastore.CreditCardPaymentMethod,
		Amount:       rand.Float64(),
		PaidAt:       paidAt,
		ChargeID:     uuid.New().String(),
		InvoiceID:    int(rand.Int31()),
		Status:       status,
		CreditCardID: &creditCardID,
	}
}

func generateNotification(patientID uint) datastore.Notification {
	return datastore.Notification{
		Title:     uuid.NewString(),
		Body:      uuid.NewString(),
		IsRead:    rand.Float32() > 0.5,
		PatientID: patientID,
	}
}

func generateNotifications(patientID uint, n int) ([]datastore.Notification, int) {
	notifications := make([]datastore.Notification, n, n)
	readCount := 0
	for i := 0; i < n; i++ {
		notifications[i] = generateNotification(patientID)
		if notifications[i].IsRead {
			readCount++
		}
	}
	return notifications, readCount
}

func generatePatientWithNotifications(n int) ([]*datastore.Patient, []int) {
	patients := make([]*datastore.Patient, n, n)
	readCounts := make([]int, n, n)
	for i := 0; i < n; i++ {
		patients[i] = generatePatient()
		patients[i].Notification, readCounts[i] = generateNotifications(patients[i].ID, rand.Intn(20)+2)
	}
	return patients, readCounts
}
