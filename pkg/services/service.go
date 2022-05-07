// Package services 业务服务包，封装业务处理相关代码
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-07 11:47
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package services

import gfErrors "github.com/yidejia/gofw/pkg/errors"

// Service 业务服务基类
type Service struct {
}

// ErrorBadRequest 返回请求格式不正确错误
func (svc *Service) ErrorBadRequest(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorBadRequest(message...)
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
func (svc *Service) ErrorUnprocessableEntity(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorUnprocessableEntity(message...)
}

// ErrorLocked 返回资源已锁定错误
func (svc *Service) ErrorLocked(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorLocked(message...)
}

// ErrorInternal 返回系统内部错误
func (svc *Service) ErrorInternal(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorInternal(message...)
}
