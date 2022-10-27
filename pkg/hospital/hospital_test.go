package hospital_test

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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

	Context("FindPatientByID", func() {
		It("should found the patient", func() {
			patient, err := graphQLClient.FindPatientByID(ctx, "HN-846103")
			Expect(err).To(BeNil())
			Expect(patient.Id).To(Equal("HN-846103"))
			Expect(patient.NationalId).To(HaveValue(Equal("1108182787046")))
		})

		It("should return nil if not found", func() {
			patient, err := graphQLClient.FindPatientByID(ctx, uuid.NewString())
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
			It("should return invoice with no error", func() {
				invoice, err := graphQLClient.FindInvoiceByID(ctx, 1)
				Expect(err).To(BeNil())
				Expect(invoice.Id).To(Equal(1))
				Expect(invoice.AppointmentID).To(Equal("2"))
				Expect(invoice.PatientID).To(Equal("HN-414878"))
				Expect(invoice.Total).To(BeEquivalentTo(3922590))
			})
			It("should return invoice with discount", func() {
				invoice, err := graphQLClient.FindInvoiceByID(ctx, 2)
				Expect(err).To(BeNil())
				Expect(invoice.Id).To(Equal(2))
				Expect(invoice.Total).To(BeEquivalentTo(2748108 - 50000))
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

	Context("CountAppointmentsWithFilters", func() {
		var (
			filters *hospital.ListAppointmentsFilters
			count   int
			err     error
		)
		BeforeEach(func() {
			dID := "12"
			filters = &hospital.ListAppointmentsFilters{DoctorID: &dID, Status: hospital.AppointmentStatusCompleted}
		})
		JustBeforeEach(func() {
			count, err = graphQLClient.CountAppointmentsWithFilters(ctx, filters)
		})
		When("doctor has one or more appointments that matched the filter", func() {
			It("should return 2", func() {
				Expect(err).To(BeNil())
				Expect(count).To(Equal(2))
			})
		})
		When("doctor no appointment that matched the filter", func() {
			BeforeEach(func() {
				filters.Status = hospital.AppointmentStatusCancelled
			})
			It("should return 2", func() {
				Expect(err).To(BeNil())
				Expect(count).To(Equal(0))
			})
		})
	})

	Context("parseListAppointmentsFiltersToAppointmentWhereInput", func() {

	})

	Context("ListAppointmentsByDoctorIDWithFilters", func() {
		var (
			doctorID     = "12"
			filters      *hospital.ListAppointmentsFilters
			appointments []*hospital.AppointmentOverview
			err          error
			take, skip   int
		)
		BeforeEach(func() {
			take = 10
			skip = 0
			filters = &hospital.ListAppointmentsFilters{Status: hospital.AppointmentStatusCompleted, DoctorID: &doctorID}
		})
		JustBeforeEach(func() {
			appointments, err = graphQLClient.ListAppointmentsWithFilters(ctx, filters, take, skip)
		})

		When("take and skip is set", func() {
			BeforeEach(func() {
				take = 1
				skip = 1
			})
			It("should return one appointment with the first one skipped", func() {
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(1))
				Expect(appointments[0].Id).To(Equal("37"))
			})
		})

		Context("target patient or doctor", func() {
			When("PatientID and DoctorID are not set in filters", func() {
				BeforeEach(func() {
					filters = &hospital.ListAppointmentsFilters{Status: hospital.AppointmentStatusCompleted}
				})
				It("should return error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("neither PatientID nor DoctorID is supplied"))
				})
			})
			When("Only PatientID is set", func() {
				var patientID = "HN-285237"
				BeforeEach(func() {
					filters = &hospital.ListAppointmentsFilters{Status: hospital.AppointmentStatusCompleted, PatientID: &patientID}
				})
				It("should return the appointments of that patient", func() {
					Expect(err).To(BeNil())
					Expect(appointments).To(HaveLen(2))
					Expect(appointments[0].Id).To(Equal("50"))
					Expect(appointments[1].Id).To(Equal("98"))
				})
			})
			When("Only DoctorID is set", func() {
				It("should return the appointments of that patient", func() {
					Expect(err).To(BeNil())
					Expect(appointments).To(HaveLen(2))
					Expect(appointments[0].Id).To(Equal("38"))
					Expect(appointments[1].Id).To(Equal("37"))
				})
			})
			When("Both PatientID and DoctorID are set", func() {
				BeforeEach(func() {
					patientID := "HN-314696"
					doctorID := "29"
					filters = &hospital.ListAppointmentsFilters{Status: hospital.AppointmentStatusCancelled, PatientID: &patientID, DoctorID: &doctorID}
				})
				It("should find the appointment that belong to them", func() {
					Expect(err).To(BeNil())
					Expect(appointments).To(HaveLen(1))
					Expect(appointments[0].Id).To(Equal("25"))
				})
			})
		})

		Context("there is no text or date filters", func() {
			When("target status is scheduled", func() {
				BeforeEach(func() {
					filters.Status = hospital.AppointmentStatusScheduled
				})
				It("should return list of scheduled appointment in ascending order", func() {
					Expect(err).To(BeNil())
					Expect(appointments).To(HaveLen(2))
					testhelper.AssertListOfAppointments(appointments, hospital.AppointmentStatusScheduled, testhelper.ASC)
				})
			})
			When("target status is completed", func() {
				It("should return list of completed appointment in ascending order", func() {
					Expect(err).To(BeNil())
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
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(2))
				Expect(appointments[0].Id).To(Equal("38"))
				Expect(appointments[1].Id).To(Equal("37"))
			})
		})
		When("there is date range filter", func() {
			BeforeEach(func() {
				dID := "11"
				filters.DoctorID = &dID
				st := time.Date(2022, 4, 27, 0, 0, 0, 0, time.UTC)
				et := time.Date(2022, 9, 7, 0, 0, 0, 0, time.UTC)
				filters.StartDate = &st
				filters.EndDate = &et
			})
			It("should return appointment on 2022-09-07 UTC", func() {
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(2))
				Expect(appointments[0].Id).To(Equal("53"))
				Expect(appointments[1].Id).To(Equal("50"))
			})
		})
		When("there are date and text filter", func() {
			BeforeEach(func() {
				st := time.Date(2022, 9, 7, 10, 0, 0, 0, time.UTC)
				dt := st.Add(time.Hour * 24)
				text := "562380"
				filters.Text = &text
				filters.StartDate = &st
				filters.EndDate = &dt
			})
			It("should return appointment on 2022-09-07 UTC with patient number that contain 562380", func() {
				Expect(err).To(BeNil())
				Expect(appointments).To(HaveLen(1))
				Expect(appointments[0].Id).To(Equal("37"))
			})
		})
	})

	Context("FindDoctorAppointmentByID", func() {
		var (
			appointmentID int
			appointment   *hospital.DoctorAppointment
		)
		JustBeforeEach(func() {
			var err error
			appointment, err = graphQLClient.FindDoctorAppointmentByID(ctx, appointmentID)
			Expect(err).To(BeNil())
		})

		When("appointment is not found", func() {
			BeforeEach(func() {
				appointmentID = 1283923472
			})
			It("should return nil with no error", func() {
				Expect(appointment).To(BeNil())
			})
		})
		When("appointment is found", func() {
			BeforeEach(func() {
				appointmentID = 26
			})
			It("should return appointment with data", func() {
				Expect(appointment).ToNot(BeNil())
				Expect(appointment.Id).To(Equal(fmt.Sprintf("%d", appointmentID)))
			})
		})
	})
})
