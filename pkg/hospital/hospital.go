package hospital

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"time"
)

type SystemClient interface {
	FindPatientByGovCredential(ctx context.Context, cred string) (*Patient, error)
	AssertDoctorCredential(ctx context.Context, username, password string) (bool, error)
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
	BirthDate    time.Time `json:"birthDate"`
	BloodType    BloodType `json:"bloodType"`
	CreatedAt    time.Time `json:"createdAt"`
	Firstname_en string    `json:"firstname_en"`
	Firstname_th string    `json:"firstname_th"`
	Height       float64   `json:"height"`
	Id           string    `json:"id"`
	Initial_en   string    `json:"initial_en"`
	Initial_th   string    `json:"initial_th"`
	Lastname_en  string    `json:"lastname_en"`
	Lastname_th  string    `json:"lastname_th"`
	NationalId   string    `json:"nationalId"`
	Nationality  string    `json:"nationality"`
	PassportId   string    `json:"passportId"`
	PhoneNumber  string    `json:"phoneNumber"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Weight       float64   `json:"weight"`
}

type Doctor struct {
	CreatedAt    time.Time `json:"createdAt"`
	Firstname_en string    `json:"firstname_en"`
	Firstname_th string    `json:"firstname_th"`
	Id           string    `json:"id"`
	Initial_en   string    `json:"initial_en"`
	Initial_th   string    `json:"initial_th"`
	Lastname_en  string    `json:"lastname_en"`
	Lastname_th  string    `json:"lastname_th"`
	Password     string    `json:"password"`
	Position     string    `json:"position"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Username     string    `json:"username"`
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
	resp, err := getDoctor(ctx, c.client, &DoctorWhereInput{Username: &StringFilter{Equals: username}})
	if err != nil {
		return nil, err
	}
	return (*Doctor)(resp.Doctor), nil
}
