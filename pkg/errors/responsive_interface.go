// Package errors 错误包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-07 10:27
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package errors

// ResponsiveError 可响应错误接口
type ResponsiveError interface {
	// HttpStatus 返回 http 状态码
	HttpStatus() int
	// Message 返回错误信息
	Message() string
	// Error 返回内部错误对象
	Error() error
	// Errors 返回多个错误信息的映射
	Errors() map[string][]string
}

// IsBadRequest 是请求资源格式错误
func IsBadRequest(err interface{}) bool {
	_, ok := err.(*errorBadRequest)
	return ok
}

// IsUnauthorized 是用户未授权错误
func IsUnauthorized(err interface{}) bool {
	_, ok := err.(*errorUnauthorized)
	return ok
}

// IsForbidden 是无权请求错误
func IsForbidden(err interface{}) bool {
	_, ok := err.(*errorForbidden)
	return ok
}

// IsNotFound 是请求资源不存在错误
func IsNotFound(err interface{}) bool {
	_, ok := err.(*errorNotFound)
	return ok
}

// IsMethodNotAllowed 是请求方法不允许错误
func IsMethodNotAllowed(err interface{}) bool {
	_, ok := err.(*errorMethodNotAllowed)
	return ok
}

// IsUnprocessableEntity 是资源无法处理错误
func IsUnprocessableEntity(err interface{}) bool {
	_, ok := err.(*errorUnprocessableEntity)
	return ok
}

// IsLocked 是资源已锁定错误
func IsLocked(err interface{}) bool {
	_, ok := err.(*errorLocked)
	return ok
}

// IsInternal 是内部错误
func IsInternal(err interface{}) bool {
	_, ok := err.(*errorInternal)
	return ok
}

// IsCustom 是自定义错误
func IsCustom(err interface{}) bool {
	_, ok := err.(*errorCustom)
	return ok
}
