package hospital_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	"math/rand"
	"time"
)

var _ = Describe("Hospital Client", func() {

	var (
		mockCtrl      *gomock.Controller
		graphQLClient *hospital.GraphQLClient
	)

	BeforeEach(func() {
		c := hospital.Config{HospitalSysEndpoint: "http://localhost:30821/graphql"}
		mockCtrl = gomock.NewController(GinkgoT())
		graphQLClient = hospital.NewGraphQLClient(&c)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("FindPatientByGovCredential", func() {
		It("should find patient by passport ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "JN848321")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-162462"))
		})

		It("should find patient by national ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(context.Background(), "4671253551800")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-427845"))
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
				assertion, err := graphQLClient.AssertDoctorCredential(context.Background(), "Anthony23", "not-password")
				Expect(err).To(BeNil())
				Expect(assertion).To(BeFalse())
			})
		})

		When("doctor credential is valid", func() {
			It("should return true", func() {
				assertion, err := graphQLClient.AssertDoctorCredential(context.Background(), "Anthony23", "password")
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
				doctor, err := graphQLClient.FindDoctorByUsername(context.Background(), "Elias_Wolf")
				Expect(err).To(BeNil())
				Expect(doctor.Id).To(Equal("5"))
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
				Expect(invoice.Id).To(Equal(1))
				Expect(invoice.AppointmentID).To(Equal("2"))
				Expect(invoice.PatientID).To(Equal("HN-414878"))
			})
		})
	})

	Context("ListAppointmentsByPatientID", func() {
		When("no appointment is found", func() {
			It("should return empty slice with no error", func() {
				appointments, err := graphQLClient.ListAppointmentsByPatientID(context.Background(), "HN-something", time.Now())
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(0))
			})
		})
		When("appointment(s) is/are found", func() {
			It("should return scheduled appointments", func() {
				appointments, err := graphQLClient.ListAppointmentsByPatientID(context.Background(), "HN-129512", time.Unix(1659211832, 0))
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(1))
			})
			It("should return appointments from started of 2023", func() {
				appointments, err := graphQLClient.ListAppointmentsByPatientID(context.Background(), "HN-853857", time.Unix(1672506001, 0))
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(3))
			})
		})
	})

	Context("FindAppointmentByID", func() {
		When("appointment is not found", func() {
			It("should return nil with no error", func() {
				appointment, err := graphQLClient.FindAppointmentByID(context.Background(), int(rand.Int31()))
				Expect(err).To(BeNil())
				Expect(appointment).To(BeNil())
			})
		})
		When("appointment has no invoice and prescriptions", func() {
			It("should return appointment with nil on invoice and zero length prescriptions", func() {
				appointment, err := graphQLClient.FindAppointmentByID(context.Background(), 1)
				Expect(err).To(BeNil())
				Expect(appointment).ToNot(BeNil())
				Expect(appointment.Prescriptions).To(HaveLen(0))
				Expect(appointment.Invoice).To(BeNil())
			})
		})
		When("appointment has invoice and prescriptions", func() {
			It("should return appointment with invoice and non zero length prescriptions", func() {
				appointment, err := graphQLClient.FindAppointmentByID(context.Background(), 2)
				Expect(err).To(BeNil())
				Expect(appointment).ToNot(BeNil())
				Expect(appointment.Prescriptions).ToNot(HaveLen(0))
				Expect(appointment.Invoice).ToNot(BeNil())
			})
		})
	})
})
