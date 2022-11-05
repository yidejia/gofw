package requests

import (
	"fmt"
	"sort"
	"time"

	"github.com/yidejia/gofw/pkg/helpers"

	"github.com/thedevsaddam/govalidator"

	"github.com/yidejia/gofw/pkg/logger"

	"github.com/spf13/cast"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/hash"
)

// SignAble 可签名请求接口
type SignAble interface {
	// ParamsToSign 返回用于签名验证的请求参数
	ParamsToSign() map[string]interface{}
}

// AppSecretFunc 根据应用 key 获取应用密钥函数
type AppSecretFunc func(appKey string) (string, error)

// SignOptions 签名选项
type SignOptions struct {
	AppKey              string        // 应用 key
	AppSecret           string        // 签名密钥
	AppSecretFunc       AppSecretFunc // 根据应用 key 获取应用密码函数
	Timestamp           int64         // 时间戳
	RandomStr           string        // 随机字符串
	ErrorMessage        string        // 签名操作内部错误信息
	InvalidErrorMessage string        // 签名无效错误信息
	ExpireTime          int64         // 签名有效期，单位分钟
	ExpiredErrorMessage string        // 签名过期错误信息
	Lazy                bool          // 延迟初始化
	LazyOptions         []SignOption  // 延迟初始化的选项
}

// SignOption 签名选项设置函数
type SignOption func(*SignOptions)

// SignRequest 签名请求
type SignRequest struct {
	Request
	AppKey    string `json:"app_key" form:"app_key" valid:"app_key"`          // 应用 key
	Timestamp int64  `json:"timestamp" form:"timestamp" valid:"timestamp"`    // 时间戳
	TS        int64  `json:"ts" form:"ts" valid:"ts"`                         // 兼容的时间戳字段
	RandomStr string `json:"random_str" form:"random_str" valid:"random_str"` // 随机字符串
	Sign      string `json:"sign" form:"sign" valid:"sign"`                   // 签名
}

// appSecretFunc 获取应用密钥的函数
var appSecretFunc = func(appKey string) (string, error) {
	// 是当前应用就返回自己的密钥
	if appKey == config.Get("app.key") {
		return config.Get("app.secret"), nil
	}
	return "", nil
}

// SetAppSecretFunc 设置获取应用密钥的函数
func SetAppSecretFunc(f AppSecretFunc) {
	appSecretFunc = f
}

