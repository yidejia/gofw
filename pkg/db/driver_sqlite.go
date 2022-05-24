package db

import (
	"fmt"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/file"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Sqlite 数据库驱动
type Sqlite struct {
}

func NewSqlite() *Sqlite {
	return &Sqlite{}
}

func (s *Sqlite) Connect(name string) gorm.Dialector {
	// 初始化 sqlite
	dbFile := config.Get(fmt.Sprintf("database.connections.%v.database", name))
	// 数据库文件不存在，返回无效配置
	if !file.Exists(dbFile) {
		return nil
	}
	dbConfig := sqlite.Open(dbFile)
	return dbConfig
}

func (s *Sqlite) ConnectToSlave(name string, host string) gorm.Dialector {
	// 数据库文件不存在，返回无效配置
	if !file.Exists(host) {
		return nil
	}
	// 将 host 当文件路径直接打开即可
	dbConfig := sqlite.Open(host)
	return dbConfig
}

func (s *Sqlite) DeleteAllTables(connection string) error {

	dbFile := config.Get(fmt.Sprintf("database.connections.%v.database", connection))
	// 数据库文件不存在，直接返回
	if !file.Exists(dbFile) {
		return nil
	}

	var tables []string

	_db := Connection(connection)

	// 读取所有数据表
	err := _db.Select(&tables, "SELECT name FROM sqlite_master WHERE type='table'").Error
	if err != nil {
		return err
	}

	// 删除所有表
	for _, table := range tables {
		err = _db.Migrator().DropTable(table)
		if err != nil {
			return err
		}
	}
	return err
}
