package cache_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"github.com/synthia-telemed/backend-api/test/container"
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
	container.Terminate
	cache.Config
}

var _ = BeforeSuite(func() {
	redisContainer = setupRedisContainer()
})

var _ = AfterSuite(func() {
	Expect(redisContainer.Terminate()).To(Succeed())
})

func setupRedisContainer() Redis {
	req := testcontainers.ContainerRequest{
		Image:        "redis:6-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	testCon, err := container.NewTestContainer(req, "6379")
	Expect(err).To(BeNil())
	return Redis{
		Config: cache.Config{
			Endpoint: fmt.Sprintf("%s:%s", testCon.Host, testCon.Port.Port()),
		},
		Terminate: testCon.Terminate,
	}
}
