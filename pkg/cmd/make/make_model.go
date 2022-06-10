package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/file"
	"os"
)

// CmdMakeModel 生成模型文件命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:08
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeModel = &cobra.Command{
	Use:   "model",
	Short: "Crate model file",
	Example: "go run main.go make model user -c user/user",
	Run:   runMakeModel,
	Args:  cobra.MinimumNArgs(1), // 至少传 1 个参数
}

func init() {
	CmdMakeModel.Flags().BoolP("force" ,"f", false, "Force files to be created in the models root directory")
}

func runMakeModel(cmd *cobra.Command, args []string) {

	// 获取是否强制执行选项
	force := parseForceFlag(cmd, args)
	// 获取注释
	comment := parseCommentFlag(cmd, args, force)
	// 获取名称
	path, name, pkgName := parseNameParam(cmd, args)
	if len(pkgName) == 0 {
		pkgName = "models"
	}

	// 模型目录不存在
	if !file.Exists("app/models") {
		if err := os.Mkdir("app/models", os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create models folder: %s", err.Error()))
		}
	}

	// 完整目标文件目录
	var fullPath string
	if len(path) > 0 {
		fullPath = fmt.Sprintf("app/models/%s", path)
	} else {
		if !force {
			console.Exit("The model should belong to its own package. If you want to continue creating, please set the \"-f\" option.")
		}
		fullPath = "app/models"
	}

	// 目标目录不存在
	if !file.Exists(fullPath) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create model parent folder: %s", err.Error()))
		}
	}

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(name, comment, pkgName)
	// 拼接目标文件路径
	filePath := fmt.Sprintf("%s/%s.go", fullPath, model.PackageName)
	// 从模板中创建文件并进行变量替换
	createFileFromStub(filePath, "model", model)
}