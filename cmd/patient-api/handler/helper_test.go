package handler_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
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
		IsDefault:   false,
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
		IsDefault:   false,
		Last4Digits: pCard.Last4Digits,
		Brand:       pCard.Brand,
		PatientID:   patientID,
		CardID:      pCard.ID,
		Name:        name,
	}
	return pCard, dCard
}
