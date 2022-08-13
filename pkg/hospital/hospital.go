package hospital

import (
	"context"
	"github.com/Khan/genqlient/graphql"
)

type SystemClient interface {
	FindPatientByGovCredential(ctx context.Context, cred string) (*getPatientPatient, error)
}

type GraphQLClient struct {
	client graphql.Client
}

func NewGraphQLClient(endpoint string) *GraphQLClient {
	return &GraphQLClient{
		client: graphql.NewClient(endpoint, nil),
	}
}

func (c GraphQLClient) FindPatientByGovCredential(ctx context.Context, cred string) (*getPatientPatient, error) {
	resp, err := getPatient(ctx, c.client, &PatientWhereInput{
		OR: []*PatientWhereInput{
			{NationalId: &StringNullableFilter{Equals: cred, Mode: QueryModeDefault}},
			{PassportId: &StringNullableFilter{Equals: cred, Mode: QueryModeDefault}},
		},
	})

	if err != nil || resp.GetPatient() == nil {
		return nil, err
	}

	return resp.Patient, nil
}
