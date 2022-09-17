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
			du := time.Millisecond * 10
			Expect(client.Set(ctx, key, value, du)).To(Succeed())
			time.Sleep(du * 10)
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

	Context("MultipleSet", func() {
		It("set multiple key-value", func() {
			kv := map[string]string{
				"k1": uuid.NewString(),
				"k2": uuid.NewString(),
			}
			Expect(client.MultipleSet(ctx, kv)).To(Succeed())
			values, err := redisClient.MGet(ctx, "k1", "k2").Result()
			Expect(err).To(BeNil())
			for _, v := range values {
				Expect(v).ToNot(BeNil())
			}
		})
	})

	Context("MultipleGet", func() {
		It("get list of string values", func() {
			kv := []string{"k1", uuid.NewString(), "k2", uuid.NewString()}
			Expect(redisClient.MSet(ctx, kv).Err()).To(Succeed())
			values, err := client.MultipleGet(ctx, "k1", "k2", "non-exist")
			Expect(err).To(BeNil())
			Expect(values).To(HaveLen(3))
			Expect(values).To(Equal([]string{kv[1], kv[3], ""}))
		})
	})

	Context("Hash Data", func() {
		It("set the key with given fields and values", func() {
			key := uuid.New().String()
			kv := map[string]string{
				"val1": uuid.New().String(),
				"val2": uuid.New().String(),
			}
			Expect(client.HashSet(ctx, key, kv)).To(Succeed())
			for field, val := range kv {
				v, err := redisClient.HGet(ctx, key, field).Result()
				Expect(err).To(BeNil())
				Expect(v).To(Equal(val))
			}
		})

		It("get the field", func() {
			key := uuid.New().String()
			field := uuid.New().String()
			val := uuid.New().String()
			Expect(redisClient.HSet(ctx, key, field, val).Err()).To(Succeed())
			v, err := client.HashGet(ctx, key, field)
			Expect(err).To(BeNil())
			Expect(v).To(Equal(val))
		})

		It("get empty string when get non-existed field", func() {
			key := uuid.New().String()
			Expect(redisClient.HSet(ctx, key, uuid.New().String(), uuid.New().String()).Err()).To(Succeed())
			v, err := client.HashGet(ctx, key, uuid.New().String())
			Expect(err).To(BeNil())
			Expect(v).To(BeEmpty())
		})

		It("get empty string when get non-existed key", func() {
			v, err := client.HashGet(ctx, uuid.New().String(), uuid.New().String())
			Expect(err).To(BeNil())
			Expect(v).To(BeEmpty())
		})
	})

	Context("Delete key(s)", func() {
		It("delete a key", func() {
			key := uuid.NewString()
			Expect(redisClient.Set(ctx, key, uuid.NewString(), 0).Err()).To(Succeed())
			Expect(client.Delete(ctx, key)).To(Succeed())
			Expect(redisClient.Get(ctx, key).Err()).To(Equal(redis.Nil))
		})

		It("delete multiple keys", func() {
			keys := []string{"k1", "k2"}
			Expect(redisClient.MSet(ctx, keys[0], uuid.NewString(), keys[1], uuid.NewString()).Err()).To(Succeed())
			Expect(client.Delete(ctx, keys...)).To(Succeed())
			count, err := redisClient.Exists(ctx, keys...).Result()
			Expect(err).To(BeNil())
			Expect(count).To(Equal(int64(0)))
		})

		It("delete a key of hash data", func() {
			key := uuid.NewString()
			field := uuid.NewString()
			Expect(redisClient.HSet(ctx, key, field, uuid.NewString()).Err()).To(Succeed())
			Expect(client.Delete(ctx, key)).To(Succeed())
			Expect(redisClient.HGet(ctx, key, field).Err()).To(Equal(redis.Nil))
		})

		It("return no error when delete non-existed key", func() {
			Expect(client.Delete(ctx, uuid.NewString())).To(Succeed())
		})
	})
})
