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
	Short: "Create a command, should be snake_case, exmaple: make cmd buckup_database -n \"backup database\"",
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func init() {
	CmdMakeCMD.Flags().StringP("comment" ,"c", "", "add comment for cmd struct")
	CmdMakeCMD.Run = runMakeCMD
}

func runMakeCMD(cmd *cobra.Command, args []string) {

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0])
	// 获取注释
	comment, err := CmdMakeCMD.Flags().GetString("comment")
	if err != nil {
		console.Exit(err.Error())
	}
	if len(comment) == 0 {
		console.Exit("missing comments for cmd struct")
	}

	// 目标目录不存在
	if !file.Exists("app/cmd") {
		_ = os.Mkdir("app/cmd", 0644)
	}

	// 拼接目标文件路径
	filePath := fmt.Sprintf("app/cmd/%s.go", model.PackageName)

	// 从模板中创建文件（做好变量替换）
	createFileFromStub(filePath, "cmd", model, map[string]string{"{{Comment}}":comment})

	// 友好提示
	console.Success("command name:" + model.PackageName)
	console.Success("command variable name: cmd.Cmd" + model.StructName)
	console.Warning("please edit main.go's app.Commands slice to register command")
}