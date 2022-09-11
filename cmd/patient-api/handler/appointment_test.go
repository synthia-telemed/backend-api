package handler_test

import (
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
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_clock"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
	"go.uber.org/zap"
	"math/rand"
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

		patient *datastore.Patient
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockPaymentDataStore = mock_datastore.NewMockPaymentDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockClock = mock_clock.NewMockClock(mockCtrl)
		h = handler.NewAppointmentHandler(mockPatientDataStore, mockPaymentDataStore, mockHospitalSysClient, mockClock, zap.NewNop().Sugar())
		patient = testhelper.GeneratePatient()
		c.Set("Patient", patient)
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
				mockPatientDataStore.EXPECT().FindByID(id).Return(nil, errors.New("err")).Times(1)
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

	Context("ListAppointments", func() {
		BeforeEach(func() {
			handlerFunc = h.ListAppointments
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
				c.Set("Patient", testhelper.GenerateCreditCard())
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
				scheduled = testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusScheduled, n)
				cancelled = testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusCancelled, n)
				completed = testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusCompleted, n)
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
				testhelper.AssertListOfAppointments(res.Scheduled, hospital.AppointmentStatusScheduled, testhelper.ASC)
				testhelper.AssertListOfAppointments(res.Completed, hospital.AppointmentStatusCompleted, testhelper.DESC)
				testhelper.AssertListOfAppointments(res.Cancelled, hospital.AppointmentStatusCancelled, testhelper.DESC)
			})
		})
	})

	Context("GetAppointment", func() {
		BeforeEach(func() {
			handlerFunc = h.GetAppointment
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
				c.Set("Patient", testhelper.GenerateCreditCard())
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("AppointmentID is not provided", func() {
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("AppointmentID is not integer", func() {
			BeforeEach(func() {
				c.AddParam("appointmentID", "hi")
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
		When("hospital client FindAppointmentByID error", func() {
			BeforeEach(func() {
				id := int(rand.Int31())
				c.AddParam("appointmentID", fmt.Sprintf("%d", id))
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), id).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("appointment is not found", func() {
			BeforeEach(func() {
				id := int(rand.Int31())
				c.AddParam("appointmentID", fmt.Sprintf("%d", id))
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), id).Return(nil, nil).Times(1)
			})
			It("should return 404", func() {
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})
		When("appointment is not own by the patient", func() {
			BeforeEach(func() {
				appointment, id := testhelper.GenerateAppointment("not-patient-id", hospital.AppointmentStatusScheduled)
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), id).Return(appointment, nil).Times(1)
			})
			It("should return 403", func() {
				Expect(rec.Code).To(Equal(http.StatusForbidden))
			})
		})

		When("appointment is found and it's completed with find payment error", func() {
			BeforeEach(func() {
				appointment, id := testhelper.GenerateAppointment(patient.RefID, hospital.AppointmentStatusCompleted)
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), id).Return(appointment, nil).Times(1)
				mockPaymentDataStore.EXPECT().FindLatestByInvoiceIDAndStatus(appointment.Invoice.Id, datastore.SuccessPaymentStatus).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("appointment is found and it's completed", func() {
			var (
				payment     *datastore.Payment
				appointment *hospital.Appointment
			)
			BeforeEach(func() {
				var id int
				appointment, id = testhelper.GenerateAppointment(patient.RefID, hospital.AppointmentStatusCompleted)
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), id).Return(appointment, nil).Times(1)
				now := time.Now()
				card := testhelper.GenerateCreditCard()
				payment = &datastore.Payment{
					ID:           uint(rand.Uint32()),
					CreatedAt:    now,
					UpdatedAt:    now,
					Method:       datastore.CreditCardPaymentMethod,
					Amount:       rand.Float64() * 10000,
					PaidAt:       &now,
					ChargeID:     uuid.New().String(),
					InvoiceID:    appointment.Invoice.Id,
					Status:       datastore.SuccessPaymentStatus,
					CreditCard:   card,
					CreditCardID: &card.ID,
				}
				mockPaymentDataStore.EXPECT().FindLatestByInvoiceIDAndStatus(appointment.Invoice.Id, datastore.SuccessPaymentStatus).Return(payment, nil).Times(1)
			})
			It("should return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.GetAppointmentResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Id).To(Equal(appointment.Id))
				Expect(res.Payment).ToNot(BeNil())
			})
		})
		When("appointment is found and it's not completed", func() {
			var (
				appointment *hospital.Appointment
			)
			BeforeEach(func() {
				var id int
				appointment, id = testhelper.GenerateAppointment(patient.RefID, hospital.AppointmentStatusScheduled)
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), id).Return(appointment, nil).Times(1)
			})
			It("should return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.GetAppointmentResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Id).To(Equal(appointment.Id))
				Expect(res.Payment).To(BeNil())
			})
		})
	})

})
