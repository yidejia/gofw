package cmd

import (
	"github.com/yidejia/gofw/pkg/console"

	"github.com/spf13/cobra"
)

// CmdVersion 输出框架版本命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-25 19:15
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Show Gofw Framework version",
	Run:   runVersion,
	Args:  cobra.ExactArgs(0),
}

func runVersion(cmd *cobra.Command, args []string) {
	console.Success("Gofw Framework 0.1.14")
}
