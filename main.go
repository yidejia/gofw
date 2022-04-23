package main

import (
	"fmt"
	"github.com/yidejia/gofw/cmd"
	"github.com/yidejia/gofw/console"
	"os"

	"github.com/spf13/cobra"
)

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