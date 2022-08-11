package hospital

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"net/http"
)

type SystemClient interface {
	FindPatientByNationalID() error
}

type GraphQLClient struct {
	client graphql.Client
}

func NewGraphQLClient(endpoint string) *GraphQLClient {
	return &GraphQLClient{
		client: graphql.NewClient(endpoint, http.DefaultClient),
	}
}

func (c GraphQLClient) FindPatientByNationalID(nationalID string) (*datastore.Patient, error) {
	resp, err := getPatient(context.Background(), c.client, PatientWhereInput{
		NationalId: StringNullableFilter{Equals: nationalID},
	})
	if err != nil || resp == nil {
		return nil, err
	}

	return &datastore.Patient{
		RefID:       resp.Patient.Id,
		BirthDate:   resp.Patient.BirthDate,
		BloodType:   datastore.BloodType(resp.Patient.BloodType),
		FirstnameEn: resp.Patient.Firstname_en,
		FirstnameTh: resp.Patient.Firstname_th,
		InitialEn:   resp.Patient.Initial_en,
		InitialTh:   resp.Patient.Initial_th,
		LastnameEn:  resp.Patient.Lastname_en,
		LastnameTh:  resp.Patient.Lastname_th,
		NationalID:  resp.Patient.NationalId,
		Nationality: resp.Patient.Nationality,
		PhoneNumber: resp.Patient.PhoneNumber,
		Weight:      float32(resp.Patient.Weight),
		Height:      float32(resp.Patient.Height),
	}, nil
}
