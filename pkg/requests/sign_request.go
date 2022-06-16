package requests

import (
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/thedevsaddam/govalidator"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	gfErrors "github.com/yidejia/gofw/pkg/errors"
	"github.com/yidejia/gofw/pkg/hash"
	"github.com/yidejia/gofw/pkg/helpers"
	"github.com/yidejia/gofw/pkg/maptool"
	"time"
)

// Signatable 请求可签名接口
type Signatable interface {
	// ParamsToSign 返回用于签名验证的请求参数
	ParamsToSign() map[string]interface{}
}

// SignAppSecretReader 签名密钥读取器接口
type SignAppSecretReader interface {
	// ReadAppSecret 读取签名密钥
	ReadAppSecret(appKey string) (appSecret string, err gfErrors.ResponsiveError)
}

// SignRequest 签名请求
type SignRequest struct {
	Request
	AppKey    string `json:"app_key" form:"app_key" valid:"app_key"`          // 应用 key
	RandomStr string `json:"random_str" form:"random_str" valid:"random_str"` // 随机字符串
	Timestamp int64 `json:"timestamp" form:"timestamp" valid:"timestamp"`     // 时间戳
	Sign      string `json:"sign" form:"sign" valid:"sign"`                   // 签名
}

func NewSignRequest() *SignRequest {
	return &SignRequest{}
}

// ReadAppSecret 读取签名密钥
// 框架默认实现，一般应用给自己的请求签名时使用，给其他应用的请求签名时，需要另外实现
func (req *SignRequest) ReadAppSecret(appKey string) (appSecret string, err gfErrors.ResponsiveError) {
	appSecret = config.Get("app.secret")
	return
}

// Validate 验证请求
func (req *SignRequest) Validate(extra ...interface{}) map[string][]string  {

	rules := govalidator.MapData{
		"app_key":    []string{"required", "min:2"},
		"random_str": []string{"required", "len:10"},
		"timestamp":  []string{"required", "digits:10"},
		"sign":       []string{"required", "len:32"},
	}

	messages := govalidator.MapData{
		"app_key": []string{
			"required:应用名为必填项",
			"min:应用名长度需大于 2",
		},
		"random_str": []string{
			"required:随机字符串为必填项",
			"min:随机字符串长度需等于 10",
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

// ValidateSign 验证请求签名
func (req *SignRequest) ValidateSign(params map[string]interface{}, sign string, errs map[string][]string) map[string][]string {
	return ValidateSign(params, sign, errs)
}

// makeParamString 生成参数字符串
func makeParamString(params map[string]interface{}, secretReader ...SignAppSecretReader) (paramString string, newParams map[string]interface{}, err gfErrors.ResponsiveError) {

	// 按顺序拼接参数名和参数值
	paramString = ""
	// 对参数按字典序排序
	paramNames := maptool.SortIndictOrder(params)

	if len(paramNames) > 0 {

		appKey, ok := params["app_key"]
		if !ok || appKey == "" {
			// 未设置应用 key，说明是当前应用发起的签名，默认读取应用配置里的数据
			appKey = config.Get("app.key")
		}

		var appSecret string
		// 设置了应用密钥读取器
		if len(secretReader) > 0 {
			if appSecret, err = secretReader[0].ReadAppSecret(cast.ToString(appKey)); err != nil {
				return
			}
		} else {
			if appSecret, err = NewSignRequest().ReadAppSecret(cast.ToString(appKey)); err != nil {
				return
			}
		}

		randomStr, ok := params["random_str"]
		if !ok || randomStr == "" {
			// 未设置随机字符串，说明是当前应用发起的签名，随机生成一个10位字符串
			randomStr = helpers.RandomString(10)
		}

		timestamp, ok := params["timestamp"]
		if !ok || cast.ToInt64(timestamp) <= 0 {
			// 未设置时间戳，说明是当前应用发起的签名，获取当前时间戳
			timestamp = app.TimenowInTimezone().Unix()
		}

		// 填充标准参数
		newParams = make(map[string]interface{})
		newParams["app_key"] = appKey
		newParams["random_str"] = randomStr
		newParams["timestamp"] = timestamp

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
func MakeSign(params map[string]interface{}, secretReader ...SignAppSecretReader) (sign string, newParams map[string]interface{}, err gfErrors.ResponsiveError) {

	sign = ""

	var paramString string
	paramString, newParams, err = makeParamString(params, secretReader...)
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
func CheckSign(params map[string]interface{}, sign string, secretReader ...SignAppSecretReader) (ok bool, newParams map[string]interface{}, err gfErrors.ResponsiveError) {

	if len(params) == 0 {
		ok = false
		err = gfErrors.NewErrorBadRequest(errors.New("请求参数不能为空"), "请求参数不能为空")
		return
	}

	paramString, newParams, err := makeParamString(params, secretReader...)
	if err != nil {
		ok = false
		return
	}

	if paramString == "" {
		ok = false
		err = gfErrors.NewErrorInternal(errors.New("检查请求签名失败"), "检查请求签名失败")
		return
	}

	return hash.Md5(hash.Md5(paramString)) == sign, newParams, err
}

// ValidateSign 验证请求签名
func ValidateSign(params map[string]interface{}, sign string, errs map[string][]string) map[string][]string {
	ok, _, err := CheckSign(params, sign)
	if err != nil {
		errs["sign"] = append(errs["sign"], err.Message())
	}
	if !ok {
		errs["sign"] = append(errs["sign"], "请求签名无效")
	}
	// 请求签名 15 分钟内有效
	if cast.ToInt64(params["timestamp"]) < app.TimenowInTimezone().Add(-(time.Duration(15) * time.Minute)).Unix() {
		errs["sign"] = append(errs["sign"], "请求签名已过期")
	}
	return errs
}
