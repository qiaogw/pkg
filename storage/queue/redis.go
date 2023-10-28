package queue

import (
	"github.com/go-admin-team/redisqueue/v2"
	"github.com/qiaogw/pkg/storage"
	"github.com/redis/go-redis/v9"
)

// NewRedis 创建并初始化一个使用 Redis 作为消息队列实现的 Redis 对象。
func NewRedis(
	producerOptions *redisqueue.ProducerOptions,
	consumerOptions *redisqueue.ConsumerOptions,
) (*Redis, error) {
	var err error
	r := &Redis{}
	r.producer, err = r.newProducer(producerOptions)
	if err != nil {
		return nil, err
	}
	r.consumer, err = r.newConsumer(consumerOptions)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Redis 是一个使用 Redis 作为消息队列实现的 Queue 对象。
type Redis struct {
	client   *redis.Client
	consumer *redisqueue.Consumer
	producer *redisqueue.Producer
}

// String 返回队列的标识，用于标识它的类型。
func (Redis) String() string {
	return "redis"
}

// newConsumer 使用指定的选项创建一个新的 redisqueue.Consumer。
func (r *Redis) newConsumer(options *redisqueue.ConsumerOptions) (*redisqueue.Consumer, error) {
	if options == nil {
		options = &redisqueue.ConsumerOptions{}
	}
	return redisqueue.NewConsumerWithOptions(options)
}

// newProducer 使用指定的选项创建一个新的 redisqueue.Producer。
func (r *Redis) newProducer(options *redisqueue.ProducerOptions) (*redisqueue.Producer, error) {
	if options == nil {
		options = &redisqueue.ProducerOptions{}
	}
	return redisqueue.NewProducerWithOptions(options)
}

// Append 将消息添加到队列中。
func (r *Redis) Append(message storage.Messager) error {
	err := r.producer.Enqueue(&redisqueue.Message{
		ID:     message.GetID(),
		Stream: message.GetStream(),
		Values: message.GetValues(),
	})
	return err
}

// Register 注册一个消费者函数用于处理特定队列中的消息。
func (r *Redis) Register(name string, f storage.ConsumerFunc) {
	r.consumer.Register(name, func(message *redisqueue.Message) error {
		m := new(Message)
		m.SetValues(message.Values)
		m.SetStream(message.Stream)
		m.SetID(message.ID)
		return f(m)
	})
}

// Run 启动队列并开始消费消息。
func (r *Redis) Run() {
	r.consumer.Run()
}

// Shutdown 停止队列的运行。
func (r *Redis) Shutdown() {
	r.consumer.Shutdown()
}
