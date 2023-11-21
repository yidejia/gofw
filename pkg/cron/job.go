package cron

import cronPKG "github.com/robfig/cron/v3"

// Job 定时任务接口
type Job interface {
	cronPKG.Job
	// Spec 返回定时任务执行时间格式
	Spec() string
	// Name 返回定时任务名称
	Name() string
}
