package nsq

import (
	"fmt"
	"time"

	"github.com/spf13/cast"

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

	// 收集 NSQ lookupd 节点地址
	var addresses []string
	lookupds := config.GetInterface("queue.nsq.lookupds")
	if _lookupds, ok := lookupds.([]map[string]interface{}); ok {
		for _, lookupd := range _lookupds {
			addresses = append(addresses, fmt.Sprintf("%s:%d", cast.ToString(lookupd["host"]), cast.ToInt(lookupd["port"])))
		}
	}
	if len(addresses) == 0 {
		panic("No NSQ lookupd node")
	}

	if err = consumer.ConnectToNSQLookupds(addresses); err != nil {
		panic("NSQ consumer connect to lookupd failed:" + err.Error())
	}
}
