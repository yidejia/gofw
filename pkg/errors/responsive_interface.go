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
