package hospital_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestHospital(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hospital Suite")
}

var compose *testcontainers.LocalDockerCompose

var _ = BeforeSuite(func() {
	setupTestHospitalSystem()
})

var _ = AfterSuite(func() {
	Expect(compose.Down().Error).To(BeNil())
})

func setupTestHospitalSystem() {
	id := uuid.New().String()
	compose = testcontainers.NewLocalDockerCompose([]string{"docker-compose.test.yaml"}, id)
	setupComposeService(compose, "postgres", wait.ForLog("database system is ready to accept connections").WithOccurrence(2))
	setupComposeService(compose, "rabbitmq", wait.ForLog("Ready to start client connection listeners"))
	setupComposeService(compose, "hospital-sys", wait.ForLog("Nest application successfully started"))

}

func setupComposeService(compose *testcontainers.LocalDockerCompose, service string, wait wait.Strategy) {
	err := compose.
		WithCommand([]string{"up", "-d", service}).
		WaitForService(service, wait).
		Invoke().Error
	Expect(err).To(BeNil())
}
