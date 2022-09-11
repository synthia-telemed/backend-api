package handler_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/doctor-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
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
		doctor                *datastore.Doctor
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockDoctorDataStore = mock_datastore.NewMockDoctorDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		h = handler.NewAppointmentHandler(mockPatientDataStore, mockDoctorDataStore, mockHospitalSysClient, mockCacheClient, mockClock, mockIDGenerator, zap.NewNop().Sugar())
		doctor = testhelper.GenerateDoctor()
	})
	JustBeforeEach(func() {
		handlerFunc(c)
	})

	Context("ParseDoctor", func() {
		BeforeEach(func() {
			handlerFunc = h.ParseDoctor
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

	Context("AuthorizedDoctorToAppointment", func() {
		var (
			appointment   *hospital.Appointment
			appointmentID int
		)

		BeforeEach(func() {
			handlerFunc = h.AuthorizedDoctorToAppointment
			appointment, appointmentID = testhelper.GenerateAppointment("", doctor.RefID, hospital.AppointmentStatusScheduled)
		})
		When("appointment ID is not provided", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", "")
			})
			It("should return 400 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrAppointmentIDMissing)
			})
		})
		When("appointment ID is invalid", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", "non-int")
			})
			It("should return 400 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrAppointmentIDInvalid)
			})
		})
		When("find appointment by ID error", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), appointmentID).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 404 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("appointment is not found", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), appointmentID).Return(nil, nil).Times(1)
			})
			It("should return 404 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusNotFound))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrAppointmentNotFound)
			})
		})
		When("Doctor is not set in context", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), appointmentID).Return(appointment, nil).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("Doctor in the context is not datastore.Doctor", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), appointmentID).Return(appointment, nil).Times(1)
				c.Set("Doctor", "anything")
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("Doctor doesn't own the appointment", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", appointment.Id)
				c.Set("Doctor", doctor)
				a, _ := testhelper.GenerateAppointment("", uuid.NewString(), hospital.AppointmentStatusScheduled)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), appointmentID).Return(a, nil).Times(1)
			})
			It("should return 403 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusForbidden))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrForbidden)
			})
		})
		When("Doctor own the appointment", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", appointment.Id)
				c.Set("Doctor", doctor)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), appointmentID).Return(appointment, nil).Times(1)
			})
			It("should set the appointment to context", func() {
				rawApp, existed := c.Get("Appointment")
				Expect(existed).To(BeTrue())
				app, ok := rawApp.(*hospital.Appointment)
				Expect(ok).To(BeTrue())
				Expect(app).To(Equal(appointment))
			})
		})
	})
})
