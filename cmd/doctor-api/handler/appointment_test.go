package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/doctor-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/cache"
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
	"strings"
	"time"
)

var _ = Describe("Doctor Appointment Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.AppointmentHandler
		handlerFunc gin.HandlerFunc

		mockDoctorDataStore      *mock_datastore.MockDoctorDataStore
		mockAppointmentDataStore *mock_datastore.MockAppointmentDataStore
		mockPatientDataStore     *mock_datastore.MockPatientDataStore
		mockHospitalSysClient    *mock_hospital_client.MockSystemClient
		mockCacheClient          *mock_cache_client.MockClient
		mockClock                *mock_clock.MockClock
		mockIDGenerator          *mock_id.MockGenerator
		doctor                   *datastore.Doctor
		appointment              *hospital.Appointment
		appointmentID            int
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockDoctorDataStore = mock_datastore.NewMockDoctorDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockAppointmentDataStore = mock_datastore.NewMockAppointmentDataStore(mockCtrl)
		mockClock = mock_clock.NewMockClock(mockCtrl)
		mockCacheClient = mock_cache_client.NewMockClient(mockCtrl)
		mockIDGenerator = mock_id.NewMockGenerator(mockCtrl)
		h = handler.NewAppointmentHandler(mockAppointmentDataStore, mockPatientDataStore, mockDoctorDataStore, mockHospitalSysClient, mockCacheClient, mockClock, mockIDGenerator, zap.NewNop().Sugar())
		doctor = testhelper.GenerateDoctor()
		appointment, appointmentID = testhelper.GenerateAppointment("", doctor.RefID, hospital.AppointmentStatusScheduled, false)
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
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
			appointment, appointmentID = testhelper.GenerateAppointment("", doctor.RefID, hospital.AppointmentStatusScheduled, false)
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
				a, _ := testhelper.GenerateAppointment("", uuid.NewString(), hospital.AppointmentStatusScheduled, false)
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

	Context("InitAppointmentRoom", func() {
		BeforeEach(func() {
			handlerFunc = h.InitAppointmentRoom
			c.Set("Doctor", doctor)
			c.Set("Appointment", appointment)
		})

		When("appointment doesn't have schedule status", func() {
			BeforeEach(func() {
				appointment.Status = hospital.AppointmentStatusCompleted
			})
			It("should return 400 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrInitNonScheduledAppointment)
			})
		})
		When("init room earlier than 10 minutes of the appointment time", func() {
			BeforeEach(func() {
				mockClock.EXPECT().Now().Return(appointment.StartDateTime.Add(-time.Minute * 10).Add(-time.Second))
			})
			It("should return 400 with not time yet error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrNotTimeYet)
			})
		})
		When("init room later than 3 hours of the appointment time", func() {
			BeforeEach(func() {
				mockClock.EXPECT().Now().Return(appointment.StartDateTime.Add(time.Hour * 3).Add(time.Second))
			})
			It("should return 400 with not time yet error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrNotTimeYet)
			})
		})
		When("get current appointment of the doctor from cache error", func() {
			BeforeEach(func() {
				mockClock.EXPECT().Now().Return(appointment.StartDateTime.Add(time.Minute)).Times(1)
				mockCacheClient.EXPECT().Get(gomock.Any(), cache.CurrentDoctorAppointmentIDKey(doctor.ID), false).Return("", testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		Context("doctor is currently in a room", func() {
			var clockTime time.Time
			BeforeEach(func() {
				clockTime = appointment.StartDateTime.Add(time.Minute)
				mockClock.EXPECT().Now().Return(clockTime).Times(1)
			})
			When("doctor is in another room", func() {
				BeforeEach(func() {
					mockCacheClient.EXPECT().Get(gomock.Any(), cache.CurrentDoctorAppointmentIDKey(doctor.ID), false).Return(uuid.NewString(), nil).Times(1)
				})
				It("should return 400 with error", func() {
					Expect(rec.Code).To(Equal(http.StatusBadRequest))
					testhelper.AssertErrorResponseBody(rec.Body, handler.ErrDoctorInAnotherRoom)
				})
			})
			Context("doctor is in the designated room", func() {
				BeforeEach(func() {
					mockCacheClient.EXPECT().Get(gomock.Any(), cache.CurrentDoctorAppointmentIDKey(doctor.ID), false).Return(appointment.Id, nil).Times(1)
				})
				When("get room ID from cache error", func() {
					BeforeEach(func() {
						mockCacheClient.EXPECT().Get(gomock.Any(), cache.AppointmentRoomIDKey(appointment.Id), false).Return("", testhelper.MockError).Times(1)
					})
					It("should return 500", func() {
						Expect(rec.Code).To(Equal(http.StatusInternalServerError))
					})
				})
				When("successfully get room ID from cache", func() {
					var roomID string
					BeforeEach(func() {
						roomID = uuid.NewString()
						mockCacheClient.EXPECT().Get(gomock.Any(), cache.AppointmentRoomIDKey(appointment.Id), false).Return(roomID, nil).Times(1)
					})
					It("should return 201 with room ID", func() {
						Expect(rec.Code).To(Equal(http.StatusCreated))
						var res handler.InitAppointmentRoomResponse
						Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
						Expect(res.RoomID).To(Equal(roomID))
					})
				})
			})
			Context("doctor is not in any room", func() {
				BeforeEach(func() {
					mockCacheClient.EXPECT().Get(gomock.Any(), cache.CurrentDoctorAppointmentIDKey(doctor.ID), false).Return("", nil).Times(1)
				})
				When("room id generation error", func() {
					BeforeEach(func() {
						mockIDGenerator.EXPECT().GenerateRoomID().Return("", testhelper.MockError).Times(1)
					})
					It("should return 500", func() {
						Expect(rec.Code).To(Equal(http.StatusInternalServerError))
					})
				})
				Context("successful generated roomID", func() {
					var (
						patient *datastore.Patient
						roomID  string
						kv      map[string]string
					)
					BeforeEach(func() {
						roomID = uuid.NewString()
						mockIDGenerator.EXPECT().GenerateRoomID().Return(roomID, nil).Times(1)
						kv = map[string]string{
							cache.CurrentDoctorAppointmentIDKey(doctor.ID): appointment.Id,
							cache.AppointmentRoomIDKey(appointment.Id):     roomID,
						}
						patient = testhelper.GeneratePatient()
						appointment.PatientID = patient.RefID
					})

					When("set current appointment of the doctor and room ID of appointment to cache error", func() {
						BeforeEach(func() {
							mockCacheClient.EXPECT().MultipleSet(gomock.Any(), kv).Return(testhelper.MockError).Times(1)
						})
						It("should return 500", func() {
							Expect(rec.Code).To(Equal(http.StatusInternalServerError))
						})
					})
					When("find patient by ID error", func() {
						BeforeEach(func() {
							mockCacheClient.EXPECT().MultipleSet(gomock.Any(), kv).Return(nil).Times(1)
							mockPatientDataStore.EXPECT().FindByRefID(appointment.PatientID).Return(nil, testhelper.MockError).Times(1)
						})
						It("should return 500", func() {
							Expect(rec.Code).To(Equal(http.StatusInternalServerError))
						})
					})
					When("set room information to cache error", func() {
						BeforeEach(func() {
							mockCacheClient.EXPECT().MultipleSet(gomock.Any(), kv).Return(nil).Times(1)
							mockPatientDataStore.EXPECT().FindByRefID(appointment.PatientID).Return(patient, nil).Times(1)
							info := map[string]string{
								"PatientID":     fmt.Sprintf("%d", patient.ID),
								"DoctorID":      fmt.Sprintf("%d", doctor.ID),
								"AppointmentID": appointment.Id,
								"StartedAt":     clockTime.Format(time.RFC3339),
							}
							mockCacheClient.EXPECT().HashSet(gomock.Any(), cache.RoomInfoKey(roomID), info).Return(testhelper.MockError).Times(1)
						})
						It("should return 500", func() {
							Expect(rec.Code).To(Equal(http.StatusInternalServerError))
						})
					})
					When("successfully set room info to cache", func() {
						BeforeEach(func() {
							mockCacheClient.EXPECT().MultipleSet(gomock.Any(), kv).Return(nil).Times(1)
							mockPatientDataStore.EXPECT().FindByRefID(appointment.PatientID).Return(patient, nil).Times(1)
							info := map[string]string{
								"PatientID":     fmt.Sprintf("%d", patient.ID),
								"DoctorID":      fmt.Sprintf("%d", doctor.ID),
								"AppointmentID": appointment.Id,
								"StartedAt":     clockTime.Format(time.RFC3339),
							}
							mockCacheClient.EXPECT().HashSet(gomock.Any(), cache.RoomInfoKey(roomID), info).Return(nil).Times(1)
						})
						It("should return 201 with room ID", func() {
							Expect(rec.Code).To(Equal(http.StatusCreated))
							var res handler.InitAppointmentRoomResponse
							Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
							Expect(res.RoomID).To(Equal(roomID))
						})
					})
				})
			})
		})
	})

	Context("CompleteAppointment", func() {
		var (
			roomID                   string
			getCurrentAppointmentKey string
			getRoomIDKey             string
			getRoomInfoKey           string
			now                      time.Time
			duration                 time.Duration
			startedTime              time.Time
			req                      *handler.CompleteAppointmentRequest
		)

		BeforeEach(func() {
			handlerFunc = h.CompleteAppointment
			c.Set("Doctor", doctor)
			roomID = uuid.NewString()
			getCurrentAppointmentKey = cache.CurrentDoctorAppointmentIDKey(doctor.ID)
			getRoomIDKey = cache.AppointmentRoomIDKey(appointment.Id)
			getRoomInfoKey = cache.RoomInfoKey(roomID)
			now = time.Now()
			duration = (time.Minute * 30) + (time.Second * 10)
			startedTime = now.Add(-duration).Round(time.Second)
			req = &handler.CompleteAppointmentRequest{Status: hospital.SettableAppointmentStatusCompleted}
			body, err := json.Marshal(req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("post", "/", bytes.NewReader(body))
		})

		When("request body is invalid", func() {
			BeforeEach(func() {
				body := fmt.Sprintf(`{"status": "%s"}`, hospital.AppointmentStatusScheduled)
				c.Request = httptest.NewRequest("post", "/", strings.NewReader(body))
			})
			It("should return 400 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrInvalidRequestBody)
			})
		})
		When("get current appointment ID from cache error", func() {
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), getCurrentAppointmentKey, false).Return("", testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("current appointment ID is not found", func() {
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), getCurrentAppointmentKey, false).Return("", nil).Times(1)
			})
			It("should return 400 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrDoctorNotInRoom)
			})
		})
		When("get room ID from cache error", func() {
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), getCurrentAppointmentKey, false).Return(appointment.Id, nil).Times(1)
				mockCacheClient.EXPECT().Get(gomock.Any(), getRoomIDKey, false).Return("", testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("get started time from cache error", func() {
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), getCurrentAppointmentKey, false).Return(appointment.Id, nil).Times(1)
				mockCacheClient.EXPECT().Get(gomock.Any(), getRoomIDKey, false).Return(roomID, nil).Times(1)
				mockCacheClient.EXPECT().HashGet(gomock.Any(), getRoomInfoKey, "StartedAt").Return("", testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("parsing started time error", func() {
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), getCurrentAppointmentKey, false).Return(appointment.Id, nil).Times(1)
				mockCacheClient.EXPECT().Get(gomock.Any(), getRoomIDKey, false).Return(roomID, nil).Times(1)
				mockCacheClient.EXPECT().HashGet(gomock.Any(), getRoomInfoKey, "StartedAt").Return(startedTime.Format(time.Stamp), nil).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		Context("getting information and parse from cache success", func() {
			var dbAppointment *datastore.Appointment
			BeforeEach(func() {
				mockCacheClient.EXPECT().Get(gomock.Any(), getCurrentAppointmentKey, false).Return(appointment.Id, nil).Times(1)
				mockCacheClient.EXPECT().Get(gomock.Any(), getRoomIDKey, false).Return(roomID, nil).Times(1)
				mockCacheClient.EXPECT().HashGet(gomock.Any(), getRoomInfoKey, "StartedAt").Return(startedTime.Format(time.RFC3339), nil).Times(1)
				mockClock.EXPECT().Now().Return(now).Times(1)
				dbAppointment = &datastore.Appointment{
					RefID:       appointment.Id,
					Duration:    duration.Seconds(),
					StartedTime: startedTime.UTC(),
				}
			})
			When("save appointment to db error", func() {
				BeforeEach(func() {
					mockAppointmentDataStore.EXPECT().Create(dbAppointment).Return(testhelper.MockError).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
			When("delete appointment and room information in cache error", func() {
				BeforeEach(func() {
					mockAppointmentDataStore.EXPECT().Create(dbAppointment).Return(nil).Times(1)
					mockCacheClient.EXPECT().Delete(gomock.Any(), gomock.InAnyOrder([]string{getRoomInfoKey, getRoomIDKey, getCurrentAppointmentKey})).Return(testhelper.MockError).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
			When("set status of appointment in hospital sys to complete error", func() {
				BeforeEach(func() {
					mockAppointmentDataStore.EXPECT().Create(dbAppointment).Return(nil).Times(1)
					mockCacheClient.EXPECT().Delete(gomock.Any(), gomock.InAnyOrder([]string{getRoomInfoKey, getRoomIDKey, getCurrentAppointmentKey})).Return(nil).Times(1)
					mockHospitalSysClient.EXPECT().SetAppointmentStatus(gomock.Any(), appointmentID, req.Status).Return(testhelper.MockError).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
			When("no error occurred", func() {
				BeforeEach(func() {
					mockAppointmentDataStore.EXPECT().Create(dbAppointment).Return(nil).Times(1)
					mockCacheClient.EXPECT().Delete(gomock.Any(), gomock.InAnyOrder([]string{getRoomInfoKey, getRoomIDKey, getCurrentAppointmentKey})).Return(nil).Times(1)
					mockHospitalSysClient.EXPECT().SetAppointmentStatus(gomock.Any(), appointmentID, req.Status).Return(nil).Times(1)
				})
				It("should return 200 with duration", func() {
					Expect(rec.Code).To(Equal(http.StatusOK))
					var res handler.CompleteAppointmentResponse
					Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
					Expect(res).To(Equal(handler.CompleteAppointmentResponse{Duration: duration.Seconds()}))
				})
			})
		})
	})

	Context("TodayAppointment", func() {
		var (
			now time.Time
		)

		BeforeEach(func() {
			handlerFunc = h.TodayAppointment
			c.Set("Doctor", doctor)

			loc, _ := time.LoadLocation("Asia/Bangkok")
			now = time.Now().In(loc)
			mockClock.EXPECT().Now().Return(now).Times(1)
		})

		When("list appointments by doctor id graphQL error", func() {
			BeforeEach(func() {
				mockHospitalSysClient.EXPECT().ListAppointmentsByDoctorID(gomock.Any(), doctor.RefID, now).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("no error occurred", func() {
			var (
				n int
			)
			BeforeEach(func() {
				n = 3
				scheduled := testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusScheduled, n)
				cancelled := testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusCancelled, n)
				completed := testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusCompleted, n)
				appointments := make([]*hospital.AppointmentOverview, n*3)
				for i := 0; i < n; i++ {
					appointments[i*3+0] = scheduled[i]
					appointments[i*3+1] = cancelled[i]
					appointments[i*3+2] = completed[i]
				}
				categorized := &hospital.CategorizedAppointment{
					Completed: completed,
					Scheduled: scheduled,
					Cancelled: cancelled,
				}
				mockHospitalSysClient.EXPECT().ListAppointmentsByDoctorID(gomock.Any(), doctor.RefID, now).Return(appointments, nil).Times(1)
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
})
