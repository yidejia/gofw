package middleware

import (
	gfErrs "github.com/yidejia/gofw/pkg/errors"
	"github.com/yidejia/gofw/pkg/limiter"
	"github.com/yidejia/gofw/pkg/logger"
	"github.com/yidejia/gofw/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// LimitIP 全局限流中间件，针对 IP 进行限流，需要配置 redis 才能使用
// @author 余海坚 haijianyu10@qq.com
// @created 2022-09-02 17:30
// @copyright © 2010-2022 广州伊的家网络科技有限公司
//
// limit 为格式化字符串，如 "5-S"，示例:
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
func LimitIP(limit string, variables ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 针对 IP 限流
		key := limiter.GetKeyIP(c)
		if ok := limitHandler(c, key, limit, variables...); !ok {
			return
		}
		c.Next()
	}
}

// LimitRoute 路由限流中间件，针对路由进行限流，需要配置 redis 才能使用
// @author 余海坚 haijianyu10@qq.com
// @created 2022-09-06 11:47
// @copyright © 2010-2022 广州伊的家网络科技有限公司
//
// limit 为格式化字符串，如 "5-S"，示例:
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
func LimitRoute(limit string, variables ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 针对单个路由，增加访问次数
		c.Set("limiter-once", false)
		// 针对 IP + 路由进行限流
		key := limiter.GetKeyRouteWithIP(c)
		if ok := limitHandler(c, key, limit, variables...); !ok {
			return
		}
		c.Next()
	}
}

// LimitKey 自定义 key 限流中间件，针对自定义 key 进行限流，需要配置 redis 才能使用，在控制器中结合业务逻辑使用
// @author 余海坚 haijianyu10@qq.com
// @created 2022-09-06 16:35
// @copyright © 2010-2022 广州伊的家网络科技有限公司
//
// limit 为格式化字符串，如 "5-S"，示例:
// * 5 reqs/second: "5-S"
// * 10 reqs/minute: "10-M"
// * 1000 reqs/hour: "1000-H"
// * 2000 reqs/day: "2000-D"
func LimitKey(c *gin.Context, key string, limit string, variables ...interface{}) bool {
	c.Set("limiter-once", false)
	return limitHandler(c, key, limit, variables...)
}

// limitHandler 对请求进行限流处理
func limitHandler(c *gin.Context, key string, limit string, variables ...interface{}) bool {

	// 获取超额的情况
	rate, _err := limiter.CheckRate(c, key, limit)
	if _err != nil {
		logger.LogIf(_err)
		response.InternalError(c, _err)
		return false
	}

	// ---- 设置标头信息-----
	// X-RateLimit-Limit :10000 最大访问次数
	// X-RateLimit-Remaining :9993 剩余的访问次数
	// X-RateLimit-Reset :1513784506 到该时间点，访问次数会重置为 X-RateLimit-Limit
	c.Header("X-RateLimit-Limit", cast.ToString(rate.Limit))
	c.Header("X-RateLimit-Remaining", cast.ToString(rate.Remaining))
	c.Header("X-RateLimit-Reset", cast.ToString(rate.Reset))

	// 超额
	if rate.Reached {

		// 请求上下文标记请求频率已超标
		c.Set("limiter-reached", true)

		if len(variables) > 1 {
			// 调用方同时设置了错误和元数据
			err := variables[0].(gfErrs.ResponsiveError)
			meta := variables[1].(map[string]interface{})
			response.AbortWithError(c, err, meta)
		} else if len(variables) > 0 {
			// 调用方只设置了错误
			err := variables[0].(gfErrs.ResponsiveError)
			response.AbortWithError(c, err)
		} else {
			// 默认只返回调用超额了错误
			response.TooManyRequests(c, "接口请求太频繁")
		}

		return false
	}

	return true
}
