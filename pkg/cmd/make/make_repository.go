package make

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/yidejia/gofw/pkg/console"
    "github.com/yidejia/gofw/pkg/file"
    "os"
)

// CmdMakeRepository 生成数据仓库文件命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-09 10:32
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeRepository = &cobra.Command{
    Use:   "repo",
    Short:  "Create repository file",
    Example: "go run main.go make repo user -c user/user",
    Run: runMakeRepository,
    Args:  cobra.MinimumNArgs(1), // 至少传 1 个参数
}

func init() {
    CmdMakeRepository.Flags().BoolP("force" ,"f", false, "Force files to be created in the repositories root directory")
}

func runMakeRepository(cmd *cobra.Command, args []string) {

    // 获取是否强制执行选项
    force := parseForceFlag(cmd, args)
    // 获取注释
    comment := parseCommentFlag(cmd, args, force)
    // 获取名称
    path, name, pkgName := parseNameParam(cmd, args)
    if len(pkgName) == 0 {
        pkgName = "repositories"
    }

    // 仓库目录不存在
    if !file.Exists("app/repositories") {
        if err := os.Mkdir("app/repositories", os.ModePerm); err != nil {
            console.Exit(fmt.Sprintf("failed to create repositories folder: %s", err.Error()))
        }
    }

    // 完整目标文件目录
    var fullPath string
    if len(path) > 0 {
        fullPath = fmt.Sprintf("app/repositories/%s", path)
    } else {
        if !force {
            console.Exit("The repository should belong to its own package. If you want to continue creating, please set the \"-f\" option.")
        }
        fullPath = "app/repositories"
    }

    // 目标目录不存在
    if !file.Exists(fullPath) {
        if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
            console.Exit(fmt.Sprintf("failed to create repository parent folder: %s", err.Error()))
        }
    }

    // 格式化模型名称，返回一个 Model 对象
    model := makeModelFromString(name, comment, pkgName)
    // 拼接目标文件路径
    filePath := fmt.Sprintf("%s/%s_repository.go", fullPath, model.PackageName)
    // 从模板中创建文件并进行变量替换
    createFileFromStub(filePath, "repository", model)
}