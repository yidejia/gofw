package cmd

import (
	"fmt"
	"github.com/yidejia/gofw/pkg/cache"
	"github.com/yidejia/gofw/pkg/console"

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

// CmdCacheForget 删除指定键缓存命令
var CmdCacheForget = &cobra.Command{
	Use:   "forget",
	Short: "Delete redis key, example: cache forget cache-key",
	Run:   runCacheForget,
}

// forget 命令的选项
var cacheKey string

func init() {
	// 注册 cache 命令的子命令
	CmdCache.AddCommand(CmdCacheClear, CmdCacheForget)

	// 设置 cache forget 命令的选项
	CmdCacheForget.Flags().StringVarP(&cacheKey, "key", "k", "", "KEY of the cache")
	_ = CmdCacheForget.MarkFlagRequired("key")
}

func runCacheClear(cmd *cobra.Command, args []string) {
	cache.Flush()
	console.Success("Cache cleared.")
}

func runCacheForget(cmd *cobra.Command, args []string) {
	cache.Forget(cacheKey)
	console.Success(fmt.Sprintf("Cache key [%s] deleted.", cacheKey))
}