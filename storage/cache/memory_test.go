package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemory()

	// 测试 Set 和 Get
	err := cache.Set("key1", "value1", 10)
	assert.NoError(t, err)

	val, err := cache.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	// 测试过期
	err = cache.Expire("key1", time.Second)
	assert.NoError(t, err)

	time.Sleep(time.Second)

	val, err = cache.Get("key1")
	assert.Error(t, err) // 应该报错，因为键已过期

	// 测试 HashGet 和 HashDel
	err = cache.Set("hash:key", "hashValue", 10)
	assert.NoError(t, err)

	hashVal, err := cache.HashGet("hash:", "key")
	assert.NoError(t, err)
	assert.Equal(t, "hashValue", hashVal)

	err = cache.HashDel("hash:", "key")
	assert.NoError(t, err)

	hashVal, err = cache.HashGet("hash:", "key")
	assert.Error(t, err) // 应该报错，因为键已被删除

	// 测试 Increase 和 Decrease
	err = cache.Set("counter", 5, 10)
	assert.NoError(t, err)

	err = cache.Increase("counter")
	assert.NoError(t, err)

	counterVal, err := cache.Get("counter")
	assert.NoError(t, err)
	assert.Equal(t, "6", counterVal)

	err = cache.Decrease("counter")
	assert.NoError(t, err)

	counterVal, err = cache.Get("counter")
	assert.NoError(t, err)
	assert.Equal(t, "5", counterVal)

	// 测试 Del
	err = cache.Del("counter")
	assert.NoError(t, err)

	counterVal, err = cache.Get("counter")
	assert.Error(t, err) // 应该报错，因为键已被删除
}
