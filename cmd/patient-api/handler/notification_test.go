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
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Notification Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.NotificationHandler
		handlerFunc gin.HandlerFunc
		patientID   uint

		mockNotificationDataStore *mock_datastore.MockNotificationDataStore
		mockPatientDataStore      *mock_datastore.MockPatientDataStore
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockNotificationDataStore = mock_datastore.NewMockNotificationDataStore(mockCtrl)
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		h = handler.NewNotificationHandler(mockNotificationDataStore, mockPatientDataStore, zap.NewNop().Sugar())
		patientID = uint(rand.Uint32())
		c.Set("UserID", patientID)
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("ListNotifications", func() {
		BeforeEach(func() {
			handlerFunc = h.ListNotifications
		})

		When("list notification from datastore error", func() {
			BeforeEach(func() {
				mockNotificationDataStore.EXPECT().ListLatest(patientID).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("no error occurred", func() {
			var notifications []datastore.Notification
			BeforeEach(func() {
				notifications, _ = testhelper.GenerateNotifications(patientID, 5)
				mockNotificationDataStore.EXPECT().ListLatest(patientID).Return(notifications, nil).Times(1)
			})
			It("should return list of notifications", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res []datastore.Notification
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res).To(Equal(notifications))
			})
		})
	})

	Context("CountUnRead", func() {
		BeforeEach(func() {
			handlerFunc = h.CountUnRead
		})

		When("count unread notification error", func() {
			BeforeEach(func() {
				mockNotificationDataStore.EXPECT().CountUnRead(patientID).Return(0, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("no error occurred", func() {
			var count int
			BeforeEach(func() {
				count = rand.Int()
				mockNotificationDataStore.EXPECT().CountUnRead(patientID).Return(count, nil).Times(1)
			})
			It("should the count", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				var res handler.CountUnReadNotificationResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Count).To(Equal(count))
			})
		})
	})

	Context("AuthorizedPatientToNotification", func() {
		BeforeEach(func() {
			handlerFunc = h.AuthorizedPatientToNotification

		})

		When("notification id is not set in param", func() {
			It("should return 400 with ErrInvalidNotificationID", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrInvalidNotificationID)
			})
		})
		When("notification id in param in not a unsigned integer", func() {
			BeforeEach(func() {
				c.AddParam("id", "not-uint")
			})
			It("should return 400 with ErrInvalidNotificationID", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrInvalidNotificationID)
			})
		})

		Context("id param is properly set", func() {
			var notification datastore.Notification
			BeforeEach(func() {
				notification = testhelper.GenerateNotification(patientID)
				c.AddParam("id", fmt.Sprintf("%d", notification.ID))
			})

			When("find notification by id error", func() {
				BeforeEach(func() {
					mockNotificationDataStore.EXPECT().FindByID(notification.ID).Return(nil, testhelper.MockError).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
			When("notification is not found", func() {
				BeforeEach(func() {
					mockNotificationDataStore.EXPECT().FindByID(notification.ID).Return(nil, nil).Times(1)
				})
				It("should return 404", func() {
					Expect(rec.Code).To(Equal(http.StatusNotFound))
					testhelper.AssertErrorResponseBody(rec.Body, handler.ErrNotificationNotFound)
				})
			})
			When("notification is owned by the patient", func() {
				BeforeEach(func() {
					otherNotification := testhelper.GenerateNotification(uint(rand.Uint32()))
					mockNotificationDataStore.EXPECT().FindByID(notification.ID).Return(&otherNotification, nil).Times(1)
				})
				It("should return 403", func() {
					Expect(rec.Code).To(Equal(http.StatusForbidden))
					testhelper.AssertErrorResponseBody(rec.Body, handler.ErrForbidden)
				})
			})
			When("no error", func() {
				BeforeEach(func() {
					mockNotificationDataStore.EXPECT().FindByID(notification.ID).Return(&notification, nil).Times(1)
				})
				It("should set the notification to context", func() {
					Expect(rec.Code).To(Equal(http.StatusOK))
					rawNotification, exists := c.Get("Notification")
					Expect(exists).To(BeTrue())
					retrievedNotification, ok := rawNotification.(*datastore.Notification)
					Expect(ok).To(BeTrue())
					Expect(retrievedNotification).To(Equal(&notification))
				})
			})
		})
	})

	Context("Read notification", func() {
		var notification datastore.Notification
		BeforeEach(func() {
			handlerFunc = h.Read
			notification = testhelper.GenerateNotification(patientID)
			c.Set("Notification", &notification)
		})

		When("set read status error", func() {
			BeforeEach(func() {
				mockNotificationDataStore.EXPECT().SetAsRead(notification.ID).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("no error", func() {
			BeforeEach(func() {
				mockNotificationDataStore.EXPECT().SetAsRead(notification.ID).Return(nil).Times(1)
			})
			It("should return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("Read all notifications", func() {
		BeforeEach(func() {
			handlerFunc = h.ReadAll
		})

		When("set all notification read status error", func() {
			BeforeEach(func() {
				mockNotificationDataStore.EXPECT().SetAllAsRead(patientID).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("no error", func() {
			BeforeEach(func() {
				mockNotificationDataStore.EXPECT().SetAllAsRead(patientID).Return(nil).Times(1)
			})
			It("should return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("Set notification token", func() {
		var (
			req     *handler.SetNotificationTokenRequest
			patient *datastore.Patient
		)
		BeforeEach(func() {
			handlerFunc = h.SetNotificationToken
			patient = testhelper.GeneratePatient()
			req = &handler.SetNotificationTokenRequest{Token: uuid.NewString()}
			body, err := json.Marshal(req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("GET", "/", bytes.NewReader(body))
		})

		When("request body is invalid", func() {
			BeforeEach(func() {
				c.Request = httptest.NewRequest("GET", "/", strings.NewReader(`{"not-token": "wasd"}`))
			})
			It("should return 400 with error", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrInvalidRequestBody)
			})
		})
		When("Patient is not set in context", func() {
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("Patient in the context is not *datastore.Patient", func() {
			BeforeEach(func() {
				c.Set("Patient", &datastore.Payment{})
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("patient struct is properly set", func() {
			var targetPatient datastore.Patient
			BeforeEach(func() {
				c.Set("Patient", patient)
				targetPatient = *patient
				targetPatient.NotificationToken = req.Token
			})

			When("save patient error", func() {
				BeforeEach(func() {
					mockPatientDataStore.EXPECT().Save(&targetPatient).Return(testhelper.MockError).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
			When("no error", func() {
				BeforeEach(func() {
					mockPatientDataStore.EXPECT().Save(&targetPatient).Return(nil).Times(1)
				})
				It("should return 200", func() {
					Expect(rec.Code).To(Equal(http.StatusOK))
				})
			})
		})
	})
})
