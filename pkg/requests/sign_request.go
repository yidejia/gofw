package requests

import (
	"errors"
	"fmt"
	"time"

	"github.com/thedevsaddam/govalidator"

	"github.com/yidejia/gofw/pkg/logger"

	"github.com/spf13/cast"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/hash"
	"github.com/yidejia/gofw/pkg/helpers"
	"github.com/yidejia/gofw/pkg/maptool"
)

// Signable 请求可签名接口
type Signable interface {
	// ParamsToSign 返回用于签名验证的请求参数
	ParamsToSign() map[string]interface{}
}

// SignSecretFunc 根据应用 key 获取签名密钥函数
type SignSecretFunc func(appKey string) (string, error)

// SignOptions 签名选项
type SignOptions struct {
	Secret              string         // 签名密钥
	SecretFunc          SignSecretFunc // 根据应用 key 获取签名密码函数
	ErrorMessage        string         // 签名操作内部错误信息
	InvalidErrorMessage string         // 签名无效错误信息
	ExpireTime          int64          // 签名有效期，单位分钟
	ExpiredErrorMessage string         // 签名过期错误信息
}

// SignOption 签名选项设置函数
type SignOption func(*SignOptions)

// SignRequest 签名请求
type SignRequest struct {
	Request
	AppKey    string `json:"app_key" form:"app_key" valid:"app_key"`          // 应用 key
	RandomStr string `json:"random_str" form:"random_str" valid:"random_str"` // 随机字符串
	Timestamp int64  `json:"timestamp" form:"timestamp" valid:"timestamp"`    // 时间戳
	Sign      string `json:"sign" form:"sign" valid:"sign"`                   // 签名
}

// signSecretFunc 获取请求签名密钥的函数
var signSecretFunc = func(appKey string) (string, error) {
	// 是当前应用就返回自己的密钥
	if appKey == config.Get("app.key") {
		return config.Get("app.secret"), nil
	}
	return "", nil
}

// SetSignSecretFunc 设置获取请求签名密钥的函数
func SetSignSecretFunc(f SignSecretFunc) {
	signSecretFunc = f
}

// Validate 验证请求
func (req *SignRequest) Validate(extra ...interface{}) map[string][]string {

	rules := govalidator.MapData{
		"app_key":    []string{"required", "min:2"},
		"random_str": []string{"required", "len:10"},
		"timestamp":  []string{"required", "digits:10"},
		"sign":       []string{"required", "len:32"},
	}

	messages := govalidator.MapData{
		"app_key": []string{
			"required:应用 key 为必填项",
			"min:应用 key 长度需大于 2",
		},
		"random_str": []string{
			"required:随机字符串为必填项",
			"len:随机字符串长度需等于 10",
		},
		"timestamp": []string{
			"required:时间戳为必填项",
			"digits:时间戳需为10位整数",
		},
		"sign": []string{
			"required:请求签名为必填项",
			"len:请求签名长度需要为 32",
		},
	}

	return req.ValidateStruct(req, rules, messages)
}

// ParamsToSign 返回用于签名验证的请求参数
func (req *SignRequest) ParamsToSign() map[string]interface{} {
	return map[string]interface{}{
		"app_key":    req.AppKey,
		"random_str": req.RandomStr,
		"timestamp":  req.Timestamp,
	}
}

// NewSignRequest 新建签名请求
func NewSignRequest() *SignRequest {
	return &SignRequest{}
}

// NewSignOptions 新建签名选项
func (req *SignRequest) NewSignOptions(options ...SignOption) *SignOptions {
	signOptions := &SignOptions{
		Secret:              config.Get("app.secret"), // 默认使用当前应用自己的密钥进行签名
		SecretFunc:          signSecretFunc,
		ErrorMessage:        "",                                                       // 默认直接返回内部错误信息，可能过于技术语言，信息要直达普通用户时最好自定义
		InvalidErrorMessage: "请求签名无效",                                                 // 默认签名无效时的请求信息，可能过于技术语言，信息要直达普通用户时最好自定义
		ExpireTime:          cast.ToInt64(config.Get("app.api_sign_expire_time", 15)), // 签名有效期默认15分钟
		ExpiredErrorMessage: "请求签名已过期",                                                // 使用默认错误信息
	}
	// 调用选项函数设置各个选项
	for _, option := range options {
		option(signOptions)
	}
	return signOptions
}

// WithSecret 设置签名密钥
func (req *SignRequest) WithSecret(secret string) SignOption {
	return func(options *SignOptions) {
		options.Secret = secret
	}
}

// WithSecretFunc 设置根据应用 key 获取签名密钥函数
func (req *SignRequest) WithSecretFunc(secretFunc SignSecretFunc) SignOption {
	return func(options *SignOptions) {
		options.SecretFunc = secretFunc
	}
}

// WithErrorMessage 设置签名操作内部错误信息
func (req *SignRequest) WithErrorMessage(message string) SignOption {
	return func(options *SignOptions) {
		options.ErrorMessage = message
	}
}

