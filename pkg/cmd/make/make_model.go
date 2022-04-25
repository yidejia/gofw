package make

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CmdMakeModel 生成模型文件的命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:08
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeModel = &cobra.Command{
	Use:   "model",
	Short: "Crate model file, example: make model user",
	Run:   runMakeModel,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeModel(cmd *cobra.Command, args []string) {

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0])

	// 确保模型的目录存在，例如 `app/models/user`
	dir := fmt.Sprintf("app/models/%s/", model.PackageName)
	// os.MkdirAll 会确保父目录和子目录都会创建，第二个参数是目录权限，使用 0777
	_ = os.MkdirAll(dir, os.ModePerm)

	// 替换变量
	createFileFromStub(dir+model.PackageName+"_model.go", "model/model", model)
	createFileFromStub(dir+model.PackageName+"_util.go", "model/model_util", model)
	createFileFromStub(dir+model.PackageName+"_hooks.go", "model/model_hooks", model)
}