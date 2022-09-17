package id_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/id"
)

var _ = Describe("NanoID", func() {
	var (
		generator id.Generator
	)

	BeforeEach(func() {
		generator = id.NewNanoID()
	})

	Context("GenerateRoomID", func() {
		It("should generate id with length of 21", func() {
			id, err := generator.GenerateRoomID()
			Expect(err).To(BeNil())
			Expect(id).To(HaveLen(21))
		})
	})
})
