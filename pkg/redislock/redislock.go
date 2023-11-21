// Package redislock redis 分布式锁包
// @author 余海坚 haijianyu10@qq.com
// @created 2023-11-21 15:24
// @copyright © 2010-2023 广州伊的家网络科技有限公司
package redislock

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/yidejia/gofw/pkg/config"
)

// MutexConfig 互斥锁配置
type MutexConfig struct {
	KeyPrefix string  // redis key 前缀
	Factor    float64 // 时间误差系数
}

// Lock 分布锁管理器
type Lock struct {
	sync        *redsync.Redsync // 分布式锁
	mutexConfig *MutexConfig     // 分布式锁配置
}

// Mutex 资源访问互斥体
type Mutex struct {
	Mutex *redsync.Mutex
}

var once sync.Once

// internalLock 内部使用的分布锁管理器单例
var internalLock *Lock

// InitWithConfig 加载配置并初始化任务调度器
func InitWithConfig() {
	once.Do(func() {
		internalLock = &Lock{}
		// 创建 redis 连接池实现分布式锁
		pool := &redis.Pool{
			MaxIdle:     5,
			IdleTimeout: 30 * time.Second,
			Dial: func() (redis.Conn, error) {
				if config.Get("redis.cron.password") == "" {
					return redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Get("redis.lock.host"), config.Get("redis.lock.port")))
				} else {
					return redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Get("redis.lock.host"), config.Get("redis.lock.port")), redis.DialPassword(config.Get("redis.lock.password")))
				}
			},
		}
		internalLock.sync = redsync.New([]redsync.Pool{pool})
		// 加载分布式锁配置
		internalLock.mutexConfig = &MutexConfig{
			KeyPrefix: config.Get("app.name") + ":lock:",
			Factor:    config.GetFloat64("redis.lock.mutex_factor"),
		}
	})
}

// Acquire 申请锁
func Acquire(name string, expireTime time.Duration, tries int) (*Mutex, error) {
	mutex := internalLock.sync.NewMutex(internalLock.mutexConfig.KeyPrefix+name, redsync.SetExpiry(expireTime), redsync.SetTries(tries))
	if err := mutex.Lock(); err != nil {
		return nil, err
	}
	return &Mutex{Mutex: mutex}, nil
}

// Release 释放锁
func Release(mutex *Mutex) (bool, error) {
	return mutex.Mutex.Unlock()
}
