// Package bootstrap 启动程序功能
package bootstrap

import (
	"fmt"
	"github.com/yidejia/gofw/pkg/cache"
	"github.com/yidejia/gofw/pkg/config"
)

// SetupCache 初始化缓存
func SetupCache() {

	// 初始化缓存专用的 redis client, 使用专属缓存 DB
	rds := cache.NewRedisStore(
		fmt.Sprintf("%v:%v", config.GetString("redis.cache.host"), config.GetString("redis.cache.port")),
		config.GetString("redis.cache.username"),
		config.GetString("redis.cache.password"),
		config.GetInt("redis.cache.database"),
	)

	cache.InitWithCacheStore(rds)
}