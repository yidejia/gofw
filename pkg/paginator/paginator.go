// Package paginator 处理分页逻辑
package paginator

import (
	"github.com/yidejia/gofw/pkg/logger"
	"gorm.io/gorm"
	"math"
)

// Paging 分页数据
type Paging struct {
	CurrentPage int    // 当前页
	PerPage     int    // 每页条数
	TotalPage   int    // 总页数
	TotalCount  int64  // 总条数
}

// Paginator 分页操作类
type Paginator struct {
	query *gorm.DB     // db query 指针
	Page       int    // 当前页
	PerPage    int    // 每页条数
	TotalCount int64  // 总条数
	TotalPage  int    // 总页数 = TotalCount/PerPage
	Offset     int    // 数据库读取数据时 Offset 的值
}

// Paginate 分页
// @param db GORM 查询句柄，用以查询数据集和获取数据总数
// @param models 模型数组，传址获取数据
// @param page 当前页码
// @param PerPage 每页条数
// 用法:
//  query := database.DB.Model(Model{}).Where("category_id = ?", modelId)
//  var models []Model
//  paging := paginator.Paginate(query, &models, page, perPage)
func Paginate(db *gorm.DB, models interface{}, page, perPage int) Paging {

	// 初始化 Paginator 实例
	p := &Paginator{
		query: db,
		Page: page,
		PerPage: perPage,
	}
	p.TotalCount = p.getTotalCount()
	p.TotalPage = p.getTotalPage()
	p.Offset = (p.Page - 1) * p.PerPage

	// 查询数据库
	err := p.query.
		Limit(p.PerPage).
		Offset(p.Offset).
		Find(models).
		Error

	// 数据库出错
	if err != nil {
		logger.LogIf(err)
		return Paging{}
	}

	return Paging{
		CurrentPage: p.Page,
		PerPage:     p.PerPage,
		TotalPage:   p.TotalPage,
		TotalCount:  p.TotalCount,
	}
}

// getTotalCount 返回的是数据库里的条数
func (p *Paginator) getTotalCount() int64 {
	var count int64
	if err := p.query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// getTotalPage 计算总页数
func (p Paginator) getTotalPage() int {
	if p.TotalCount == 0 {
		return 0
	}
	nums := int64(math.Ceil(float64(p.TotalCount) / float64(p.PerPage)))
	if nums == 0 {
		nums = 1
	}
	return int(nums)
}