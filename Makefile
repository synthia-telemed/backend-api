proto:
	protoc --go_out=./pkg/services/token/proto --go_opt=paths=source_relative \
        --go-grpc_out=./pkg/services/token/proto --go-grpc_opt=paths=source_relative \
        --proto_path=pkg/services/token/proto \
        --validate_out="lang=go:." \
        pkg/services/token/proto/token.proto