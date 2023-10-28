package queue

import (
	"github.com/qiaogw/pkg/storage"
	"sync"
	"time"

	"github.com/google/uuid"
)

// queue 代表一个消息通道。
type queue chan storage.Messager

// NewMemory 创建并初始化一个使用内存作为消息队列的 Memory 对象。
func NewMemory(poolNum uint) *Memory {
	return &Memory{
		queue:   new(sync.Map),
		PoolNum: poolNum,
	}
}

// Memory 是一个使用内存作为消息队列实现的 Queue 对象。
type Memory struct {
	queue   *sync.Map
	wait    sync.WaitGroup
	mutex   sync.RWMutex
	PoolNum uint
}

// String 返回队列的标识，用于标识它的类型。
func (*Memory) String() string {
	return "memory"
}

// makeQueue 创建一个新的消息通道。
func (m *Memory) makeQueue() queue {
	if m.PoolNum <= 0 {
		return make(queue)
	}
	return make(queue, m.PoolNum)
}

// Append 将消息添加到队列中。
func (m *Memory) Append(message storage.Messager) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	memoryMessage := new(Message)
	memoryMessage.SetID(message.GetID())
	memoryMessage.SetStream(message.GetStream())
	memoryMessage.SetValues(message.GetValues())

	v, ok := m.queue.Load(message.GetStream())

	if !ok {
		v = m.makeQueue()
		m.queue.Store(message.GetStream(), v)
	}

	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(message.GetStream(), q)
	}
	go func(gm storage.Messager, gq queue) {
		gm.SetID(uuid.New().String())
		gq <- gm
	}(memoryMessage, q)
	return nil
}

// Register 注册一个消费者函数用于处理特定队列中的消息。
func (m *Memory) Register(name string, f storage.ConsumerFunc) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, ok := m.queue.Load(name)
	if !ok {
		v = m.makeQueue()
		m.queue.Store(name, v)
	}
	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(name, q)
	}
	go func(out queue, gf storage.ConsumerFunc) {
		var err error
		for message := range q {
			err = gf(message)
			if err != nil {
				if message.GetErrorCount() < 3 {
					message.SetErrorCount(message.GetErrorCount() + 1)
					// 每次间隔时长放大
					i := time.Second * time.Duration(message.GetErrorCount())
					time.Sleep(i)
					out <- message
				}
				err = nil
			}
		}
	}(q, f)
}

// Run 启动队列并等待任务完成。
func (m *Memory) Run() {
	m.wait.Add(1)
	m.wait.Wait()
}

// Shutdown 停止队列的运行。
func (m *Memory) Shutdown() {
	m.wait.Done()
}
