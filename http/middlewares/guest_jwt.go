package middlewares

import (
	"github.com/yidejia/gofw/jwt"
	"github.com/yidejia/gofw/response"

	"github.com/gin-gonic/gin"
)

// GuestJWT 强制使用游客身份访问
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 15:35
// @copyright © 2010-2022 广州伊的家网络科技有限公司
func GuestJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(c.GetHeader("Authorization")) > 0 {

			// 解析 token 成功，说明登录成功了
			_, err := jwt.NewJWT().ParserToken(c)
			if err == nil {
				response.Unauthorized(c, "请使用游客身份访问")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}