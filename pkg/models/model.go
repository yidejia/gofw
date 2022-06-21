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
	"github.com/yidejia/gofw/pkg/helpers"
	"strings"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

// JSONTime 用于 JSON 数据的时间
type JSONTime struct {
	time.Time
}

// MarshalJSON 实现编码 JSON数据接口
func (t JSONTime) MarshalJSON() ([]byte, error) {
	// 自定义 JSON 数据中的时间格式
	formatted := fmt.Sprintf("\"%v\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

// UnmarshalJSON 实现解码 JSON 数据接口
func (t *JSONTime) UnmarshalJSON(data []byte) error {

	if string(data) == "null" {
		return nil
	}

	var err error

	str := string(data)
	// 去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	_time, err := time.Parse(TimeFormat, timeStr)
	*t = JSONTime{Time: _time}

	return err
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

// String 将时间对象转换成字符串时调用
func (t *JSONTime) String() string {
	return fmt.Sprintf("hhh:%v", t.Time.String())
}

// Model 模型基类
type Model struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 通用时间戳
type CommonTimestampsField struct {
	CreatedAt JSONTime `gorm:"column:created_at;type:timestamp NULL;comment:创建时间;" json:"created_at,omitempty"`
	UpdatedAt JSONTime `gorm:"column:updated_at;type:timestamp NULL;comment:更新时间;" json:"updated_at,omitempty"`
}

// TimeToString 时间字段转换字符串
func (m CommonTimestampsField) TimeToString(field string) string {
	if field == "created_at" {
		createdAt, _ := m.CreatedAt.MarshalJSON()
		return  strings.Trim(string(createdAt), "\"")
	} else if field == "updated_at" {
		updatedAt, _ := m.UpdatedAt.MarshalJSON()
		return strings.Trim(string(updatedAt), "\"")
	} else {
		return ""
	}
}

// DeletedAtTimestampsField 删除时间戳
// 一般用于软删除
type DeletedAtTimestampsField struct {
	DeletedAt *JSONTime `gorm:"column:deleted_at;type:timestamp NULL;index;comment:删除时间;" json:"deleted_at,omitempty"`
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
func (m Model) ToMap(fields ...string) map[string]interface{} {
	_map := helpers.StructToMap(m, fields...)
	if helpers.SearchStringInSlice(fields, "id") > 0 {
		_map["id"] = m.ID
	}
	return _map
}