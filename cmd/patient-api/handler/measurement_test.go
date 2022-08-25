package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/cmd/patient-api/handler"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/synthia-telemed/backend-api/test/mock_datastore"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var _ = Describe("Measurement Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *handler.MeasurementHandler
		handlerFunc gin.HandlerFunc
		patientID   uint

		mockMeasurementDataStore *mock_datastore.MockMeasurementDataStore
	)

	BeforeEach(func() {
		mockCtrl, rec, c = initHandlerTest()
		patientID = uint(rand.Uint32())
		c.Set("UserID", patientID)

		mockMeasurementDataStore = mock_datastore.NewMockMeasurementDataStore(mockCtrl)
		h = handler.NewMeasurementHandler(mockMeasurementDataStore, zap.NewNop().Sugar())
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("CreateBloodPressure", func() {
		var req *handler.BloodPressureRequest
		BeforeEach(func() {
			handlerFunc = h.CreateBloodPressure
			req = &handler.BloodPressureRequest{
				DateTime:  time.Now(),
				Systolic:  uint(rand.Uint32()),
				Diastolic: uint(rand.Uint32()),
				Pulse:     uint(rand.Uint32()),
			}
			body, err := json.Marshal(&req)
			Expect(err).To(BeNil())
			c.Request = httptest.NewRequest("GET", "/", bytes.NewReader(body))
		})

		When("Request body is invalid", func() {
			BeforeEach(func() {
				c.Request = httptest.NewRequest("GET", "/", strings.NewReader(`{ "invalid": "error ja" }`))
			})
			It("should return 400", func() {
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})

		When("create blood pressure error", func() {
			BeforeEach(func() {
				mockMeasurementDataStore.EXPECT().CreateBloodPressure(gomock.Any()).Return(errors.New("error")).Times(1)
			})
			It("should return 500", func() {
				Expect(rec.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		When("no error", func() {
			BeforeEach(func() {
				mockMeasurementDataStore.EXPECT().CreateBloodPressure(gomock.Any()).Return(nil).Times(1)
			})
			It("should return 201 wth blood pressure", func() {
				Expect(rec.Code).To(Equal(http.StatusCreated))
				var bp datastore.BloodPressure
				Expect(json.Unmarshal(rec.Body.Bytes(), &bp)).To(Succeed())
				Expect(bp.DateTime).To(Equal(req.DateTime.UTC()))
				Expect(bp.PatientID).To(Equal(patientID))
				Expect(bp.Diastolic).To(Equal(req.Diastolic))
				Expect(bp.Systolic).To(Equal(req.Systolic))
				Expect(bp.Pulse).To(Equal(req.Pulse))
			})
		})

		Context("CreateGlucose", func() {
			var req *handler.GlucoseRequest
			BeforeEach(func() {
				handlerFunc = h.CreateGlucose
				b := false
				req = &handler.GlucoseRequest{
					DateTime:     time.Now(),
					Value:        uint(rand.Uint32()),
					IsBeforeMeal: &b,
				}
				body, err := json.Marshal(&req)
				Expect(err).To(BeNil())
				c.Request = httptest.NewRequest("GET", "/", bytes.NewReader(body))
			})

			When("Request body is invalid", func() {
				BeforeEach(func() {
					c.Request = httptest.NewRequest("GET", "/", strings.NewReader(`{ "invalid": "error ja" }`))
				})
				It("should return 400", func() {
					Expect(rec.Code).To(Equal(http.StatusBadRequest))
				})
			})

			When("create glucose error", func() {
				BeforeEach(func() {
					mockMeasurementDataStore.EXPECT().CreateGlucose(gomock.Any()).Return(errors.New("error")).Times(1)
				})
				It("should return 500", func() {
					Expect(rec.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			When("no error", func() {
				BeforeEach(func() {
					mockMeasurementDataStore.EXPECT().CreateGlucose(gomock.Any()).Return(nil).Times(1)
				})
				It("should return 201 wth blood pressure", func() {
					Expect(rec.Code).To(Equal(http.StatusCreated))
					var g datastore.Glucose
					Expect(json.Unmarshal(rec.Body.Bytes(), &g)).To(Succeed())
					Expect(g.DateTime).To(Equal(req.DateTime.UTC()))
					Expect(g.PatientID).To(Equal(patientID))
					Expect(g.Value).To(Equal(req.Value))
					Expect(g.IsBeforeMeal).To(Equal(*req.IsBeforeMeal))
				})
			})
		})
	})
})
