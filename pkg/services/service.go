// Package services 业务服务包，封装业务处理相关代码
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-07 11:47
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package services

import (
	gfErrors "github.com/yidejia/gofw/pkg/errors"
	"github.com/yidejia/gofw/pkg/logger"
)

// Service 业务服务基类
type Service struct {
}

// ErrorBadRequest 返回请求格式不正确错误
// 没有内部错误对象需要返回时，err 可以设置为 nil
func (svc *Service) ErrorBadRequest(err error, message ...string) gfErrors.ResponsiveError {
	logger.LogIf(err)
	return gfErrors.NewErrorBadRequest(err, message...)
}

// ErrorUnauthorized 返回用户未授权错误
func (svc *Service) ErrorUnauthorized(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorUnauthorized(message...)
}

// ErrorForbidden 返回无权访问错误
func (svc *Service) ErrorForbidden(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorForbidden(message...)
}

// ErrorNotFound 返回资源不存在错误
func (svc *Service) ErrorNotFound(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorNotFound(message...)
}

// ErrorMethodNotAllowed 返回请求方法不允许错误
func (svc *Service) ErrorMethodNotAllowed(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorMethodNotAllowed(message...)
}

// ErrorUnprocessableEntity 返回请求方法不允许错误
// 不需要返回多个错误信息映射时，errors 可以设置为 nil
func (svc *Service) ErrorUnprocessableEntity(errors map[string][]string, message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorUnprocessableEntity(errors, message...)
}

// ErrorLocked 返回资源已锁定错误
func (svc *Service) ErrorLocked(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorLocked(message...)
}

// ErrorInternal 返回系统内部错误
// 没有内部错误对象需要返回时，err 可以设置为 nil
func (svc *Service) ErrorInternal(err error, message ...string) gfErrors.ResponsiveError {
	logger.LogIf(err)
	return gfErrors.NewErrorInternal(err, message...)
}
