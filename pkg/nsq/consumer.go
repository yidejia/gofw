package nsq

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cast"

	"github.com/nsqio/go-nsq"
	"github.com/yidejia/gofw/pkg/config"
)

// consumers 消息消费者映射表
var consumers sync.Map

// RegisterConsumer 注册 NSQ 消息消费者
func RegisterConsumer(addr, topic, channel string, handler nsq.Handler, concurrency int) {

	conf := nsq.NewConfig()
	conf.LookupdPollInterval = 1 * time.Second
	conf.MaxInFlight = 1000
	consumer, err := nsq.NewConsumer(topic, channel, conf)
	if err != nil {
		panic("New NSQ consumer failed:" + err.Error())
	}

	consumer.AddConcurrentHandlers(handler, concurrency)

	// 调用方指定了 topic 所在节点地址，直接连接就可以了
	if len(addr) > 0 {
		if err = consumer.ConnectToNSQD(addr); err != nil {
			panic(fmt.Sprintf("NSQ consume consumer to NSQD %s failed: %s", addr, err.Error()))
		}
		// 缓存消费者
		consumers.Store(
			fmt.Sprintf("%s_%s_%s", addr, topic, channel),
			consumer,
		)
		return
	}

	// 调用方未指定 topic 所在节点地址，需要收集 NSQ lookupd 节点地址，通过连接 lookupd 节点，自动发现 topic 所在的 nsqd 节点
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
		panic("NSQ consumer connect to lookupd failed: " + err.Error())
	}

	// 缓存消费者
	consumers.Store(
		fmt.Sprintf("%s_%s", topic, channel),
		consumer,
	)
}

// StopConsumers 停止 NSQ 消息消费者
func StopConsumers() {
	consumers.Range(func(key, value interface{}) bool {
		if c, ok := value.(*nsq.Consumer); ok {
			c.Stop()
		}
		return true
	})
}
