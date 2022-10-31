package hospital

import (
	"context"
	"errors"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"sort"
	"strconv"
	"time"
)

type SystemClient interface {
	FindPatientByID(ctx context.Context, id string) (*Patient, error)
	FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error)
	AssertDoctorCredential(ctx context.Context, username, password string) (bool, error)
	FindDoctorByUsername(ctx context.Context, username string) (*Doctor, error)
	FindInvoiceByID(ctx context.Context, id int) (*InvoiceOverview, error)
	PaidInvoice(ctx context.Context, id int) error
	ListAppointmentsByPatientID(ctx context.Context, patientID string, since time.Time) ([]*AppointmentOverview, error)
	ListAppointmentsByDoctorID(ctx context.Context, doctorID string, date time.Time) ([]*AppointmentOverview, error)
	ListAppointmentsWithFilters(ctx context.Context, filters *ListAppointmentsFilters, take, skip int) ([]*AppointmentOverview, error)
	CountAppointmentsWithFilters(ctx context.Context, filters *ListAppointmentsFilters) (int, error)
	FindAppointmentByID(ctx context.Context, appointmentID int) (*Appointment, error)
	FindDoctorAppointmentByID(ctx context.Context, appointmentID int) (*DoctorAppointment, error)
	SetAppointmentStatus(ctx context.Context, appointmentID int, status SettableAppointmentStatus) error
	CategorizeAppointmentByStatus(apps []*AppointmentOverview) *CategorizedAppointment
}
type Config struct {
	HospitalSysEndpoint string `env:"HOSPITAL_SYS_ENDPOINT,required"`
}

type GraphQLClient struct {
	client graphql.Client
}

func NewGraphQLClient(config *Config) *GraphQLClient {
	return &GraphQLClient{
		client: graphql.NewClient(config.HospitalSysEndpoint, nil),
	}
}

func (c GraphQLClient) FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error) {
	where := &PatientWhereInput{
		OR: []*PatientWhereInput{
			{NationalId: &StringNullableFilter{Equals: &cred}},
			{PassportId: &StringNullableFilter{Equals: &cred}},
		},
	}
	return c.getAndParsePatient(ctx, where)
}

func (c GraphQLClient) FindPatientByID(ctx context.Context, id string) (*Patient, error) {
	where := &PatientWhereInput{Id: &StringFilter{Equals: &id}}
	return c.getAndParsePatient(ctx, where)
}

