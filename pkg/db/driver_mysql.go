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

func (my *Mysql) DeleteAllTables(connection string) error {

	dbname := CurrentDatabase(connection)
	var tables []string

	_db := Connection(connection)

	// 读取所有数据表
	err := _db.Table("information_schema.tables").
		Where("table_schema = ?", dbname).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// 暂时关闭外键检测
	_db.Exec("SET foreign_key_checks = 0;")

	// 删除所有表
	for _, table := range tables {
		err = _db.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// 开启 MySQL 外键检测
	_db.Exec("SET foreign_key_checks = 1;")
	return nil
}
