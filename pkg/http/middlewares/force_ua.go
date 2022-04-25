// Package middlewares Gin 中间件
package middlewares

import (
	"errors"
	"github.com/yidejia/gofw/pkg/response"

	"github.com/gin-gonic/gin"
)

// ForceUA 中间件，强制请求必须附带 User-Agent 标头
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-25 10:50
// @copyright © 2010-2022 广州伊的家网络科技有限公司
func ForceUA() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取 User-Agent 标头信息
		if len(c.Request.Header["User-Agent"]) == 0 {
			response.BadRequest(c, errors.New("User-Agent 标头未找到"), "请求必须附带 User-Agent 标头")
			return
		}

		c.Next()
	}
}