package handler_test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"github.com/synthia-telemed/backend-api/test/mock_hospital_client"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Info Handler", func() {
	var (
		mockCtrl        *gomock.Controller
		c               *gin.Context
		rec             *httptest.ResponseRecorder
		h               *handler.InfoHandler
		handlerFunc     gin.HandlerFunc
		patient         *datastore.Patient
		hospitalPatient *hospital.Patient

		mockPatientDataStore  *mock_datastore.MockPatientDataStore
		mockhospitalSysClient *mock_hospital_client.MockSystemClient
	)

	BeforeEach(func() {
		mockCtrl, rec, c = testhelper.InitHandlerTest()
		patient = testhelper.GeneratePatient()
		hospitalPatient = testhelper.GenerateHospitalPatient()
		mockPatientDataStore = mock_datastore.NewMockPatientDataStore(mockCtrl)
		mockhospitalSysClient = mock_hospital_client.NewMockSystemClient(mockCtrl)
		h = handler.NewInfoHandler(mockPatientDataStore, mockhospitalSysClient, zap.NewNop().Sugar())
		c.Set("Patient", patient)
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("ParseHospitalPatientInfo", func() {
		BeforeEach(func() {
			handlerFunc = h.ParseHospitalPatientInfo
		})

		When("Patient is not set in context", func() {
			BeforeEach(func() {
				c.Set("Patient", nil)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("Patient in context is not a type of datastore.Patient", func() {
			BeforeEach(func() {
				c.Set("Patient", testhelper.GenerateDoctor())
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("find patient by id error in hospital system client", func() {
			BeforeEach(func() {
				mockhospitalSysClient.EXPECT().FindPatientByID(gomock.Any(), patient.RefID).Return(nil, testhelper.MockError).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})
		When("patient info is not found", func() {
			BeforeEach(func() {
				mockhospitalSysClient.EXPECT().FindPatientByID(gomock.Any(), patient.RefID).Return(nil, nil).Times(1)
			})
			It("should return 400 with ErrPatientNotFound", func() {
				Expect(rec.Code).To(Equal(http.StatusNotFound))
				testhelper.AssertErrorResponseBody(rec.Body, handler.ErrPatientNotFound)
			})
		})
		When("patient info is found", func() {
			BeforeEach(func() {
				mockhospitalSysClient.EXPECT().FindPatientByID(gomock.Any(), patient.RefID).Return(hospitalPatient, nil).Times(1)
			})
			It("should return 200 with name in EN and TH", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				rawInfo, exist := c.Get("PatientInfo")
				Expect(exist).To(BeTrue())
				patientInfo, ok := rawInfo.(*hospital.Patient)
				Expect(ok).To(BeTrue())
				Expect(patientInfo).To(Equal(hospitalPatient))
			})
		})
	})

	Context("GetName", func() {
		BeforeEach(func() {
			handlerFunc = h.GetName
			c.Set("PatientInfo", hospitalPatient)
		})

		It("should return 200 with name in EN and TH", func() {
			Expect(rec.Code).To(Equal(http.StatusOK))
			var res *handler.GetNameResponse
			Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
			Expect(res.EN).To(Equal(hospitalPatient.NameEN))
			Expect(res.TH).To(Equal(hospitalPatient.NameTH))
		})
	})

	Context("GetPatientInfo", func() {
		BeforeEach(func() {
			handlerFunc = h.GetPatientInfo
			c.Set("PatientInfo", hospitalPatient)
		})

		It("should return 200 with patient information", func() {
			Expect(rec.Code).To(Equal(http.StatusOK))
			var res *hospital.Patient
			Expect(json.Unmarshal(rec.Body.Bytes(), &res)).To(Succeed())
			Expect(res.NationalId).To(Equal(hospitalPatient.NationalId))
		})
	})
})
