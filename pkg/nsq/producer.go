// Package nsq 定制包
// @author 余海坚 haijianyu10@qq.com
// @created 2023-11-22 17:53
// @copyright © 2010-2023 广州伊的家网络科技有限公司
package nsq

import (
	"fmt"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/yidejia/gofw/pkg/config"
)

// Producer 消息生产者
// 对 NSQ 消息生产者进行自定义封装
// @author 余海坚 haijianyu10@qq.com
// @created 2023-12-01 16:47
// @copyright © 2010-2023 广州伊的家网络科技有限公司
type Producer struct {
	NSQProducer *nsq.Producer
}

// producers 消息生产者映射表
var producers sync.Map

var once sync.Once

// InitProducersWithConfig 加载配置并初始化 NSQ 消息生产者
func InitProducersWithConfig() {

	once.Do(func() {

		// 初始化消息生产者并缓存
		for name := range config.GetStringMap("queue.nsq.producers") {

			producer, err := nsq.NewProducer(
				config.Get(fmt.Sprintf("queue.nsq.producers.%s.addr", name)),
				nsq.NewConfig(),
			)
			if err != nil {
				panic(fmt.Sprintf("init NSQ producer %s failed: %s", name, err.Error()))
			}

			if err = producer.Ping(); err != nil {
				panic(fmt.Sprintf("NSQ producer %s ping failed: %s", name, err.Error()))
			}

			// 缓存消息生产者，方便通过生产者名称快捷获取
			producers.Store(name, &Producer{NSQProducer: producer})
		}
	})
}

// ConnectProducer 连接消息生产者
func ConnectProducer(name ...string) *Producer {

	var _name string
	if len(name) > 0 {
		_name = name[0]
	} else {
		_name = "default"
	}

	producer, ok := producers.Load(_name)
	if !ok {
		panic(fmt.Sprintf("NSQ producer %s not exists", _name))
	}

	return producer.(*Producer)
}

// Publish 发布消息
func (p *Producer) Publish(topic, message string) error {
	return p.NSQProducer.Publish(topic, []byte(message))
}

// DeferredPublish 延迟发布消息
func (p *Producer) DeferredPublish(topic, message string, delay time.Duration) error {
	return p.NSQProducer.DeferredPublish(topic, delay, []byte(message))
}

// StopProducers 停止 NSQ 消息生产者
func StopProducers() {
	producers.Range(func(key, value interface{}) bool {
		if p, ok := value.(*Producer); ok {
			p.NSQProducer.Stop()
		}
		return true
	})
}
