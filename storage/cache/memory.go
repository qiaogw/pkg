package cache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cast"
)

// item 表示缓存中的项，包含值和过期时间。
type item struct {
	Value   string
	Expired time.Time
}

// NewMemory 创建并初始化一个使用内存作为缓存实现的 Memory 对象。
func NewMemory() *Memory {
	return &Memory{
		items: new(sync.Map),
	}
}

// Memory 是一个使用内存作为缓存实现的 Cache 对象。
type Memory struct {
	items *sync.Map
	mutex sync.RWMutex
}

// String 返回 Cache 的标识，用于标识它的类型。
func (*Memory) String() string {
	return "memory"
}

// connect 暂时无实际操作，用于保持接口一致性。
func (m *Memory) connect() {
}

// Get 从缓存中获取指定 key 的值。
func (m *Memory) Get(key string) (string, error) {
	item, err := m.getItem(key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, nil
}

// getItem 从缓存中获取指定 key 的项，如果项不存在或已过期，则返回 nil。
func (m *Memory) getItem(key string) (*item, error) {
	var err error
	i, ok := m.items.Load(key)
	if !ok {
		return nil, nil
	}
	switch i.(type) {
	case *item:
		item := i.(*item)
		if item.Expired.Before(time.Now()) {
			// 过期，删除
			_ = m.del(key)
			return nil, nil
		}
		return item, nil
	default:
		err = fmt.Errorf("value of %s type error", key)
		return nil, err
	}
}

// Set 将指定 key 的值设置为给定的 val，并设置过期时间。
func (m *Memory) Set(key string, val interface{}, expire int) error {
	s, err := cast.ToStringE(val)
	if err != nil {
		return err
	}
	item := &item{
		Value:   s,
		Expired: time.Now().Add(time.Duration(expire) * time.Second),
	}
	return m.setItem(key, item)
}

// setItem 将指定 key 的项设置为给定的 item。
func (m *Memory) setItem(key string, item *item) error {
	m.items.Store(key, item)
	return nil
}

// Del 删除缓存中的指定 key。
func (m *Memory) Del(key string) error {
	return m.del(key)
}

// del 删除缓存中的指定 key。
func (m *Memory) del(key string) error {
	m.items.Delete(key)
	return nil
}

// HashGet 从指定的哈希表中获取指定 key 的值。
func (m *Memory) HashGet(hk, key string) (string, error) {
	item, err := m.getItem(hk + key)
	if err != nil || item == nil {
		return "", err
	}
	return item.Value, err
}

// HashDel 从指定的哈希表中删除指定 key。
func (m *Memory) HashDel(hk, key string) error {
	return m.del(hk + key)
}

// Increase 将指定 key 的值增加 num。
func (m *Memory) Increase(key string) error {
	return m.calculate(key, 1)
}

// Decrease 将指定 key 的值减少 num。
func (m *Memory) Decrease(key string) error {
	return m.calculate(key, -1)
}

// calculate 对指定 key 的值进行增减操作。
func (m *Memory) calculate(key string, num int) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return err
	}
	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	var n int
	n, err = cast.ToIntE(item.Value)
	if err != nil {
		return err
	}
	n += num
	item.Value = strconv.Itoa(n)
	return m.setItem(key, item)
}

// Expire 设置指定 key 的过期时间。
func (m *Memory) Expire(key string, dur time.Duration) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return err
	}
	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	item.Expired = time.Now().Add(dur)
	return m.setItem(key, item)
}
