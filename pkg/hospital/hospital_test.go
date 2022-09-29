package hospital_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/hospital"
	testhelper "github.com/synthia-telemed/backend-api/test/helper"
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
				appointments, err := graphQLClient.ListAppointmentsByDoctorID(ctx, "24", time.Date(2022, 9, 9, 10, 3, 2, 0, time.UTC))
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(0))
			})
		})
		When("appointment are found", func() {
			It("should return appointments on that date", func() {
				appointments, err := graphQLClient.ListAppointmentsByDoctorID(ctx, "9", time.Date(2022, 9, 7, 13, 43, 0, 0, time.UTC))
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
		When("appointment has invoice without discount and prescriptions", func() {
			It("should return appointment with invoice and non zero length prescriptions", func() {
				appointment, err := graphQLClient.FindAppointmentByID(ctx, 2)
				Expect(err).To(BeNil())
				Expect(appointment).ToNot(BeNil())
				Expect(appointment.Prescriptions).To(HaveLen(4))
				Expect(appointment.Invoice).ToNot(BeNil())
				Expect(appointment.Invoice.InvoiceItems).To(HaveLen(10))
			})
		})
		When("appointment has invoice with discount and prescriptions", func() {
			It("should return appointment with invoice, discount and non zero length prescriptions", func() {
				appointment, err := graphQLClient.FindAppointmentByID(ctx, 4)
				Expect(err).To(BeNil())
				Expect(appointment).ToNot(BeNil())
				Expect(appointment.Prescriptions).To(HaveLen(1))
				Expect(appointment.Invoice).ToNot(BeNil())
				Expect(appointment.Invoice.InvoiceItems).To(HaveLen(9))
				Expect(appointment.Invoice.InvoiceDiscounts).To(HaveLen(1))
				Expect(appointment.Invoice.InvoiceDiscounts[0].Name).To(Equal("Social Security"))
				Expect(appointment.Invoice.InvoiceDiscounts[0].Amount).To(BeEquivalentTo(50000))
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

	Context("CategorizeAppointmentByStatus", func() {
		It("should categorized appointment by status", func() {
			categorized := hospital.CategorizedAppointment{
				Completed: testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusCompleted, 3),
				Scheduled: testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusScheduled, 2),
				Cancelled: testhelper.GenerateAppointmentOverviews(hospital.AppointmentStatusCancelled, 3),
			}

			appointments := make([]*hospital.AppointmentOverview, 0)
			appointments = append(appointments, categorized.Completed...)
			appointments = append(appointments, categorized.Scheduled...)
			appointments = append(appointments, categorized.Cancelled...)
			res := graphQLClient.CategorizeAppointmentByStatus(appointments)
			hospital.ReverseSlice(categorized.Scheduled)
			Expect(res.Completed).To(Equal(categorized.Completed))
			Expect(res.Scheduled).To(Equal(categorized.Scheduled))
			Expect(res.Cancelled).To(Equal(categorized.Cancelled))
		})
	})

	Context("ReverseSlice", func() {
		It("should reserve the order of elements", func() {
			s := []int{1, 2, 3, 4, 5}
			hospital.ReverseSlice(s)
			Expect(s).To(Equal([]int{5, 4, 3, 2, 1}))
		})
	})

	Context("ListAppointmentsByDoctorIDWithFilters", func() {
		var (
			doctorID     = "12"
			filters      *hospital.ListAppointmentsByDoctorIDFilters
			appointments []*hospital.AppointmentOverview
			err          error
		)
		BeforeEach(func() {
			filters = &hospital.ListAppointmentsByDoctorIDFilters{Status: hospital.AppointmentStatusCompleted}
		})
		JustBeforeEach(func() {
			appointments, err = graphQLClient.ListAppointmentsByDoctorIDWithFilters(ctx, doctorID, filters)
			Expect(err).To(BeNil())
		})

		Context("there is no text or date filters", func() {
			When("target status is scheduled", func() {
				BeforeEach(func() {
					filters.Status = hospital.AppointmentStatusScheduled
				})
				It("should return list of scheduled appointment in ascending order", func() {
					Expect(appointments).To(HaveLen(2))
					testhelper.AssertListOfAppointments(appointments, hospital.AppointmentStatusScheduled, testhelper.ASC)
				})
			})
			When("target status is completed", func() {
				It("should return list of completed appointment in ascending order", func() {
					Expect(appointments).To(HaveLen(2))
					testhelper.AssertListOfAppointments(appointments, hospital.AppointmentStatusCompleted, testhelper.DESC)
				})
			})
		})
		When("there is text filters", func() {
			BeforeEach(func() {
				text := "lar"
				filters.Text = &text
			})
			It("should return appointment that patient name has 'lar'", func() {
				Expect(appointments).To(HaveLen(2))
				Expect(appointments[0].Id).To(Equal("38"))
				Expect(appointments[1].Id).To(Equal("37"))
			})
		})
		When("there is date filter", func() {
			BeforeEach(func() {
				date := time.Date(2022, 9, 7, 10, 0, 0, 0, time.UTC)
				filters.Date = &date
			})
			It("should return appointment on 2022-09-07 UTC", func() {
				Expect(appointments).To(HaveLen(2))
				Expect(appointments[0].Id).To(Equal("38"))
				Expect(appointments[1].Id).To(Equal("37"))
			})
		})
		When("there are date and text filter", func() {
			BeforeEach(func() {
				date := time.Date(2022, 9, 7, 10, 0, 0, 0, time.UTC)
				text := "562380"
				filters.Text = &text
				filters.Date = &date
			})
			It("should return appointment on 2022-09-07 UTC with patient number that contain 562380", func() {
				Expect(appointments).To(HaveLen(1))
				Expect(appointments[0].Id).To(Equal("37"))
			})
		})
	})
})
