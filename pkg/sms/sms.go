// Package sms 发送短信
package sms

import (
	"fmt"
	"github.com/yidejia/gofw/pkg/config"
	"sync"
)

// Message 是短信的结构体
type Message struct {
	Template string
	Data     map[string]string
	Content  string
}

// SMS 是我们发送短信的操作类
type SMS struct {
	Driver Driver
}

// once 单例模式
var once sync.Once

// internalSMS 内部使用的 SMS 对象
var internalSMS *SMS

var drivers = make(map[string]DriverFunc)

// NewSMS 单例模式获取
func NewSMS() *SMS {
	once.Do(func() {
		driverName := config.Get("sms.driver")
		driverFunc, ok := drivers[driverName]
		if ok {
			internalSMS = &SMS{
				Driver: driverFunc(),
			}
		} else {
			panic(fmt.Sprintf("sms does not supported %s driver", driverName))
		}
	})
	return internalSMS
}

// RegisterDriver 注册短信驱动
func RegisterDriver(driverName string, driverFunc DriverFunc) {
	drivers[driverName] = driverFunc
}

// GetMessageTemplate 获取消息模板
func (sms *SMS) GetMessageTemplate() string {
	return sms.Driver.GetMessageTemplate()
}

// HandleVerifyCode 处理验证码
func (sms *SMS) HandleVerifyCode(code string) map[string]string {
	return sms.Driver.HandleVerifyCode(code)
}

// Send 发送短信
func (sms *SMS) Send(phone string, message Message) bool {
	return sms.Driver.Send(phone, message, sms.Driver.ReadConfig())
}

// SendString 发送纯文本短信
func (sms *SMS) SendString(phone string, message string) bool {
	return sms.Send(phone, Message{Content: message})
}