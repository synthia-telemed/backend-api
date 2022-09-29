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
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/pkg/payment"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_clock"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
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
		mockPaymentDataStore    *mock_datastore.MockPaymentDataStore
		mockhospitalSysClient   *mock_hospital_client.MockSystemClient
		mockClock               *mock_clock.MockClock
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		patientID = uint(rand.Uint32())
		customerID = uuid.New().String()
		c.Set("UserID", patientID)
		c.Set("CustomerID", customerID)

		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockCreditCardDataStore = mock_datastore.NewMockCreditCardDataStore(mockCtrl)
		mockPaymentDataStore = mock_datastore.NewMockPaymentDataStore(mockCtrl)
		mockPaymentClient = mock_payment.NewMockClient(mockCtrl)
		mockhospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockClock = mock_clock.NewMockClock(mockCtrl)
		h = handler.NewPaymentHandler(mockPaymentClient, mockPatientDataStore, mockCreditCardDataStore, mockhospitalSysClient, mockPaymentDataStore, mockClock, zap.NewNop().Sugar())
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

			req = &handler.AddCreditCardRequest{CardToken: uuid.New().String(), Name: "test_card", IsDefault: true}
			reqBody, err := json.Marshal(&req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
		})

		When("card_token is not present in request body", func() {
			BeforeEach(func() {
				c.Request = httptest.NewRequest("POST", "/", nil)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("Count card by patient ID error", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().Count(patientID).Return(0, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("Patient has maximum number of card", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().Count(patientID).Return(5, nil).Times(1)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("add credit card to Omise error", func() {
			BeforeEach(func() {
				mockCreditCardDataStore.EXPECT().Count(patientID).Return(0, nil).Times(1)
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
				pCard, dCard = testhelper.GeneratePaymentAndDataStoreCard(patientID, req.Name, true)
				mockPaymentClient.EXPECT().AddCreditCard(customerID, req.CardToken).Return(pCard, nil).Times(1)
				mockCreditCardDataStore.EXPECT().Create(dCard).Return(nil).Times(1)
			})

			When("it's the first credit card and set as not default", func() {
				BeforeEach(func() {
					req.IsDefault = false
					reqBody, err := json.Marshal(&req)
					Expect(err).To(BeNil())
					c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
					mockCreditCardDataStore.EXPECT().Count(patientID).Return(0, nil).Times(1)
				})
				It("should return 201 with credit card set as default", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
				})
			})
			When("it's the first credit card and set as default", func() {
				BeforeEach(func() {
					mockCreditCardDataStore.EXPECT().Count(patientID).Return(0, nil).Times(1)
				})
				It("should return 201 with credit card set as default", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
				})
			})

			When("patient already has some cards and set new card as default", func() {
				BeforeEach(func() {
					mockCreditCardDataStore.EXPECT().Count(patientID).Return(3, nil).Times(1)
					mockCreditCardDataStore.EXPECT().SetAllToNonDefault(patientID).Return(nil).Times(1)
				})
				It("should return 201 with credit card set as default", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
				})
			})
			When("patient already has some cards and set new card as not default", func() {
				BeforeEach(func() {
					req.IsDefault = false
					reqBody, err := json.Marshal(&req)
					Expect(err).To(BeNil())
					c.Request = httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))

					mockCreditCardDataStore.EXPECT().Count(patientID).Return(3, nil).Times(1)
					dCard.IsDefault = false
				})
				It("should return 201 with credit card set as non default", func() {
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
				mockCreditCardDataStore.EXPECT().FindByPatientID(patientID).Return(nil, testhelper.MockError).Times(1)
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
				cards = testhelper.GenerateCreditCards(3)
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
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(nil, testhelper.MockError).Times(1)
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
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return("", testhelper.MockError).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("save customerID error", func() {
				BeforeEach(func() {
					mockPaymentClient.EXPECT().CreateCustomer(patientID).Return(customerID, nil).Times(1)
					pp := &datastore.Patient{PaymentCustomerID: &customerID}
					mockPatientDataStore.EXPECT().Save(pp).Return(testhelper.MockError).Times(1)
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
				mockCreditCardDataStore.EXPECT().FindByID(cardID).Return(nil, testhelper.MockError).Times(1)
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

	Context("SetCreditCardIsDefault", func() {
		var (
			card *datastore.CreditCard
			req  *handler.SetCreditCardIsDefaultRequest
		)

		BeforeEach(func() {
			handlerFunc = h.SetCreditCardIsDefault
			card = testhelper.GenerateCreditCard()
			req = &handler.SetCreditCardIsDefaultRequest{IsDefault: true}
			body, err := json.Marshal(req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("PATCH", "/", bytes.NewReader(body))
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
		When("set default to false", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				req.IsDefault = false
				body, err := json.Marshal(req)
				Expect(err).To(BeNil())
				c.Request = httptest.NewRequest("PATCH", "/", bytes.NewReader(body))
				mockCreditCardDataStore.EXPECT().SetIsDefault(card.ID, req.IsDefault).Return(nil).Times(1)
			})
			It("should set credit card default to false and return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
		When("set all card to non-default error", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				mockCreditCardDataStore.EXPECT().SetAllToNonDefault(card.PatientID).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("set is default status error", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				mockCreditCardDataStore.EXPECT().SetAllToNonDefault(card.PatientID).Return(nil).Times(1)
				mockCreditCardDataStore.EXPECT().SetIsDefault(card.ID, req.IsDefault).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("set is default to true with no error", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				mockCreditCardDataStore.EXPECT().SetAllToNonDefault(card.PatientID).Return(nil).Times(1)
				mockCreditCardDataStore.EXPECT().SetIsDefault(card.ID, req.IsDefault).Return(nil).Times(1)
			})
			It("should set all other cards to non-default, set specific card to true, and return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("DeleteCreditCard", func() {
		var (
			card *datastore.CreditCard
		)

		BeforeEach(func() {
			handlerFunc = h.DeleteCreditCard
			card = testhelper.GenerateCreditCard()
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
				mockPaymentClient.EXPECT().RemoveCreditCard(customerID, card.CardID).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("remove credit card from data store failed", func() {
			BeforeEach(func() {
				c.Set("CreditCard", card)
				mockCreditCardDataStore.EXPECT().Delete(card.ID).Return(testhelper.MockError).Times(1)
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
			It("should return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("ParseAndVerifyUnpaidInvoiceOwnership", func() {
		var invoiceID int
		BeforeEach(func() {
			handlerFunc = h.ParseAndVerifyUnpaidInvoiceOwnership
			invoiceID = int(rand.Int31())
		})
		When("invoiceID is invalid", func() {
			BeforeEach(func() {
				c.AddParam("invoiceID", "not-int-ka")
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("find invoice by ID with hospital sys client error", func() {
			BeforeEach(func() {
				c.AddParam("invoiceID", fmt.Sprintf("%d", invoiceID))
				mockhospitalSysClient.EXPECT().FindInvoiceByID(gomock.Any(), invoiceID).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("invoice is not found", func() {
			BeforeEach(func() {
				c.AddParam("invoiceID", fmt.Sprintf("%d", invoiceID))
				mockhospitalSysClient.EXPECT().FindInvoiceByID(gomock.Any(), invoiceID).Return(nil, nil).Times(1)
			})
			It("should return 404", func() {
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})
		When("invoice is paid", func() {
			BeforeEach(func() {
				c.AddParam("invoiceID", fmt.Sprintf("%d", invoiceID))
				i := &hospital.InvoiceOverview{PatientID: uuid.New().String(), Paid: true}
				mockhospitalSysClient.EXPECT().FindInvoiceByID(gomock.Any(), invoiceID).Return(i, nil).Times(1)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("find patient by ID error", func() {
			BeforeEach(func() {
				c.AddParam("invoiceID", fmt.Sprintf("%d", invoiceID))
				p := &datastore.Patient{ID: patientID, RefID: uuid.New().String()}
				i := &hospital.InvoiceOverview{PatientID: p.RefID}
				mockhospitalSysClient.EXPECT().FindInvoiceByID(gomock.Any(), invoiceID).Return(i, nil).Times(1)
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("patient's refID is not equal to patient ID in invoice", func() {
			BeforeEach(func() {
				c.AddParam("invoiceID", fmt.Sprintf("%d", invoiceID))
				p := &datastore.Patient{ID: patientID, RefID: uuid.New().String()}
				i := &hospital.InvoiceOverview{PatientID: uuid.New().String()}
				mockhospitalSysClient.EXPECT().FindInvoiceByID(gomock.Any(), invoiceID).Return(i, nil).Times(1)
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(p, nil).Times(1)
			})
			It("should return 403", func() {
				Expect(rec.Code).To(Equal(http.StatusForbidden))
			})
		})
		When("no error occurred", func() {
			BeforeEach(func() {
				c.AddParam("invoiceID", fmt.Sprintf("%d", invoiceID))
				p := &datastore.Patient{ID: patientID, RefID: uuid.New().String()}
				i := &hospital.InvoiceOverview{Id: invoiceID, PatientID: p.RefID}
				mockhospitalSysClient.EXPECT().FindInvoiceByID(gomock.Any(), invoiceID).Return(i, nil).Times(1)
				mockPatientDataStore.EXPECT().FindByID(patientID).Return(p, nil).Times(1)
			})
			It("set invoice to the context", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				i, ok := c.Get("Invoice")
				Expect(ok).To(BeTrue())
				invoice, ok := i.(*hospital.InvoiceOverview)
				Expect(ok).To(BeTrue())
				Expect(invoice.Id).To(Equal(invoiceID))
			})
		})
	})

	Context("PayInvoiceWithCreditCard", func() {
		var (
			creditCard   *datastore.CreditCard
			invoice      *hospital.InvoiceOverview
			invoiceIDStr string
		)
		BeforeEach(func() {
			handlerFunc = h.PayInvoiceWithCreditCard
			creditCard = testhelper.GenerateCreditCard()
			invoice = testhelper.GenerateHospitalInvoice(false)
			invoiceIDStr = fmt.Sprintf("%d", invoice.Id)
			c.Set("CreditCard", creditCard)
			c.Set("Invoice", invoice)
		})
		When("pay with credit card error", func() {
			BeforeEach(func() {
				mockPaymentClient.EXPECT().PayWithCreditCard(customerID, creditCard.CardID, invoiceIDStr, int(invoice.Total*100)).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("hospital sys client PaidInvoice error", func() {
			BeforeEach(func() {
				p := testhelper.GeneratePayment(true)
				mockPaymentClient.EXPECT().PayWithCreditCard(customerID, creditCard.CardID, invoiceIDStr, int(invoice.Total*100)).Return(p, nil).Times(1)
				now := time.Now()
				mockClock.EXPECT().NowPointer().Return(&now).Times(1)
				mockhospitalSysClient.EXPECT().PaidInvoice(gomock.Any(), invoice.Id).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("create payment in datastore error", func() {
			BeforeEach(func() {
				p := testhelper.GeneratePayment(true)
				mockPaymentClient.EXPECT().PayWithCreditCard(customerID, creditCard.CardID, invoiceIDStr, int(invoice.Total*100)).Return(p, nil).Times(1)
				now := time.Now()
				mockClock.EXPECT().NowPointer().Return(&now).Times(1)
				mockhospitalSysClient.EXPECT().PaidInvoice(gomock.Any(), invoice.Id).Return(nil).Times(1)
				mockPaymentDataStore.EXPECT().Create(gomock.Any()).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		Context("no error occurred", func() {
			var (
				paymentCharge *payment.Payment
				paymentData   *datastore.Payment
			)
			When("payment failed", func() {
				BeforeEach(func() {
					paymentCharge = testhelper.GeneratePayment(false)
					mockPaymentClient.EXPECT().PayWithCreditCard(customerID, creditCard.CardID, invoiceIDStr, int(invoice.Total*100)).Return(paymentCharge, nil).Times(1)
					paymentData = testhelper.GenerateDataStorePayment(datastore.CreditCardPaymentMethod, datastore.FailedPaymentStatus, invoice, paymentCharge, creditCard)
					mockClock.EXPECT().NowPointer().Return(paymentData.PaidAt).Times(1)
					mockPaymentDataStore.EXPECT().Create(paymentData).Return(nil).Times(1)
				})
				It("should return 201 with failure message", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
					var res handler.PayInvoiceWithCreditCardResponse
					Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
					Expect(res.FailureMessage).To(Equal(paymentCharge.FailureMessage))
					Expect(res.Status).To(Equal(datastore.FailedPaymentStatus))
				})
			})
			When("payment success", func() {
				BeforeEach(func() {
					paymentCharge = testhelper.GeneratePayment(true)
					mockPaymentClient.EXPECT().PayWithCreditCard(customerID, creditCard.CardID, invoiceIDStr, int(invoice.Total*100)).Return(paymentCharge, nil).Times(1)
					mockhospitalSysClient.EXPECT().PaidInvoice(gomock.Any(), invoice.Id).Return(nil).Times(1)
					paymentData = testhelper.GenerateDataStorePayment(datastore.CreditCardPaymentMethod, datastore.SuccessPaymentStatus, invoice, paymentCharge, creditCard)
					mockClock.EXPECT().NowPointer().Return(paymentData.PaidAt).Times(1)
					mockPaymentDataStore.EXPECT().Create(paymentData).Return(nil).Times(1)
				})
				It("should return 201 with success message", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
					var res handler.PayInvoiceWithCreditCardResponse
					Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
					Expect(res.FailureMessage).To(BeNil())
					Expect(res.Status).To(Equal(datastore.SuccessPaymentStatus))
				})
			})
		})
	})
})
