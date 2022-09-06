package cache_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

var (
	redisContainer Redis
)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache Suite")
}

type Redis struct {
	cache.Config
	terminate func(ctx context.Context) error
}

var _ = BeforeSuite(func() {
	redisContainer = setupRedisContainer()
})

var _ = AfterSuite(func() {
	Expect(redisContainer.terminate(context.Background())).To(Succeed())
})

func setupRedisContainer() Redis {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:6-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	Expect(err).To(BeNil())
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "6379")
	return Redis{
		Config: cache.Config{
			Endpoint: fmt.Sprintf("%s:%s", host, port.Port()),
		},
		terminate: container.Terminate,
	}
}
