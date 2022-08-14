package cache_test

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"math/rand"
)

var _ = Describe("Cache Suite", func() {
	var (
		redis  *miniredis.Miniredis
		client cache.Client
	)

	BeforeEach(func() {
		redis = miniredis.RunT(GinkgoT())
		client = cache.NewRedisClient(&cache.Config{
			Endpoint: redis.Addr(),
		})
	})

	Context("Basic Get and Set", func() {
		var (
			key   string
			value string
		)

		BeforeEach(func() {
			rand.Seed(GinkgoRandomSeed())
			key = fmt.Sprintf("key-%d", rand.Int())
			value = fmt.Sprintf("value-%d", rand.Int())
		})

		It("set the value", func() {
			err := client.Set(context.Background(), key, value, 0)
			Expect(err).To(BeNil())

			retrievedValue, err := redis.Get(key)
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(Equal(value))
		})

		It("get the value and not delete", func() {
			err := redis.Set(key, value)
			Expect(err).To(BeNil())

			retrievedValue, err := client.Get(context.Background(), key, false)
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(Equal(value))

			val, err := redis.Get(key)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(value))
		})

		It("get the value and delete", func() {
			err := redis.Set(key, value)
			Expect(err).To(BeNil())

			retrievedValue, err := client.Get(context.Background(), key, true)
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(Equal(value))

			_, err = redis.Get(key)
			Expect(err).To(Equal(miniredis.ErrKeyNotFound))
		})

		It("return empty string if key does not exist", func() {
			retrievedValue, err := client.Get(context.Background(), key, false)
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(BeEmpty())
		})

	})

})
