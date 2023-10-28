package queue

import (
	"fmt"
	"github.com/go-admin-team/redisqueue/v2"
	"github.com/qiaogw/pkg/storage"
	"log"
	"sync"
	"testing"
	"time"
)

func TestMemory_AppendAndRegister(t *testing.T) {
	memory := NewMemory(0)

	var receivedMessages []storage.Messager

	// 模拟一个消费者函数，用于处理收到的消息
	consumerFunc := func(message storage.Messager) error {
		receivedMessages = append(receivedMessages, message)
		return nil
	}

	// 注册消费者函数到队列中
	memory.Register("test_queue", consumerFunc)

	// 创建一条测试消息
	message := &Message{
		Message: redisqueue.Message{
			ID:     "1", // 使用 "1" 作为消息的预期 ID
			Stream: "test_queue",
			Values: map[string]interface{}{
				"key": "value",
			},
		},
	}

	// 将消息添加到队列
	err := memory.Append(message)
	if err != nil {
		t.Errorf("Append() error = %v", err)
	}

	// 等待一段时间以确保消息被处理
	time.Sleep(100 * time.Millisecond)

	// 检查是否正确接收和处理了消息

	if len(receivedMessages) != 1 {
		t.Errorf("Received %d messages, expected 1", len(receivedMessages))
	}

	receivedMessage := receivedMessages[0]

	if receivedMessage.GetStream() != message.GetStream() {
		t.Errorf("Received message stream %s, expected %s", receivedMessage.GetStream(), message.GetStream())
	}

	receivedValues := receivedMessage.GetValues()
	expectedValue := "value"
	if receivedValues["key"] != expectedValue {
		t.Errorf("Received message value %s, expected %s", receivedValues["key"], expectedValue)
	}
}

func TestMemory_Append(t *testing.T) {
	// 定义测试用例的数据结构
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
	}
	type args struct {
		name    string
		message storage.Messager
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test01",
			fields{},
			args{
				name: "test",
				message: &Message{
					Message: redisqueue.Message{
						ID:     "",
						Stream: "test",
						Values: map[string]interface{}{
							"key": "value",
						},
					},
					ErrorCount: 3,              // 初始化错误计数
					mux:        sync.RWMutex{}, // 初始化互斥锁
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory(100) // 创建一个内存消息队列
			if err := m.Append(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Register(t *testing.T) {
	log.SetFlags(19)
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
	}
	type args struct {
		name string
		f    storage.ConsumerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test01",
			fields{},
			args{
				name: "test",
				f: func(message storage.Messager) error {
					fmt.Println(message.GetValues())
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory(100)            // 创建一个内存消息队列
			m.Register(tt.name, tt.args.f) // 注册消费者函数到队列中
			if err := m.Append(&Message{redisqueue.Message{
				Stream: "test",
				Values: map[string]interface{}{
					"key": "value",
				},
			}, 3, sync.RWMutex{}}); err != nil {
				t.Error(err)
				return
			}
			go func() {
				m.Run() // 启动消息队列
			}()
			time.Sleep(3 * time.Second) // 等待一段时间以让消息被处理
			m.Shutdown()                // 停止消息队列的运行
		})
	}
}
