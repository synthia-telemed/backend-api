package datastore_test

import (
	"fmt"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/datastore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/rand"
	"os"
)

var _ = Describe("Patient Datastore", Ordered, func() {

	var (
		db               *gorm.DB
		patientDataStore datastore.PatientDataStore
	)

	BeforeAll(func() {
		_ = godotenv.Load(".test.env")
		var err error
		db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		var err error
		patientDataStore, err = datastore.NewGormPatientDataStore(db)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		Expect(db.Migrator().DropTable(&datastore.Patient{})).To(BeNil())
	})

	Context("Finding", func() {

		var users []datastore.Patient

		BeforeEach(func() {
			users = generateUsers(10)
			db.Create(&users)
		})

		It("should find user by ID", func() {
			user := users[0]
			foundUser, err := patientDataStore.FindByID(user.ID)
			Expect(err).To(BeNil())
			Expect(foundUser.ID).To(Equal(user.ID))
		})
	})

})

func generateUsers(num int) []datastore.Patient {
	rand.Seed(GinkgoRandomSeed())
	users := make([]datastore.Patient, num)
	for i := 0; i < num; i++ {
		users[i] = datastore.Patient{
			RefID: fmt.Sprintf("HN-%d", rand.Uint32()),
		}
	}
	return users
}
