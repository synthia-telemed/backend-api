package handler_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/doctor-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_cache_client"
	"github.com/synthia-telemed/backend-api/test/mock_clock"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
	"github.com/synthia-telemed/backend-api/test/mock_id"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Doctor Appointment Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.AppointmentHandler
		handlerFunc gin.HandlerFunc

		mockDoctorDataStore   *mock_datastore.MockDoctorDataStore
		mockPatientDataStore  *mock_datastore.MockPatientDataStore
		mockHospitalSysClient *mock_hospital_client.MockSystemClient
		mockCacheClient       *mock_cache_client.MockClient
		mockClock             *mock_clock.MockClock
		mockIDGenerator       *mock_id.MockGenerator
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockDoctorDataStore = mock_datastore.NewMockDoctorDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		h = handler.NewAppointmentHandler(mockPatientDataStore, mockDoctorDataStore, mockHospitalSysClient, mockCacheClient, mockClock, mockIDGenerator, zap.NewNop().Sugar())
	})
	JustBeforeEach(func() {
		handlerFunc(c)
	})

	Context("ParseDoctor", func() {
		var doctor *datastore.Doctor
		BeforeEach(func() {
			handlerFunc = h.ParseDoctor
			doctor = testhelper.GenerateDoctor()
			c.Set("UserID", doctor.ID)
		})

		When("find doctor by ID error", func() {
			BeforeEach(func() {
				mockDoctorDataStore.EXPECT().FindByID(doctor.ID).Return(nil, nil).Times(1)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrDoctorNotFound)
			})
		})
		When("doctor is not found", func() {
			BeforeEach(func() {
				mockDoctorDataStore.EXPECT().FindByID(doctor.ID).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("doctor is found", func() {
			BeforeEach(func() {
				mockDoctorDataStore.EXPECT().FindByID(doctor.ID).Return(doctor, nil).Times(1)
			})
			It("should set the doctor to context", func() {
				rawDoc, existed := c.Get("Doctor")
				Expect(existed).To(BeTrue())
				doc, ok := rawDoc.(*datastore.Doctor)
				Expect(ok).To(BeTrue())
				Expect(doc).To(Equal(doctor))
			})
		})
	})
})
