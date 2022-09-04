package handler_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"math/rand"
	"net/http/httptest"
	"time"
)

func initHandlerTest() (*gomock.Controller, *httptest.ResponseRecorder, *gin.Context) {
	rand.Seed(GinkgoRandomSeed())
	mockCtrl := gomock.NewController(GinkgoT())
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	return mockCtrl, rec, c
}

func generatePatient() *datastore.Patient {
	return &datastore.Patient{
		ID:        uint(rand.Uint32()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RefID:     uuid.New().String(),
	}
}

func generateCreditCard() *datastore.CreditCard {
	return &datastore.CreditCard{
		ID:          uint(rand.Uint32()),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Last4Digits: fmt.Sprintf("%d", rand.Intn(10000)),
		Brand:       "Visa",
		PatientID:   uint(rand.Uint32()),
		CardID:      uuid.New().String(),
	}
}

func generateCreditCards(n int) []datastore.CreditCard {
	cards := make([]datastore.CreditCard, n)
	for i := 0; i < n; i++ {
		cards[i] = *generateCreditCard()
	}
	return cards
}

func generatePaymentAndDataStoreCard(patientID uint, name string) (*payment.Card, *datastore.CreditCard) {
	pCard := &payment.Card{
		ID:          uuid.New().String(),
		Last4Digits: fmt.Sprintf("%d", rand.Intn(10000)),
		Brand:       "MasterCard",
	}
	dCard := &datastore.CreditCard{
		Last4Digits: pCard.Last4Digits,
		Brand:       pCard.Brand,
		PatientID:   patientID,
		CardID:      pCard.ID,
		Name:        name,
	}
	return pCard, dCard
}

func generatePayment(isSuccess bool) *payment.Payment {
	var (
		failure *string
	)
	if !isSuccess {
		f := "failed to charge"
		failure = &f
	}
	return &payment.Payment{
		ID:             uuid.New().String(),
		Amount:         rand.Int(),
		Currency:       "THB",
		Paid:           isSuccess,
		Success:        isSuccess,
		FailureCode:    failure,
		FailureMessage: failure,
	}
}

func generateHospitalInvoice(paid bool) *hospital.InvoiceOverview {
	return &hospital.InvoiceOverview{
		Id:            rand.Int(),
		CreatedAt:     time.Now(),
		Paid:          paid,
		Total:         rand.Float64(),
		AppointmentID: uuid.New().String(),
		PatientID:     uuid.New().String(),
	}
}

func generateDataStorePayment(method datastore.PaymentMethod, status datastore.PaymentStatus, i *hospital.InvoiceOverview, p *payment.Payment, c *datastore.CreditCard) *datastore.Payment {
	var paidAt *time.Time
	if status != datastore.PendingPaymentStatus || method == datastore.CreditCardPaymentMethod {
		now := time.Now()
		paidAt = &now
	}
	return &datastore.Payment{
		Method:       method,
		Amount:       i.Total,
		PaidAt:       paidAt,
		ChargeID:     p.ID,
		InvoiceID:    i.Id,
		Status:       status,
		CreditCard:   c,
		CreditCardID: &c.ID,
	}
}

func generateAppointmentOverviews(status hospital.AppointmentStatus, n int) []*hospital.AppointmentOverview {
	apps := make([]*hospital.AppointmentOverview, n, n)
	for i := 0; i < n; i++ {
		apps[i] = generateAppointmentOverview(status)
	}
	return apps
}

func generateAppointmentOverview(status hospital.AppointmentStatus) *hospital.AppointmentOverview {
	return &hospital.AppointmentOverview{
		Id:        uuid.New().String(),
		DateTime:  time.Now(),
		PatientId: uuid.New().String(),
		Status:    status,
		Doctor: hospital.DoctorOverview{
			FullName:      uuid.New().String(),
			Position:      uuid.New().String(),
			ProfilePicURL: uuid.New().String(),
		},
	}
}

type Ordering string

const (
	DESC Ordering = "DESC"
	ASC  Ordering = "ASC"
)

func assertListOfAppointments(apps []*hospital.AppointmentOverview, status hospital.AppointmentStatus, order Ordering) {
	prevTime := apps[0].DateTime
	for i := 1; i < len(apps); i++ {
		a := apps[i]
		Expect(a.Status).To(Equal(status))
		if order == DESC {
			Expect(a.DateTime.After(prevTime)).To(BeTrue())
		} else {
			Expect(a.DateTime.Before(prevTime)).To(BeTrue())
		}
	}
}
