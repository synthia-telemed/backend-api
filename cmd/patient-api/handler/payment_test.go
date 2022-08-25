package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
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
	"strings"
)

var _ = Describe("Payment Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.PaymentHandler
		handlerFunc gin.HandlerFunc
		patientID   uint

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
		c.Set("UserID", patientID)

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
			patient   *datastore.Patient
			cardToken string
		)

		BeforeEach(func() {
			handlerFunc = h.AddCreditCard
			cardToken = "token"
			patient = &datastore.Patient{PaymentCustomerID: nil}
			reqBody := strings.NewReader(fmt.Sprintf(`{"card_token": "%s"}`, cardToken))
			c.Request = httptest.NewRequest("POST", "/", reqBody)
		})

		When("card_token is not present in request body", func() {
			BeforeEach(func() {
				c.Request = httptest.NewRequest("POST", "/", nil)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("find patient by id error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(nil, errors.New("error")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("haven't added any credit card", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(patient, nil).Times(1)
			})

			When("payment client create customer error", func() {
				BeforeEach(func() {
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return("", errors.New("payment client error")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("patient update error", func() {
				BeforeEach(func() {
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return("cus_id", nil).Times(1)
					mockPatientDataStore.EXPECT().Save(gomock.Any()).Return(errors.New("saving error")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("successfully added credit card", func() {
				BeforeEach(func() {
					cusID := "cus_id"
					p := &datastore.Patient{PaymentCustomerID: &cusID}
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return(cusID, nil).Times(1)
					mockPatientDataStore.EXPECT().Save(p).Return(nil).Times(1)
					mockPaymentClient.EXPECT().AddCreditCard(cusID, cardToken).Return(nil).Times(1)
				})
				It("should return 201", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
				})
			})
		})

		When("already added credit card", func() {
			BeforeEach(func() {
				cusID := "cus_id"
				p := &datastore.Patient{PaymentCustomerID: &cusID}
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(p, nil).Times(1)
				mockPaymentClient.EXPECT().AddCreditCard(cusID, cardToken).Return(nil).Times(1)
			})
			It("should return 201", func() {
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})
		})

		When("add credit card error", func() {
			BeforeEach(func() {
				cusID := "cus_id"
				p := &datastore.Patient{PaymentCustomerID: &cusID}
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(p, nil).Times(1)
				mockPaymentClient.EXPECT().AddCreditCard(cusID, cardToken).Return(errors.New("error")).Times(1)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Context("Get patient's credit cards", func() {
		BeforeEach(func() {
			handlerFunc = h.GetCreditCards
			c.Request = httptest.NewRequest("GET", "/", nil)
		})

		When("find patient by id error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(nil, errors.New("error")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("patient has no credit cards", func() {
			BeforeEach(func() {
				patient := &datastore.Patient{PaymentCustomerID: nil}
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(patient, nil).Times(1)
			})
			It("should return 200 with empty list", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.String()).To(Equal(`[]`))
			})
		})

		When("patient has at least one credit card", func() {
			var cards []payment.Card
			BeforeEach(func() {
				cards = generatePaymentCard(3)
				cusID := "cus_id"
				patient := &datastore.Patient{PaymentCustomerID: &cusID}
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(patient, nil).Times(1)
				mockPaymentClient.EXPECT().ListCards(cusID).Return(cards, nil)
			})
			It("should return 200 with list of cards", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var c []payment.Card
				Expect(json.Unmarshal(rec.Body.Bytes(), &c)).To(Succeed())
				Expect(c).To(Equal(cards))
			})
		})

		When("payment client error when getting list of cards", func() {
			BeforeEach(func() {
				cusID := "cus_id"
				patient := &datastore.Patient{PaymentCustomerID: &cusID}
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(patient, nil).Times(1)
				mockPaymentClient.EXPECT().ListCards(cusID).Return(nil, errors.New("error"))
			})
			It("should return 200 with list of cards", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

func generatePaymentCard(n int) []payment.Card {
	cards := make([]payment.Card, n)
	for i := 0; i < n; i++ {
		cards[i] = payment.Card{
			ID:         fmt.Sprintf("id-%d", rand.Int()),
			LastDigits: fmt.Sprintf("%d", rand.Intn(10000)),
			Brand:      "Visa",
		}
	}
	return cards
}
