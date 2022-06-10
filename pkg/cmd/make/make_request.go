package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/file"
	"os"
)

// CmdMakeRequest 生成请求文件命令
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

func init() {
	CmdMakeRequest.Flags().BoolP("force" ,"f", false, "Force files to be created in the services root directory")
}

func runMakeRequest(cmd *cobra.Command, args []string) {

	// 获取是否强制执行选项
	force := parseForceFlag(cmd, args)
	// 获取注释
	comment := parseCommentFlag(cmd, args, force)
	// 获取名称
	path, name, pkgName := parseNameParam(cmd, args)
	if len(pkgName) == 0 {
		pkgName = "requests"
	}

	// 模型目录不存在
	if !file.Exists("app/requests") {
		if err := os.Mkdir("app/requests", os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create requests folder: %s", err.Error()))
		}
	}

	// 完整目标文件目录
	var fullPath string
	if len(path) > 0 {
		fullPath = fmt.Sprintf("app/requests/%s", path)
	} else {
		if !force {
			console.Exit("The request should belong to its own package. If you want to continue creating, please set the \"-f\" option.")
		}
		fullPath = "app/requests"
	}

	// 目标目录不存在
	if !file.Exists(fullPath) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create request parent folder: %s", err.Error()))
		}
	}

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(name, comment, pkgName)
	// 拼接目标文件路径
	filePath := fmt.Sprintf("%s/%s_request.go", fullPath, model.PackageName)
	// 从模板中创建文件并进行变量替换
	createFileFromStub(filePath, "request", model)
}