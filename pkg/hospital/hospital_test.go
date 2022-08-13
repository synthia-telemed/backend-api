package hospital_test

import (
	"context"
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
		mockCtrl = gomock.NewController(GinkgoT())
		graphQLClient = hospital.NewGraphQLClient(&hospital.Config{
			HospitalSysEndpoint: "https://hospital-mock.synthia.tech/graphql",
		})
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("FindPatientByGovCredential", func() {
		It("should find patient by passport ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "QH226189")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-106186"))
		})

		It("should find patient by national ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "6514582729055")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-127801"))
		})

		It("should return nil when patient not found", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "not-exist-national-id")
			Expect(err).To(BeNil())
			Expect(patient).To(BeNil())
		})
	})
})
