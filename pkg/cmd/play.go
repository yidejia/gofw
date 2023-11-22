package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/console"
)

var CmdPlay = &cobra.Command{
	Use:   "play",
	Short: "Likes the Go Playground, but running at our application context",
	Run:   runPlay,
}

// runPlay 调试时运行测试代码
func runPlay(cmd *cobra.Command, args []string) {
	// TODO 可以在这里输入代码，在终端进行验证
	console.Success("Hello " + config.Get("app.name") + ", Application was initialized successfully.")
	console.Warning("Next, GoFW expect you to use code to change the world.")
}
