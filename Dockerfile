FROM golang:1.18-alpine as builder
WORKDIR /app
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o patient-api cmd/patient-api/main.go

FROM alpine:3
WORKDIR /app
COPY ./ ./
COPY --from=builder /app/patient-api ./bin/patient-api
ENTRYPOINT ["/app/bin/patient-api"]