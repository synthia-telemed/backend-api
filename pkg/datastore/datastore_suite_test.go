package datastore_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestDatastore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datastore Suite")
}

type terminateFunc func()
type PostgresDB struct {
	datastore.Config
	terminate terminateFunc
}

var (
	postgres PostgresDB
)

var _ = BeforeSuite(func() {
	postgres = setupPostgresDBContainer()
})

var _ = AfterSuite(func() {
	postgres.terminate()
})

func setupPostgresDBContainer() PostgresDB {
	config := datastore.Config{
		User:     "postgres",
		Password: "postgres",
		Name:     "synthia",
		SSLMode:  "disable",
		TimeZone: "Asia/Bangkok",
	}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     config.User,
			"POSTGRES_PASSWORD": config.Password,
			"POSTGRES_DB":       config.Name,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	Expect(err).To(BeNil())
	config.Host, _ = container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	config.Port = port.Int()

	terminateFunc := func() {
		Expect(container.Terminate(ctx)).To(Succeed())
	}

	return PostgresDB{
		Config:    config,
		terminate: terminateFunc,
	}
}
