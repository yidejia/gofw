// Package requests 请求包
// 处理请求数据和表单验证
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 16:54
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"github.com/yidejia/gofw/pkg/auth"
	"github.com/yidejia/gofw/pkg/db"
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

// MapConverter 映射转换器接口
// 一般用在更新模型时，从请求中提取部分字段，用于更新模型的部分字段
type MapConverter interface {
	// ToMap 将请求结构体转换成映射
	ToMap() map[string]interface{}
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

// Bind 绑定请求数据到结构体
func Bind(c *gin.Context, data interface{}) bool {
	// 解析请求，支持 JSON 数据、表单请求和 URL Query
	if err := c.ShouldBind(data); err != nil {
		response.BadRequest(c, err, "请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。")
		return false
	}
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