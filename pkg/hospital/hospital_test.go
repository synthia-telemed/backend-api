package hospital_test

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
)

var _ = Describe("Hospital Client", func() {

	var (
		mockCtrl      *gomock.Controller
		graphQLClient *hospital.GraphQLClient
	)

	BeforeEach(func() {
		var c hospital.Config
		Expect(env.Parse(&c)).To(Succeed())
		mockCtrl = gomock.NewController(GinkgoT())
		graphQLClient = hospital.NewGraphQLClient(&c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("FindPatientByGovCredential", func() {
		It("should find patient by passport ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "SO629265")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-525661"))
		})

		It("should find patient by national ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "3089074169079")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-937553"))
		})

		It("should return nil when patient not found", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "not-exist-national-id")
			Expect(err).To(BeNil())
			Expect(patient).To(BeNil())
		})
	})
})
