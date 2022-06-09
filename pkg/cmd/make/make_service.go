package make

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/yidejia/gofw/pkg/console"
    "github.com/yidejia/gofw/pkg/file"
    "os"
)

// CmdMakeService 生成服务文件命令
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-09 15:15
// @copyright © 2010-2022 广州伊的家网络科技有限公司
var CmdMakeService = &cobra.Command{
    Use:     "service",
    Short:   "Create service file",
    Example: "go run main.go make service user -c user",
    Run:     runMakeService,
    Args:    cobra.MinimumNArgs(1), // 至少传 1 个参数
}

func init() {
    CmdMakeService.Flags().BoolP("force" ,"f", false, "Force files to be created in the services root directory")
}

func runMakeService(cmd *cobra.Command, args []string) {

    // 获取注释
    comment := parseCommentFlag(cmd, args)
    // 获取名称
    path, name, pkgName := parseNameParam(cmd, args)
    if len(pkgName) == 0 {
        pkgName = "services"
    }
    // 获取是否强制执行选项
    force := parseForceFlag(cmd, args)

    // 模型目录不存在
    if !file.Exists("app/services") {
        if err := os.Mkdir("app/services", os.ModePerm); err != nil {
            console.Exit(fmt.Sprintf("failed to create services folder: %s", err.Error()))
        }
    }

    // 完整目标文件目录
    var fullPath string
    if len(path) > 0 {
        fullPath = fmt.Sprintf("app/services/%s", path)
    } else {
        if !force {
            console.Exit("The service should belong to its own package. If you want to continue creating, please set the \"-f\" option.")
        }
        fullPath = "app/services"
    }

    // 目标目录不存在
    if !file.Exists(fullPath) {
        if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
            console.Exit(fmt.Sprintf("failed to create service parent folder: %s", err.Error()))
        }
    }

    // 格式化模型名称，返回一个 Model 对象
    model := makeModelFromString(name, comment, pkgName)
    // 拼接目标文件路径
    filePath := fmt.Sprintf("%s/%s_service.go", fullPath, model.PackageName)
    // 从模板中创建文件并进行变量替换
    createFileFromStub(filePath, "service", model)
}