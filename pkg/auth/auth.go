// Package auth 用户认证包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 11:31
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package auth

import (
	"github.com/gin-gonic/gin"
	gfErrs "github.com/yidejia/gofw/pkg/errors"
)

// userResolver 获取用户模型的函数
var userResolver UserResolver

// SetUserResolver 设置获取用户模型的函数
func SetUserResolver(_userResolver UserResolver) {
	userResolver = _userResolver
}

// ResolveUser 获取用户模型
func ResolveUser(id uint64) (user Authenticate, err gfErrs.ResponsiveError) {
	if userResolver != nil {
		return userResolver(id)
	}
	return nil, gfErrs.NewErrorNotFound("用户不存在")
}

// CurrentUID 从 gin.context 中获取当前登录用户 ID
func CurrentUID(c *gin.Context) string {
	return c.GetString("user_id")
}

// CurrentUser 获取当前登录用户
func CurrentUser(c *gin.Context) Authenticate {

	_user, ok := c.Get("user")
	if !ok {
		return nil
	}

	authUser, ok := _user.(Authenticate)
	if !ok {
		return nil
	}

	return authUser
}
