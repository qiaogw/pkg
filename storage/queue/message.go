package queue

import (
	"github.com/go-admin-team/redisqueue/v2"
	"sync"

	"github.com/qiaogw/pkg/storage"
)

// Message 是一个包装了 redisqueue.Message 的结构，用于添加错误计数功能。
type Message struct {
	redisqueue.Message              // 内嵌 redisqueue.Message，继承其属性
	ErrorCount         int          // 错误计数
	mux                sync.RWMutex // 读写锁，用于保护并发访问
}

// GetID 获取消息的唯一标识符。
func (m *Message) GetID() string {
	return m.ID
}

// GetStream 获取消息所在的流名称。
func (m *Message) GetStream() string {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.Stream
}

// GetValues 获取消息的键值对数据。
func (m *Message) GetValues() map[string]interface{} {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.Values
}

// SetID 设置消息的唯一标识符。
func (m *Message) SetID(id string) {
	m.ID = id
}

// SetStream 设置消息所在的流名称。
func (m *Message) SetStream(stream string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.Stream = stream
}

// SetValues 设置消息的键值对数据。
func (m *Message) SetValues(values map[string]interface{}) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.Values = values
}

// GetPrefix 获取消息的前缀。
func (m *Message) GetPrefix() (prefix string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.Values == nil {
		return
	}
	v, _ := m.Values[storage.PrefixKey]
	prefix, _ = v.(string)
	return
}

// SetPrefix 设置消息的前缀。
func (m *Message) SetPrefix(prefix string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.Values == nil {
		m.Values = make(map[string]interface{})
	}
	m.Values[storage.PrefixKey] = prefix
}

// SetErrorCount 设置消息的错误计数。
func (m *Message) SetErrorCount(count int) {
	m.ErrorCount = count
}

// GetErrorCount 获取消息的错误计数。
func (m *Message) GetErrorCount() int {
	return m.ErrorCount
}
