package make

import (
	"fmt"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/console"

	"github.com/spf13/cobra"
)

// CmdMakeMigration 生成数据库迁移文件命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:57
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeMigration = &cobra.Command{
	Use:   "migration",
	Short: "Create a migration file",
	Example: "go run main.go make migration create_users_table -c create_users_table",
	Run:   runMakeMigration,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeMigration(cmd *cobra.Command, args []string) {

	// 获取注释
	comment, _ := cmd.Flags().GetString("comment")
	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0], comment, "")

	// 日期格式化
	timeStr := app.TimenowInTimezone().Format("2006_01_02_150405")
	fileName := timeStr + "_" + model.PackageName
	filePath := fmt.Sprintf("database/migrations/%s.go", fileName)

	createFileFromStub(filePath, "migration", model, map[string]string{"{{FileName}}": fileName})

	console.Success("Migration file created，after modify it, use `migrate up` to migrate database.")
}