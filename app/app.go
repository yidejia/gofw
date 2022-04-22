// Package app 应用信息
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 16:59
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package app

import (
	"github.com/yidejia/gofw/config"
	"time"
)

func IsLocal() bool {
	return config.Get("app.env") == "local"
}

func IsProduction() bool {
	return config.Get("app.env") == "production"
}

func IsTesting() bool {
	return config.Get("app.env") == "testing"
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation(config.GetString("app.timezone"))
	return time.Now().In(chinaTimezone)
}