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
	FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error)
	AssertDoctorCredential(ctx context.Context, username, password string) (bool, error)
	FindDoctorByUsername(ctx context.Context, username string) (*Doctor, error)
	FindInvoiceByID(ctx context.Context, id int) (*InvoiceOverview, error)
	PaidInvoice(ctx context.Context, id int) error
	ListAppointmentsByPatientID(ctx context.Context, patientID string, since time.Time) ([]*AppointmentOverview, error)
	ListAppointmentsByDoctorID(ctx context.Context, doctorID string, date time.Time) ([]*AppointmentOverview, error)
	ListAppointmentsWithFilters(ctx context.Context, filters *ListAppointmentsFilters) ([]*AppointmentOverview, error)
	FindAppointmentByID(ctx context.Context, appointmentID int) (*Appointment, error)
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

type Name struct {
	FullName  string
	Initial   string
	Firstname string
	Lastname  string
}

func NewName(init, first, last string) *Name {
	return &Name{
		FullName:  parseFullName(init, first, last),
		Initial:   init,
		Firstname: first,
		Lastname:  last,
	}
}

type Patient struct {
	BirthDate   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PassportId  *string
	NameEN      *Name
	NameTH      *Name
	NationalId  *string
	Id          string
	Nationality string
	PhoneNumber string
	BloodType   BloodType
	Height      float64
	Weight      float64
}

type Doctor struct {
	Id            string
	NameEN        *Name
	NameTH        *Name
	Position      string
	ProfilePicURL string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Username      string
}

type InvoiceOverview struct {
	CreatedAt     time.Time
	AppointmentID string
	PatientID     string
	Id            int
	Total         float64
	Paid          bool
}

type AppointmentOverview struct {
	Id            string            `json:"id"`
	Detail        string            `json:"detail"`
	StartDateTime time.Time         `json:"start_date_time"`
	EndDateTime   time.Time         `json:"end_date_time"`
	Status        AppointmentStatus `json:"status"`
	Doctor        DoctorOverview    `json:"doctor"`
	Patient       PatientOverview   `json:"patient"`
}
type DoctorOverview struct {
	ID            string `json:"id"`
	FullName      string `json:"full_name"`
	Position      string `json:"position"`
	ProfilePicURL string `json:"profile_pic_url"`
}
type PatientOverview struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type Appointment struct {
	Id              string            `json:"id"`
	PatientID       string            `json:"patient_id"`
	StartDateTime   time.Time         `json:"start_date_time"`
	EndDateTime     time.Time         `json:"end_date_time"`
	NextAppointment *time.Time        `json:"next_appointment"`
	Detail          string            `json:"detail"`
	Status          AppointmentStatus `json:"status"`
	Doctor          DoctorOverview    `json:"doctor"`
	Invoice         *Invoice          `json:"invoice"`
	Prescriptions   []*Prescription   `json:"prescriptions"`
}
type Invoice struct {
	InvoiceItems     []*InvoiceItem     `json:"invoice_items"`
	InvoiceDiscounts []*InvoiceDiscount `json:"invoice_discounts"`
	Id               int                `json:"id"`
	Total            float64            `json:"total"`
	Paid             bool               `json:"paid"`
}
type InvoiceItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
type InvoiceDiscount struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}
type Prescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PictureURL  string `json:"picture_url"`
	Amount      int    `json:"amount"`
}

func (c GraphQLClient) FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error) {
	resp, err := getPatient(ctx, c.client, &PatientWhereInput{
		OR: []*PatientWhereInput{
			{NationalId: &StringNullableFilter{Equals: &cred}},
			{PassportId: &StringNullableFilter{Equals: &cred}},
		},
	})

	if err != nil || resp.GetPatient() == nil {
		return nil, err
	}

	p := resp.GetPatient()
	return &Patient{
		Id:          p.Id,
		NameEN:      NewName(p.Initial_en, p.Firstname_en, p.Lastname_en),
		NameTH:      NewName(p.Initial_th, p.Firstname_th, p.Lastname_th),
		BirthDate:   p.BirthDate,
		BloodType:   p.BloodType,
		Height:      p.Height,
		Weight:      p.Weight,
		NationalId:  p.NationalId,
		Nationality: p.Nationality,
		PassportId:  p.PassportId,
		PhoneNumber: p.PhoneNumber,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
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
	Text      *string    `json:"text"`
	Date      *time.Time `json:"date"`
	DoctorID  *string
	PatientID *string
	Status    AppointmentStatus `json:"status" binding:"required,enum" enums:"CANCELLED,COMPLETED,SCHEDULED"`
}

func (c GraphQLClient) ListAppointmentsWithFilters(ctx context.Context, filters *ListAppointmentsFilters) ([]*AppointmentOverview, error) {
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
	if filters.Date != nil {
		startDateTime := time.Date(filters.Date.Year(), filters.Date.Month(), filters.Date.Day(), 0, 0, 0, 0, filters.Date.Location())
		endDateTime := startDateTime.Add(24 * time.Hour)
		where.StartDateTime = &DateTimeFilter{Gte: &startDateTime, Lt: &endDateTime}
	}

	order := SortOrderDesc
	if filters.Status == AppointmentStatusScheduled {
		order = SortOrderAsc
	}
	resp, err := getAppointments(ctx, c.client, where, []*AppointmentOrderByWithRelationInput{{StartDateTime: &order}})
	if err != nil {
		return nil, err
	}
	return c.parseHospitalAppointmentToAppointmentOverview(resp.Appointments), nil
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
				ID:       a.Patient.Id,
				FullName: parseFullName(a.Patient.Initial_en, a.Patient.Firstname_en, a.Patient.Lastname_en),
			},
		}
	}
	return appointments
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
