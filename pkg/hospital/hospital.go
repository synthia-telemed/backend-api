package hospital

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"strconv"
	"time"
)

type SystemClient interface {
	FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error)
	AssertDoctorCredential(ctx context.Context, username, password string) (bool, error)
	FindDoctorByUsername(ctx context.Context, username string) (*Doctor, error)
	FindInvoiceByID(ctx context.Context, id int) (*Invoice, error)
	PaidInvoice(ctx context.Context, id int) error
	ListAppointmentsByPatientID(ctx context.Context, patientID string) ([]*Appointment, error)
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

type Invoice struct {
	CreatedAt     time.Time
	Id            int
	Paid          bool
	Total         float64
	AppointmentID string
	PatientID     string
}

type Appointment struct {
	Id        string
	DateTime  time.Time
	PatientId string
	Status    AppointmentStatus
	Doctor    DoctorBasicInfo
}

type DoctorBasicInfo struct {
	Initial   string
	Firstname string
	Lastname  string
	Position  string
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

func (c GraphQLClient) FindInvoiceByID(ctx context.Context, id int) (*Invoice, error) {
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
	return &Invoice{
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

func (c GraphQLClient) ListAppointmentsByPatientID(ctx context.Context, patientID string) ([]*Appointment, error) {

	resp, err := getAppointments(ctx, c.client, &AppointmentWhereInput{
		PatientId: &StringFilter{Equals: patientID, Mode: QueryModeDefault},
	})
	if err != nil {
		return nil, err
	}
	appointments := make([]*Appointment, len(resp.Appointments))
	for i, a := range resp.Appointments {
		appointments[i] = &Appointment{
			Id:        a.Id,
			DateTime:  a.DateTime,
			PatientId: a.PatientId,
			Status:    a.Status,
			Doctor: DoctorBasicInfo{
				Initial:   a.Doctor.Initial_en,
				Firstname: a.Doctor.Firstname_en,
				Lastname:  a.Doctor.Lastname_en,
				Position:  a.Doctor.Position,
			},
		}
	}
	return appointments, nil
}
