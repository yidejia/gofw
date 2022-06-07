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
	comment := parseCmdComment(cmd, args)
	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0], comment)

	// 命令目录不存在
	if !file.Exists("app/cmd") {
		if err := os.Mkdir("app/cmd", 0644); err != nil {
			panic(fmt.Sprintf("failed to create cmd folder: %s", err.Error()))
		}
	}

	// 拼接目标文件路径
	filePath := fmt.Sprintf("app/cmd/%s.go", model.PackageName)

	// 从模板中创建文件并进行变量替换
	createFileFromStub(filePath, "cmd", model)

	// 友好提示
	console.Success("command name:" + model.PackageName)
	console.Success("command variable name: cmd.Cmd" + model.StructName)
	console.Warning("please edit main.go's app.Commands slice to register command")
}