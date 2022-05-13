// Package user 用户模块包
package user

import (
	"github.com/yidejia/gofw/app/models/user"
	repos "github.com/yidejia/gofw/app/repositories"
	"github.com/yidejia/gofw/pkg/database"
	gfErrors "github.com/yidejia/gofw/pkg/errors"
)

// UserRepository 用户数据仓库
type UserRepository struct {
	repos.Repository
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) Create(user *user.User) (err gfErrors.ResponsiveError)  {
	if _err := database.DB.Create(user).Error; _err != nil {
		err = repo.ErrorInternal(_err,"创建用户失败")
	}
	return
}
