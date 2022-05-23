package db

import (
	"errors"
	"fmt"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"sync"
	"time"
)

// Connector 数据库连接器接口
type Connector interface {
	// Connection 返回数据库连接名
	Connection() string
}

// connections 数据库连接映射表
var connections sync.Map

// InitWithConfig 加载配置初始化数据库连接
func InitWithConfig() {

	var dbConfig gorm.Dialector
	// 数据库驱动名
	var driverName string
	// 数据库驱动
	var driver Driver

	// 遍历数据库连接配置建立数据库连接池
	for name, _ := range config.GetStringMap("database.connections") {

		driverName = config.Get(fmt.Sprintf("database.connections.%v.driver", name))
		switch driverName {
		case "mysql":
			driver = NewMysql()
		case "sqlite":
			driver = NewSqlite()
		default:
			panic(errors.New(fmt.Sprintf("database connection %v not supported %v driver", name, driverName)))
		}

		dbConfig = driver.Connect(name)
		// 生成的配置无效，继续处理下一个连接
		if dbConfig == nil {
			continue
		}
		// 使用 gorm.Open 连接数据库
		db, err := gorm.Open(dbConfig, &gorm.Config{
			Logger: logger.NewGormLogger(),
		})
		// 处理错误
		if err != nil {
			panic(fmt.Sprintf("open database %v error: %v", name, err.Error()))
		}

		// 设置了只读从库
		if read := config.GetStringSlice(fmt.Sprintf("database.connections.%v.read", name)); len(read) > 0 {

			var replicas []gorm.Dialector
			var readDBConfig gorm.Dialector
			for _, host := range read {
				readDBConfig = driver.ConnectToSlave(name, host)
				if readDBConfig == nil {
					continue
				}
				replicas = append(replicas, readDBConfig)
			}

			err = db.Use(dbresolver.Register(dbresolver.Config{
				Replicas: replicas,
				// 负载均衡策略
				Policy: dbresolver.RandomPolicy{},
			}))
			if err != nil {
				panic(fmt.Sprintf("settup database %v slave error: %v", name, err.Error()))
			}
		}

		// 获取底层的 sqlDB
		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Sprintf("get sqlDB %v error: %v", name, err.Error()))
		}
		// 测试数据库连接
		if err = sqlDB.Ping(); err != nil {
			panic(fmt.Sprintf("failed to connect to DB %v, error: %v", name, err.Error()))
		}

		// 设置最大连接数
		sqlDB.SetMaxOpenConns(config.GetInt(fmt.Sprintf("database.connections.%v.max_open_connections", name)))
		// 设置最大空闲连接数
		sqlDB.SetMaxIdleConns(config.GetInt(fmt.Sprintf("database.connections.%v.max_idle_connections", name)))
		// 设置每个连接的过期时间
		sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt(fmt.Sprintf("database.connections.%v.max_life_seconds", name))) * time.Second)

		// 缓存数据库连接
		connections.Store(name, db)
	}
}

// Connection 通过连接名获取数据库连接实例
// 基于 Model 进行数据库操作时，优先使用下面的 DB 方法，只有需要显式使用特定数据库连接时才使用这个方法
func Connection(name string) *gorm.DB {
	db, ok := connections.Load(name)
	if !ok {
		panic(fmt.Sprintf("DB Connection %v not exists", name))
	}
	return db.(*gorm.DB)
}

// DB 通过实现数据库连接器接口的对象获取数据库连接实例
// Model 基类已经默认实现了接口，可以直接使用，不使用默认数据库连接的 Model 需要重新实现连接器接口
func DB(connector Connector) *gorm.DB {
	return Connection(connector.Connection())
}
