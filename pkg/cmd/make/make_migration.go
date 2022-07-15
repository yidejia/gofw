package make

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/file"
)

// CmdMakeMigration 生成数据库迁移文件命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:57
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeMigration = &cobra.Command{
	Use:     "migration",
	Short:   "Create a migration file",
	Example: "go run main.go make migration create_users_table -c create_users_table",
	Run:     runMakeMigration,
	Args:    cobra.MinimumNArgs(1), // 至少传 1 个参数
}

func init() {
	CmdMakeMigration.Flags().StringP("model", "m", "", "Create table for a model")
}

func runMakeMigration(cmd *cobra.Command, args []string) {

	// 获取注释
	comment, _ := cmd.Flags().GetString("comment")
	name := ""
	forModel, err := cmd.Flags().GetString("model")
	if err != nil {
		console.Exit(fmt.Sprintf("Get model flage error: %s", err.Error()))
	}
	if forModel != "" {
		name = forModel
	}

	// 迁移目录不存在
	if !file.Exists("database/migrations") {
		if err := os.Mkdir("database/migrations", os.ModePerm); err != nil {
			console.Exit(fmt.Sprintf("failed to create migrations folder: %s", err.Error()))
		}
	}

	// 日期格式化
	timeStr := app.TimeNowInTimezone().Format("2006_01_02_150405")
	var fileName string
	variables := make(map[string]string)
	// 格式化模型名称，返回一个 Model 对象
	var model Model
	// 设置了模型名称
	if name != "" {
		model = makeModelFromString(name, comment, "")
		fileName = fmt.Sprintf("%s_create_%s_table", timeStr, model.TableName)
		// 启用表名注释那段代码
		variables["{{CommentOutTableName}}"] = ""
		// 表名注释
		variables["{{TableNameComment}}"] = model.Comment
		// 迁移说明
		variables["{{Instruction}}"] = fmt.Sprintf("创建%s表", variables["{{TableNameComment}}"])
	} else {
		// 未设置模型名称时，默认为 user 作为示例
		name = "user"
		model = makeModelFromString(name, comment, "")
		// 根据输入参数生成文件名
		fileName = timeStr + "_" + args[0]
		// 注释掉表名注释那段代码
		variables["{{CommentOutTableName}}"] = "// "
		// 表名注释固定为“用户”，作为示例
		variables["{{TableNameComment}}"] = "用户"
		// 迁移说明就是注释
		variables["{{Instruction}}"] = model.Comment
	}

	filePath := fmt.Sprintf("database/migrations/%s.go", fileName)

	variables["{{FileName}}"] = fileName

	createFileFromStub(filePath, "migration", model, variables)

	console.Success("Migration file created，after modify it, use `migrate up` to migrate database.")
}
