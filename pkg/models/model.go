// Package models 模型通用属性和方法
package models

import (
	"github.com/spf13/cast"
	"github.com/yidejia/gofw/pkg/config"
	"time"
)

// Model 模型基类
type Model struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
}

// GetStringID 获取 ID 的字符串格式
func (m Model) GetStringID() string {
	return cast.ToString(m.ID)
}

// Connection 获取模型对应的数据库连接
func (m Model) Connection() string {
	// 返回默认的数据库连接
	return config.Get("database.default")
}

// ModelName 模型名称
func (m Model) ModelName() string {
	return "模型"
}