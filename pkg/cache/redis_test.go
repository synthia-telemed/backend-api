package cache_test

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/synthia-telemed/backend-api/pkg/cache"
	"math/rand"
	"time"
)

var _ = Describe("Cache Suite", func() {
	var (
		redisClient *redis.Client
		client      cache.Client
		ctx         context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		redisClient = redis.NewClient(&redis.Options{Addr: redisContainer.Endpoint})
		Expect(redisClient.Ping(ctx).Err()).To(Succeed())
		client = cache.NewRedisClient(&redisContainer.Config)

	})

	Context("Basic Get and Set", func() {
		var (
			key   string
			value string
		)

		BeforeEach(func() {
			rand.Seed(GinkgoRandomSeed())
			key = uuid.New().String()
			value = uuid.New().String()
		})

		It("set the value", func() {
			Expect(client.Set(ctx, key, value, 0)).To(Succeed())

			retrievedValue, err := redisClient.Get(ctx, key).Result()
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(Equal(value))
		})

		It("set value with expiration time", func() {
			du := time.Millisecond
			Expect(client.Set(ctx, key, value, du)).To(Succeed())
			time.Sleep(du)
			Expect(redisClient.Get(ctx, key).Err()).To(Equal(redis.Nil))
		})

		It("get the value and not delete", func() {
			Expect(redisClient.Set(ctx, key, value, 0).Err()).To(Succeed())

			retrievedValue, err := client.Get(ctx, key, false)
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(Equal(value))

			val, err := redisClient.Get(ctx, key).Result()
			Expect(err).To(BeNil())
			Expect(val).To(Equal(value))
		})

		It("get the value and delete", func() {
			Expect(redisClient.Set(ctx, key, value, 0).Err()).To(Succeed())

			retrievedValue, err := client.Get(ctx, key, true)
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(Equal(value))

			Expect(redisClient.Get(ctx, key).Err()).To(Equal(redis.Nil))
		})

		It("return empty string if key does not exist", func() {
			retrievedValue, err := client.Get(ctx, key, false)
			Expect(err).To(BeNil())
			Expect(retrievedValue).To(BeEmpty())
		})

	})

})
