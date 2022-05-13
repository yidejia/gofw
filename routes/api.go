// Package routes 注册路由
package routes

import (
	"github.com/yidejia/gofw/pkg/config"

	"github.com/yidejia/gofw/app/http/controllers/api/v1/user"
	gfMiddlewares "github.com/yidejia/gofw/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	// 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
	var v1 *gin.RouterGroup
	// 未配置 api 域名时，需要在 url 路径中增加 api 的前缀路径
	if len(config.Get("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}

	// 全局限流中间件：每小时限流。这里是所有 API （根据 IP）请求加起来。
	// 作为参考 Github API 每小时最多 60 个请求（根据 IP）。
	// 测试时，可以调高一点。
	v1.Use(gfMiddlewares.LimitIP("200-H"))
	{
		// 用户路由组
		usersGroup := v1.Group("/users")
		// 限流中间件：每小时限流，作为参考 Github API 每小时最多 60 个请求（根据 IP）
		// 测试时，可以调高一点
		usersGroup.Use(gfMiddlewares.LimitIP("1000-H"))
		{
			// TODO 注册用户路由，就按照示例注册其他模块的路由，一般一个模块对应一个路由组
			uc := new(user.UsersController)
			// 创建用户
			usersGroup.POST("", uc.Store)
			// 用户列表
			usersGroup.GET("", uc.Index)
		}
	}
}