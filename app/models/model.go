// Package models 应用模型包，封装业务实体定义相关代码
package models

import (
	gfModels "github.com/yidejia/gofw/pkg/models"
)

// Model 应用模型基类，内嵌了框架模型基类，可以根据应用需要进行扩展
type Model struct {
	gfModels.Model
}
