proto:
	protoc --go_out=./pkg/services/token --go_opt=paths=source_relative \
        --go-grpc_out=./pkg/services/token --go-grpc_opt=paths=source_relative \
        --proto_path=pkg/services/token \
        --validate_out="lang=go:." \
        pkg/services/token/token.proto
