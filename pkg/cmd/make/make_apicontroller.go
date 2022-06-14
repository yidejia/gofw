package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/file"
	"os"
)

// CmdMakeAPIController 生成 api 控制器文件命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:17
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeAPIController = &cobra.Command{
	Use:     "api",
	Short:   "Create api controller",
	Example: "go run main.go make api-ctr v1/user -c user",
	Run:     runMakeAPIController,
	Args:    cobra.MinimumNArgs(1), // 至少传 1 个参数
}

func init() {
	CmdMakeAPIController.Flags().BoolP("force" ,"f", false, "Force files to be created in the api-ctr root directory")
	CmdMakeAPIController.Flags().Uint8P("version", "v", 1, "Set api-ctr version")
}

func runMakeAPIController(cmd *cobra.Command, args []string) {

	// 获取是否强制执行选项
	force := parseForceFlag(cmd, args)
	// 获取注释
	comment := parseCommentFlag(cmd, args, force)
	// 接口默认版本号
	versionName := "v1"
	// 获取接口版本号
	version, err := cmd.Flags().GetUint8("version")
	if err == nil && version > 1 {
		versionName = fmt.Sprintf("v%d", version)
	}
	// 获取名称
	path, name, pkgName := parseNameParam(cmd, args)
	if len(pkgName) == 0 {
		// 默认是对应版本包
		pkgName = versionName
	}

	// 控制器目录不存在
	if !file.Exists("app/http/controllers/api") {
		if err := os.MkdirAll("app/http/controllers/api", os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create api-ctr folder: %s", err.Error()))
		}
	}

	// 完整目标文件目录
	var fullPath string
	if len(path) > 0 {
		fullPath = fmt.Sprintf("app/http/controllers/api/%s/%s", versionName, path)
	} else {
		if !force {
			console.Exit("The service should belong to its own package. If you want to continue creating, please set the \"-f\" option.")
		}
		fullPath = fmt.Sprintf("app/http/controllers/api/%s", versionName)
	}

	// 目标目录不存在
	if !file.Exists(fullPath) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create api-ctr parent folder: %s", err.Error()))
		}
	}

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(name, comment, pkgName)
	// 拼接目标文件路径
	filePath := fmt.Sprintf("%s/%s_controller.go", fullPath, model.VariableNamePlural)
	// 从模板中创建文件并进行变量替换
	createFileFromStub(filePath, "apicontroller", model, map[string]string{"{{Version}}":versionName})
}