package hospital_test

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"math/rand"
)

var _ = Describe("Hospital Client", func() {

	var (
		mockCtrl      *gomock.Controller
		graphQLClient *hospital.GraphQLClient
	)

	BeforeEach(func() {
		rand.Seed(GinkgoRandomSeed())
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

	Context("AssertDoctorCredential", func() {
		When("doctor's username is not found", func() {
			It("should return false", func() {
				assertion, err := graphQLClient.AssertDoctorCredential(context.Background(), "not-exist-doctor", "password")
				Expect(err).To(BeNil())
				Expect(assertion).To(BeFalse())
			})
		})

		When("doctor credential is invalid", func() {
			It("should return false", func() {
				assertion, err := graphQLClient.AssertDoctorCredential(context.Background(), "Roma40", "not-password")
				Expect(err).To(BeNil())
				Expect(assertion).To(BeFalse())
			})
		})

		When("doctor credential is valid", func() {
			It("should return true", func() {
				assertion, err := graphQLClient.AssertDoctorCredential(context.Background(), "Rickie_Ward29", "password")
				Expect(err).To(BeNil())
				Expect(assertion).To(BeTrue())
			})
		})
	})

	Context("FindDoctorByUsername", func() {
		When("doctor is not found", func() {
			It("should return nil with no error", func() {
				doctor, err := graphQLClient.FindDoctorByUsername(context.Background(), "awdasdwasdwad")
				Expect(err).To(BeNil())
				Expect(doctor).To(BeNil())
			})
		})

		When("doctor is found", func() {
			It("should return doctor", func() {
				doctor, err := graphQLClient.FindDoctorByUsername(context.Background(), "Jonathon_Kshlerin19")
				Expect(err).To(BeNil())
				Expect(doctor.Id).To(Equal("10"))
			})
		})
	})

	Context("FindInvoiceByID", func() {
		When("invoice not found", func() {
			It("should return nil with no error", func() {
				invoice, err := graphQLClient.FindInvoiceByID(context.Background(), int(rand.Int31()))
				Expect(err).To(BeNil())
				Expect(invoice).To(BeNil())
			})
		})
		When("invoice is found", func() {
			It("should invoice with no error", func() {
				invoice, err := graphQLClient.FindInvoiceByID(context.Background(), 1)
				Expect(err).To(BeNil())
				Expect(invoice.Id).To(Equal("1"))
				Expect(invoice.AppointmentID).To(Equal("9"))
				Expect(invoice.PatientID).To(Equal("HN-265555"))
			})
		})
	})
})
