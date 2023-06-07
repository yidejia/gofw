// Package sms 短信包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 18:02
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package sms

// Driver 短信接口
type Driver interface {
	// ReadConfig 读取配置
	ReadConfig() map[string]string
	// GetMessageTemplate 获取消息模板，不支持的短信驱动直接返回空字符串即可
	GetMessageTemplate() string
	// HandleVerifyCode 处理验证码，不支持的短信驱动直接返回空映射即可
	HandleVerifyCode(code string) map[string]string
	// BeforeSend 发送短信前
	BeforeSend(phone string, message *Message, config map[string]string)
	// Send 发送短信
	Send(phone string, message *Message, config map[string]string) bool
}

// DriverFunc 返回短信驱动的函数
type DriverFunc func() Driver
