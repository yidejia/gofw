// Package services 应用业务服务包，封装业务处理逻辑相关代码
package services

import (
	gfSvc "github.com/yidejia/gofw/pkg/services"
)

// Service 应用业务服务基类，内嵌了框架业务服务基类，可以根据应用需要进行扩展
type Service struct {
	gfSvc.Service
}
