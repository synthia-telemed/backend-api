package server_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/server"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Server handler", func() {
	var (
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		h           *server.GinHandler
		handlerFunc gin.HandlerFunc
	)

	BeforeEach(func() {
		_, rec, c = testhelper.InitHandlerTest()
		h = &server.GinHandler{Logger: zap.NewNop().Sugar()}
	})

	JustBeforeEach(func() {
		handlerFunc(c)
	})

	Context("ParseUserID", func() {
		BeforeEach(func() {
			handlerFunc = h.ParseUserID
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		})

		When("X-USER-ID is not present", func() {
			It("should return 401", func() {
				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
				Expect(rec.Body.String()).To(Equal(`{"message":"Missing user ID"}`))
			})
		})

		When("X-USER-ID is invalid", func() {
			BeforeEach(func() {
				c.Request.Header.Set("X-USER-ID", "invalid")
			})
			It("should return 401", func() {
				Expect(rec.Code).To(Equal(http.StatusUnauthorized))
				Expect(rec.Body.String()).To(Equal(`{"message":"Invalid user ID"}`))
			})
		})

		When("X-USER-ID is valid", func() {
			BeforeEach(func() {
				c.Request.Header.Set("X-USER-ID", "99")
			})
			It("should return 200", func() {
				Expect(rec.Code).To(Equal(http.StatusOK))
				id, ok := c.Get("UserID")
				Expect(ok).To(BeTrue())
				Expect(id).To(Equal(uint(99)))
			})
		})
	})
})
