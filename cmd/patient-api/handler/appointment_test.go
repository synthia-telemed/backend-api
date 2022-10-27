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
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_cache_client"
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

		mockPatientDataStore     *mock_datastore.MockPatientDataStore
		mockPaymentDataStore     *mock_datastore.MockPaymentDataStore
		mockAppointmentDataStore *mock_datastore.MockAppointmentDataStore
		mockHospitalSysClient    *mock_hospital_client.MockSystemClient
		mockCacheClient          *mock_cache_client.MockClient
		mockClock                *mock_clock.MockClock

		patient *datastore.Patient
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockPaymentDataStore = mock_datastore.NewMockPaymentDataStore(mockCtrl)
		mockAppointmentDataStore = mock_datastore.NewMockAppointmentDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockCacheClient = mock_cache_client.NewMockClient(mockCtrl)
		mockClock = mock_clock.NewMockClock(mockCtrl)
		h = handler.NewAppointmentHandler(mockPatientDataStore, mockPaymentDataStore, mockAppointmentDataStore, mockHospitalSysClient, mockCacheClient, mockClock, zap.NewNop().Sugar())
		patient = testhelper.GeneratePatient()
		c.Set("Patient", patient)
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
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
				categorized  *hospital.CategorizedAppointment
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
				categorized = &hospital.CategorizedAppointment{
					Completed: completed,
					Scheduled: scheduled,
					Cancelled: cancelled,
				}
				mockHospitalSysClient.EXPECT().CategorizeAppointmentByStatus(appointments).Return(categorized)
			})
			It("should return 200 with list of appointments group by status", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res hospital.CategorizedAppointment
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Completed).To(HaveLen(n))
				Expect(res.Cancelled).To(HaveLen(n))
				Expect(res.Scheduled).To(HaveLen(n))
			})
		})
	})

	Context("AuthorizedPatientToAppointment", func() {
		BeforeEach(func() {
			handlerFunc = h.AuthorizedPatientToAppointment
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
				appointment, id := testhelper.GenerateAppointment("not-patient-id", "", hospital.AppointmentStatusScheduled, false)
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), id).Return(appointment, nil).Times(1)
			})
			It("should return 403", func() {
				Expect(rec.Code).To(Equal(http.StatusForbidden))
			})
		})
		When("appointment is found and own by patient", func() {
			var (
				appointment   *hospital.Appointment
				appointmentID int
			)
			BeforeEach(func() {
				appointment, appointmentID = testhelper.GenerateAppointment(patient.RefID, "", hospital.AppointmentStatusScheduled, false)
				c.AddParam("appointmentID", appointment.Id)
				mockHospitalSysClient.EXPECT().FindAppointmentByID(gomock.Any(), appointmentID).Return(appointment, nil).Times(1)
			})
			It("should set appointment to context", func() {
				rawApp, exist := c.Get("Appointment")
				Expect(exist).To(BeTrue())
				app, ok := rawApp.(*hospital.Appointment)
				Expect(ok).To(BeTrue())
				Expect(app).To(Equal(appointment))
			})
		})
	})

	Context("GetAppointment", func() {
		var (
			appointment   *hospital.Appointment
			appointmentDS *datastore.Appointment
		)

		BeforeEach(func() {
			handlerFunc = h.GetAppointment
			appointment, _ = testhelper.GenerateAppointment(patient.RefID, "", hospital.AppointmentStatusCompleted, true)
			appointmentDS = testhelper.GenerateDataStoreAppointment(appointment.Id)
			c.Set("Appointment", appointment)
		})

		When("appointment is found and it's scheduled", func() {
			BeforeEach(func() {
				appointment, _ = testhelper.GenerateAppointment(patient.RefID, "", hospital.AppointmentStatusScheduled, false)
				c.Set("Appointment", appointment)
			})
			It("should return appointment", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.GetAppointmentResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Id).To(Equal(appointment.Id))
				Expect(res.Payment).To(BeNil())
				Expect(res.Duration).To(BeNil())
			})
		})

		When("completed appointment is found but find local appointment error", func() {
			BeforeEach(func() {
				mockAppointmentDataStore.EXPECT().FindByRefID(appointment.Id).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("completed appointment is found but find payment error", func() {
			BeforeEach(func() {
				mockAppointmentDataStore.EXPECT().FindByRefID(appointment.Id).Return(appointmentDS, nil).Times(1)
				mockPaymentDataStore.EXPECT().FindLatestByInvoiceIDAndStatus(appointment.Invoice.Id, datastore.SuccessPaymentStatus).Return(nil, errors.New("err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("completed appointment is found and it's paid", func() {
			var (
				payment *datastore.Payment
			)
			BeforeEach(func() {
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
				mockAppointmentDataStore.EXPECT().FindByRefID(appointment.Id).Return(appointmentDS, nil).Times(1)
				mockPaymentDataStore.EXPECT().FindLatestByInvoiceIDAndStatus(appointment.Invoice.Id, datastore.SuccessPaymentStatus).Return(payment, nil).Times(1)
			})
			It("should return 200 with payment info", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.GetAppointmentResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Id).To(Equal(appointment.Id))
				Expect(res.Payment).ToNot(BeNil())
				Expect(res.Duration).ToNot(BeNil())
			})
		})
		When("completed appointment is found and it's have not been paid paid", func() {
			BeforeEach(func() {
				mockAppointmentDataStore.EXPECT().FindByRefID(appointment.Id).Return(appointmentDS, nil).Times(1)
				appointment.Invoice.Paid = false
			})
			It("should return 200 with payment is null", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.GetAppointmentResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Id).To(Equal(appointment.Id))
				Expect(res.Payment).To(BeNil())
				Expect(res.Duration).ToNot(BeNil())
			})
		})
	})

	Context("GetAppointmentRoomID", func() {
		var appointment *hospital.Appointment
		BeforeEach(func() {
			handlerFunc = h.GetAppointmentRoomID
			appointment, _ = testhelper.GenerateAppointment(patient.RefID, "", hospital.AppointmentStatusScheduled, false)
			c.Set("Appointment", appointment)
		})
		When("get roomID from cache error", func() {
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), cache.AppointmentRoomIDKey(appointment.Id), false).Return("", testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("get roomID from cache error", func() {
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), cache.AppointmentRoomIDKey(appointment.Id), false).Return("", nil).Times(1)
			})
			It("should return 404 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusNotFound))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrRoomIDNotFound)
			})
		})
		When("roomID of the appointment is found", func() {
			var roomID string
			BeforeEach(func() {
				roomID = uuid.NewString()
				mockCacheClient.EXPECT().Get(gomock.Any(), cache.AppointmentRoomIDKey(appointment.Id), false).Return(roomID, nil).Times(1)
			})
			It("should return 200 with roomID", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.GetAppointmentRoomIDResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.RoomID).To(Equal(roomID))
			})
		})
	})

	Context("GetNextScheduledAppointment", func() {
		var (
			appointment *hospital.AppointmentOverview
			where       *hospital.ListAppointmentsFilters
		)
		BeforeEach(func() {
			handlerFunc = h.GetNextScheduledAppointment
			appointment = testhelper.GenerateAppointmentOverview(hospital.AppointmentStatusScheduled)
			where = &hospital.ListAppointmentsFilters{
				PatientID: &patient.RefID,
				Status:    hospital.AppointmentStatusScheduled,
			}
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
		When("List appointment with filter query error", func() {
			BeforeEach(func() {
				mockHospitalSysClient.EXPECT().ListAppointmentsWithFilters(gomock.Any(), where, 1, 0).Return(nil, testhelper.MockError)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("Patient has no scheduled appointment", func() {
			BeforeEach(func() {
				mockHospitalSysClient.EXPECT().ListAppointmentsWithFilters(gomock.Any(), where, 1, 0).Return([]*hospital.AppointmentOverview{}, nil)
			})
			It("should return empty body with status of 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.String()).To(BeEmpty())
			})
		})
		When("Patient has at one or more scheduled appointment", func() {
			BeforeEach(func() {
				anotherAppointment := testhelper.GenerateAppointmentOverview(hospital.AppointmentStatusScheduled)
				mockHospitalSysClient.EXPECT().ListAppointmentsWithFilters(gomock.Any(), where, 1, 0).Return([]*hospital.AppointmentOverview{appointment, anotherAppointment}, nil)
			})
			It("should return empty body with status of 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res hospital.AppointmentOverview
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Id).To(Equal(appointment.Id))
			})
		})
	})

})
