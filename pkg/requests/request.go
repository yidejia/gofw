// Package requests 请求包
// 处理请求数据和表单验证
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 16:54
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"github.com/yidejia/gofw/pkg/auth"
	"github.com/yidejia/gofw/pkg/db"
	gfJSON "github.com/yidejia/gofw/pkg/json"
	"github.com/yidejia/gofw/pkg/maptool"
	"github.com/yidejia/gofw/pkg/response"
)

// Validatable 可验证接口
// 实现这个接口的对象可调用自身方法验证自己的数据
type Validatable interface {
	// Validate 对数据进行验证
	Validate(extra ...interface{}) map[string][]string
}

// ModelConverter 模型转换器接口
// 一般用在创建模型时，从请求中提取数据快速生成模型数据
type ModelConverter interface {
	// ToModel 将请求结构体转换成模型
	ToModel() db.IModel
}

// Request 请求基类
type Request struct {
}

// CurrentUID 从 gin.context 中获取当前登录用户 ID
func (req *Request) CurrentUID(c *gin.Context) string {
	return auth.CurrentUID(c)
}

// CurrentUser 获取当前登录用户
func (req *Request) CurrentUser(c *gin.Context) (user auth.Authenticate) {
	return auth.CurrentUser(c)
}

// Bind 绑定请求数据到结构体
func (req *Request) Bind(c *gin.Context, data interface{}) bool {
	return Bind(c, data)
}

// BindMap 绑定请求数据到映射
func (req *Request) BindMap(c *gin.Context, data *map[string]interface{}) bool {
	return BindMap(c, data)
}

// Validate 数据验证
func (req *Request) Validate(data Validatable, extra ...interface{}) map[string][]string {
	return Validate(data, extra...)
}

// BindAndValidate 绑定请求数据到结构体并进行数据验证
func (req *Request) BindAndValidate(c *gin.Context, data Validatable, extra ...interface{}) bool {
	return BindAndValidate(c, data, extra...)
}

// ValidateStruct 验证结构体
func (req *Request) ValidateStruct(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	return ValidateStruct(data, rules, messages)
}

// ValidateFile 验证文件
func (req *Request) ValidateFile(c *gin.Context, data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	return ValidateFile(c, data, rules, messages)
}

// MergeValidateErrors 合并验证错误信息
func (req *Request) MergeValidateErrors(errs map[string][]string, moreErrs ...map[string][]string) map[string][]string {
	return MergeValidateErrors(errs, moreErrs...)
}

// MergeParams 合并请求参数
func (req *Request) MergeParams(params map[string]interface{}, moreParams ...map[string]interface{}) map[string]interface{} {
	return MergeParams(params, moreParams...)
}

// ToMap 将请求结构体转换成映射
func (req *Request) ToMap(data interface{}) map[string]interface{} {
	return maptool.StructToMap(data)
}

// ToMapOnly 将请求结构体转换成映射并只返回指定的键值对
func (req *Request) ToMapOnly(data interface{}, keys ...string) map[string]interface{} {

	m := req.ToMap(data)

	if len(keys) > 0 {

		nm := map[string]interface{}{}

		for _, key := range keys {
			for k, v := range m {
				if k == key {
					nm[k] = v
				}
			}
		}

		m = nm
	}

	return m
}

// ToMapExcept 将请求结构体转换成映射并排除指定的键值对
func (req *Request) ToMapExcept(data interface{}, keys ...string) map[string]interface{} {

	m := req.ToMap(data)

	if len(keys) > 0 {

		nm := map[string]interface{}{}

		for _, key := range keys {
			for k, v := range m {
				if k != key {
					nm[k] = v
				}
			}
		}

		m = nm
	}

	return m
}

// LogRequest 记录请求日志
func (req *Request) LogRequest(c *gin.Context, reqData interface{}) error {
	return LogRequest(c, reqData)
}

// GetRequestLog 获取请求日志
func (req *Request) GetRequestLog(c *gin.Context) string {
	return GetRequestLog(c)
}

// LogResponse 记录响应日志
func (req *Request) LogResponse(c *gin.Context, respData interface{}) error {
	return LogResponse(c, respData)
}

// GetResponseLog 获取响应日志
func (req *Request) GetResponseLog(c *gin.Context) string {
	return GetResponseLog(c)
}

// Bind 绑定请求数据到结构体
func Bind(c *gin.Context, data interface{}) bool {
	// 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(data); err != nil {
		response.BadRequest(c, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return false
	}
	// 绑定成功
	return true
}

