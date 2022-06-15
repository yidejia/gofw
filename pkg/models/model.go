// Package models 模型包
// 定义业务实体数据结构，应该重点关注结构定义，对模型的操作应该在数据仓库类为实现
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:37
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/spf13/cast"
	"github.com/yidejia/gofw/pkg/config"
	"time"
)

// JSONTime 用于 JSON 数据的时间
type JSONTime struct {
	time.Time
}

// MarshalJSON 定义 JSON 数据中的时间格式
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// Value 往数据库插入数据时调用
func (t JSONTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	// 判断给定时间是否和默认零时间的时间戳相同
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 将 time.Time 值转换成 JSONTime 值
func (t *JSONTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JSONTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// Model 模型基类
type Model struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 通用时间戳
type CommonTimestampsField struct {
	CreatedAt JSONTime `gorm:"column:created_at;type:timestamp NULL;comment:'创建时间';" json:"created_at,omitempty"`
	UpdatedAt JSONTime `gorm:"column:updated_at;type:timestamp NULL;comment:'更新时间';" json:"updated_at,omitempty"`
}

// DeletedAtTimestampsField 删除时间戳
// 一般用于软删除
type DeletedAtTimestampsField struct {
	DeletedAt *JSONTime `gorm:"column:deleted_at;type:timestamp NULL;index;comment:'删除时间';" json:"deleted_at,omitempty"`
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

// ToMap 将模型转换成映射
func (m Model) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id": m.ID,
	}
}