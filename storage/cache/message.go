package cache

import (
	"github.com/go-admin-team/redisqueue/v2"
	"github.com/qiaogw/pkg/storage"
)

// Message 是一个封装了 redisqueue.Message 的自定义消息类型。
type Message struct {
	redisqueue.Message
}

// GetID 返回消息的唯一标识符。
func (m *Message) GetID() string {
	return m.ID
}

// GetStream 返回消息所属的流名称。
func (m *Message) GetStream() string {
	return m.Stream
}

// GetValues 返回消息的键值对数据。
func (m *Message) GetValues() map[string]interface{} {
	return m.Values
}

// SetID 设置消息的唯一标识符。
func (m *Message) SetID(id string) {
	m.ID = id
}

// SetStream 设置消息所属的流名称。
func (m *Message) SetStream(stream string) {
	m.Stream = stream
}

// SetValues 设置消息的键值对数据。
func (m *Message) SetValues(values map[string]interface{}) {
	m.Values = values
}

// GetPrefix 返回消息的前缀信息。
func (m *Message) GetPrefix() (prefix string) {
	if m.Values == nil {
		return
	}
	v, _ := m.Values[storage.PrefixKey]
	prefix, _ = v.(string)
	return
}

// SetPrefix 设置消息的前缀信息。
func (m *Message) SetPrefix(prefix string) {
	if m.Values == nil {
		m.Values = make(map[string]interface{})
	}
	m.Values[storage.PrefixKey] = prefix
}
