package requests

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/thedevsaddam/govalidator"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/hash"
	"github.com/yidejia/gofw/pkg/maptool"
	"time"
)

// SignRequest 签名请求
type SignRequest struct {
	Request
	App string `json:"app" form:"app" valid:"app"`
	Timestamp int `json:"timestamp" form:"timestamp" valid:"timestamp"`
	Sign string `json:"sign" form:"sign" valid:"sign"`
}

// Validate 验证请求
func (req *SignRequest) Validate(extra ...interface{}) map[string][]string  {

	rules := govalidator.MapData{
		"app":              []string{"required", "min:1"},
		"timestamp":        []string{"required", "digits:10"},
		"sign":             []string{"required", "len:60"},
	}

	messages := govalidator.MapData{
		"app": []string{
			"required:应用名为必填项",
			"min:应用名长度需大于 1",
		},
		"timestamp": []string{
			"required:时间戳为必填项",
			"digits:时间戳需为10位整数",
		},
		"sign": []string{
			"required:请求签名为必填项",
			"len:请求签名长度需要为 60",
		},
	}

	return req.ValidateStruct(req, rules, messages)
}

// ToMap 从请求中提取数据生成映射
func (req *SignRequest) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"app":       req.App,
		"timestamp": req.Timestamp,
	}
}

// ValidateSign 验证请求签名
func (req *SignRequest) ValidateSign(params map[string]interface{}, sign string, errs map[string][]string) map[string][]string {
	return ValidateSign(params, sign, errs)
}

// makeParamString 生成参数字符串
func makeParamString(params map[string]interface{}) (paramString string) {
	// 对参数名按字典序排序
	paramNames := maptool.SortIndictOrder(params)
	// 按顺序拼接参数名和参数值值
	paramString = ""
	if len(paramNames) > 0 {
		isFirstParam := true
		for _, paramName := range paramNames {
			if isFirstParam {
				paramString = paramString + fmt.Sprintf("%s=%s", paramName, cast.ToString(params[paramName]))
				isFirstParam = false
			} else {
				paramString = paramString + fmt.Sprintf("&%s=%s", paramName, cast.ToString(params[paramName]))
			}
		}
	}
	return
}

// MakeSign 生成请求签名
func MakeSign(params map[string]interface{}) (sign string) {
	sign = ""
	paramString := makeParamString(params)
	if paramString != "" {
		sign = hash.BcryptHash(paramString)
	}
	return
}

// CheckSign 检查请求签名
func CheckSign(params map[string]interface{}, sign string) bool {
	if len(params) == 0 {
		return false
	}
	paramString := makeParamString(params)
	return hash.BcryptCheck(paramString, sign)
}

// ValidateSign 验证请求签名
func ValidateSign(params map[string]interface{}, sign string, errs map[string][]string) map[string][]string {
	if !CheckSign(params, sign) {
		errs["sign"] = append(errs["sign"], "请求签名无效")
	}
	// 请求签名 15 分钟内有效
	if (cast.ToInt64(params["timestamp"]) + cast.ToInt64(time.Minute * 15)) < app.TimenowInTimezone().Unix() {
		errs["sign"] = append(errs["sign"], "请求签名已过期")
	}
	return errs
}
