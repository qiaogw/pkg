package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// NewRedis 创建一个使用 Redis 作为缓存实现的 Cache 对象。
func NewRedis(client *redis.Client, options *redis.Options) (*Redis, error) {
	if client == nil {
		client = redis.NewClient(options)
	}
	r := &Redis{
		client: client,
	}
	err := r.connect()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Redis 是一个使用 Redis 作为缓存实现的 Cache 对象。
type Redis struct {
	client *redis.Client
}

// String 返回 Cache 的标识，用于标识它的类型。
func (*Redis) String() string {
	return "redis"
}

// connect 测试与 Redis 的连接。
func (r *Redis) connect() error {
	_, err := r.client.Ping(context.TODO()).Result()
	return err
}

// Get 从缓存中获取指定 key 的值。
func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.TODO(), key).Result()
}

// Set 将指定 key 的值设置为给定的 val，并设置过期时间。
func (r *Redis) Set(key string, val interface{}, expire int) error {
	return r.client.Set(context.TODO(), key, val, time.Duration(expire)*time.Second).Err()
}

// Del 删除缓存中的指定 key。
func (r *Redis) Del(key string) error {
	return r.client.Del(context.TODO(), key).Err()
}

// HashGet 从指定的哈希表中获取指定 key 的值。
func (r *Redis) HashGet(hk, key string) (string, error) {
	return r.client.HGet(context.TODO(), hk, key).Result()
}

// HashDel 从指定的哈希表中删除指定 key。
func (r *Redis) HashDel(hk, key string) error {
	return r.client.HDel(context.TODO(), hk, key).Err()
}

// Increase 对指定 key 的值进行递增操作。
func (r *Redis) Increase(key string) error {
	return r.client.Incr(context.TODO(), key).Err()
}

// Decrease 对指定 key 的值进行递减操作。
func (r *Redis) Decrease(key string) error {
	return r.client.Decr(context.TODO(), key).Err()
}

// Expire 设置指定 key 的过期时间。
func (r *Redis) Expire(key string, dur time.Duration) error {
	return r.client.Expire(context.TODO(), key, dur).Err()
}

// GetClient 获取原始的 Redis 客户端。
func (r *Redis) GetClient() *redis.Client {
	return r.client
}
