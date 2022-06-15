// Package repositories 数据仓库包
// 封装数据库操作相关代码
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-07 16:17
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package repositories

import (
	"errors"
	"fmt"
	"github.com/yidejia/gofw/pkg/db"
	gfErrors "github.com/yidejia/gofw/pkg/errors"
	"gorm.io/gorm"
)

// Repository 数据仓库基类
type Repository struct {
}

// NewErrorNotFound 生成资源不存在错误
func (repo *Repository) NewErrorNotFound(message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorNotFound(message...)
}

// NewErrorInternal 生成系统内部错误
// 没有内部错误对象需要返回时，err 可以设置为 nil
func (repo *Repository) NewErrorInternal(err error, message ...string) gfErrors.ResponsiveError {
	return gfErrors.NewErrorInternal(err, message...)
}

// NewError 自动生成合适的错误
// 一般用于获取单个模块的场景，自动判断是模型不存在还是查询出错
func (repo *Repository) NewError(err error, iModel db.IModel, message ...string) gfErrors.ResponsiveError {
	// 模型不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if iModel != nil {
			return repo.NewErrorNotFound(iModel.ModelName() + "不存在")
		} else {
			return repo.NewErrorNotFound(message...)
		}
	} else {
		// 查询出错
		if iModel != nil {
			// 设置了日志消息
			if len(message) > 1 {
				return repo.NewErrorInternal(err, fmt.Sprintf("获取%s失败", iModel.ModelName()), message[0], message[1])
			} else {
				return repo.NewErrorInternal(err, fmt.Sprintf("获取%s失败", iModel.ModelName()))
			}
		} else {
			return repo.NewErrorInternal(err, message...)
		}
	}
}
