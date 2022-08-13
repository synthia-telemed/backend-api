package handler_test

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"github.com/synthia-telemed/backend-api/test/mock_cache_client"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
	"github.com/synthia-telemed/backend-api/test/mock_sms_client"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var _ = Describe("Auth Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.AuthHandler
		handlerFunc gin.HandlerFunc

		mockPatientDataStore  *mock_datastore.MockPatientDataStore
		mockHospitalSysClient *mock_hospital_client.MockSystemClient
		mockSmsClient         *mock_sms_client.MockClient
		mockCacheClient       *mock_cache_client.MockClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		rec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(rec)

		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockSmsClient = mock_sms_client.NewMockClient(mockCtrl)
		mockCacheClient = mock_cache_client.NewMockClient(mockCtrl)
		h = handler.NewAuthHandler(mockPatientDataStore, mockHospitalSysClient, mockSmsClient, mockCacheClient, zap.NewNop().Sugar())
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	Context("Signin", func() {
		BeforeEach(func() {
			handlerFunc = h.Signin
		})

		When("request body is valid", func() {
			BeforeEach(func() {
				reqBody := strings.NewReader(`{"not-credential": "1234567890"}`)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", reqBody)
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("patient is not found", func() {
			BeforeEach(func() {
				reqBody := strings.NewReader(`{"credential": "1234567890"}`)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", reqBody)
				mockHospitalSysClient.EXPECT().FindPatientByGovCredential(context.Background(), "1234567890").Return(nil, nil).Times(1)
			})
			It("should return 404", func() {
				Expect(rec.Code).To(Equal(http.StatusNotFound))
			})
		})

		When("patient is found", func() {
			p := &hospital.Patient{Id: "HN-1234", PhoneNumber: "0812223330"}
			BeforeEach(func() {
				reqBody := strings.NewReader(`{"credential": "1234567890"}`)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", reqBody)
				mockHospitalSysClient.EXPECT().FindPatientByGovCredential(context.Background(), "1234567890").Return(p, nil).Times(1)
				mockCacheClient.EXPECT().Set(context.Background(), gomock.Any(), p.Id, time.Minute*10).Return(nil).Times(1)
				mockSmsClient.EXPECT().Send(p.PhoneNumber, gomock.Any()).Return(nil).Times(1)
			})

			It("should return 201 with phone number", func() {
				Expect(rec.Code).To(Equal(http.StatusCreated))
				var res handler.LoginResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(BeNil())
				Expect(res.PhoneNumber).To(Equal("081***3330"))
			})
		})

	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

})
