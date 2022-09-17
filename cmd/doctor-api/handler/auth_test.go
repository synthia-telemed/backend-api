package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/doctor-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
	"github.com/synthia-telemed/backend-api/test/mock_token_service"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Doctor Auth Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.AuthHandler
		handlerFunc gin.HandlerFunc

		mockDoctorDataStore   *mock_datastore.MockDoctorDataStore
		mockHospitalSysClient *mock_hospital_client.MockSystemClient
		mockTokenService      *mock_token_service.MockService
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		mockDoctorDataStore = mock_datastore.NewMockDoctorDataStore(mockCtrl)
		mockHospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		mockTokenService = mock_token_service.NewMockService(mockCtrl)
		h = handler.NewAuthHandler(mockHospitalSysClient, mockTokenService, mockDoctorDataStore, zap.NewNop().Sugar())
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Signin", func() {
		var req handler.SigninRequest
		BeforeEach(func() {
			handlerFunc = h.Signin
			req = handler.SigninRequest{Username: "doctor-a", Password: "password"}
			reqBody, err := json.Marshal(req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("post", "/", bytes.NewReader(reqBody))
		})

		When("request body is invalid", func() {
			BeforeEach(func() {
				c.Request = httptest.NewRequest("post", "/", strings.NewReader(`{"not-username": "sethanantp"}`))
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("doctor credential is invalid", func() {
			BeforeEach(func() {
				mockHospitalSysClient.EXPECT().AssertDoctorCredential(gomock.Any(), req.Username, req.Password).Return(false, nil).Times(1)
			})
			It("should return 401", func() {
				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		When("hospitalSysClient.AssertDoctorCredential error", func() {
			BeforeEach(func() {
				mockHospitalSysClient.EXPECT().AssertDoctorCredential(gomock.Any(), req.Username, req.Password).Return(false, errors.New("some-err")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("doctor credential is valid", func() {
			var (
				token       string
				queryDoctor *hospital.Doctor
			)

			BeforeEach(func() {
				queryDoctor = &hospital.Doctor{Id: fmt.Sprintf("doc-%d", rand.Int())}
				token = "token"
				mockHospitalSysClient.EXPECT().AssertDoctorCredential(gomock.Any(), req.Username, req.Password).Return(true, nil).Times(1)
			})

			When("no error occurred", func() {
				BeforeEach(func() {
					mockHospitalSysClient.EXPECT().FindDoctorByUsername(gomock.Any(), req.Username).Return(queryDoctor, nil).Times(1)
					mockDoctorDataStore.EXPECT().FindOrCreate(&datastore.Doctor{RefID: queryDoctor.Id}).Return(nil).Times(1)
					mockTokenService.EXPECT().GenerateToken(uint64(0), "Doctor").Return(token, nil).Times(1)
				})
				It("should return 201 with token", func() {
					var res handler.SigninResponse
					Expect(rec.Code).To(Equal(http.StatusCreated))
					Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
					Expect(res.Token).To(Equal(token))
				})
			})

			When("hospitalSysClient.FindDoctorByUsername error", func() {
				BeforeEach(func() {
					mockHospitalSysClient.EXPECT().FindDoctorByUsername(gomock.Any(), req.Username).Return(nil, errors.New("err")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
			When("doctorDataStore.FindOrCreate error", func() {
				BeforeEach(func() {
					mockHospitalSysClient.EXPECT().FindDoctorByUsername(gomock.Any(), req.Username).Return(queryDoctor, nil).Times(1)
					mockDoctorDataStore.EXPECT().FindOrCreate(&datastore.Doctor{RefID: queryDoctor.Id}).Return(errors.New("err")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
			When("tokenService.GenerateToken error", func() {
				BeforeEach(func() {
					mockHospitalSysClient.EXPECT().FindDoctorByUsername(gomock.Any(), req.Username).Return(queryDoctor, nil).Times(1)
					mockDoctorDataStore.EXPECT().FindOrCreate(&datastore.Doctor{RefID: queryDoctor.Id}).Return(nil).Times(1)
					mockTokenService.EXPECT().GenerateToken(uint64(0), "Doctor").Return("", errors.New("err")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})

})
