package cmd

import (
	"github.com/yidejia/gofw/console"
	"github.com/yidejia/gofw/helpers"

	"github.com/spf13/cobra"
)

// CmdKey  key 命令，生成随机字符串，可用来设置我们的 APP_KEY 环境变量的值
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 16:40
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdKey = &cobra.Command{
	Use:   "key",
	Short: "Generate App Key, will print the generated Key",
	Run:   runKeyGenerate,
	Args:  cobra.NoArgs, // 不允许传参
}

func runKeyGenerate(cmd *cobra.Command, args []string) {
	console.Success("---")
	console.Success("App Key:")
	console.Success(helpers.RandomString(32))
	console.Success("---")
	console.Warning("please go to .env file to change the APP_KEY option")
}