// Package auth 用户认证包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 11:31
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package auth

import "github.com/gin-gonic/gin"

// CurrentUID 从 gin.context 中获取当前登录用户 ID
func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}
