// Package make 封装常用结构体和文件生成命令的包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:00
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package make

import (
	"embed"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/file"
	"github.com/yidejia/gofw/pkg/str"
)

// Model 模板文件填充模型，参数解释
type Model struct {
	TableName          string // 生成模型结构体时对应的数据表名
	StructName         string // 结构体名称
	StructNamePlural   string // 结构体名称复数形式
	VariableName       string // 结构体变量名
	VariableNamePlural string // 结构体变量名复数形式
	PackageName        string // 结构体所属包名
	Comment            string // 结构体注释
	ModuleComment      string // 模块注释
	CustomPackageName  string // 自定义包名，设置后模板填充数据时优先使用自定义包名
}

// stubsFS 方便我们后面打包这些 .stub 为后缀名的文件
//go:embed stubs
var stubsFS embed.FS

// CmdMake 生成模板文件命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:00
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMake = &cobra.Command{
	Use:   "make",
	Short: "Generate file and code",
}

func init() {
	// 注册 make 的子命令
	CmdMake.AddCommand(
		CmdMakeCMD,
		CmdMakeModel,
		CmdMakeRepository,
		CmdMakeService,
		CmdMakeAPIController,
		CmdMakeRequest,
		CmdMakeMigration,
		CmdMakeFactory,
		CmdMakeSeeder,
	)

	for _, subCmd := range CmdMake.Commands() {
		// 为子命令设置注释选项
		subCmd.Flags().StringP("comment", "c", "", "set comment for struct or file")
	}
}

// parseCommentFlag 解析命令注释选项
func parseCommentFlag(cmd *cobra.Command, args []string, force bool) string {

	// 检查是否设置了注释
	comment, err := cmd.Flags().GetString("comment")
	if err != nil {
		console.Exit(fmt.Sprintf("Get comment flag error: %s", err.Error()))
	}

	if len(comment) == 0 {
		console.Exit("Missing comment for struct or file, please use \"-c\" to set comment flag")
	}

	comments := strings.Split(comment, "/")
	if len(comments) < 2 && !force {
		console.Exit("Comments should contain the module name, like -c user/user")
	}

	return comment
}

// parseForceFlag 解析强制执行选项
func parseForceFlag(cmd *cobra.Command, args []string) bool {
	// 检查是否设置了强制执行选项
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		console.Exit(fmt.Sprintf("Get force flage error: %s", err.Error()))
	}
	return force
}

// parseNameParam 解析命令名称参数
func parseNameParam(cmd *cobra.Command, args []string) (path string, name string, pkgName string) {
	names := strings.Split(args[0], "/")
	namesLen := len(names)
	if namesLen == 1 {
		path = ""
		pkgName = ""
		name = names[0]
	} else {
		name = names[namesLen-1]
		pkgName = names[namesLen-2]
		path = strings.Join(names[0:namesLen-1], "/")
	}
	return
}

// makeModelFromString 格式化用户输入的内容并设置好模板文件的填充数据
func makeModelFromString(name string, comment string, pkgName string) Model {
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
	comments := strings.Split(comment, "/")
	// 包含模块注释
	if len(comments) > 1 {
		model.ModuleComment = comments[0]
		model.Comment = comments[1]
	} else {
		// 只包含结构体和文件注释
		model.Comment = comment
	}
	model.CustomPackageName = pkgName
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
	replaces["{{ModuleComment}}"] = model.ModuleComment
	replaces["{{CustomPackageName}}"] = model.CustomPackageName
	replaces["{{Author}}"] = config.Get("app.developer", "")
	replaces["{{AuthorEmail}}"] = config.Get("app.developer_email", "")
	replaces["{{CreatedDataTime}}"] = app.TimeNowInTimezone().Format("2006-01-02 15:04")
	replaces["{{CopyrightFromYear}}"] = config.GetString("app.copyright_from_year", app.TimeNowInTimezone().Format("2006"))
	replaces["{{CopyrightToYear}}"] = app.TimeNowInTimezone().Format("2006")
	replaces["{{AuthorCompany}}"] = config.Get("app.developer_company", "非常牛逼有限公司")
	replaces["{{AppName}}"] = config.Get("app.name")

	if replaces["{{Author}}"] == "" {
		console.Exit("Please set the developer name in the .env file with DEVELOPER=")
	}

	if replaces["{{AuthorEmail}}"] == "" {
		console.Exit("Please set the developer email in the .env file with DEVELOPER_EMAIL=")
	}

	// 设置默认版本号
	if replaces["{{Version}}"] == "" {
		replaces["{{Version}}"] = "v1"
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