// WithInvalidErrorMessage 设置签名无效时的错误信息
func (req *SignRequest) WithInvalidErrorMessage(message string) SignOption {
	return func(options *SignOptions) {
		options.InvalidErrorMessage = message
	}
}

// WithExpireTime 设置签名有效期，单位分钟
func (req *SignRequest) WithExpireTime(expireTime int64) SignOption {
	return func(options *SignOptions) {
		options.ExpireTime = expireTime
	}
}

// WithExpiredErrorMessage 设置签名过期错误信息
func (req *SignRequest) WithExpiredErrorMessage(message string) SignOption {
	return func(options *SignOptions) {
		options.ExpiredErrorMessage = message
	}
}

// ValidateSign 验证请求签名
func (req *SignRequest) ValidateSign(params map[string]interface{}, sign string, errs map[string][]string, options *SignOptions) map[string][]string {
	return ValidateSign(params, sign, errs, options)
}

// makeParamString 生成参数字符串
func makeParamString(params map[string]interface{}, options *SignOptions) (paramString string, newParams map[string]interface{}, err error) {

	// 按顺序拼接参数名和参数值
	paramString = ""
	// 对参数名按字典序排序
	paramNames := maptool.SortIndictOrder(params)

	if len(paramNames) > 0 {

		appKey, ok := params["app_key"]
		if !ok || appKey == "" {
			// 未设置应用 key，说明是当前应用发起的签名，默认读取应用配置里的数据
			appKey = config.Get("app.key")
		}

		var appSecret string
		// 设置了获取签名密钥函数
		if options.SecretFunc != nil {
			if appSecret, err = options.SecretFunc(cast.ToString(appKey)); err != nil {
				return
			}
		} else if len(options.Secret) > 0 {
			// 设置了签名密钥
			appSecret = options.Secret
		} else if appKey == config.Get("app.key") {
			// 是当前应用时，可以使用自己的密钥
			appSecret = config.Get("app.secret")
		} else {
			err = errors.New("没有设置签名密钥")
			return
		}

		randomStr, ok := params["random_str"]
		if !ok || randomStr == "" {
			// 未设置随机字符串，说明是当前应用发起的签名，随机生成一个10位字符串
			randomStr = helpers.RandomString(10)
		}

		timestamp, ok := params["timestamp"]
		if !ok || cast.ToInt64(timestamp) <= 0 {
			// 未设置时间戳，说明是当前应用发起的签名，获取当前时间戳
			timestamp = app.TimeNowInTimezone().Unix()
		}

		// 填充标准参数
		newParams = map[string]interface{}{
			"app_key":    appKey,
			"random_str": randomStr,
			"timestamp":  timestamp,
		}

		paramString = fmt.Sprintf("app_key=%s&app_secret=%s&random_str=%s&timestamp=%d", appKey, appSecret, randomStr, timestamp)

		for _, paramName := range paramNames {
			// 收集原来的请求参数
			newParams[paramName] = params[paramName]
			// 过滤已经完成拼接的参数
			if paramName != "app_key" && paramName != "app_secret" && paramName != "random_str" && paramName != "timestamp" {
				paramString = paramString + fmt.Sprintf("&%s=%s", paramName, cast.ToString(params[paramName]))
			}
		}
	}

	return
}

// MakeSign 生成请求签名
func MakeSign(params map[string]interface{}, options *SignOptions) (sign string, newParams map[string]interface{}, err error) {

	sign = ""

	var paramString string
	paramString, newParams, err = makeParamString(params, options)
	if err != nil {
		return
	}

	if paramString != "" {
		sign = hash.Md5(hash.Md5(paramString))
		newParams["sign"] = sign
	}

	return
}

// CheckSign 检查请求签名
func CheckSign(params map[string]interface{}, sign string, options *SignOptions) (ok bool, newParams map[string]interface{}, err error) {

	if len(params) == 0 {
		ok = false
		err = errors.New("请求参数不能为空")
		return
	}

	paramString, newParams, err := makeParamString(params, options)
	if err != nil {
		ok = false
		return
	}

	if paramString == "" {
		ok = false
		err = errors.New("检查请求签名失败")
		return
	}

	return hash.Md5(hash.Md5(paramString)) == sign, newParams, err
}

// ValidateSign 验证请求签名
func ValidateSign(params map[string]interface{}, sign string, errs map[string][]string, options *SignOptions) map[string][]string {
	ok, _, err := CheckSign(params, sign, options)
	if err != nil {
		if len(options.ErrorMessage) > 0 {
			logger.ErrorString("请求签名", "验证请求签名", err.Error())
			errs["sign"] = append(errs["sign"], options.ErrorMessage)
		} else {
			errs["sign"] = append(errs["sign"], err.Error())
		}
	}
	if !ok {
		errs["sign"] = append(errs["sign"], options.InvalidErrorMessage)
	}
	// 请求签名 15 分钟内有效
	if cast.ToInt64(params["timestamp"]) < app.TimeNowInTimezone().Add(-(time.Duration(options.ExpireTime) * time.Minute)).Unix() {
		errs["sign"] = append(errs["sign"], options.ExpiredErrorMessage)
	}
	return errs
}
