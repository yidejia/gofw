package migrate

import (
	"database/sql"

	"gorm.io/gorm"
)

// migrationFunc 定义 up 和 down 回调方法的类型
type migrationFunc func(gorm.Migrator, *sql.DB) error

// migrationFiles 所有的迁移文件数组
var migrationFiles []MigrationFile

// MigrationFile 代表着单个迁移文件
type MigrationFile struct {
	Up       migrationFunc
	Down     migrationFunc
	FileName string // 迁移文件名
	Connection string // 执行迁移文件时使用的数据库连接名
}

// Add 新增一个迁移文件，所有的迁移文件都需要调用此方法来注册
func Add(name string, up migrationFunc, down migrationFunc, connection ...string) {
	// 指定了数据库连接
	if len(connection) > 0 {
		migrationFiles = append(migrationFiles, MigrationFile{
			FileName: name,
			Connection: connection[0],
			Up:       up,
			Down:     down,
		})
	} else {
		// 使用默认数据库连接
		migrationFiles = append(migrationFiles, MigrationFile{
			FileName: name,
			Connection: "",
			Up:       up,
			Down:     down,
		})
	}
}

// getMigrationFile 通过迁移文件的名称来获取到 MigrationFile 对象
func getMigrationFile(name string) MigrationFile {
	for _, mFile := range migrationFiles {
		if name == mFile.FileName {
			return mFile
		}
	}
	return MigrationFile{}
}

// isNotMigrated 判断迁移是否已执行
func (mFile MigrationFile) isNotMigrated(migrations []Migration) bool {
	for _, migration := range migrations {
		if migration.Migration == mFile.FileName {
			return false
		}
	}
	return true
}