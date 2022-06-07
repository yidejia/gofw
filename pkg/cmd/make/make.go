// Package make 封装常用结构体和文件生成命令的包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:00
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package make

import (
	"embed"
	"fmt"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/file"
	"github.com/yidejia/gofw/pkg/str"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// Model 模板文件填充模型，参数解释
// 单个词，用户命令传参，以 User 模型为例：
//  - user
//  - User
//  - users
//  - Users
// 整理好的数据：
// {
//     "TableName": "users",
//     "StructName": "User",
//     "StructNamePlural": "Users"
//     "VariableName": "user",
//     "VariableNamePlural": "users",
//     "PackageName": "user"
// }
// -
// 两个词或者以上，用户命令传参，以 TopicComment 模型为例：
//  - topic_comment
//  - topic_comments
//  - TopicComment
//  - TopicComments
// 整理好的数据：
// {
//     "TableName": "topic_comments",
//     "StructName": "TopicComment",
//     "StructNamePlural": "TopicComments"
//     "VariableName": "topicComment",
//     "VariableNamePlural": "topicComments",
//     "PackageName": "topic_comment"
// }
type Model struct {
	TableName          string // 生成模型结构体时对应的数据表名
	StructName         string // 结构体名称
	StructNamePlural   string // 结构体名称复数形式
	VariableName       string // 结构体变量名
	VariableNamePlural string // 结构体变量名复数形式
	PackageName        string // 结构体所属包名
	Comment            string // 结构体注释
}

// stubsFS 方便我们后面打包这些 .stub 为后缀名的文件
//go:embed stubs
var stubsFS embed.FS

// CmdMake 说明 cobra 命令
var CmdMake = &cobra.Command{
	Use:   "make",
	Short: "Generate file and code",
}

func init() {
	// 注册 make 的子命令
	CmdMake.AddCommand(
		CmdMakeCMD,
		CmdMakeModel,
		CmdMakeAPIController,
		CmdMakeRequest,
		CmdMakeMigration,
		CmdMakeFactory,
		CmdMakeSeeder,
	)

	for _, subCmd := range CmdMake.Commands() {
		// 为子命令设置注释选项
		subCmd.Flags().StringP("comment" ,"c", "", "set comment for struct or file")
	}
}

// parseCmdComment 解析并提取命令注释
func parseCmdComment(cmd *cobra.Command, args []string) string {
	// 检查是否设置了注释
	comment, err := cmd.Flags().GetString("comment")
	if err != nil {
		console.Exit(err.Error())
	}
	if len(comment) == 0 {
		console.Exit("Missing comment for struct or file, please use \"-c\" to set comment flag")
	}
	return comment
}

// makeModelFromString 格式化用户输入的内容并设置好模板文件的填充数据
func makeModelFromString(name string, comment string) Model {
	model := Model{}
	// 结构体名为驼峰单数形式
	model.StructName = str.Singular(strcase.ToCamel(name))
	// 结构体名的复数
	model.StructNamePlural = str.Plural(model.StructName)
	// 结构体为模型时，模型对应的数据表名为结构体名复数的蛇形形式
	model.TableName = str.Snake(model.StructNamePlural)
	// 结构体变量名为结构体名小驼峰形式
	model.VariableName = str.LowerCamel(model.StructName)
	// 结构体变量名复数为结构体系名复数的小驼峰形式
	model.VariableNamePlural = str.LowerCamel(model.StructNamePlural)
	model.PackageName = str.Snake(model.StructName)
	// 注释一般为中文，开发环境未安装中文语言包时可以设置为英文
	model.Comment = comment
	return model
}

// createFileFromStub 读取 stub 文件并进行变量替换
// 最后一个选项可选，如若传参，应传 map[string]string 类型，作为附加的变量搜索替换
func createFileFromStub(filePath string, stubName string, model Model, variables ...interface{}) {

	// 目标文件已存在
	if file.Exists(filePath) {
		console.Exit(filePath + " already exists!")
	}

	// 读取 stub 模板文件
	modelData, err := stubsFS.ReadFile("stubs/" + stubName + ".stub")
	if err != nil {
		console.Exit(err.Error())
	}
	modelStub := string(modelData)

	// 最后一个参数可选，支持自定义变量替换
	replaces := make(map[string]string)
	if len(variables) > 0 {
		replaces = variables[0].(map[string]string)
	}

	// 添加默认的替换变量
	replaces["{{VariableName}}"] = model.VariableName
	replaces["{{VariableNamePlural}}"] = model.VariableNamePlural
	replaces["{{StructName}}"] = model.StructName
	replaces["{{StructNamePlural}}"] = model.StructNamePlural
	replaces["{{PackageName}}"] = model.PackageName
	replaces["{{TableName}}"] = model.TableName
	replaces["{{Comment}}"] = model.Comment
	replaces["{{Author}}"] = config.Get("app.developer", "")
	replaces["{{AuthorEmail}}"] = config.Get("app.developer_email", "")
	replaces["{{CreatedDataTime}}"] = app.TimenowInTimezone().Format("2006-01-02 15:04")
	replaces["{{CopyrightToYear}}"] = app.TimenowInTimezone().Format("2006")

	if len(replaces["{{Author}}"]) == 0 {
		console.Exit("Please set the developer name in the .env file DEVELOPER=")
	}

	if len(replaces["{{AuthorEmail}}"]) == 0 {
		console.Exit("Please set the developer email in the .env file DEVELOPER_EMAIL=")
	}

	// 对模板内容做变量替换
	for search, replace := range replaces {
		modelStub = strings.ReplaceAll(modelStub, search, replace)
	}

	// 存储到目标文件中
	err = file.Put([]byte(modelStub), filePath)
	if err != nil {
		console.Exit(err.Error())
	}

	// 提示成功
	console.Success(fmt.Sprintf("[%s] created.", filePath))
}