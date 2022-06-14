// Package maptool 映射工具包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-14 16:40
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package maptool

import (
	"sort"
)

// SortIndictOrder 按字典序排序
func SortIndictOrder(mapData map[string]interface{}) (keys []string) {
	keys = make([]string, 0)
	for key, _ := range mapData {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}
