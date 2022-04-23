package cmd

import (
	"github.com/yidejia/gofw/cache"
	"github.com/yidejia/gofw/console"

	"github.com/spf13/cobra"
)

// CmdCache 缓存命令
var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Cache management",
}

// CmdCacheClear 清空缓存命令
var CmdCacheClear = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache",
	Run:   runCacheClear,
}

func init() {
	// 注册 cache 命令的子命令
	CmdCache.AddCommand(CmdCacheClear)
}

func runCacheClear(cmd *cobra.Command, args []string) {
	cache.Flush()
	console.Success("Cache cleared.")
}