// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
	gfMiddlewares "github.com/yidejia/gofw/pkg/http/middlewares"
	"github.com/yidejia/gofw/routes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRoute 路由初始化
func SetupRoute(router *gin.Engine) {

	// 注册全局中间件
	registerGlobalMiddleWare(router)

	//  注册 API 路由
	routes.RegisterAPIRoutes(router)

	// TODO 如果注册了其它路由，请在这里调用注册函数，注册函数统一放在 routes 包，如果分模块注册路由，请在 routes 包下创建 api 包存放路由文件
	// 示例：
	// routes.RegisterUserAPIRoutes(router)
	// routes.RegisterWebRoutes(router)
	// routes.RegisterUserWebRoutes(router)
	// routes.RegisterAdminRoutes(router)

	//  配置 404 路由
	setup404Handler(router)
}

// registerGlobalMiddleWare 注册全局中间件
func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(
		gfMiddlewares.Logger(),
		gfMiddlewares.Recovery(),
		gfMiddlewares.ForceUA(),
		// TODO 在这里注册其它中间件
	)
}

// setup404Handler 配置 404 路由
func setup404Handler(router *gin.Engine) {
	// 处理 404 请求
	router.NoRoute(func(c *gin.Context) {
		// 获取标头信息的 Accept 信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是 HTML 的话
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			// 默认返回 JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})
}