// Package app 应用包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 16:59
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package app

import (
	"time"

	"github.com/yidejia/gofw/pkg/config"
)

// ExitHandler 应用退出处理函数
type ExitHandler func() error

// ErrExitHandler 应用退出错误处理函数
type ErrExitHandler func(moduleName string, processName string, err error)

// ExitProcess 应用退出处理工序
type ExitProcess struct {
	Name       string         // 工序名称
	Handler    ExitHandler    // 工序处理函数
	ErrHandler ErrExitHandler // 工序错误处理函数
}

// exitProcesses 应用退出处理工序集合
var exitProcesses = make(map[string][]*ExitProcess)

// RegisterExitHandler 注册应用退出处理函数
// @param moduleName 模块名
// @param processName 处理工序名称
// @param errHandler 可选处理工序错误处理函数
func RegisterExitHandler(moduleName string, processName string, handler ExitHandler, errHandler ...ErrExitHandler) {
	p := &ExitProcess{Name: processName, Handler: handler}
	// 设置了错误处理函数
	if len(errHandler) > 0 {
		p.ErrHandler = errHandler[0]
	}
	exitProcesses[moduleName] = append(exitProcesses[moduleName], p)
}

// Exit 退出应用
// @param errHandler 可选处理工序错误处理函数
func Exit(errHandler ...ErrExitHandler) {
	// 退出应用前调用各模块注册的处理函数进行清理和资源回收工作
	for moduleName, processes := range exitProcesses {
		for _, process := range processes {
			if err := process.Handler(); err != nil {
				// 调用模块自身退出错误处理函数
				if process.ErrHandler != nil {
					process.ErrHandler(moduleName, process.Name, err)
				} else if len(errHandler) > 0 {
					// 调用应用退出错误处理函数
					errHandler[0](moduleName, process.Name, err)
				}
			}
		}
	}
}

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
