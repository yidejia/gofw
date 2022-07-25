// Package services 业务服务包
// 封装业务处理相关代码
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-07 11:47
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package services

import (
	gfErrs "github.com/yidejia/gofw/pkg/errors"
)

// Service 业务服务基类
type Service struct {
}

// NewErrorBadRequest 生成请求格式不正确错误
// 没有内部错误对象需要返回时，err 可以设置为 nil
func (svc *Service) NewErrorBadRequest(err error, message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorBadRequest(err, message...)
}

// NewErrorUnauthorized 生成用户未授权错误
func (svc *Service) NewErrorUnauthorized(message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorUnauthorized(message...)
}

// NewErrorForbidden 生成无权访问错误
func (svc *Service) NewErrorForbidden(message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorForbidden(message...)
}

// NewErrorNotFound 生成资源不存在错误
func (svc *Service) NewErrorNotFound(message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorNotFound(message...)
}

// NewErrorMethodNotAllowed 生成请求方法不允许错误
func (svc *Service) NewErrorMethodNotAllowed(message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorMethodNotAllowed(message...)
}

// NewErrorUnprocessableEntity 生成请求方法不允许错误
// 不需要返回多个错误信息映射时，errors 可以设置为 nil
func (svc *Service) NewErrorUnprocessableEntity(errors map[string][]string, message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorUnprocessableEntity(errors, message...)
}

// NewErrorLocked 生成资源已锁定错误
func (svc *Service) NewErrorLocked(message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorLocked(message...)
}

// NewErrorInternal 生成系统内部错误
// 没有内部错误对象需要返回时，err 可以设置为 nil
func (svc *Service) NewErrorInternal(err error, message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorInternal(err, message...)
}

// NewErrorServiceUnavailable 生成服务不可用错误
func (svc *Service) NewErrorServiceUnavailable(message ...string) gfErrs.ResponsiveError {
	return gfErrs.NewErrorServiceUnavailable(message...)
}
