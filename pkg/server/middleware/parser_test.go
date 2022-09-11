package middleware_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/server/middleware"

	"net/http"
	"net/http/httptest"
)

var _ = Describe("Parser Middleware", func() {
	var (
		c           *gin.Context
		rec         *httptest.ResponseRecorder
		handlerFunc gin.HandlerFunc
	)

	BeforeEach(func() {
		rec = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
		handlerFunc = middleware.ParseUserID
	})

	JustBeforeEach(func() {
		handlerFunc(c)
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
