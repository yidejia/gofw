package middleware

import (
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/yidejia/gofw/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/yidejia/gofw/pkg/response"
	"go.uber.org/zap"
)

// Recovery 使用 zap.Error() 来记录 Panic 和 call stack
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 17:18
// @copyright © 2010-2022 广州伊的家网络科技有限公司
func Recovery() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				// 获取用户的请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, true)

				// 链接中断，客户端中断连接为正常行为，不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 链接中断的情况
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					_ = c.Error(err.(error))
					c.Abort()
					// 链接已断开，无法写状态码
					return
				}

				// 如果不是链接中断，就开始记录堆栈信息
				logger.Error("recovery from panic",
					zap.Time("time", time.Now()),               // 记录时间
					zap.Any("error", err),                      // 记录错误信息
					zap.String("request", string(httpRequest)), // 请求信息
					zap.Stack("stacktrace"),                    // 调用堆栈信息
				)

				// 中断处理并返回内部系统错误
				response.InternalError(c, nil, "系统错误，请稍候再试")
			}
		}()
		c.Next()
	}
}
