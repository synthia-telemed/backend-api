package handler_test

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/test/mock_clock"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Appointment Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.AppointmentHandler
		handlerFunc gin.HandlerFunc

		mockPatientDataStore  *mock_datastore.MockPatientDataStore
		mockPaymentDataStore  *mock_datastore.MockPaymentDataStore
		mockHospitalSysClient *mock_hospital_client.MockSystemClient
		mockClock             *mock_clock.MockClock
	)

	BeforeEach(func() {
		mockCtrl, rec, c = initHandlerTest()
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockPaymentDataStore = mock_datastore.NewMockPaymentDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockClock = mock_clock.NewMockClock(mockCtrl)
		h = handler.NewAppointmentHandler(mockPatientDataStore, mockPaymentDataStore, mockHospitalSysClient, mockClock, zap.NewNop().Sugar())
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("ListAppointments", func() {
		var (
			patient *datastore.Patient
		)
		BeforeEach(func() {
			handlerFunc = h.ListAppointments
			patient = generatePatient()
			c.Set("Patient", patient)
		})
		When("Patient struct is not set", func() {
			BeforeEach(func() {
				c.Set("Patient", nil)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("Patient struct parsing error", func() {
			BeforeEach(func() {
				c.Set("Patient", generateCreditCard())
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("list appointments from hospital client error", func() {
			BeforeEach(func() {
				mockClock.EXPECT().Now().Return(time.Now()).Times(1)
				mockHospitalSysClient.EXPECT().ListAppointmentsByPatientID(gomock.Any(), patient.RefID, gomock.Any()).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("no error occurred", func() {
			var (
				n            int
				appointments []*hospital.AppointmentOverview
				scheduled    []*hospital.AppointmentOverview
				cancelled    []*hospital.AppointmentOverview
				completed    []*hospital.AppointmentOverview
			)
			BeforeEach(func() {
				mockClock.EXPECT().Now().Return(time.Now()).Times(1)
				n = 3
				scheduled = generateAppointmentOverviews(hospital.AppointmentStatusScheduled, n)
				cancelled = generateAppointmentOverviews(hospital.AppointmentStatusCancelled, n)
				completed = generateAppointmentOverviews(hospital.AppointmentStatusCompleted, n)
				appointments = make([]*hospital.AppointmentOverview, n*3)
				for i := 0; i < n; i++ {
					appointments[i*3+0] = scheduled[i]
					appointments[i*3+1] = cancelled[i]
					appointments[i*3+2] = completed[i]
				}
				mockHospitalSysClient.EXPECT().ListAppointmentsByPatientID(gomock.Any(), patient.RefID, gomock.Any()).Return(appointments, nil).Times(1)
			})
			It("should return 200 with list of appointments group by status", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.ListAppointmentsResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Completed).To(HaveLen(n))
				assertListOfAppointments(res.Scheduled, hospital.AppointmentStatusScheduled, ASC)
				assertListOfAppointments(res.Completed, hospital.AppointmentStatusCompleted, DESC)
				assertListOfAppointments(res.Cancelled, hospital.AppointmentStatusCancelled, DESC)
			})
		})
	})

})
