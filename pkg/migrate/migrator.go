// Package migrate 处理数据库迁移
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 17:37
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package migrate

import (
	"database/sql"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/console"
	"github.com/yidejia/gofw/pkg/db"
	"github.com/yidejia/gofw/pkg/file"
	"io/ioutil"
	"os"

	"gorm.io/gorm"
)

// Migrator 数据迁移操作类
type Migrator struct {
	Folder   string // 迁移文件存放目录
	DB       *gorm.DB
	Migrator gorm.Migrator
}

// Migration 对应数据的 migrations 表里的一条数据
type Migration struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	Migration string `gorm:"type:varchar(255);not null;unique;"`
	Batch     int
}

// NewMigrator 创建 Migrator 实例，用以执行迁移操作
func NewMigrator() *Migrator {

	// 默认数据库连接
	_db := db.Connection()
	// 工作目录
	wd, _ := os.Getwd()

	// 初始化必要属性
	migrator := &Migrator{
		Folder:   wd + "/" + config.Get("database.migration_folder"),
		DB:       _db,
		Migrator: _db.Migrator(),
	}
	// migrations 不存在的话就创建它
	migrator.createMigrationsTable()

	return migrator
}

// 创建 migrations 表
func (migrator *Migrator) createMigrationsTable() {

	migration := Migration{}

	// 不存在才创建
	if !migrator.Migrator.HasTable(&migration) {
		_ = migrator.Migrator.CreateTable(&migration)
	}
}

// Up 执行所有未迁移过的文件
func (migrator *Migrator) Up() {

	// 读取所有迁移文件，确保按照时间排序
	migrateFiles := migrator.readAllMigrationFiles()

	// 获取当前批次的值
	batch := migrator.getBatch()

	// 获取所有迁移数据
	var migrations []Migration
	migrator.DB.Find(&migrations)

	// 可以通过此值来判断数据库是否已是最新
	ran := false

	// 对迁移文件进行遍历，如果没有执行过，就执行 up 回调
	for _, mFile := range migrateFiles {

		// 对比文件名称，看是否已经运行过
		if mFile.isNotMigrated(migrations) {
			migrator.runUpMigration(mFile, batch)
			ran = true
		}
	}

	if !ran {
		console.Success("database is up to date.")
	}
}

// Rollback 回滚上一个操作
func (migrator *Migrator) Rollback() {

	// 获取最后一批次的迁移数据
	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)
	var migrations []Migration
	migrator.DB.Where("batch = ?", lastMigration.Batch).Order("id DESC").Find(&migrations)

	// 回滚最后一批次的迁移
	if !migrator.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to rollback.")
	}
}

// 回退迁移，按照倒序执行迁移的 down 方法
func (migrator *Migrator) rollbackMigrations(migrations []Migration) bool {

	// 标记是否真的有执行了迁移回退的操作
	ran := false

	for _, _migration := range migrations {

		// 友好提示
		console.Warning("rollback " + _migration.Migration)

		// 执行迁移文件的 down 方法
		mFile := getMigrationFile(_migration.Migration)

		if mFile.Down != nil {

			DB, sqlDB := migrator.getDBFromMigrationFile(mFile)

			mFile.Down(DB.Migrator(), sqlDB)
		}

		ran = true

		// 回退成功了就删除掉这条记录
		migrator.DB.Delete(&_migration)

		// 打印运行状态
		console.Success("finish " + mFile.FileName)
	}
	return ran
}

// 获取当前这个批次的值
func (migrator *Migrator) getBatch() int {

	// 默认为 1
	batch := 1

	// 取最后执行的一条迁移数据
	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)

	// 如果有值的话，加一
	if lastMigration.ID > 0 {
		batch = lastMigration.Batch + 1
	}
	return batch
}

// 从文件目录读取文件，保证正确的时间排序
func (migrator *Migrator) readAllMigrationFiles() []MigrationFile {

	// 读取 database/migrations/ 目录下的所有文件
	// 默认是会按照文件名称进行排序
	files, err := ioutil.ReadDir(migrator.Folder)
	console.ExitIf(err)

	var migrateFiles []MigrationFile
	for _, f := range files {

		// 去除文件后缀 .go
		fileName := file.FileNameWithoutExtension(f.Name())

		// 通过迁移文件的名称获取『MigrationFile』对象
		mFile := getMigrationFile(fileName)

		// 加个判断，确保迁移文件可用，再放进 migrateFiles 数组中
		if len(mFile.FileName) > 0 {
			migrateFiles = append(migrateFiles, mFile)
		}
	}

	// 返回排序好的『MigrationFile』数组
	return migrateFiles
}

// 从迁移文件中获取数据库实例
func (migrator *Migrator) getDBFromMigrationFile(mFile MigrationFile) (*gorm.DB, *sql.DB) {

	var _db *gorm.DB

	if len(mFile.Connection) > 0 {
		_db = db.Connection(mFile.Connection)
	} else {
		_db = db.Connection()
	}

	sqlDB, _ := _db.DB()

	return _db, sqlDB
}

// 执行迁移，执行迁移的 up 方法
func (migrator *Migrator) runUpMigration(mFile MigrationFile, batch int) {

	// 执行 up 函数的 SQL
	if mFile.Up != nil {

		DB, sqlDB := migrator.getDBFromMigrationFile(mFile)

		// 友好提示
		console.Warning("migrating " + mFile.FileName)
		// 执行 up 函数
		mFile.Up(DB.Migrator(), sqlDB)
		// 提示已迁移了哪个文件
		console.Success("migrated " + mFile.FileName)
	}

	// 入库
	err := migrator.DB.Create(&Migration{Migration: mFile.FileName, Batch: batch}).Error
	console.ExitIf(err)
}

// Reset 回滚所有迁移
func (migrator *Migrator) Reset() {

	var migrations []Migration

	// 按照倒序读取所有迁移文件
	migrator.DB.Order("id DESC").Find(&migrations)

	// 回滚所有迁移
	if !migrator.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to reset.")
	}
}

// Refresh 回滚所有迁移，并运行所有迁移
func (migrator *Migrator) Refresh() {

	// 回滚所有迁移
	migrator.Reset()

	// 再次执行所有迁移
	migrator.Up()
}

// Fresh Drop 所有的表并重新运行所有迁移
func (migrator *Migrator) Fresh(connection ...string) {

	var err error
	// 删除指定数据库所有表
	if len(connection) > 0 {
		err = db.DeleteAllTables(connection...)
	} else {
		// 未指定数据库连接时只删除默认数据库的所有数据表
		err = db.DeleteAllTables(config.Get("database.default"))
	}
	console.ExitIf(err)
	console.Success("clearup database")

	// 重新创建 migrates 表
	migrator.createMigrationsTable()
	console.Success("[migrations] table created.")

	// 重新调用 up 命令
	migrator.Up()
}