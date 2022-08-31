package handler_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
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
		CreatedAt:      time.Now(),
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
		paidAt = &p.CreatedAt
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
