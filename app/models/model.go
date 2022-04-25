// Package models 存放模型的包
package models

import (
	gfModels "github.com/yidejia/gofw/pkg/models"
	"time"
)

// Model 应用模型基类，内嵌了框架的模型基类，可以根据应用的需要进行扩展
type Model struct {
	gfModels.Model
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
}
