// Package repositories 应用数据仓库包，封装数据操作相关代码
package repositories

import (
	gfRepos "github.com/yidejia/gofw/pkg/repositories"
)

// Repository 应用数据仓库基类，内嵌了框架数据仓库基类，可以根据应用需要进行扩展
type Repository struct {
	gfRepos.Repository
}
