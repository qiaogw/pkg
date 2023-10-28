package locker

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/bsm/redislock"
)

// NewRedis 创建并初始化一个使用 Redis 实现的 Locker。
func NewRedis(c *redis.Client) *Redis {
	return &Redis{
		client: c,
	}
}

// Redis 是一个使用 Redis 实现的 Locker。
type Redis struct {
	client *redis.Client
	mutex  *redislock.Client
}

// String 返回该 Locker 的标识，用于识别它的类型。
func (Redis) String() string {
	return "redis"
}

// Lock 尝试获取一个分布式锁。key 是锁的唯一标识，ttl 是锁的超时时间（秒），
// options 是锁的配置选项。如果获取锁成功，将返回一个锁对象；如果获取失败，将返回错误。
func (r *Redis) Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error) {
	// 如果 r.mutex 为 nil，创建一个 redislock.Client 并初始化。
	if r.mutex == nil {
		r.mutex = redislock.New(r.client)
	}
	// 使用 Obtain 方法尝试获取锁，传入 key、锁的超时时间和配置选项。
	return r.mutex.Obtain(context.TODO(), key, time.Duration(ttl)*time.Second, options)
}
