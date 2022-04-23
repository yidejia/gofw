// Package cmd 存放程序的所有子命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 16:29
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package cmd

import (
	"github.com/yidejia/gofw/helpers"
	"os"

	"github.com/spf13/cobra"
)

// Env 存储全局选项 --env 的值
var Env string

// RegisterGlobalFlags 注册全局选项（flag）
func RegisterGlobalFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVarP(&Env, "env", "e", "", "load .env file, example: --env=testing will use .env.testing file")
}

// RegisterDefaultCmd 注册默认命令
func RegisterDefaultCmd(rootCmd *cobra.Command, subCmd *cobra.Command) {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	firstArg := helpers.FirstElement(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use && firstArg != "-h" && firstArg != "--help" {
		args := append([]string{subCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}