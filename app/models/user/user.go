// Package user 存放用户模块相关模型的包
package user

import (
	"github.com/yidejia/gofw/app/models"
)

// User 用户模型，根据应用需要对模型字段进行定制
type User struct {
	models.Model

	Name string `json:"name,omitempty"`
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	models.CommonTimestampsField
}
