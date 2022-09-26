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
		ctx           context.Context
	)

	BeforeEach(func() {
		c := hospital.Config{HospitalSysEndpoint: "http://localhost:30821/graphql"}
		mockCtrl = gomock.NewController(GinkgoT())
		graphQLClient = hospital.NewGraphQLClient(&c)
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("FindPatientByGovCredential", func() {
		It("should find patient by passport ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(ctx, "JN848321")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-162462"))
		})

		It("should find patient by national ID", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(ctx, "4671253551800")
			Expect(err).To(BeNil())
			Expect(patient).ToNot(BeNil())
			Expect(patient.Id).To(Equal("HN-427845"))
		})

		It("should return nil when patient not found", func() {
			patient, err := graphQLClient.FindPatientByGovCredential(ctx, "not-exist-national-id")
			Expect(err).To(BeNil())
			Expect(patient).To(BeNil())
		})
	})

	Context("AssertDoctorCredential", func() {
		When("doctor's username is not found", func() {
			It("should return false", func() {
				assertion, err := graphQLClient.AssertDoctorCredential(ctx, "not-exist-doctor", "password")
				Expect(err).To(BeNil())
				Expect(assertion).To(BeFalse())
			})
		})

		When("doctor credential is invalid", func() {
			It("should return false", func() {
				assertion, err := graphQLClient.AssertDoctorCredential(ctx, "Anthony23", "not-password")
				Expect(err).To(BeNil())
				Expect(assertion).To(BeFalse())
			})
		})

		When("doctor credential is valid", func() {
			It("should return true", func() {
				assertion, err := graphQLClient.AssertDoctorCredential(ctx, "Anthony23", "password")
				Expect(err).To(BeNil())
				Expect(assertion).To(BeTrue())
			})
		})
	})

	Context("FindDoctorByUsername", func() {
		When("doctor is not found", func() {
			It("should return nil with no error", func() {
				doctor, err := graphQLClient.FindDoctorByUsername(ctx, "awdasdwasdwad")
				Expect(err).To(BeNil())
				Expect(doctor).To(BeNil())
			})
		})

		When("doctor is found", func() {
			It("should return doctor", func() {
				doctor, err := graphQLClient.FindDoctorByUsername(ctx, "Elias_Wolf")
				Expect(err).To(BeNil())
				Expect(doctor.Id).To(Equal("5"))
			})
		})
	})

	Context("FindInvoiceByID", func() {
		When("invoice not found", func() {
			It("should return nil with no error", func() {
				invoice, err := graphQLClient.FindInvoiceByID(ctx, int(rand.Int31()))
				Expect(err).To(BeNil())
				Expect(invoice).To(BeNil())
			})
		})
		When("invoice is found", func() {
			It("should invoice with no error", func() {
				invoice, err := graphQLClient.FindInvoiceByID(ctx, 1)
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
				appointments, err := graphQLClient.ListAppointmentsByPatientID(ctx, "HN-something", time.Now())
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(0))
			})
		})
		When("appointment(s) is/are found", func() {
			It("should return scheduled appointments", func() {
				appointments, err := graphQLClient.ListAppointmentsByPatientID(ctx, "HN-129512", time.Unix(1659211832, 0))
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(1))
			})
			It("should return appointments from started of 2023", func() {
				appointments, err := graphQLClient.ListAppointmentsByPatientID(ctx, "HN-853857", time.Unix(1672506001, 0))
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(3))
			})
		})
	})

	Context("ListAppointmentsByDoctorID", func() {
		When("no appointment is found", func() {
			It("should return empty slice with no error", func() {
				appointments, err := graphQLClient.ListAppointmentsByDoctorID(ctx, 24, time.Date(2022, 9, 9, 10, 3, 2, 0, time.UTC))
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(0))
			})
		})
		When("appointment are found", func() {
			It("should return appointments on that date", func() {
				appointments, err := graphQLClient.ListAppointmentsByDoctorID(ctx, 9, time.Date(2022, 9, 7, 13, 43, 0, 0, time.UTC))
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(3))
			})
		})
	})

	Context("FindAppointmentByID", func() {
		When("appointment is not found", func() {
			It("should return nil with no error", func() {
				appointment, err := graphQLClient.FindAppointmentByID(ctx, int(rand.Int31()))
				Expect(err).To(BeNil())
				Expect(appointment).To(BeNil())
			})
		})
		When("appointment has no invoice and prescriptions", func() {
			It("should return appointment with nil on invoice and zero length prescriptions", func() {
				appointment, err := graphQLClient.FindAppointmentByID(ctx, 1)
				Expect(err).To(BeNil())
				Expect(appointment).ToNot(BeNil())
				Expect(appointment.Prescriptions).To(HaveLen(0))
				Expect(appointment.Invoice).To(BeNil())
			})
		})
		When("appointment has invoice and prescriptions", func() {
			It("should return appointment with invoice and non zero length prescriptions", func() {
				appointment, err := graphQLClient.FindAppointmentByID(ctx, 2)
				Expect(err).To(BeNil())
				Expect(appointment).ToNot(BeNil())
				Expect(appointment.Prescriptions).To(HaveLen(4))
				Expect(appointment.Invoice).ToNot(BeNil())
				Expect(appointment.Invoice.InvoiceItems).To(HaveLen(10))
			})
		})
	})

	Context("PaidInvoice", func() {
		It("should set the paid status to true", func() {
			Expect(graphQLClient.PaidInvoice(ctx, 18)).To(Succeed())
			invoice, err := graphQLClient.FindInvoiceByID(ctx, 18)
			Expect(err).To(BeNil())
			Expect(invoice.Paid).To(BeTrue())
		})
	})

	DescribeTable("Set appointment status", func(appID int, status hospital.SettableAppointmentStatus, expectedStatus hospital.AppointmentStatus) {
		Expect(graphQLClient.SetAppointmentStatus(ctx, appID, status)).To(Succeed())
		appointment, err := graphQLClient.FindAppointmentByID(ctx, appID)
		Expect(err).To(BeNil())
		Expect(appointment.Status).To(Equal(expectedStatus))
	},
		Entry("set status of appointment to complete", 34, hospital.SettableAppointmentStatusCompleted, hospital.AppointmentStatusCompleted),
		Entry("set status of appointment to cancelled", 86, hospital.SettableAppointmentStatusCancelled, hospital.AppointmentStatusCancelled),
	)
})
