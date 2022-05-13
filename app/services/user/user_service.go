package user

import (
	"github.com/yidejia/gofw/app/models/user"
	userRepos "github.com/yidejia/gofw/app/repositories/user"
	svcs "github.com/yidejia/gofw/app/services"
	gfErrors "github.com/yidejia/gofw/pkg/errors"
)

// UserService 用户服务
type UserService struct {
	svcs.Service
	repo *userRepos.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: userRepos.NewUserRepository(),
	}
}

func (svc *UserService) Create() (user.User, gfErrors.ResponsiveError)  {

	_user := user.User{}

	svc.repo.Create(&_user)
	_user.ID = 1

	if _user.ID > 0 {
		return _user, nil
	} else {
		return _user, svc.ErrorInternal("创建用户失败")
	}
}
