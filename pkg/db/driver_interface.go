package db

import "gorm.io/gorm"

// Driver 数据库驱动接口
type Driver interface {

	// Connect 连接数据库
	// @param name 数据库连接名
	Connect(name string) gorm.Dialector

	// ConnectToSlave 连接到只读从库
	// @param name 数据库连接名
	// @param host 从库地址
	ConnectToSlave(name string, host string) gorm.Dialector

	// DeleteAllTables 删除所有数据表
	DeleteAllTables(connection string) error
}
