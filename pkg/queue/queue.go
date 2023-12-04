// Package queue 队列包
// @author 余海坚 haijianyu10@qq.com
// @created 2023-11-23 15:47
// @copyright © 2010-2023 广州伊的家网络科技有限公司
package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/yidejia/gofw/pkg/nsq"

	nsqPKG "github.com/nsqio/go-nsq"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/logger"
)

// Job 队列任务接口
type Job interface {
	// JobName 返回队列任务名称
	JobName() string
	// OnJobQueue 执行任务的队列
	OnJobQueue() string
	// NewJob 新建任务用于将队列消息绑定到结构体上
	NewJob() Job
	// HandleJob 处理任务
	HandleJob(job Job) error
}

// jobs 已注册的队列任务
var jobs = make(map[string]Job)

// RegisterJob 注册队列任务
func RegisterJob(job Job) {
	jobs[job.JobName()] = job
}

// GetJob 根据任务名获取队列任务
func GetJob(name string) Job {
	if job, ok := jobs[name]; ok {
		return job
	}
	return nil
}

// dispatchJob 内部分发队列任务
func dispatchJob(job Job) (producer *nsq.Producer, topic, message string, err error) {

	// 以 JSON 格式序列化任务
	var messageByte []byte
	messageByte, err = json.Marshal(job)
	if err != nil {
		err = errors.New("序列化任务失败：" + err.Error())
		return
	}

	// 发布消息的生产者
	producer = ConnectProducer(job.OnJobQueue())
	// 发布消息的 topic
	topic = fmt.Sprintf("%s_queue_%s", config.Get("app.name"), job.OnJobQueue())
	// 发布的消息，使用『@@』分割任务名和任务消息
	message = fmt.Sprintf("%s@@%s", job.JobName(), string(messageByte))

	return
}

// DispatchJob 分发队列任务
func DispatchJob(job Job) error {

	producer, topic, message, err := dispatchJob(job)
	if err != nil {
		return err
	}

	// 发布消息到队列
	if err = producer.Publish(topic, message); err != nil {
		return errors.New("发布消息到队列失败：" + err.Error())
	}

	return nil
}

// DispatchJobDelay 延迟分发队列任务
func DispatchJobDelay(job Job, delay time.Duration) error {

	producer, topic, message, err := dispatchJob(job)
	if err != nil {
		return err
	}

	// 延迟发布消息到队列
	if err = producer.DeferredPublish(topic, message, delay); err != nil {
		return errors.New("延迟发布消息到队列失败：" + err.Error())
	}

	return nil
}

// JobHandler 队列任务处理器
type JobHandler struct {
}

// HandleMessage 处理 NSQ 队列消息
func (h *JobHandler) HandleMessage(message *nsqPKG.Message) error {

	if len(message.Body) == 0 {
		// 返回 nil 将自动发送一个 FIN 命令到 NSQ 标识消息已处理
		return nil
	}

	// 标记消息已处理
	message.Finish()

	// 分割任务名和任务参数
	_message := string(message.Body)
	messages := strings.Split(_message, "@@")
	// 消息格式不正确，只记录日志
	if len(messages) < 2 {
		logger.ErrorString("队列任务", "处理 NSQ 队列消息-分割任务名和任务参数", fmt.Sprintf("消息格式不正确：%s", _message))
		return nil
	}

	job := GetJob(messages[0])
	// 任务不存在，只记录日志
	if job == nil {
		logger.ErrorString("队列任务", "处理 NSQ 队列消息-根据任务名获取任务", fmt.Sprintf("任务[%s]不存在", messages[0]))
		return nil
	}

	// 任务存入消息队列中时以 JSON 格式序列化，现在反序列化消息到任务结构体
	_job := job.NewJob()
	err := json.Unmarshal([]byte(messages[1]), _job)
	// 反序列化消息失败
	if err != nil {
		logger.ErrorString("队列任务", "处理 NSQ 队列消息-反序列化消息", fmt.Sprintf("反序列化消息 %s 失败", messages[1]))
		return nil
	}

	// 处理任务
	return job.HandleJob(_job)
}

// producers 消息生产者映射表
var producers sync.Map

// InitWithConfig 加载配置初始化队列
func InitWithConfig() {
	// 注册队列消息消费者
	if consumers := config.GetStringMap("queue.job"); len(consumers) > 0 {

		appName := config.Get("app.name")
		var topic string
		jobHandler := &JobHandler{}

		for _queue := range consumers {

			// 获取队列的消息生产者
			producer := nsq.ConnectProducer(config.Get(fmt.Sprintf("queue.job.%s.producer", _queue)))
			// 缓存消息生产者，方便通过队列名快捷获取
			producers.Store(_queue, producer)

			topic = fmt.Sprintf("%s_queue_%s", appName, _queue)

			consumerAddr := ""
			// 不启用 NSQ lookupd 节点时，只能直连消费者节点
			if !config.GetBool("queue.nsq.enable_lookupd") {
				consumerAddr = config.Get(fmt.Sprintf("queue.job.%s.consumer.addr", _queue))
				if len(consumerAddr) == 0 {
					panic(fmt.Sprintf("queue %s can not lookup consumer addr", _queue))
				}
			}

			nsq.RegisterConsumer(consumerAddr, topic, "1", jobHandler, config.GetInt(fmt.Sprintf("queue.job.%s.processes", _queue)))
		}
	}
}

// ConnectProducer 连接消息生产者
func ConnectProducer(name ...string) *nsq.Producer {

	var _name string
	if len(name) > 0 {
		_name = name[0]
	} else {
		_name = "default"
	}

	producer, ok := producers.Load(_name)
	if !ok {
		panic(fmt.Sprintf("[Queue] NSQ producer %s not exists", _name))
	}

	return producer.(*nsq.Producer)
}
