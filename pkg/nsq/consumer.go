package nsq

import (
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/yidejia/gofw/pkg/config"
)

// RegisterConsumer 注册 NSQ 消息消费者
func RegisterConsumer(topic, channel string, handler nsq.Handler, concurrency int) {

	conf := nsq.NewConfig()
	conf.LookupdPollInterval = 1 * time.Second
	conf.MaxInFlight = 1000
	consumer, err := nsq.NewConsumer(topic, channel, conf)
	if err != nil {
		panic("New NSQ consumer failed:" + err.Error())
	}

	consumer.AddConcurrentHandlers(handler, concurrency)

	if err = consumer.ConnectToNSQD(fmt.Sprintf(
		"%s:%d",
		config.Get("queue.nsq.host"),
		config.GetInt("queue.nsq.port"),
	)); err != nil {
		panic("NSQ consumer connect to NSQD failed:" + err.Error())
	}
}
