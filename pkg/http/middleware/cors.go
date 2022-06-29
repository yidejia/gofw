package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-16 09:53
// @copyright © 2010-2022 广州伊的家网络科技有限公司
func Cors() gin.HandlerFunc {

	return func(c *gin.Context) {

		allowOrigin := "*"
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			allowOrigin = origin
		}

		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, token, X-Requested-With, Origin, Access-Control-Request-Headers, SERVER_NAME, Access-Control-Allow-Headers, Accept, Connection, User-Agent, Cookie")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 允许放行 OPTIONS 请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
