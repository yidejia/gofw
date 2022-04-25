// Package auth 用户认证包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 11:31
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yidejia/gofw/pkg/logger"
)

// userResolver 获取用户模型的函数
var userResolver UserResolver

// SetUserResolver 设置获取用户模型的函数
func SetUserResolver(_userResolver UserResolver)  {
	userResolver = _userResolver
}

// ResolveUser 获取用户模型
func ResolveUser(id string) (user Authenticate) {
	if userResolver != nil {
		return userResolver(id)
	}
	return user
}

// CurrentUID 从 gin.context 中获取当前登录用户 ID
func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}

// CurrentUser 获取当前登录用户
func CurrentUser(c *gin.Context) (user Authenticate) {
	_user, ok := c.MustGet("current_user").(Authenticate)
	if !ok {
		logger.LogIf(errors.New("无法获取用户"))
		return user
	}
	return _user
}
