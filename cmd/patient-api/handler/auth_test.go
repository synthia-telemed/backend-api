package handler_test

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_cache_client"
	"github.com/synthia-telemed/backend-api/test/mock_clock"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
	"github.com/synthia-telemed/backend-api/test/mock_sms_client"
	"github.com/synthia-telemed/backend-api/test/mock_token_service"
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
		mockTokenService      *mock_token_service.MockService
		mockClock             *mock_clock.MockClock
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()

		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockSmsClient = mock_sms_client.NewMockClient(mockCtrl)
		mockCacheClient = mock_cache_client.NewMockClient(mockCtrl)
		mockTokenService = mock_token_service.NewMockService(mockCtrl)
		mockClock = mock_clock.NewMockClock(mockCtrl)
		h = handler.NewAuthHandler(mockPatientDataStore, mockHospitalSysClient, mockSmsClient, mockCacheClient, mockTokenService, mockClock, zap.NewNop().Sugar())
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
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
			var otpExpiredTime time.Time
			BeforeEach(func() {
				reqBody := strings.NewReader(`{"credential": "1234567890"}`)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", reqBody)
				p := &hospital.Patient{Id: "HN-1234", PhoneNumber: "0812223330"}
				mockHospitalSysClient.EXPECT().FindPatientByGovCredential(context.Background(), "1234567890").Return(p, nil).Times(1)
				mockCacheClient.EXPECT().Set(gomock.Any(), gomock.Any(), p.Id, time.Minute*10).Return(nil).Times(1)
				mockSmsClient.EXPECT().Send(p.PhoneNumber, gomock.Any()).Return(nil).Times(1)
				now := time.Now()
				mockClock.EXPECT().Now().Return(now).Times(1)
				otpExpiredTime = now.Add(10 * time.Minute)
			})

			It("should return 201 with phone number", func() {
				Expect(rec.Code).To(Equal(http.StatusCreated))
				var res handler.SigninResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.PhoneNumber).To(Equal("081***3330"))
				Expect(res.ExpiredAt.Equal(otpExpiredTime)).To(BeTrue())
			})
		})

	})

	Context("OTP Verification", func() {
		BeforeEach(func() {
			handlerFunc = h.VerifyOTP
		})

		When("request body is valid", func() {
			BeforeEach(func() {
				reqBody := strings.NewReader(`{"not-otp": "123456"}`)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", reqBody)
			})

			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("OTP is invalid or expired", func() {
			BeforeEach(func() {
				reqBody := strings.NewReader(`{"otp": "123456"}`)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", reqBody)
				mockCacheClient.EXPECT().Get(gomock.Any(), "123456", true).Return("", nil).Times(1)
			})

			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("OTP is valid", func() {
			BeforeEach(func() {
				reqBody := strings.NewReader(`{"otp": "123456"}`)
				c.Request, _ = http.NewRequest(http.MethodPost, "/", reqBody)
				mockCacheClient.EXPECT().Get(gomock.Any(), "123456", true).Return("HN-1234", nil).Times(1)
				mockPatientDataStore.EXPECT().FindOrCreate(gomock.Any()).Return(nil).Times(1)
				mockTokenService.EXPECT().GenerateToken(uint64(0), "Patient").Return("token", nil).Times(1)
			})

			It("should return 201 with token", func() {
				Expect(rec.Code).To(Equal(http.StatusCreated))
				var res handler.VerifyOTPResponse
				Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
				Expect(res.Token).To(Equal("token"))
			})
		})
	})

	Context("SignOut", func() {
		var (
			patient        *datastore.Patient
			updatedPatient *datastore.Patient
		)
		BeforeEach(func() {
			handlerFunc = h.SignOut
			patient = testhelper.GeneratePatient()
			c.Set("Patient", patient)
			p := *patient
			p.NotificationToken = ""
			updatedPatient = &p
		})

		When("update patient's notification token error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().Save(updatedPatient).Return(testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("no error", func() {
			BeforeEach(func() {
				mockPatientDataStore.EXPECT().Save(updatedPatient).Return(nil).Times(1)
			})
			It("should return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})

	})
})
