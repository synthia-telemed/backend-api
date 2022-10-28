package handler_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Patient Gin Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.PatientGinHandler
		handlerFunc gin.HandlerFunc

		mockPatientDataStore *mock_datastore.MockPatientDataStore
		patient              *datastore.Patient
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		patientGinHandler := handler.NewPatientGinHandler(mockPatientDataStore, zap.NewNop().Sugar())
		h = &patientGinHandler
		patient = testhelper.GeneratePatient()
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("ParsePatient", func() {
		BeforeEach(func() {
			handlerFunc = h.ParsePatient
			c.Set("Patient", nil)
		})

		When("find patient by ID error", func() {
			BeforeEach(func() {
				id := uint(rand.Uint32())
				c.Set("UserID", id)
				mockPatientDataStore.EXPECT().FindByID(id).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("patient is not found", func() {
			BeforeEach(func() {
				id := uint(rand.Uint32())
				c.Set("UserID", id)
				mockPatientDataStore.EXPECT().FindByID(id).Return(nil, nil).Times(1)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("patient patient is found", func() {
			BeforeEach(func() {
				id := uint(rand.Uint32())
				c.Set("UserID", id)
				mockPatientDataStore.EXPECT().FindByID(id).Return(patient, nil).Times(1)
			})
			It("should set the patient", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				p, exist := c.Get("Patient")
				Expect(exist).To(BeTrue())
				Expect(p).To(Equal(patient))
			})
		})
	})
})