// Validate 验证请求
func (req *SignRequest) Validate(extra ...interface{}) map[string][]string {

	// 时间戳进行参数名兼容处理
	if req.Timestamp == 0 && req.TS > 0 {
		req.Timestamp = req.TS
	}

	rules := govalidator.MapData{
		"app_key":    []string{"required", "min:2"},
		"timestamp":  []string{"required", "digits:10"},
		"random_str": []string{"required", "len:10"},
		"sign":       []string{"required", "len:32"},
	}

	messages := govalidator.MapData{
		"app_key": []string{
			"required:应用 key 为必填项",
			"min:应用 key 长度需大于 2",
		},
		"timestamp": []string{
			"required:时间戳为必填项",
			"digits:时间戳需为10位整数",
		},
		"random_str": []string{
			"required:随机字符串为必填项",
			"len:随机字符串长度需等于 10",
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
		"timestamp":  req.Timestamp,
		"random_str": req.RandomStr,
	}
}

// NewLazySignOptions 新建延迟初始化的签名选项
func NewLazySignOptions(options ...SignOption) *SignOptions {
	return &SignOptions{
		Lazy:        true,
		LazyOptions: options,
	}
}

// NewSignOptions 新建签名选项
func NewSignOptions(options ...SignOption) *SignOptions {

	signOptions := &SignOptions{}

	// 调用选项函数设置各个选项
	for _, option := range options {
		option(signOptions)
	}

	// 设置默认时间戳
	if signOptions.Timestamp == 0 {
		signOptions.Timestamp = app.TimeNowInTimezone().Unix()
	}

	// 设置默认随机字符串
	if signOptions.RandomStr == "" {
		signOptions.RandomStr = helpers.RandomString(10)
	}

	// 设置默认应用 key
	if signOptions.AppKey == "" {
		signOptions.AppKey = config.Get("app.key")
	}

	// 设置默认应用密钥
	if signOptions.AppSecret == "" && signOptions.AppKey != "" {
		var err error
		// 使用选项设置的函数
		if signOptions.AppSecretFunc != nil {
			if signOptions.AppSecret, err = signOptions.AppSecretFunc(signOptions.AppKey); err != nil {
				logger.ErrorString("新建签名选项", "设置默认应用密钥(signOptions.AppSecretFunc)", err.Error())
			}
		} else if appSecretFunc != nil {
			// 使用应用默认的函数
			if signOptions.AppSecret, err = appSecretFunc(signOptions.AppKey); err != nil {
				logger.ErrorString("新建签名选项", "设置默认应用密钥(appSecretFunc)", err.Error())
			}
		} else {
			// 默认直接返回当前应用密钥
			signOptions.AppSecret = config.Get("app.secret")
		}
	}

	// 设置签名无效时的默认提示信息，可能过于技术语言，信息要直达普通用户时最好自定义
	if signOptions.InvalidErrorMessage == "" {
		signOptions.InvalidErrorMessage = "请求签名无效"
	}

	// 签名有效期默认15分钟
	if signOptions.ExpireTime == 0 {
		signOptions.ExpireTime = cast.ToInt64(config.Get("app.api_sign_expire_time", 15))
	}

	// 设置签名已过期默认提示信息
	if signOptions.ExpiredErrorMessage == "" {
		signOptions.ExpiredErrorMessage = "请求签名已过期"
	}

	return signOptions
}

// WithAppKey 设置应用 key
func WithAppKey(appKey string) SignOption {
	return func(options *SignOptions) {
		options.AppKey = appKey
	}
}

// WithAppSecret 设置应用密钥
func WithAppSecret(appSecret string) SignOption {
	return func(options *SignOptions) {
		options.AppSecret = appSecret
	}
}

// WithAppSecretFunc 设置根据应用 key 获取应用密钥函数
func WithAppSecretFunc(appSecretFunc AppSecretFunc) SignOption {
	return func(options *SignOptions) {
		options.AppSecretFunc = appSecretFunc
	}
}

// WithTimestamp 设置时间戳
func WithTimestamp(timestamp int64) SignOption {
	return func(options *SignOptions) {
		options.Timestamp = timestamp
	}
}

// WithRandomStr 设置随机字符串
func WithRandomStr(randomStr string) SignOption {
	return func(options *SignOptions) {
		options.RandomStr = randomStr
	}
}

// WithErrorMessage 设置签名操作内部错误信息
func WithErrorMessage(message string) SignOption {
	return func(options *SignOptions) {
		options.ErrorMessage = message
	}
}

// WithInvalidErrorMessage 设置签名无效时的错误信息
func WithInvalidErrorMessage(message string) SignOption {
	return func(options *SignOptions) {
		options.InvalidErrorMessage = message
	}
}

// WithExpireTime 设置签名有效期，单位分钟
func WithExpireTime(expireTime int64) SignOption {
	return func(options *SignOptions) {
		options.ExpireTime = expireTime
	}
}

// WithExpiredErrorMessage 设置签名过期错误信息
func WithExpiredErrorMessage(message string) SignOption {
	return func(options *SignOptions) {
		options.ExpiredErrorMessage = message
	}
}

// ValidateSign 验证请求签名
func (req *SignRequest) ValidateSign(params map[string]interface{}, sign string, options *SignOptions, errs map[string][]string, extra ...interface{}) map[string][]string {
	// 验证签名标准参数
	errs = req.MergeValidateErrors(errs, req.Validate(extra...))
	return ValidateSign(params, sign, options, errs)
}

// makeParamString 生成参数字符串
func makeParamString(params map[string]interface{}) (paramString string) {

	// 对映射键按字典序排序
	keys := make([]string, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {

		// 签名预设字段不参与拼接，它们将固定在字符串前缀部分
		if k == "app_key" || k == "app_secret" || k == "timestamp" || k == "random_str" || k == "sign" {
			continue
		}

		v := params[k]
		m, isMap := v.(map[string]interface{})
		s, isSlice := v.([]interface{})

		// 映射值是映射，递归处理
		if isMap {
			if len(m) > 0 {
				newParams := make(map[string]interface{}, len(m))
				for nk, mv := range m {
					newParams[fmt.Sprintf("%s:%s", k, nk)] = mv
				}
				paramString += makeParamString(newParams)
			}
		} else if isSlice {
			// 映射值是切片，将切片转换成映射后，递归处理
			if len(s) > 0 {
				newParams := make(map[string]interface{}, len(s))
				for i, sv := range s {
					newParams[fmt.Sprintf("%s:%d", k, i)] = sv
				}
				paramString += makeParamString(newParams)
			}
		} else {
			// 拼接映射键和值
			if k != "" && v != nil {
				paramString += fmt.Sprintf("&%s=%v", k, v)
			}
		}
	}

	return
}

// createSign 内部生成参数签名
func createSign(params map[string]interface{}, timestamp int64, randomStr, appKey, appSecret string) (sign string) {
	sign = hash.Md5(
		hash.Md5(
			fmt.Sprintf("app_key=%s&app_secret=%s&ts=%d&random_str=%s", appKey, appSecret, timestamp, randomStr) + makeParamString(params),
		),
	)
	return
}

// MakeSign 生成请求签名
func MakeSign(params map[string]interface{}, options *SignOptions) (sign string) {
	sign = createSign(params, options.Timestamp, options.RandomStr, options.AppKey, options.AppSecret)
	params["app_key"] = options.AppKey
	params["timestamp"] = options.Timestamp
	params["random_str"] = options.RandomStr
	params["sign"] = sign
	return
}

// CheckSign 检查请求签名
func CheckSign(params map[string]interface{}, sign string, options *SignOptions) bool {

	// 当前签名选项启用了延迟初始化，这个机制方便于自动从请求参数中提取签名所需参数
	if options.Lazy {
		// 从参数中提取签名选项
		options.LazyOptions = append(options.LazyOptions, WithTimestamp(cast.ToInt64(params["timestamp"])))
		options.LazyOptions = append(options.LazyOptions, WithRandomStr(cast.ToString(params["random_str"])))
		options.LazyOptions = append(options.LazyOptions, WithAppKey(cast.ToString(params["app_key"])))
		// 真正初始化签名选项
		*options = *NewSignOptions(options.LazyOptions...)
	}

	return MakeSign(params, options) == sign
}

// ValidateSign 验证请求签名
func ValidateSign(params map[string]interface{}, sign string, options *SignOptions, errs map[string][]string) map[string][]string {

	// 签名无效
	if ok := CheckSign(params, sign, options); !ok {
		errs["sign"] = append(errs["sign"], options.InvalidErrorMessage)
	}

	// 签名已过期
	if cast.ToInt64(params["timestamp"]) < app.TimeNowInTimezone().Add(-(time.Duration(options.ExpireTime) * time.Minute)).Unix() {
		errs["sign"] = append(errs["sign"], options.ExpiredErrorMessage)
	}

	return errs
}
