package make

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CmdMakeRequest 生成 请求文件的命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:25
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeRequest = &cobra.Command{
	Use:   "req",
	Short: "Create request file",
	Example: "go run main.go make req user -c user",
	Run:   runMakeRequest,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeRequest(cmd *cobra.Command, args []string) {

	// 获取注释
	comment, _ := cmd.Flags().GetString("comment")
	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0], comment)

	// 拼接目标文件路径
	filePath := fmt.Sprintf("app/requests/%s_request.go", model.PackageName)

	// 基于模板创建文件（做好变量替换）
	createFileFromStub(filePath, "request", model)
}