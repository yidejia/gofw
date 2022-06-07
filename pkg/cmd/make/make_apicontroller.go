package make

import (
	"fmt"
	"github.com/yidejia/gofw/pkg/console"
	"strings"

	"github.com/spf13/cobra"
)

// CmdMakeAPIController 生成 api 控制器文件的命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:17
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeAPIController = &cobra.Command{
	Use:   "api-ctr",
	Short: "Create api controller",
	Example: "go run main.go make api-ctr v1/user -c user",
	Run:   runMakeAPIController,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeAPIController(cmd *cobra.Command, args []string) {

	// 处理参数，要求附带 API 版本（v1 或者 v2）
	array := strings.Split(args[0], "/")
	if len(array) != 2 {
		console.Exit("api controller name format: v1/user")
	}

	// apiVersion 用来拼接目标路径
	// name 用来生成 cmd.Model 实例
	apiVersion, name := array[0], array[1]
	// 获取注释
	comment, _ := cmd.Flags().GetString("comment")
	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(name, comment)

	// 组建目标目录
	filePath := fmt.Sprintf("app/http/controllers/api/%s/%s_controller.go", apiVersion, model.TableName)

	// 基于模板创建文件（做好变量替换）
	createFileFromStub(filePath, "apicontroller", model)
}