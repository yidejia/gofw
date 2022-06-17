package sms

import (
	"encoding/json"
	aliyunsmsclient "github.com/KenmyZhang/aliyun-communicate"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/logger"
)

// Aliyun 实现 sms.Driver interface
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 18:03
// @copyright © 2010-2022 广州伊的家网络科技有限公司
type Aliyun struct{}

func NewAliyun() *Aliyun {
	return &Aliyun{}
}

// ReadConfig 读取配置
func (s *Aliyun) ReadConfig() map[string]string {
	return config.GetStringMapString("sms.drivers.aliyun")
}

// GetMessageTemplate 获取消息模板
func (s *Aliyun) GetMessageTemplate() string {
	return config.GetString("sms.drivers.aliyun.template_code")
}

// HandleVerifyCode 处理验证码
func (s *Aliyun) HandleVerifyCode(code string) map[string]string {
	return map[string]string{"code": code}
}

// Send 实现 sms.Driver interface 的 Send 方法
func (s *Aliyun) Send(phone string, message Message, config map[string]string) bool {

	smsClient := aliyunsmsclient.New("http://dysmsapi.aliyuncs.com/")

	templateParam, err := json.Marshal(message.Data)
	if err != nil {
		logger.ErrorString("短信[阿里云]", "解析绑定错误", err.Error())
		return false
	}

	logger.DebugJSON("短信[阿里云]", "配置信息", config)

	result, err := smsClient.Execute(
		config["access_key_id"],
		config["access_key_secret"],
		phone,
		config["sign_name"],
		message.Template,
		string(templateParam),
	)

	logger.DebugJSON("短信[阿里云]", "请求内容", smsClient.Request)
	logger.DebugJSON("短信[阿里云]", "接口响应", result)

	if err != nil {
		logger.ErrorString("短信[阿里云]", "发信失败", err.Error())
		return false
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		logger.ErrorString("短信[阿里云]", "解析响应 JSON 错误", err.Error())
		return false
	}

	if result.IsSuccessful() {
		logger.DebugString("短信[阿里云]", "发信成功", "")
		return true
	} else {
		logger.ErrorString("短信[阿里云]", "服务商返回错误", string(resultJSON))
		return false
	}
}