// Package database 数据库包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 16:43
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/yidejia/gofw/pkg/config"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// DB 对象
var DB *gorm.DB
var SQLDB *sql.DB

// Connect 连接数据库
func Connect(dbConfig gorm.Dialector, _logger gormLogger.Interface) {

	// 使用 gorm.Open 连接数据库
	var err error
	DB, err = gorm.Open(dbConfig, &gorm.Config{
		Logger: _logger,
	})
	// 处理错误
	if err != nil {
		fmt.Println(err.Error())
	}

	// 获取底层的 sqlDB
	SQLDB, err = DB.DB()
	if err != nil {
		fmt.Println(err.Error())
	}
}

// CurrentDatabase 获取当前数据库名称
func CurrentDatabase() (dbname string) {
	dbname = DB.Migrator().CurrentDatabase()
	return
}

// DeleteAllTables 删除所有表
func DeleteAllTables() error {
	var err error
	switch config.Get("database.connection") {
	case "mysql":
		err = deleteMySQLTables()
	case "sqlite":
		err = deleteAllSqliteTables()
	default:
		panic(errors.New("database connection not supported"))
	}

	return err
}

// deleteAllSqliteTables 删除所有 Sqlite 表
func deleteAllSqliteTables() error {
	var tables []string

	// 读取所有数据表
	err := DB.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}

	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return err
}

// deleteMySQLTables 删除所有 MySQL 表
func deleteMySQLTables() error {
	dbname := CurrentDatabase()
	var tables []string

	// 读取所有数据表
	err := DB.Table("information_schema.tables").
		Where("table_schema = ?", dbname).
		Pluck("table_name", &tables).
		Error
	if err != nil {
		return err
	}

	// 暂时关闭外键检测
	DB.Exec("SET foreign_key_checks = 0;")

	// 删除所有表
	for _, table := range tables {
		err := DB.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}

	// 开启 MySQL 外键检测
	DB.Exec("SET foreign_key_checks = 1;")
	return nil
}

// TableName 模型对应的数据表名
func TableName(obj interface{}) string {
	stmt := &gorm.Statement{DB: DB}
	_ = stmt.Parse(obj)
	return stmt.Schema.Table
}