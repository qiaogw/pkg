package queue

import (
	json "github.com/json-iterator/go"
	"github.com/nsqio/go-nsq"
	"github.com/qiaogw/pkg/storage"
)

// NewNSQ 创建并初始化一个使用 NSQ 作为消息队列实现的 NSQ 对象。
func NewNSQ(addresses []string, cfg *nsq.Config, channelPrefix string) (*NSQ, error) {
	n := &NSQ{
		addresses:     addresses,
		cfg:           cfg,
		channelPrefix: channelPrefix,
	}
	var err error
	n.producer, err = n.newProducer()
	return n, err
}

// NSQ 是一个使用 NSQ 作为消息队列实现的 Queue 对象。
type NSQ struct {
	addresses     []string
	cfg           *nsq.Config
	producer      *nsq.Producer
	consumer      *nsq.Consumer
	channelPrefix string
}

// String 返回队列的标识，用于标识它的类型。
func (NSQ) String() string {
	return "nsq"
}

// switchAddress ⚠️生产环境至少配置三个节点
func (e *NSQ) switchAddress() {
	if len(e.addresses) > 1 {
		e.addresses[0], e.addresses[len(e.addresses)-1] =
			e.addresses[1],
			e.addresses[0]
	}
}

// newProducer 使用指定的选项创建一个新的 nsq.Producer。
func (e *NSQ) newProducer() (*nsq.Producer, error) {
	if e.cfg == nil {
		e.cfg = nsq.NewConfig()
	}
	return nsq.NewProducer(e.addresses[0], e.cfg)
}

// newConsumer 使用指定的选项创建一个新的 nsq.Consumer。
func (e *NSQ) newConsumer(topic string, h nsq.Handler) (err error) {
	if e.cfg == nil {
		e.cfg = nsq.NewConfig()
	}
	if e.consumer == nil {
		e.consumer, err = nsq.NewConsumer(topic, e.channelPrefix+topic, e.cfg)
		if err != nil {
			return err
		}
	}
	e.consumer.AddHandler(h)
	err = e.consumer.ConnectToNSQDs(e.addresses)

	return err
}

// Append 将消息添加到队列中。
func (e *NSQ) Append(message storage.Messager) error {
	rb, err := json.Marshal(message.GetValues())
	if err != nil {
		return err
	}
	return e.producer.Publish(message.GetStream(), rb)
}

// Register 注册一个消费者函数用于处理特定队列中的消息。
func (e *NSQ) Register(name string, f storage.ConsumerFunc) {
	h := &nsqConsumerHandler{f}
	err := e.newConsumer(name, h)
	if err != nil {
		//目前不支持动态注册
		panic(err)
	}
}

// Run 启动队列并开始消费消息。
func (e *NSQ) Run() {
}

// Shutdown 停止队列的运行。
func (e *NSQ) Shutdown() {
	if e.producer != nil {
		e.producer.Stop()
	}
	if e.consumer != nil {
		e.consumer.Stop()
	}
}

// nsqConsumerHandler 实现了 nsq.Handler 接口，用于处理 NSQ 消息。
type nsqConsumerHandler struct {
	f storage.ConsumerFunc
}

// HandleMessage 处理 NSQ 消息。
func (e nsqConsumerHandler) HandleMessage(message *nsq.Message) error {
	m := new(Message)
	data := make(map[string]interface{})
	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		return err
	}
	m.SetValues(data)
	return e.f(m)
}
