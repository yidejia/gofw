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
	"time"

	"github.com/yidejia/gofw/pkg/nsq"

	"github.com/spf13/cast"

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
	// NewJob 新建任务用于将消息绑定到结构体上
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
func dispatchJob(job Job) (topic, message string, err error) {

	// 以 JSON 格式序列化任务
	var messageByte []byte
	messageByte, err = json.Marshal(job)
	if err != nil {
		err = errors.New("序列化任务失败：" + err.Error())
		return
	}

	// 发布消息的 topic
	topic = fmt.Sprintf("%s_queue_%s", config.Get("app.name"), job.OnJobQueue())
	// 发布的消息，使用『@@』分割任务名和任务消息
	message = fmt.Sprintf("%s@@%s", job.JobName(), string(messageByte))

	return
}

// DispatchJob 分发队列任务
func DispatchJob(job Job) error {

	topic, message, err := dispatchJob(job)
	if err != nil {
		return err
	}

	// 发布消息到队列
	if err = nsq.Publish(topic, message); err != nil {
		return errors.New("发布消息到队列失败：" + err.Error())
	}

	return nil
}

// DispatchJobDelay 延迟分发队列任务
func DispatchJobDelay(job Job, delay time.Duration) error {

	topic, message, err := dispatchJob(job)
	if err != nil {
		return err
	}

	// 延迟发布消息到队列
	if err = nsq.DeferredPublish(topic, message, delay); err != nil {
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

// InitWithConfig 加载配置初始化队列
func InitWithConfig() {
	// 注册队列消息消费者
	if consumers := config.GetStringMap("queue.job"); len(consumers) > 0 {
		appName := config.Get("app.name")
		var topic string
		jobHandler := &JobHandler{}
		for _queue, _consumer := range consumers {
			consumer := cast.ToStringMap(_consumer)
			topic = fmt.Sprintf("%s_queue_%s", appName, _queue)
			nsq.RegisterConsumer(topic, "1", jobHandler, cast.ToInt(consumer["processes"]))
		}
	}
}
