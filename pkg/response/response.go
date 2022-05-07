// Package response 响应处理工具
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 17:55
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package response

import (
	gfErrors "github.com/yidejia/gofw/pkg/errors"
	"github.com/yidejia/gofw/pkg/paginator"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Created 响应 201 和带 data 键的 JSON 数据
// 执行『新增操作』成功后调用，例如新增资源后返回新增的资源
// @param item 一个资源对象
// @param meta 附加的元数据
func Created(c *gin.Context, item interface{}, meta ...gin.H) {
	if len(meta) > 0 {
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    item,
			"meta":    meta[0],
		})
	} else {
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data":    item,
		})
	}
}

// Item 响应 200 和带 data 键的 JSON 数据
// 执行『查询或更新操作』后返回查询到或已更新的一个资源对象
// @param item 一个资源对象
// @param meta 附加的元数据
func Item(c *gin.Context, item interface{}, meta ...gin.H) {
	if len(meta) > 0 {
		JSON(c, gin.H{
			"success": true,
			"data":    item,
			"meta":    meta[0],
		})
	} else {
		JSON(c, gin.H{
			"success": true,
			"data":    item,
		})
	}
}

// Collection 响应 200 和带 data 键的 JSON 数据
// 执行『查询操作』后返回一个资源集合
// @param collection 一个资源集合
// @param meta 附加的元数据
func Collection(c *gin.Context, collection []interface{}, meta ...gin.H) {
	if len(meta) > 0 {
		JSON(c, gin.H{
			"success": true,
			"data":    collection,
			"meta":    meta[0],
		})
	} else {
		JSON(c, gin.H{
			"success": true,
			"data":    collection,
		})
	}
}

// Paginate 响应 200 和带 data 键的 JSON 数据
// 执行『查询操作』后返回一个资源集合的分页，适用于资源集合比较大的场景
// @param collection 一个资源集合
// @param paging 分页对象
// @param meta 附加的元数据
func Paginate(c *gin.Context, collection []interface{}, paging paginator.Paging, meta ...gin.H) {
	if len(meta) > 0 {
		metaData := meta[0]
		metaData["pagination"] = paging
		JSON(c, gin.H{
			"success": true,
			"data":    collection,
			"meta":    metaData,
		})
	} else {
		JSON(c, gin.H{
			"success": true,
			"data":    collection,
			"meta":    gin.H{
				"pagination": paging,
			},
		})
	}
}

// NoContent 返回一个无实体内容响应
// 执行某个『没有具体返回数据』的『变更』操作成功后调用，例如删除资源、变更资源状态，只需通过响应头判断操作是否成功
func NoContent(c *gin.Context)  {
	c.Status(http.StatusNoContent)
}

// Success 响应 200 和预设『操作成功！』的 JSON 数据
// 优先使用 NoContent 方法，对接方不方便处理响应头或单纯需要返回元数据时才折衷考虑使用此方法
// 执行某个『没有具体返回数据』的『变更』操作成功后调用，例如删除资源、变更资源状态
// @param meta 附加的元数据
func Success(c *gin.Context, meta ...gin.H) {
	if len(meta) > 0 {
		JSON(c, gin.H{
			"success": true,
			"message": "操作成功！",
			"meta":    meta[0],
		})
	} else {
		JSON(c, gin.H{
			"success": true,
			"message": "操作成功！",
		})
	}
}

// Data 响应 200 和带 data 键的 JSON 数据
// 执行『更新操作』成功后调用，例如更新话题，成功后返回已更新的话题
func Data(c *gin.Context, data interface{}) {
	JSON(c, gin.H{
		"success": true,
		"data":    data,
	})
}

// JSON 响应 200 和 JSON 数据
func JSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// AbortWithError 中断处理并返回错误
func AbortWithError(c *gin.Context, err gfErrors.ResponsiveError)  {
	c.AbortWithStatusJSON(err.HttpStatus(), gin.H{
		"success": false,
		"message": err.Error(),
	})
}

// ValidationError 处理表单验证不通过的错误，返回的 JSON 示例：
//         {
//             "errors": {
//                 "phone": [
//                     "手机号为必填项，参数名称 phone",
//                     "手机号长度必须为 11 位的数字"
//                 ]
//             },
//             "message": "请求验证不通过，具体请查看 errors"
//         }
func ValidationError(c *gin.Context, errors map[string][]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"message": "请求验证不通过，具体请查看 errors",
		"errors":  errors,
	})
}