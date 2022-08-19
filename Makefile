proto:
	protoc --go_out=./pkg/token/proto --go_opt=paths=source_relative \
        --go-grpc_out=./pkg/token/proto --go-grpc_opt=paths=source_relative \
        --proto_path=pkg/token/proto \
        --validate_out="lang=go:." \
        pkg/token/proto/token.proto

unit-test:
	ginkgo -r

unit-test-local:

	docker compose -f docker-compose.test.yml up -d
	set -a allexport; source ".env.test"; ginkgo -r; set +a allexport
	docker compose -f docker-compose.test.yml down

mockgen:
	mockgen -source=pkg/token/proto/token_grpc.pb.go -destination=test/mock_token_service/mock_token_grpc.pb.go -package mock_token_service
	mockgen -source=pkg/token/grpc.go -destination=test/mock_token_service/mock_token_grpc.go -package mock_token_service
	mockgen -source=pkg/cache/client.go -destination=test/mock_cache_client/mock_cache_client.go -package mock_cache_client
	mockgen -source=pkg/hospital/hospital.go -destination=test/mock_hospital_client/mock_hospital_client.go -package mock_hospital_client
	mockgen -source=pkg/sms/client.go -destination=test/mock_sms_client/mock_sms_client.go -package mock_sms_client
	mockgen -source=pkg/payment/client.go -destination=test/mock_payment/mock_payment.go -package mock_payment
	mockgen -source=pkg/datastore/patient.go -destination=test/mock_datastore/mock_patient_datastore.go -package mock_datastore
	mockgen -source=pkg/datastore/measurement.go -destination=test/mock_datastore/mock_measurement_datastore.go -package mock_datastore

gql-client-gen:
	genqlient ./pkg/hospital/genqlient.yaml

swagger:
	swag init --dir cmd/patient-api --parseDependency --parseInternal
