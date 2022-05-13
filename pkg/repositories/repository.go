// Package repositories 数据仓库包，封装数据操作相关代码
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-07 16:17
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package repositories

import (
	gfErrors "github.com/yidejia/gofw/pkg/errors"
)

// Repository 数据仓库基类
type Repository struct {
}

// ErrorNotFound 返回资源不存在错误
func (repo *Repository) ErrorNotFound(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorNotFound(message...)
}

// ErrorInternal 返回系统内部错误
// 没有内部错误对象需要返回时，err 可以设置为 nil
func (repo *Repository) ErrorInternal(err error, message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorInternal(err, message...)
}
