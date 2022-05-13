package bootstrap

import (
	userSvcs "github.com/yidejia/gofw/app/services/user"
	"github.com/yidejia/gofw/pkg/auth"
	gfErrors "github.com/yidejia/gofw/pkg/errors"
)

// SetupAuth 初始化授权
func SetupAuth()  {
	auth.SetUserResolver(func(userId string) (user auth.Authenticate, err gfErrors.ResponsiveError) {
		// TODO 这里实现根据用户 id 获取用户模型的逻辑，根据应用需要进行修改
		_user, err := userSvcs.NewUserService().Get(userId)
		return &_user, err
	})
}