func (c GraphQLClient) getAndParsePatient(ctx context.Context, where *PatientWhereInput) (*Patient, error) {
	resp, err := getPatient(ctx, c.client, where)
	if err != nil || resp.GetPatient() == nil {
		return nil, err
	}

	p := resp.GetPatient()
	return &Patient{
		Id:            p.Id,
		NameEN:        NewName(p.Initial_en, p.Firstname_en, p.Lastname_en),
		NameTH:        NewName(p.Initial_th, p.Firstname_th, p.Lastname_th),
		ProfilePicURL: p.ProfilePicURL,
		BirthDate:     p.BirthDate,
		BloodType:     p.BloodType,
		Height:        p.Height,
		Weight:        p.Weight,
		NationalId:    p.NationalId,
		Nationality:   p.Nationality,
		PassportId:    p.PassportId,
		PhoneNumber:   p.PhoneNumber,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, nil
}

func (c GraphQLClient) AssertDoctorCredential(ctx context.Context, username, password string) (bool, error) {
	resp, err := assertDoctorCredential(ctx, c.client, password, username)
	if err != nil {
		return false, err
	}
	return resp.AssertDoctorPassword, nil
}

func (c GraphQLClient) FindDoctorByUsername(ctx context.Context, username string) (*Doctor, error) {
	resp, err := getDoctor(ctx, c.client, &DoctorWhereInput{Username: &StringFilter{Equals: &username}})
	if err != nil || resp.GetDoctor() == nil {
		return nil, err
	}
	d := resp.GetDoctor()
	return &Doctor{
		Id:            d.Id,
		NameEN:        NewName(d.Initial_en, d.Firstname_en, d.Lastname_en),
		NameTH:        NewName(d.Initial_th, d.Firstname_th, d.Lastname_th),
		Username:      d.Username,
		Position:      d.Position,
		ProfilePicURL: d.ProfilePicURL,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}, nil
}

func (c GraphQLClient) FindInvoiceByID(ctx context.Context, id int) (*InvoiceOverview, error) {
	resp, err := getInvoice(ctx, c.client, &InvoiceWhereInput{Id: &IntFilter{Equals: &id}})
	if err != nil || resp.Invoice == nil {
		return nil, err
	}
	invoiceID, err := strconv.ParseInt(resp.Invoice.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	var discount float64
	for _, dis := range resp.Invoice.InvoiceDiscount {
		discount += dis.GetAmount()
	}
	return &InvoiceOverview{
		CreatedAt:     resp.Invoice.CreatedAt,
		Id:            int(invoiceID),
		Paid:          resp.Invoice.Paid,
		Total:         resp.Invoice.Total - discount,
		AppointmentID: resp.Invoice.Appointment.Id,
		PatientID:     resp.Invoice.Appointment.PatientId,
	}, nil
}

func (c GraphQLClient) PaidInvoice(ctx context.Context, id int) error {
	_, err := paidInvoice(ctx, c.client, float64(id))
	return err
}

func (c GraphQLClient) ListAppointmentsByPatientID(ctx context.Context, patientID string, since time.Time) ([]*AppointmentOverview, error) {
	desc := SortOrderDesc
	resp, err := getAppointments(ctx, c.client, &AppointmentWhereInput{
		PatientId:     &StringFilter{Equals: &patientID},
		StartDateTime: &DateTimeFilter{Gte: &since},
	}, []*AppointmentOrderByWithRelationInput{
		{StartDateTime: &desc},
	})
	if err != nil {
		return nil, err
	}
	return c.parseHospitalAppointmentToAppointmentOverview(resp.Appointments), nil
}

func (c GraphQLClient) ListAppointmentsByDoctorID(ctx context.Context, doctorID string, date time.Time) ([]*AppointmentOverview, error) {
	desc := SortOrderDesc
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := startTime.Add(time.Hour * 24)

	doctorIDInt64, err := strconv.ParseInt(doctorID, 10, 32)
	if err != nil {
		return nil, err
	}
	doctorIDInt := int(doctorIDInt64)

	resp, err := getAppointments(ctx, c.client, &AppointmentWhereInput{
		DoctorId:      &IntFilter{Equals: &doctorIDInt},
		StartDateTime: &DateTimeFilter{Gte: &startTime, Lt: &endTime},
	}, []*AppointmentOrderByWithRelationInput{
		{StartDateTime: &desc},
	})
	if err != nil {
		return nil, err
	}
	return c.parseHospitalAppointmentToAppointmentOverview(resp.Appointments), nil
}

type ListAppointmentsFilters struct {
	Text      *string           `json:"text" form:"text"`
	StartDate *time.Time        `json:"start_date" form:"start_date"`
	EndDate   *time.Time        `json:"end_date" form:"end_date"`
	DoctorID  *string           `swaggerignore:"true"`
	PatientID *string           `swaggerignore:"true"`
	Status    AppointmentStatus `json:"status" form:"status" binding:"required,enum" enums:"CANCELLED,COMPLETED,SCHEDULED"`
}

func (c GraphQLClient) ListAppointmentsWithFilters(ctx context.Context, filters *ListAppointmentsFilters, take, skip int) ([]*AppointmentOverview, error) {
	where, err := c.parseListAppointmentsFiltersToAppointmentWhereInput(filters)
	if err != nil {
		return nil, err
	}
	order := SortOrderDesc
	if filters.Status == AppointmentStatusScheduled {
		order = SortOrderAsc
	}
	resp, err := getAppointmentsWithPagination(ctx, c.client, where, []*AppointmentOrderByWithRelationInput{{StartDateTime: &order}}, &take, &skip)
	if err != nil {
		return nil, err
	}
	return c.parseHospitalAppointmentWithPaginationToAppointmentOverview(resp.Appointments), nil
}

func (c GraphQLClient) CountAppointmentsWithFilters(ctx context.Context, filters *ListAppointmentsFilters) (int, error) {
	where, err := c.parseListAppointmentsFiltersToAppointmentWhereInput(filters)
	if err != nil {
		return 0, err
	}
	resp, err := getAppointments(ctx, c.client, where, nil)
	if err != nil {
		return 0, err
	}
	return len(resp.Appointments), nil
}

func (c GraphQLClient) parseListAppointmentsFiltersToAppointmentWhereInput(filters *ListAppointmentsFilters) (*AppointmentWhereInput, error) {
	where := &AppointmentWhereInput{Status: &EnumAppointmentStatusFilter{Equals: &filters.Status}}
	if filters.PatientID != nil {
		where.PatientId = &StringFilter{Equals: filters.PatientID}
	} else if filters.DoctorID != nil {
		doctorIDInt64, err := strconv.ParseInt(*filters.DoctorID, 10, 32)
		if err != nil {
			return nil, err
		}
		doctorIDInt := int(doctorIDInt64)
		where.DoctorId = &IntFilter{Equals: &doctorIDInt}
	} else {
		return nil, errors.New("neither PatientID nor DoctorID is supplied")
	}
	if filters.Text != nil {
		where.OR = []*AppointmentWhereInput{
			{PatientId: &StringFilter{Contains: filters.Text}},
			{Patient: &PatientRelationFilter{Is: &PatientWhereInput{
				OR: []*PatientWhereInput{
					{Firstname_en: &StringFilter{Contains: filters.Text}},
					{Lastname_en: &StringFilter{Contains: filters.Text}},
				},
			}}},
		}
	}
	if filters.StartDate != nil && filters.EndDate != nil {
		st := filters.StartDate
		et := filters.EndDate
		startDateTime := time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
		endDateTime := time.Date(et.Year(), et.Month(), et.Day(), 23, 59, 59, 0, et.Location())
		where.StartDateTime = &DateTimeFilter{Gte: &startDateTime, Lt: &endDateTime}
	}
	return where, nil
}

func (c GraphQLClient) parseHospitalAppointmentWithPaginationToAppointmentOverview(hosApps []*getAppointmentsWithPaginationAppointmentsAppointment) []*AppointmentOverview {
	appointments := make([]*AppointmentOverview, len(hosApps))
	for i, a := range hosApps {
		appointments[i] = &AppointmentOverview{
			Id:            a.Id,
			StartDateTime: a.StartDateTime,
			EndDateTime:   a.EndDateTime,
			Status:        a.Status,
			Detail:        a.Detail,
			Doctor: DoctorOverview{
				ID:            a.Doctor.Id,
				FullName:      parseFullName(a.Doctor.Initial_en, a.Doctor.Firstname_en, a.Doctor.Lastname_en),
				Position:      a.Doctor.Position,
				ProfilePicURL: a.Doctor.ProfilePicURL,
			},
			Patient: PatientOverview{
				ID:            a.Patient.Id,
				FullName:      parseFullName(a.Patient.Initial_en, a.Patient.Firstname_en, a.Patient.Lastname_en),
				ProfilePicURL: a.Patient.ProfilePicURL,
			},
		}
	}
	return appointments
}

func (c GraphQLClient) parseHospitalAppointmentToAppointmentOverview(hosApps []*getAppointmentsAppointmentsAppointment) []*AppointmentOverview {
	appointments := make([]*AppointmentOverview, len(hosApps))
	for i, a := range hosApps {
		appointments[i] = &AppointmentOverview{
			Id:            a.Id,
			StartDateTime: a.StartDateTime,
			EndDateTime:   a.EndDateTime,
			Status:        a.Status,
			Detail:        a.Detail,
			Doctor: DoctorOverview{
				ID:            a.Doctor.Id,
				FullName:      parseFullName(a.Doctor.Initial_en, a.Doctor.Firstname_en, a.Doctor.Lastname_en),
				Position:      a.Doctor.Position,
				ProfilePicURL: a.Doctor.ProfilePicURL,
			},
			Patient: PatientOverview{
				ID:            a.Patient.Id,
				FullName:      parseFullName(a.Patient.Initial_en, a.Patient.Firstname_en, a.Patient.Lastname_en),
				ProfilePicURL: a.Patient.ProfilePicURL,
			},
		}
	}
	return appointments
}

func (c GraphQLClient) FindDoctorAppointmentByID(ctx context.Context, appointmentID int) (*DoctorAppointment, error) {
	resp, err := getDoctorAppointment(ctx, c.client, &AppointmentWhereInput{Id: &IntFilter{Equals: &appointmentID}})
	app := resp.GetAppointment()
	if err != nil || app == nil {
		return nil, err
	}
	return &DoctorAppointment{
		Id: app.GetId(),
		Patient: DoctorAppointmentPatient{
			BirthDate:     app.Patient.GetBirthDate(),
			ID:            app.Patient.GetId(),
			FullName:      parseFullName(app.Patient.GetInitial_en(), app.Patient.GetFirstname_en(), app.Patient.GetLastname_en()),
			BloodType:     app.Patient.GetBloodType(),
			Weight:        app.Patient.GetWeight(),
			Height:        app.Patient.GetHeight(),
			ProfilePicURL: app.Patient.GetProfilePicURL(),
		},
		Doctor: DoctorOverview{
			ID:            app.Doctor.GetId(),
			FullName:      parseFullName(app.Doctor.GetInitial_en(), app.Doctor.GetFirstname_en(), app.Doctor.GetLastname_en()),
			Position:      app.Doctor.GetPosition(),
			ProfilePicURL: app.Doctor.GetProfilePicURL(),
		},
		Detail:          app.GetDetail(),
		StartDateTime:   app.GetStartDateTime(),
		EndDateTime:     app.GetEndDateTime(),
		NextAppointment: app.GetNextAppointment(),
		Status:          app.GetStatus(),
	}, nil
}

func (c GraphQLClient) FindAppointmentByID(ctx context.Context, appointmentID int) (*Appointment, error) {
	resp, err := getAppointment(ctx, c.client, &AppointmentWhereInput{
		Id: &IntFilter{Equals: &appointmentID},
	})
	if err != nil || resp.GetAppointment() == nil {
		return nil, err
	}
	appointment := &Appointment{
		Id:              resp.Appointment.GetId(),
		PatientID:       resp.Appointment.GetPatientId(),
		StartDateTime:   resp.Appointment.GetStartDateTime(),
		EndDateTime:     resp.Appointment.GetEndDateTime(),
		NextAppointment: resp.Appointment.GetNextAppointment(),
		Detail:          resp.Appointment.GetDetail(),
		Status:          resp.Appointment.GetStatus(),
		Doctor: DoctorOverview{
			ID:            resp.Appointment.Doctor.GetId(),
			FullName:      parseFullName(resp.Appointment.Doctor.GetInitial_en(), resp.Appointment.Doctor.GetFirstname_en(), resp.Appointment.Doctor.GetLastname_en()),
			Position:      resp.Appointment.Doctor.GetPosition(),
			ProfilePicURL: resp.Appointment.Doctor.GetProfilePicURL(),
		},
		Invoice:       nil,
		Prescriptions: make([]*Prescription, len(resp.Appointment.GetPrescriptions())),
	}
	in := resp.Appointment.Invoice
	if in != nil {
		id, _ := strconv.ParseInt(in.GetId(), 10, 32)
		appointment.Invoice = &Invoice{
			Id:               int(id),
			Total:            in.GetTotal(),
			Paid:             in.GetPaid(),
			InvoiceItems:     make([]*InvoiceItem, len(in.GetInvoiceItems())),
			InvoiceDiscounts: make([]*InvoiceDiscount, len(in.GetInvoiceDiscount())),
		}
		for i, it := range in.InvoiceItems {
			appointment.Invoice.InvoiceItems[i] = &InvoiceItem{
				Name:     it.GetName(),
				Price:    it.GetPrice(),
				Quantity: it.GetQuantity(),
			}
		}
		for i, dis := range in.InvoiceDiscount {
			appointment.Invoice.InvoiceDiscounts[i] = &InvoiceDiscount{
				Name:   dis.Name,
				Amount: dis.Amount,
			}
		}
	}
	pre := resp.Appointment.Prescriptions
	if len(pre) != 0 {
		for i, p := range pre {
			appointment.Prescriptions[i] = &Prescription{
				Amount:      p.GetAmount(),
				Name:        p.Medicine.GetName(),
				Description: p.Medicine.GetDescription(),
				PictureURL:  p.Medicine.GetPictureURL(),
			}
		}
	}
	return appointment, nil
}

func parseFullName(init, first, last string) string {
	return fmt.Sprintf("%s %s %s", init, first, last)
}

type SettableAppointmentStatus string

const (
	SettableAppointmentStatusCancelled SettableAppointmentStatus = "CANCELLED"
	SettableAppointmentStatusCompleted SettableAppointmentStatus = "COMPLETED"
)

func (s SettableAppointmentStatus) IsValid() bool {
	switch s {
	case SettableAppointmentStatusCompleted, SettableAppointmentStatusCancelled:
		return true
	default:
		return false
	}
}

func (c GraphQLClient) SetAppointmentStatus(ctx context.Context, appointmentID int, status SettableAppointmentStatus) error {
	var s AppointmentStatus
	switch status {
	case SettableAppointmentStatusCancelled:
		s = AppointmentStatusCancelled
	case SettableAppointmentStatusCompleted:
		s = AppointmentStatusCompleted
	}
	_, err := setAppointmentStatus(ctx, c.client, float64(appointmentID), s)
	return err
}

type CategorizedAppointment struct {
	Completed []*AppointmentOverview `json:"completed"`
	Scheduled []*AppointmentOverview `json:"scheduled"`
	Cancelled []*AppointmentOverview `json:"cancelled"`
}

func (c GraphQLClient) CategorizeAppointmentByStatus(apps []*AppointmentOverview) *CategorizedAppointment {
	res := CategorizedAppointment{
		Completed: make([]*AppointmentOverview, 0),
		Scheduled: make([]*AppointmentOverview, 0),
		Cancelled: make([]*AppointmentOverview, 0),
	}
	for _, a := range apps {
		switch a.Status {
		case AppointmentStatusCancelled:
			res.Cancelled = append(res.Cancelled, a)
		case AppointmentStatusCompleted:
			res.Completed = append(res.Completed, a)
		case AppointmentStatusScheduled:
			res.Scheduled = append(res.Scheduled, a)
		}
	}
	ReverseSlice(res.Scheduled)
	return &res
}

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func (s AppointmentStatus) IsValid() bool {
	switch s {
	case AppointmentStatusCancelled, AppointmentStatusScheduled, AppointmentStatusCompleted:
		return true
	default:
		return false
	}
}
