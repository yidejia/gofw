package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/file"
	"os"
)

// CmdMakeCMD 生成命令文件的命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:07
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeCMD = &cobra.Command{
	Use:   "cmd",
	Short: "Create a command, should be snake_case",
	Example: "go run main.go make cmd backup_database -c backup_database",
	Run:   runMakeCMD,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeCMD(cmd *cobra.Command, args []string) {

	// 获取注释
	comment := parseCommentFlag(cmd, args)
	// 获取名称
	path, name, pkgName := parseNameParam(cmd, args)
	if len(pkgName) == 0 {
		pkgName = "cmd"
	}

	// 命令目录不存在
	if !file.Exists("app/cmd") {
		if err := os.Mkdir("app/cmd", os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create cmd folder: %s", err.Error()))
		}
	}

	// 完整目标文件目录
	var fullPath string
	if len(path) > 0 {
		fullPath = fmt.Sprintf("app/cmd/%s", path)
	} else {
		fullPath = "app/cmd"
	}

	// 目标目录不存在
	if !file.Exists(fullPath) {
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create cmd parent folder: %s", err.Error()))
		}
	}

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(name, comment, pkgName)

	// 拼接目标文件路径
	filePath := fmt.Sprintf("%s/%s.go", fullPath, model.PackageName)

	// 从模板中创建文件并进行变量替换
	createFileFromStub(filePath, "cmd", model)

	// 友好提示
	console.Success("command name:" + model.PackageName)
	console.Success(fmt.Sprintf("command variable name: %s.Cmd%s", model.CustomPackageName, model.StructName))
	console.Warning("please edit main.go's app.Commands slice to register command")
}