// BindMap 绑定请求数据到映射
func BindMap(c *gin.Context, data *map[string]interface{}) bool {

	// 读取请求正文内容并转换为字符串类型
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil && err != io.EOF {
		response.BadRequest(c, err, "读取请求正文内容错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return false
	}

	// 将请求体内容回填，以保证后续中间件和处理函数的正常工作
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	err = gfJSON.BindMap(string(reqBody), data)
	if err != nil && err != io.EOF {
		response.BadRequest(c, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return false
	}

	// 绑定成功
	return true
}

// Validate 数据验证
// 通过传入 extra 可变参数而不是直接传入 gin.Context 这个使用起来更方便的上下文参数进行辅助验证，是为了避免验证函数跟 gin.Context 上下文强绑定了，
// 这样将导致验证函数只能用于 http 请求的验证，无法复用到其他场景例如命令行或者 rpc 远程调用的数据验证
// @param data 待验证数据
// @param handler 验证函数
// @param extra 辅助验证的附加数据
// @return 映射 map，key 为请求参数名，value 为该参数多个错误信息组成的切片
func Validate(data Validatable, extra ...interface{}) map[string][]string {
	// 返回数据验证结果
	return data.Validate(extra...)
}

// BindAndValidate 绑定请求数据到结构体并进行数据验证
// 控制器里调用示例：
// request := requests.UserSaveRequest{}
// if ok := requests.BindAndValidate(c, &request); !ok {
//     return
// }
// 需要辅助验证的附加数据时，控制器里调用示例：
// request := requests.UserSaveRequest{}
// currentUser := auth.CurrentUser(c)
// if ok := requests.BindAndValidate(c, &request, currentUser); !ok {
//     return
// }
func BindAndValidate(c *gin.Context, data Validatable, extra ...interface{}) bool {

	if ok := Bind(c, data); !ok {
		return false
	}

	// 数据验证
	errs := Validate(data, extra...)

	// 请求过于频繁，限制客户端短时间内多次尝试破解类似验证码这种规则
	if c.GetBool("limiter-reached") {
		c.Set("limiter-reached", false)
		return false
	}

	// 验证不通过
	if len(errs) > 0 {
		response.ValidationError(c, errs, "")
		return false
	}

	return true
}

// ValidateStruct 验证结构体
func ValidateStruct(data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	// 配置选项
	opts := govalidator.Options{
		Data:          data,
		Rules:         rules,
		TagIdentifier: "valid", // 模型中的 Struct 标签标识符
		Messages:      messages,
	}
	// 开始验证
	return govalidator.New(opts).ValidateStruct()
}

// ValidateFile 验证文件
func ValidateFile(c *gin.Context, data interface{}, rules govalidator.MapData, messages govalidator.MapData) map[string][]string {
	opts := govalidator.Options{
		Request:       c.Request,
		Rules:         rules,
		Messages:      messages,
		TagIdentifier: "valid",
	}
	// 调用 govalidator 的 Validate 方法来验证文件
	return govalidator.New(opts).Validate()
}

// MergeValidateErrors 合并验证错误信息
func MergeValidateErrors(errs map[string][]string, moreErrs ...map[string][]string) map[string][]string {

	if len(moreErrs) > 0 {

		var moreErr map[string][]string
		var key string
		var value []string

		for _, moreErr = range moreErrs {
			for key, value = range moreErr {
				errs[key] = value
			}
		}
	}

	return errs
}

// MergeParams 合并请求参数
func MergeParams(params map[string]interface{}, moreParams ...map[string]interface{}) map[string]interface{} {

	if len(moreParams) > 0 {

		var moreParam map[string]interface{}
		var key string
		var value interface{}

		for _, moreParam = range moreParams {
			for key, value = range moreParam {
				params[key] = value
			}
		}
	}

	return params
}

// LogRequest 记录请求日志
func LogRequest(c *gin.Context, reqData interface{}) error {

	jsonBytes, err := json.Marshal(reqData)
	if err != nil {
		return errors.New("记录请求日志错误：" + err.Error())
	}

	c.Set("request_log", string(jsonBytes))

	return nil
}

// GetRequestLog 获取请求日志
func GetRequestLog(c *gin.Context) string {
	return c.GetString("request_log")
}

// LogResponse 记录响应日志
func LogResponse(c *gin.Context, respData interface{}) error {

	jsonBytes, err := json.Marshal(respData)
	if err != nil {
		return errors.New("记录响应日志错误：" + err.Error())
	}

	c.Set("response_log", string(jsonBytes))

	return nil
}

// GetResponseLog 获取响应日志
func GetResponseLog(c *gin.Context) string {
	return c.GetString("response_log")
}
