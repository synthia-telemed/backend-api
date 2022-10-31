proto:
	protoc --go_out=./pkg/token/proto --go_opt=paths=source_relative \
        --go-grpc_out=./pkg/token/proto --go-grpc_opt=paths=source_relative \
        --proto_path=pkg/token/proto \
        --validate_out="lang=go:." \
        pkg/token/proto/token.proto

unit-test:
	ginkgo -r

unit-test-local:
	set -a allexport; source ".env.test"; ginkgo -r; set +a allexport

mockgen:
	mockgen -source=pkg/token/proto/token_grpc.pb.go -destination=test/mock_token_service/mock_token_grpc.pb.go -package mock_token_service
	mockgen -source=pkg/token/grpc.go -destination=test/mock_token_service/mock_token_grpc.go -package mock_token_service
	mockgen -source=pkg/cache/client.go -destination=test/mock_cache_client/mock_cache_client.go -package mock_cache_client
	mockgen -source=pkg/hospital/hospital.go -destination=test/mock_hospital_client/mock_hospital_client.go -package mock_hospital_client
	mockgen -source=pkg/sms/client.go -destination=test/mock_sms_client/mock_sms_client.go -package mock_sms_client
	mockgen -source=pkg/payment/client.go -destination=test/mock_payment/mock_payment.go -package mock_payment
	mockgen -source=pkg/clock/clock.go -destination=test/mock_clock/mock_clock.go -package mock_clock
	mockgen -source=pkg/id/nanoid.go -destination=test/mock_id/mock_id.go -package mock_id
	mockgen -source=pkg/notification/client.go -destination=test/mock_notification/mock_notification.go -package mock_notification
	mockgen -source=pkg/datastore/patient.go -destination=test/mock_datastore/mock_patient_datastore.go -package mock_datastore
	mockgen -source=pkg/datastore/doctor.go -destination=test/mock_datastore/mock_doctor_datastore.go -package mock_datastore
	mockgen -source=pkg/datastore/credit_card.go -destination=test/mock_datastore/mock_credit_card.go -package mock_datastore
	mockgen -source=pkg/datastore/payment.go -destination=test/mock_datastore/mock_payment.go -package mock_datastore
	mockgen -source=pkg/datastore/appointment.go -destination=test/mock_datastore/mock_appointment.go -package mock_datastore
	mockgen -source=pkg/datastore/notification.go -destination=test/mock_datastore/mock_notification.go -package mock_datastore

gql-client-gen:
	genqlient ./pkg/hospital/genqlient.yaml
	fieldalignment -fix ./pkg/hospital/

swagger:
	swag init --parseDependency --parseInternal --dir cmd/patient-api --output cmd/patient-api/docs
	swag init --parseDependency --parseInternal --dir cmd/doctor-api --output cmd/doctor-api/docs
