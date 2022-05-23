package bootstrap

import "github.com/yidejia/gofw/pkg/redis"

// SetupRedis 初始化 Redis
func SetupRedis() {
	redis.InitWithConfig()
}