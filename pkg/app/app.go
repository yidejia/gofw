// Package app 应用包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 16:59
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package app

import (
	"time"

	"github.com/yidejia/gofw/pkg/config"
)

// IsLocal 当前运行在本地开发环境
func IsLocal() bool {
	return config.Get("app.env") == "local"
}

// IsProduction 当前运行在生产环境
func IsProduction() bool {
	return config.Get("app.env") == "production"
}

// IsTesting 当前运行在测试环境
func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

// TimeNowInTimezone 获取当前时间，支持时区
func TimeNowInTimezone() time.Time {
	timezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(timezone)
}

// URL 传参 path 拼接站点的 URL
func URL(path string) string {
	return config.Get("app.url") + path
}

// V1URL 拼接带 v1 标示 URL
func V1URL(path string) string {
	return URL("/v1/" + path)
}
