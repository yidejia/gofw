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

var once sync.Once

// internalProducer 内部使用的 NSQ 消息生产者单例
var internalProducer *nsq.Producer

// InitProducerWithConfig 加载配置并初始化 NSQ 消息生产者单例
func InitProducerWithConfig() {

	once.Do(func() {

		var err error

		internalProducer, err = nsq.NewProducer(
			fmt.Sprintf(
				"%s:%d",
				config.Get("queue.nsq.host"),
				config.GetInt("queue.nsq.port"),
			),
			nsq.NewConfig(),
		)
		if err != nil {
			panic("init nsq producer failed:" + err.Error())
		}

		if err = internalProducer.Ping(); err != nil {
			panic("nsq producer ping failed:" + err.Error())
		}
	})
}

// Publish 发布消息
func Publish(topic, message string) error {
	return internalProducer.Publish(topic, []byte(message))
}

// DeferredPublish 延迟发布消息
func DeferredPublish(topic, message string, delay time.Duration) error {
	return internalProducer.DeferredPublish(topic, delay, []byte(message))
}

// StopProducer 停止 NSQ 消息生产者单例
func StopProducer() {
	if internalProducer != nil {
		internalProducer.Stop()
	}
}
