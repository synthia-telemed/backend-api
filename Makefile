proto:
	protoc --go_out=./pkg/services/token/proto --go_opt=paths=source_relative \
        --go-grpc_out=./pkg/services/token/proto --go-grpc_opt=paths=source_relative \
        --proto_path=pkg/services/token/proto \
        --validate_out="lang=go:." \
        pkg/services/token/proto/token.proto

unit-test:
	ginkgo -r

mockgen:
	mockgen -source=pkg/services/token/proto/token_grpc.pb.go -destination=test/mock_token_service/mock_token_grpc.pb.go -package mock_token_service
	mockgen -source=pkg/cache/client.go -destination=test/mock_cache_client/mock_cache_client.go -package mock_cache_client
	mockgen -source=pkg/hospital/hospital.go -destination=test/mock_hospital_client/mock_hospital_client.go -package mock_hospital_client
	mockgen -source=pkg/sms/client.go -destination=test/mock_sms_client/mock_sms_client.go -package mock_sms_client
	mockgen -source=pkg/datastore/patient.go -destination=test/mock_datastore/mock_patient_datastore.go -package mock_datastore

gql-client-gen:
	genqlient ./pkg/hospital/genqlient.yaml