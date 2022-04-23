package make

import (
	"fmt"
	"github.com/yidejia/gofw/app"
	"github.com/yidejia/gofw/console"

	"github.com/spf13/cobra"
)

// CmdMakeMigration 生成数据库迁移文件的命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:57
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeMigration = &cobra.Command{
	Use:   "migration",
	Short: "Create a migration file, example: make migration add_users_table",
	Run:   runMakeMigration,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeMigration(cmd *cobra.Command, args []string) {

	// 日期格式化
	timeStr := app.TimenowInTimezone().Format("2006_01_02_150405")

	model := makeModelFromString(args[0])
	fileName := timeStr + "_" + model.PackageName
	filePath := fmt.Sprintf("database/migrations/%s.go", fileName)
	createFileFromStub(filePath, "migration", model, map[string]string{"{{FileName}}": fileName})
	console.Success("Migration file created，after modify it, use `migrate up` to migrate database.")
}