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
		mockCtrl, rec, c = initHandlerTest()
		patientID = uint(rand.Uint32())
		customerID = uuid.New().String()
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
			req *handler.AddCreditCardRequest
		)

		BeforeEach(func() {
			handlerFunc = h.AddCreditCard

			req = &handler.AddCreditCardRequest{CardToken: uuid.New().String(), Name: "test_card"}
			reqBody, err := json.Marshal(&req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
			c.Set("CustomerID", customerID)
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
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(generateCreditCards(5), nil).Times(1)
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
			var (
				pCard *payment.Card
				dCard *datastore.CreditCard
			)

			BeforeEach(func() {
				pCard, dCard = generatePaymentAndDataStoreCard(patientID, req.Name)
				mockPaymentClient.EXPECT().AddCreditCard(customerID, req.CardToken).Return(pCard, nil).Times(1)
				mockCreditCardDataStore.EXPECT().Create(dCard).Return(nil).Times(1)
			})

			When("it's the first credit card", func() {
				BeforeEach(func() {
					mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return([]datastore.CreditCard{}, nil).Times(1)
				})
				It("should return 201", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
				})
			})
			When("patient already has some cards", func() {
				BeforeEach(func() {
					mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(generateCreditCards(3), nil).Times(1)
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
			c.Set("CustomerID", customerID)
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
				cards = generateCreditCards(3)
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

	Context("Create or parse customerID", func() {
		BeforeEach(func() {
			handlerFunc = h.CreateOrParseCustomer
		})

		When("Find patient by ID error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("Patient doesn't have customerID", func() {
			BeforeEach(func() {
				p := &datastore.Patient{PaymentCustomerID: nil}
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(p, nil).Times(1)
			})

			When("create payment customer error", func() {
				BeforeEach(func() {
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return("", errors.New("err")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("save customerID error", func() {
				BeforeEach(func() {
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return(customerID, nil).Times(1)
					pp := &datastore.Patient{PaymentCustomerID: &customerID}
					mockPatientDataStore.EXPECT().Save(pp).Return(errors.New("err")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("no error occurred", func() {
				BeforeEach(func() {
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return(customerID, nil).Times(1)
					pp := &datastore.Patient{PaymentCustomerID: &customerID}
					mockPatientDataStore.EXPECT().Save(pp).Return(nil).Times(1)
				})
				It("should set ID to CustomerID", func() {
					id, ok := c.Get("CustomerID")
					Expect(ok).To(BeTrue())
					Expect(id).To(Equal(customerID))
				})
			})
		})

		When("patient already has customer ID", func() {
			BeforeEach(func() {
				p := &datastore.Patient{PaymentCustomerID: &customerID}
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(p, nil).Times(1)
			})
			It("should set ID to CustomerID", func() {
				id, ok := c.Get("CustomerID")
				Expect(ok).To(BeTrue())
				Expect(id).To(Equal(customerID))
			})
		})
	})

	Context("VerifyCreditCardOwnership", func() {
		var cardID uint
		BeforeEach(func() {
			handlerFunc = h.VerifyCreditCardOwnership
			cardID = uint(rand.Uint32())
		})

		When("cardID is in invalid format", func() {
			BeforeEach(func() {
				c.AddParam("cardID", "not-uint")
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("find credit card by ID error", func() {
			BeforeEach(func() {
				c.AddParam("cardID", fmt.Sprintf("%v", cardID))
				mockCreditCardDataStore.EXPECT().FindByID(cardID).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("credit card is not found", func() {
			BeforeEach(func() {
				c.AddParam("cardID", fmt.Sprintf("%v", cardID))
				mockCreditCardDataStore.EXPECT().FindByID(cardID).Return(nil, nil).Times(1)
			})
			It("should return 404", func() {
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})
		When("patient doesn't own the credit card", func() {
			BeforeEach(func() {
				c.AddParam("cardID", fmt.Sprintf("%v", cardID))
				card := &datastore.CreditCard{PatientID: uint(rand.Uint32())}
				mockCreditCardDataStore.EXPECT().FindByID(cardID).Return(card, nil).Times(1)
			})
			It("should return 403", func() {
				Expect(rec.Code).To(Equal(http.StatusForbidden))
			})
		})
	})

	Context("DeleteCreditCard", func() {
		var (
			card *datastore.CreditCard
		)

		BeforeEach(func() {
			handlerFunc = h.DeleteCreditCard
			c.Set("CustomerID", customerID)
			card = generateCreditCard()
		})

		When("credit card is not set", func() {
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("credit card parsing is failed", func() {
			BeforeEach(func() {
				c.Set("CreditCard", "just-string")
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("remove credit card from payment client failed", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				mockCreditCardDataStore.EXPECT().Delete(card.ID).Return(nil).Times(1)
				mockPaymentClient.EXPECT().RemoveCreditCard(customerID, card.CardID).Return(errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("remove credit card from data store failed", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				mockCreditCardDataStore.EXPECT().Delete(card.ID).Return(errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("no error occurred", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				mockPaymentClient.EXPECT().RemoveCreditCard(customerID, card.CardID).Return(nil).Times(1)
				mockCreditCardDataStore.EXPECT().Delete(card.ID).Return(nil).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
