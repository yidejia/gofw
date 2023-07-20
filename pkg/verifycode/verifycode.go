// Package verifycode 用以发送手机验证码和邮箱验证码
package verifycode

import (
	"fmt"
	"strings"
	"sync"

	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/helpers"
	"github.com/yidejia/gofw/pkg/logger"
	"github.com/yidejia/gofw/pkg/mail"
	"github.com/yidejia/gofw/pkg/redis"
	"github.com/yidejia/gofw/pkg/sms"
)

// VerifyCode 操作对象
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 18:08
// @copyright © 2010-2022 广州伊的家网络科技有限公司
type VerifyCode struct {
	Store Store
}

var once sync.Once
var internalVerifyCode *VerifyCode

// NewVerifyCode 单例模式获取
func NewVerifyCode() *VerifyCode {
	once.Do(func() {
		internalVerifyCode = &VerifyCode{
			Store: &RedisStore{
				RedisClient: redis.Connection(),
				// 增加前缀保持数据库整洁，出问题调试时也方便
				KeyPrefix: config.GetString("app.name") + ":verifycode:",
			},
		}
	})

	return internalVerifyCode
}

// Option 验证码选项
type Option struct {
	Scene       uint8             `json:"scene"`        // 验证码场景码
	AppKey      string            `json:"app_key"`      // 使用验证码的应用
	SMSTemplate string            `json:"sms_template"` // 短信模板
	SMSData     map[string]string `json:"sms_data"`     // 短信数据
}

// OptionFunc 验证码选项设置函数
type OptionFunc func(option *Option)

// SendWithScene 设置验证码的场景码
func SendWithScene(scene uint8) OptionFunc {
	return func(option *Option) {
		option.Scene = scene
	}
}

// SendWithAppKey 设置使用验证码的应用
func SendWithAppKey(appKey string) OptionFunc {
	return func(option *Option) {
		option.AppKey = appKey
	}
}

// SendWithSMSTemplate 设置短信模板
func SendWithSMSTemplate(template string) OptionFunc {
	return func(option *Option) {
		option.SMSTemplate = template
	}
}

// SendWithSMSData 设置短信数据
func SendWithSMSData(data map[string]string) OptionFunc {
	return func(option *Option) {
		option.SMSData = data
	}
}

// NewOption 新建验证码选项
func NewOption(optionFunS ...OptionFunc) *Option {
	option := &Option{}
	for _, optionFun := range optionFunS {
		optionFun(option)
	}
	return option
}

// SendSMS 发送短信验证码
// 调用示例：
// verifycode.NewVerifyCode().SendSMS(request.Phone, request.Content)
func (vc *VerifyCode) SendSMS(phone, content string, options ...*Option) bool {

	// 获取验证码选项
	var option *Option
	if len(options) > 0 {
		// 设置了验证码选项
		option = options[0]
	} else {
		// 生成默认验证码选项
		option = NewOption()
	}

	// 生成验证码
	var code string
	// 设置了验证码的使用场景
	if option.Scene > 0 {
		// 设置了使用验证码的应用
		if len(option.AppKey) > 0 {
			// 生成只能用于某个应用某个场景的验证码
			code = vc.generateVerifyCode(fmt.Sprintf("%s:%d:%s", option.AppKey, option.Scene, phone))
		} else {
			// 生成只能用于某个场景的验证码
			code = vc.generateVerifyCode(fmt.Sprintf("%d:%s", option.Scene, phone))
		}
	} else {
		// 生成通用型验证码
		code = vc.generateVerifyCode(phone)
	}

	// 设置了短信模板时，将验证码填充到短信模板中
	if len(option.SMSTemplate) > 0 {
		option.SMSTemplate = strings.ReplaceAll(option.SMSTemplate, "{{code}}", code)
	} else {
		// 默认将验证码填充到短信内容中
		content = strings.ReplaceAll(content, "{{code}}", code)
	}

	logger.DebugString("验证码", "发送短信验证码内容", content)

	// 方便本地和 API 自动测试
	if !app.IsProduction() &&
		(config.GetBool("verifycode.debug_mode") ||
			strings.HasPrefix(phone, config.GetString("verifycode.debug_phone_prefix"))) {
		return true
	}

	_sms := sms.NewSMS()
	// 短信消息体
	smsMessage := &sms.Message{}

	// 设置了短信模板
	if len(option.SMSTemplate) > 0 {
		smsMessage.Template = option.SMSTemplate
	} else {
		// 默认由短信驱动提供短信模板
		smsMessage.Template = _sms.GetMessageTemplate()
	}

	// 设置了短信数据
	if len(option.SMSData) > 0 {
		smsMessage.Data = option.SMSData
		// 将验证码放进短信数据中传递给短信驱动处理
		smsMessage.Data["code"] = code
	} else {
		// 默认由短信驱动提供短信数据
		smsMessage.Data = _sms.HandleVerifyCode(code)
	}

	// 设置短信内容
	smsMessage.Content = content

	// 发送短信
	return _sms.Send(phone, smsMessage)
}

// CheckAnswer 检查用户提交的验证码是否正确，key 可以是手机号或者 Email
func (vc *VerifyCode) CheckAnswer(key string, answer string) bool {

	logger.DebugJSON("验证码", "检查验证码", map[string]string{key: answer})

	// 方便开发，在非生产环境下，具备特殊前缀的手机号和特殊 Email 后缀，直接验证成功
	if !app.IsProduction() &&
		// 开启调试模式下不会真正发送验证码，直接通过验证码检查
		(config.GetBool("verifycode.debug_mode") ||
			config.Get("verifycode.debug_code") == answer ||
			strings.HasPrefix(key, config.GetString("verifycode.debug_phone_prefix")) ||
			strings.HasSuffix(key, config.GetString("verifycode.debug_email_suffix"))) {
		return true
	}

	return vc.Store.Verify(key, answer, false)
}

// generateVerifyCode 生成验证码，并放置于 Redis 中
func (vc *VerifyCode) generateVerifyCode(key string) string {

	// 生成随机码
	code := helpers.RandomNumber(config.GetInt("verifycode.code_length"))

	// 调试模式下，使用固定验证码
	if !app.IsProduction() && config.GetBool("verifycode.debug_mode") {
		code = config.GetString("verifycode.debug_code")
	}

	logger.DebugJSON("验证码", "生成验证码", map[string]string{key: code})

	// 将验证码及 KEY（邮箱或手机号）存放到 Redis 中并设置过期时间
	vc.Store.Set(key, code)

	return code
}

// SendEmail 发送邮件验证码，调用示例：
//         verifycode.NewVerifyCode().SendEmail(request.Email)
func (vc *VerifyCode) SendEmail(email string) error {

	// 生成验证码
	code := vc.generateVerifyCode(email)

	// 方便本地和 API 自动测试
	if !app.IsProduction() && strings.HasSuffix(email, config.GetString("verifycode.debug_email_suffix")) {
		return nil
	}

	content := fmt.Sprintf("<h1>您的 Email 验证码是 %v </h1>", code)
	// 发送邮件
	mail.NewMailer().Send(mail.Email{
		From: mail.From{
			Address: config.GetString("mail.from.address"),
			Name:    config.GetString("mail.from.name"),
		},
		To:      []string{email},
		Subject: "Email 验证码",
		HTML:    []byte(content),
	})

	return nil
}
