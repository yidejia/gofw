package main

import (
	"fmt"
	btsConfig "github.com/yidejia/gofw/config"
	"github.com/yidejia/gofw/pkg/cmd"
	"github.com/yidejia/gofw/pkg/console"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// 应用初始化时触发加载 config 目录下的配置信息
	btsConfig.Initialize()
}

func main() {

	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   "gofw",
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,
	}

	// 注册子命令
	rootCmd.AddCommand(
		//cmd.CmdServe,
	)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}