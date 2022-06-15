package db

import (
	"errors"
	"fmt"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strings"
	"sync"
	"time"
)

// Connector 数据库连接器接口
type Connector interface {
	// Connection 返回数据库连接名
	Connection() string
}

// IModel 模型接口
type IModel interface {
	Connector
	// GetStringID 获取 ID 的字符串格式
	GetStringID() string
	// ModelName 模型名称
	ModelName() string
	// ToMap 将模型转换成映射
	ToMap() map[string]interface{}
}

// connections 数据库连接映射表
var connections sync.Map

// NewDriver 根据数据库连接名和驱动名生成驱动实例
func NewDriver(connection string, driverName string) (driver Driver) {
	switch driverName {
	case "mysql":
		driver = NewMysql()
	case "sqlite":
		driver = NewSqlite()
	default:
		panic(errors.New(fmt.Sprintf("database connection %v not supported %v driver", connection, driverName)))
	}
	return
}

// InitWithConfig 加载配置初始化数据库连接
func InitWithConfig() {

	var dbConfig gorm.Dialector
	// 数据库驱动名
	var driverName string
	// 数据库驱动
	var driver Driver

	// 遍历数据库连接配置建立数据库连接池
	for connection, _ := range config.GetStringMap("database.connections") {

		driverName = config.Get(fmt.Sprintf("database.connections.%v.driver", connection))
		driver = NewDriver(connection, driverName)

		dbConfig = driver.Connect(connection)
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
			panic(fmt.Sprintf("open database %v error: %v", connection, err.Error()))
		}

		// 设置了只读从库
		if read := config.GetStringSlice(fmt.Sprintf("database.connections.%v.read", connection)); len(read) > 0 {

			var replicas []gorm.Dialector
			var readDBConfig gorm.Dialector
			for _, host := range read {
				readDBConfig = driver.ConnectToSlave(connection, host)
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
				panic(fmt.Sprintf("settup database %v slave error: %v", connection, err.Error()))
			}
		}

		// 获取底层的 sqlDB
		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Sprintf("get sqlDB %v error: %v", connection, err.Error()))
		}
		// 测试数据库连接
		if err = sqlDB.Ping(); err != nil {
			panic(fmt.Sprintf("failed to connect to DB %v, error: %v", connection, err.Error()))
		}

		// 设置最大连接数
		sqlDB.SetMaxOpenConns(config.GetInt(fmt.Sprintf("database.connections.%v.max_open_connections", connection)))
		// 设置最大空闲连接数
		sqlDB.SetMaxIdleConns(config.GetInt(fmt.Sprintf("database.connections.%v.max_idle_connections", connection)))
		// 设置每个连接的过期时间
		sqlDB.SetConnMaxLifetime(time.Duration(config.GetInt(fmt.Sprintf("database.connections.%v.max_life_seconds", connection))) * time.Second)

		// 缓存数据库连接
		connections.Store(connection, db)
	}
}

// Connection 通过连接名获取数据库连接实例
// 基于 Model 进行数据库操作时，优先使用下面的 DB 方法，只有需要显式使用特定数据库连接时才使用这个方法
func Connection(name ...string) *gorm.DB {

	var _name string

	// 返回指定数据库连接
	if len(name) > 0 {
		_name = name[0]
		// 清除可能存在的空格
		_name = strings.ReplaceAll(_name, " ", "")
		if _name == "" {
			// 数据库连接名为空字符串，只能返回默认数据库连接
			_name = config.GetString("database.default")
		}
	} else {
		// 返回默认数据库连接
		_name = config.GetString("database.default")
	}

	db, ok := connections.Load(_name)

	if !ok {
		panic(fmt.Sprintf("DB Connection %v not exists", _name))
	}

	return db.(*gorm.DB)
}

// DB 通过实现数据库连接器接口的对象获取数据库连接实例
// Model 基类已经默认实现了接口，可以直接使用，不使用默认数据库连接的 Model 需要重新实现连接器接口
func DB(connector Connector) *gorm.DB {
	return Connection(connector.Connection())
}

// Model 通过实现模型接口的对象获取数据库会话对象
func Model(iModel IModel) *gorm.DB {
	return DB(iModel).Model(iModel)
}

// CurrentDatabase 获取当前数据库名称
// @param connection 数据库连接名
func CurrentDatabase(connection ...string) (dbname string) {
	dbname = Connection(connection...).Migrator().CurrentDatabase()
	return
}

// DeleteAllTables 删除所有表
func DeleteAllTables(connection ...string) error {
	var driverName string
	var driver Driver
	var err error
	// 指定了数据库连接时，只删除这个数据库的所有表
	if len(connection) > 0 {
		driverName = config.Get(fmt.Sprintf("database.connections.%s.driver", connection[0]))
		driver = NewDriver(connection[0], driverName)
		return driver.DeleteAllTables(connection[0])
	} else {
		// 默认删除所有数据库连接对应的数据表
		for _connection, _ := range config.GetStringMap("database.connections") {
			driverName = config.Get(fmt.Sprintf("database.connections.%s.driver", _connection))
			driver = NewDriver(_connection, driverName)
			if err = driver.DeleteAllTables(_connection); err != nil {
				return err
			}
		}
		return nil
	}
}

// TableName 模型对应的数据表名
func TableName(iModel IModel) string {
	stmt := &gorm.Statement{DB: DB(iModel)}
	_ = stmt.Parse(iModel)
	return stmt.Schema.Table
}
