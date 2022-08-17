package handler_test

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
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

		mockPatientDataStore *mock_datastore.MockPatientDataStore
		mockPaymentClient    *mock_payment.MockClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		rec = httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ = gin.CreateTestContext(rec)
		patientID = uint(rand.Uint32())
		c.Set("patientID", patientID)

		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockPaymentClient = mock_payment.NewMockClient(mockCtrl)
		h = handler.NewPaymentHandler(mockPaymentClient, mockPatientDataStore, zap.NewNop().Sugar())
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

			When("no error", func() {
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
	})
})
