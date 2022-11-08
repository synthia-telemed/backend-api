package hospital_test

import (
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
	Expect(compose.Down().Error).To(Succeed())
})

func setupTestHospitalSystem() {
	id := uuid.New().String()
	compose = testcontainers.NewLocalDockerCompose([]string{"docker-compose.test.yaml"}, id)
	execError := compose.
		WithCommand([]string{"up", "-d"}).
		WaitForService("postgres", wait.ForLog("database system is ready to accept connections").WithOccurrence(2)).
		WaitForService("rabbitmq", wait.ForLog("Ready to start client connection listeners")).
		WaitForService("hospital-sys", wait.ForLog("Nest application successfully started").WithStartupTimeout(time.Minute*2)).
		Invoke()
	Expect(execError.Error).To(BeNil())
}
