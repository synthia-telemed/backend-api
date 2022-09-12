package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/pkg/server"
	"math/rand"
	"net/http/httptest"
	"time"
)

var MockError = errors.New("error")

func InitHandlerTest() (*gomock.Controller, *httptest.ResponseRecorder, *gin.Context) {
	rand.Seed(GinkgoRandomSeed())
	mockCtrl := gomock.NewController(GinkgoT())
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	gin.SetMode(gin.TestMode)
	return mockCtrl, rec, c
}

func GeneratePatient() *datastore.Patient {
	return &datastore.Patient{
		ID:        uint(rand.Uint32()),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RefID:     uuid.New().String(),
	}
}

func GenerateCreditCard() *datastore.CreditCard {
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

func GenerateCreditCards(n int) []datastore.CreditCard {
	cards := make([]datastore.CreditCard, n)
	for i := 0; i < n; i++ {
		cards[i] = *GenerateCreditCard()
	}
	return cards
}

func GeneratePaymentAndDataStoreCard(patientID uint, name string) (*payment.Card, *datastore.CreditCard) {
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

func GeneratePayment(isSuccess bool) *payment.Payment {
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

func GenerateHospitalInvoice(paid bool) *hospital.InvoiceOverview {
	return &hospital.InvoiceOverview{
		Id:            rand.Int(),
		CreatedAt:     time.Now(),
		Paid:          paid,
		Total:         rand.Float64(),
		AppointmentID: uuid.New().String(),
		PatientID:     uuid.New().String(),
	}
}

func GenerateDataStorePayment(method datastore.PaymentMethod, status datastore.PaymentStatus, i *hospital.InvoiceOverview, p *payment.Payment, c *datastore.CreditCard) *datastore.Payment {
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

func GenerateAppointmentOverviews(status hospital.AppointmentStatus, n int) []*hospital.AppointmentOverview {
	apps := make([]*hospital.AppointmentOverview, n, n)
	for i := 0; i < n; i++ {
		apps[i] = GenerateAppointmentOverview(status)
	}
	return apps
}

func GenerateAppointmentOverview(status hospital.AppointmentStatus) *hospital.AppointmentOverview {
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

func GenerateAppointment(patientID string, doctorID string, status hospital.AppointmentStatus) (*hospital.Appointment, int) {
	var invoice *hospital.Invoice
	if status == hospital.AppointmentStatusCompleted {
		invoice = &hospital.Invoice{
			Id:           int(rand.Int31()),
			Total:        rand.Float64() * 10000,
			Paid:         false,
			InvoiceItems: nil,
		}
	}
	id := rand.Int31()
	return &hospital.Appointment{
		Id:              fmt.Sprintf("%d", id),
		PatientID:       patientID,
		DateTime:        time.Now(),
		NextAppointment: nil,
		Detail:          uuid.New().String(),
		Status:          status,
		Doctor: hospital.DoctorOverview{
			ID:            doctorID,
			FullName:      uuid.New().String(),
			Position:      uuid.New().String(),
			ProfilePicURL: uuid.New().String(),
		},
		Invoice:       invoice,
		Prescriptions: nil,
	}, int(id)
}

type Ordering string

const (
	DESC Ordering = "DESC"
	ASC  Ordering = "ASC"
)

func AssertListOfAppointments(apps []*hospital.AppointmentOverview, status hospital.AppointmentStatus, order Ordering) {
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

func GenerateDoctor() *datastore.Doctor {
	return &datastore.Doctor{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RefID:     uuid.NewString(),
		ID:        uint(rand.Uint32()),
	}
}

func AssertErrorResponseBody(body *bytes.Buffer, expectedError *server.ErrorResponse) {
	var res server.ErrorResponse
	Expect(json.Unmarshal(body.Bytes(), &res)).To(Succeed())
	Expect(&res).To(Equal(expectedError))
}