FROM golang:1.18-alpine as patient-api-builder
WORKDIR /app
COPY ./DigiCertGlobalRootCA.crt.pem ./
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./pkg ./pkg
COPY ./cmd/patient-api ./cmd/patient-api
RUN go build -o patient-api cmd/patient-api/main.go

FROM golang:1.18-alpine as doctor-api-builder
WORKDIR /app
COPY ./DigiCertGlobalRootCA.crt.pem ./
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./pkg ./pkg
COPY ./cmd/doctor-api ./cmd/doctor-api
RUN go build -o doctor-api cmd/doctor-api/main.go

FROM alpine:3
WORKDIR /app
COPY ./ ./
COPY --from=patient-api-builder /app/patient-api ./bin/patient-api
COPY --from=doctor-api-builder /app/doctor-api ./bin/doctor-api
ENTRYPOINT ["/app/bin/patient-api"]