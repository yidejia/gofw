// Package middlewares Gin 中间件
package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yidejia/gofw/pkg/auth"
	"github.com/yidejia/gofw/pkg/jwt"
	"github.com/yidejia/gofw/pkg/response"
)

// AuthJWT 授权中间件
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 14:17
// @copyright © 2010-2022 广州伊的家网络科技有限公司
func AuthJWT() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 从请求头或者请求参数中获取 token，并验证 JWT 的准确性
		claims, err := jwt.NewJWT().ParserToken(c)
		// JWT 解析失败，有错误发生
		if err != nil {
			response.Unauthorized(c, fmt.Sprintf("无权访问：%v", err.Error()))
			return
		}

		// JWT 解析成功，设置用err
		user, gfErr := auth.ResolveUser(claims.UserID)
		if gfErr != nil {
			response.Unauthorized(c, "用户不存在，鉴权失败")
			return
		}
		// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
		c.Set("current_user_id", claims.UserID)
		c.Set("current_user_name", claims.UserName)
		c.Set("current_user", user)

		c.Next()
	}
}