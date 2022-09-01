package hospital

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"strconv"
	"time"
)

type SystemClient interface {
	FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error)
	AssertDoctorCredential(ctx context.Context, username, password string) (bool, error)
	FindDoctorByUsername(ctx context.Context, username string) (*Doctor, error)
	FindInvoiceByID(ctx context.Context, id int) (*InvoiceOverview, error)
	PaidInvoice(ctx context.Context, id int) error
	ListAppointmentsByPatientID(ctx context.Context, patientID string) ([]*AppointmentOverview, error)
	FindAppointmentByID(ctx context.Context, appointmentID int) (*Appointment, error)
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

type Patient struct {
	BirthDate    time.Time
	BloodType    BloodType
	CreatedAt    time.Time
	Firstname_en string
	Firstname_th string
	Height       float64
	Id           string
	Initial_en   string
	Initial_th   string
	Lastname_en  string
	Lastname_th  string
	NationalId   string
	Nationality  string
	PassportId   string
	PhoneNumber  string
	UpdatedAt    time.Time
	Weight       float64
}

type Doctor struct {
	CreatedAt    time.Time
	Firstname_en string
	Firstname_th string
	Id           string
	Initial_en   string
	Initial_th   string
	Lastname_en  string
	Lastname_th  string
	Password     string
	Position     string
	UpdatedAt    time.Time
	Username     string
}

type InvoiceOverview struct {
	CreatedAt     time.Time
	Id            int
	Paid          bool
	Total         float64
	AppointmentID string
	PatientID     string
}

type AppointmentOverview struct {
	Id        string            `json:"id"`
	DateTime  time.Time         `json:"date_time"`
	PatientId string            `json:"patient_id"`
	Status    AppointmentStatus `json:"status"`
	Doctor    DoctorOverview    `json:"doctor"`
}
type DoctorOverview struct {
	FullName      string `json:"full_name"`
	Position      string `json:"position"`
	ProfilePicURL string `json:"profile_pic_url"`
}

type Appointment struct {
	Id              string            `json:"id"`
	PatientID       string            `json:"patient_id"`
	DateTime        time.Time         `json:"date_time"`
	NextAppointment time.Time         `json:"next_appointment"`
	Detail          string            `json:"detail"`
	Status          AppointmentStatus `json:"status"`
	Doctor          DoctorOverview    `json:"doctor"`
	Invoice         *Invoice          `json:"invoice"`
	Prescriptions   []*Prescription   `json:"prescriptions"`
}
type Invoice struct {
	Id           int            `json:"id"`
	Total        float64        `json:"total"`
	Paid         bool           `json:"paid"`
	InvoiceItems []*InvoiceItem `json:"invoice_items"`
}
type InvoiceItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
type Prescription struct {
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

func (c GraphQLClient) FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error) {
	resp, err := getPatient(ctx, c.client, &PatientWhereInput{
		OR: []*PatientWhereInput{
			{NationalId: &StringNullableFilter{Equals: cred, Mode: QueryModeDefault}},
			{PassportId: &StringNullableFilter{Equals: cred, Mode: QueryModeDefault}},
		},
	})

	if err != nil || resp.GetPatient() == nil {
		return nil, err
	}

	return (*Patient)(resp.Patient), nil
}

func (c GraphQLClient) AssertDoctorCredential(ctx context.Context, username, password string) (bool, error) {
	resp, err := assertDoctorCredential(ctx, c.client, password, username)
	if err != nil {
		return false, err
	}
	return resp.AssertDoctorPassword, nil
}

func (c GraphQLClient) FindDoctorByUsername(ctx context.Context, username string) (*Doctor, error) {
	resp, err := getDoctor(ctx, c.client, &DoctorWhereInput{Username: &StringFilter{Equals: username, Mode: QueryModeDefault}})
	if err != nil {
		return nil, err
	}
	return (*Doctor)(resp.Doctor), nil
}

func (c GraphQLClient) FindInvoiceByID(ctx context.Context, id int) (*InvoiceOverview, error) {
	resp, err := getInvoice(ctx, c.client, &InvoiceWhereInput{Id: &IntFilter{Equals: id}})
	if err != nil {
		return nil, err
	}
	if resp.Invoice == nil {
		return nil, nil
	}
	invoiceID, err := strconv.ParseInt(resp.Invoice.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	return &InvoiceOverview{
		CreatedAt:     resp.Invoice.CreatedAt,
		Id:            int(invoiceID),
		Paid:          resp.Invoice.Paid,
		Total:         resp.Invoice.Total,
		AppointmentID: resp.Invoice.Appointment.Id,
		PatientID:     resp.Invoice.Appointment.PatientId,
	}, nil
}

func (c GraphQLClient) PaidInvoice(ctx context.Context, id int) error {
	_, err := paidInvoice(ctx, c.client, float64(id))
	return err
}

func (c GraphQLClient) ListAppointmentsByPatientID(ctx context.Context, patientID string, since time.Time) ([]*AppointmentOverview, error) {
	resp, err := getAppointments(ctx, c.client, &AppointmentWhereInput{
		PatientId: &StringFilter{Equals: patientID, Mode: QueryModeDefault},
		DateTime:  &DateTimeFilter{Gte: since},
	})
	if err != nil {
		return nil, err
	}
	appointments := make([]*AppointmentOverview, len(resp.Appointments))
	for i, a := range resp.Appointments {
		appointments[i] = &AppointmentOverview{
			Id:        a.Id,
			DateTime:  a.DateTime,
			PatientId: a.PatientId,
			Status:    a.Status,
			Doctor: DoctorOverview{
				FullName:      parseFullName(a.Doctor.Initial_en, a.Doctor.Firstname_en, a.Doctor.Lastname_en),
				Position:      a.Doctor.Position,
				ProfilePicURL: a.Doctor.ProfilePicURL,
			},
		}
	}
	return appointments, nil
}

func (c GraphQLClient) FindAppointmentByID(ctx context.Context, appointmentID int) (*Appointment, error) {
	resp, err := getAppointment(ctx, c.client, &AppointmentWhereInput{
		Id: &IntFilter{Equals: appointmentID},
	})
	if err != nil {
		return nil, err
	}
	if resp.GetAppointment() == nil {
		return nil, nil
	}
	appointment := &Appointment{
		Id:              resp.Appointment.GetId(),
		PatientID:       resp.Appointment.GetPatientId(),
		DateTime:        resp.Appointment.GetDateTime(),
		NextAppointment: resp.Appointment.GetNextAppointment(),
		Detail:          resp.Appointment.GetDetail(),
		Status:          resp.Appointment.GetStatus(),
		Doctor: DoctorOverview{
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
			Id:           int(id),
			Total:        in.GetTotal(),
			Paid:         in.GetPaid(),
			InvoiceItems: make([]*InvoiceItem, len(in.GetInvoiceItems())),
		}
		for i, it := range in.InvoiceItems {
			appointment.Invoice.InvoiceItems[i] = &InvoiceItem{
				Name:     it.GetName(),
				Price:    it.GetPrice(),
				Quantity: it.GetQuantity(),
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
			}
		}
	}
	return appointment, nil
}

func parseFullName(init, first, last string) string {
	return fmt.Sprintf("%s %s %s", init, first, last)
}
