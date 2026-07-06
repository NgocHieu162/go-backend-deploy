package cache

import (
	"fmt"
	"go-backend/internal/common/env"
	"go-backend/internal/common/helpers"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	memoryStore *persist.MemoryStore
	redisStore  *persist.RedisStore
}

func NewCache(env *env.Env) *Cache {

	// deploy restart BE là dữ liệu cache mất
	memoryStore := persist.NewMemoryStore(1 * time.Minute)

	// database phụ, dữ liệu cache muốn share cho các backend khác
	// backend restart thì không bị mất dữ liệu cache
	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     env.RedisAddr,
		Password: env.RedisPass,
	}))

	// statusCmd := redisStore.RedisClient.Ping(context.Background())
	// err := statusCmd.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("[REDIS] Connection To Redis Successfully", statusCmd.String())

	return &Cache{
		memoryStore: memoryStore,
		redisStore:  redisStore,
	}
}

func (c *Cache) MemoryCache(defaultExpire time.Duration) gin.HandlerFunc {
	return cache.CacheByRequestURI(
		c.memoryStore,
		defaultExpire,
		cache.WithCacheStrategyByRequest(customKey),
	)
}

func (c *Cache) RedisCache(defaultExpire time.Duration) gin.HandlerFunc {
	return cache.CacheByRequestURI(
		c.redisStore,
		defaultExpire,
		cache.WithCacheStrategyByRequest(customKey),
	)
}

func customKey(c *gin.Context) (bool, cache.Strategy) {
	user, err := helpers.GetUsers(c)
	key := fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.RequestURI())
	if err == nil && user != nil {
		key = fmt.Sprintf("%d:%s:%s", user.ID, c.Request.Method, c.Request.URL.RequestURI())
	}
	return true, cache.Strategy{
		CacheKey: key,
	}
}
