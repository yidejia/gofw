package db

import (
	"fmt"
	"github.com/yidejia/gofw/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Mysql 数据库驱动
type Mysql struct {
}

func NewMysql() *Mysql {
	return &Mysql{}
}

func (my *Mysql) Connect(name string) gorm.Dialector {
	// 构建 DSN 信息
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.Get(fmt.Sprintf("database.connections.%v.username", name)),
		config.Get(fmt.Sprintf("database.connections.%v.password", name)),
		config.Get(fmt.Sprintf("database.connections.%v.host", name)),
		config.Get(fmt.Sprintf("database.connections.%v.port", name)),
		config.Get(fmt.Sprintf("database.connections.%v.database", name)),
		config.Get(fmt.Sprintf("database.connections.%v.charset", name)),
	)
	dbConfig := mysql.New(mysql.Config{
		DSN: dsn,
	})
	return dbConfig
}

func (my *Mysql) ConnectToSlave(name string, host string) gorm.Dialector {
	// 构建 DSN 信息
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
		config.Get(fmt.Sprintf("database.connections.%v.username", name)),
		config.Get(fmt.Sprintf("database.connections.%v.password", name)),
		host,
		config.Get(fmt.Sprintf("database.connections.%v.port", name)),
		config.Get(fmt.Sprintf("database.connections.%v.database", name)),
		config.Get(fmt.Sprintf("database.connections.%v.charset", name)),
	)
	dbConfig := mysql.New(mysql.Config{
		DSN: dsn,
	})
	return dbConfig
}
