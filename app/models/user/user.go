// Package user 存放用户模块相关模型的包
package user

import (
	"github.com/yidejia/gofw/app/models"
	"github.com/yidejia/gofw/pkg/hash"
	gfModels "github.com/yidejia/gofw/pkg/models"
	"gorm.io/gorm"
)

// User 用户模型，根据应用需要对模型字段进行定制
type User struct {
	models.Model

	Name string `json:"name,omitempty"`
	Email    string `json:"-"`
	Phone    string `json:"-"`
	Password string `json:"-"`

	gfModels.CommonTimestampsField // 时间戳，不需要的可以删除这个嵌入
}

// BeforeSave GORM 模型钩子，在创建和更新模型前调用
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	// 对密码进行加密
	if !hash.BcryptIsHashed(user.Password) {
		user.Password = hash.BcryptHash(user.Password)
	}
	return
}

func (user *User) AuthId() uint64 {
	return user.ID
}

func (user *User) AuthName() string {
	return user.Name
}
