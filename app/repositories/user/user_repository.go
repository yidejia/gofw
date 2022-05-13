// Package user 用户模块包
package user

import (
	"github.com/yidejia/gofw/app/models/user"
	repos "github.com/yidejia/gofw/app/repositories"
	"github.com/yidejia/gofw/pkg/database"
)

// UserRepository 用户数据仓库
type UserRepository struct {
	repos.Repository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) Create(user *user.User) (err error)  {
	if err = database.DB.Create(user).Error; err != nil {
		err = repo.ErrorInternal("创建用户失败")
	}
	return
}
