package handler_test

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
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_payment"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Payment Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.PaymentHandler
		handlerFunc gin.HandlerFunc
		patientID   uint
		customerID  string

		mockPatientDataStore    *mock_datastore.MockPatientDataStore
		mockCreditCardDataStore *mock_datastore.MockCreditCardDataStore
		mockPaymentClient       *mock_payment.MockClient
	)

	BeforeEach(func() {
		rand.Seed(GinkgoRandomSeed())
		mockCtrl = gomock.NewController(GinkgoT())
		rec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(rec)
		patientID = uint(rand.Uint32())
		customerID = uuid.New().String()
		c.Set("UserID", patientID)
		c.Set("CustomerID", customerID)

		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockCreditCardDataStore = mock_datastore.NewMockCreditCardDataStore(mockCtrl)
		mockPaymentClient = mock_payment.NewMockClient(mockCtrl)
		h = handler.NewPaymentHandler(mockPaymentClient, mockPatientDataStore, mockCreditCardDataStore, zap.NewNop().Sugar())
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Add credit card", func() {
		var (
			req   *handler.AddCreditCardRequest
			pCard *payment.Card
			dCard *datastore.CreditCard
		)

		BeforeEach(func() {
			handlerFunc = h.AddCreditCard

			req = &handler.AddCreditCardRequest{CardToken: uuid.New().String(), Name: "test_card"}
			reqBody, err := json.Marshal(&req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
			pCard, dCard = generatePaymentAndDataStoreCard(patientID, req.Name)
		})

		When("card_token is not present in request body", func() {
			BeforeEach(func() {
				c.Request = httptest.NewRequest("POST", "/", nil)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("List card by patient ID error", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("Patient has maximum number of card", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(generatePaymentCards(5), nil).Times(1)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("add credit card to Omise error", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return([]datastore.CreditCard{}, nil).Times(1)
				mockPaymentClient.EXPECT().AddCreditCard(customerID, req.CardToken).Return(nil, errors.New("error")).Times(1)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("successfully added credit card", func() {
			BeforeEach(func() {
				mockPaymentClient.EXPECT().AddCreditCard(customerID, req.CardToken).Return(pCard, nil).Times(1)
				mockCreditCardDataStore.EXPECT().Create(dCard).Return(nil).Times(1)
			})

			When("it's the first credit card", func() {
				BeforeEach(func() {
					dCard.IsDefault = true
					mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return([]datastore.CreditCard{}, nil).Times(1)
				})
				It("should return 201", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
				})
			})
			When("patient already has some cards", func() {
				BeforeEach(func() {
					dCard.IsDefault = false
					mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(generatePaymentCards(3), nil).Times(1)
				})
				It("should return 201", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
				})
			})
		})
	})

	Context("Get patient's credit cards", func() {
		BeforeEach(func() {
			handlerFunc = h.GetCreditCards
			c.Request = httptest.NewRequest("GET", "/", nil)
		})

		When("query error", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("patient has no credit cards", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return([]datastore.CreditCard{}, nil).Times(1)
			})
			It("should return 200 with empty list", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.String()).To(Equal(`[]`))
			})
		})

		When("patient has at least one credit card", func() {
			var cards []datastore.CreditCard
			BeforeEach(func() {
				cards = generatePaymentCards(3)
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(cards, nil).Times(1)
			})
			It("should return 200 with list of cards", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var c []datastore.CreditCard
				Expect(json.Unmarshal(rec.Body.Bytes(), &c)).To(Succeed())
				Expect(c).To(HaveLen(len(cards)))
			})
		})
	})
})

func generatePaymentCards(n int) []datastore.CreditCard {
	cards := make([]datastore.CreditCard, n)
	for i := 0; i < n; i++ {
		cards[i] = datastore.CreditCard{
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
