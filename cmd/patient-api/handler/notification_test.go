package handler_test

import (
	"encoding/json"
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
})
