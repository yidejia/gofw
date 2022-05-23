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
