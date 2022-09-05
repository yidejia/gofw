package middleware

import (
	"net/http"

	"github.com/yidejia/gofw/pkg/limiter"
	"github.com/yidejia/gofw/pkg/logger"
	"github.com/yidejia/gofw/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// LimitIP 全局限流中间件，针对 IP 进行限流，需要配置 redis 才能使用
// limit 为格式化字符串，如 "5-S" ，示例:
//
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
//
func LimitIP(limit string, message ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 针对 IP 限流
		key := limiter.GetKeyIP(c)
		if ok := limitHandler(c, key, limit, message...); !ok {
			return
		}
		c.Next()
	}
}

// LimitPerRoute 限流中间件，用在单独的路由中
func LimitPerRoute(limit string, message ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 针对单个路由，增加访问次数
		c.Set("limiter-once", false)

		// 针对 IP + 路由进行限流
		key := limiter.GetKeyRouteWithIP(c)
		if ok := limitHandler(c, key, limit, message...); !ok {
			return
		}

		c.Next()
	}
}

// limitHandler 对请求进行限流
func limitHandler(c *gin.Context, key string, limit string, message ...string) bool {

	// 获取超额的情况
	rate, err := limiter.CheckRate(c, key, limit)
	if err != nil {
		logger.LogIf(err)
		response.InternalError(c, err)
		return false
	}

	// ---- 设置标头信息-----
	// X-RateLimit-Limit :10000 最大访问次数
	// X-RateLimit-Remaining :9993 剩余的访问次数
	// X-RateLimit-Reset :1513784506 到该时间点，访问次数会重置为 X-RateLimit-Limit
	c.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))
	c.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining))
	c.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))

	msg := "接口请求太频繁"
	// 设置了自定义消息
	if len(message) > 0 {
		msg = message[0]
	}
	// 超额
	if rate.Reached {
		// 提示用户超额了
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"message": msg,
		})
		return false
	}

	return true
}
