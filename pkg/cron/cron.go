// Package cron 定时任务包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-07-04 15:20
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package cron

import (
	"fmt"
	"sync"
	"time"

	"github.com/yidejia/gofw/pkg/logger"

	"github.com/yidejia/gofw/pkg/app"

	"github.com/go-redsync/redsync"

	"github.com/gomodule/redigo/redis"
	cronPkg "github.com/robfig/cron/v3"
	"github.com/yidejia/gofw/pkg/config"
)

// MutexConfig 互斥锁配置
type MutexConfig struct {
	KeyPrefix string  // redis key 前缀
	Factor    float64 // 时间误差系数
}

// Cron 定时任务调度器
type Cron struct {
	runner      *cronPkg.Cron    // 定时任务调度器实例
	parser      *cronPkg.Parser  // 时间格式解析器
	sync        *redsync.Redsync // 分布式锁
	mutexConfig *MutexConfig     // 分布式锁配置
}

var once sync.Once

// internalCron 内部使用 Cron 单例
var internalCron *Cron

// InitWithConfig 加载配置并初始化任务调度器
func InitWithConfig() {
	once.Do(func() {
		// 创建指定时区的任务调度器
		timezone, err := time.LoadLocation(config.Get("app.timezone"))
		if err != nil {
			panic("init cron with timezone failed:" + err.Error())
		}
		// 创建时间规格式解析器，启用秒计时单位，秒可选
		_parser := cronPkg.NewParser(
			cronPkg.SecondOptional | cronPkg.Minute | cronPkg.Hour | cronPkg.Dom | cronPkg.Month | cronPkg.Dow | cronPkg.Descriptor,
		)
		internalCron = &Cron{
			runner: cronPkg.New(
				cronPkg.WithLocation(timezone), // 指定时区
				cronPkg.WithParser(_parser),    // 指定时间格式解析器
			),
			parser: &_parser,
		}
		// 创建 redis 连接池实现分布式锁
		pool := &redis.Pool{
			MaxIdle:     5,
			IdleTimeout: 30 * time.Second,
			Dial: func() (redis.Conn, error) {
				if config.Get("redis.cron.password") == "" {
					return redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Get("redis.cron.host"), config.Get("redis.cron.port")))
				} else {
					return redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Get("redis.cron.host"), config.Get("redis.cron.port")), redis.DialPassword(config.Get("redis.cron.password")))
				}
			},
		}
		internalCron.sync = redsync.New([]redsync.Pool{pool})
		// 加载分布式锁配置
		internalCron.mutexConfig = &MutexConfig{
			KeyPrefix: config.Get("app.name") + ":cron:",
			Factor:    config.GetFloat64("redis.cron.mutex_factor"),
		}
	})
}

// lock 加锁任务，避免分布式环境下多次执行
func lock(job Job) error {
	// 解析时间规格
	schedule, err := internalCron.parser.Parse(job.Spec())
	if err != nil {
		return err
	}
	now := app.TimenowInTimezone()
	d := schedule.Next(now).Sub(now)
	d = d - time.Duration(float64(d)*internalCron.mutexConfig.Factor)
	// 创建一个指定时间后过期的互斥锁
	mutex := internalCron.sync.NewMutex(internalCron.mutexConfig.KeyPrefix+job.Name(), redsync.SetExpiry(d), redsync.SetTries(1))
	if err = mutex.Lock(); err != nil {
		return err
	}
	return nil
}

// wrapJobToFunc 将任务包装成函数，同时对任务加锁，避免分布式环境下，同一个任务在同一时间被多处调用
func wrapJobToFunc(job Job) func() {
	return func() {
		if err := lock(job); err != nil {
			logger.Warn("lock job error:" + err.Error())
			return
		}
		job.Run()
	}
}

// AddJob 添加任务
func AddJob(job Job) {
	if _, err := internalCron.runner.AddFunc(job.Spec(), wrapJobToFunc(job)); err != nil {
		panic("cron add job failed:" + err.Error())
	}
}

// Run 启动一个新协程运行定时任务
func Run() {
	go internalCron.runner.Run()
}
