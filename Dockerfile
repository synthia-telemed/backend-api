FROM golang:1.18-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o patient-api cmd/patient-api/main.go

FROM alpine:3
WORKDIR /app
COPY --from=builder /app/patient-api ./patient-api
ENTRYPOINT ["/app/patient-api"]