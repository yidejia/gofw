// Package maptool 映射工具包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-14 16:40
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package maptool

import (
	"sort"
	"strings"

	"github.com/spf13/cast"
)

// SortedMap 键名已排序映射
type SortedMap struct {
	MapData map[string]interface{}
	Keys    []string
}

// KSort 对映射按键名进行升序排序
func KSort(mapData map[string]interface{}) *SortedMap {
	sortedMap := &SortedMap{
		MapData: mapData,
		Keys:    make([]string, 0, len(mapData)),
	}
	for k, _ := range mapData {
		sortedMap.Keys = append(sortedMap.Keys, k)
	}
	sort.Strings(sortedMap.Keys)
	return sortedMap
}

// WalkByKSort 对映射按键名进行升序排序后遍历处理
func WalkByKSort(mapData map[string]interface{}, f func(k string, v interface{})) {
	sortedMap := KSort(mapData)
	for _, k := range sortedMap.Keys {
		f(k, sortedMap.MapData[k])
	}
}

// KSortRecursive 对映射按键名进行递归升序排序
func KSortRecursive(mapData map[string]interface{}) *SortedMap {
	sortedMap := &SortedMap{
		MapData: mapData,
		Keys:    make([]string, 0, len(mapData)),
	}
	for k, v := range mapData {
		if _map, ok := v.(map[string]interface{}); ok {
			sortedMap.MapData[k] = KSortRecursive(_map)
		}
		sortedMap.Keys = append(sortedMap.Keys, k)
	}
	sort.Strings(sortedMap.Keys)
	return sortedMap
}

// WalkByKSortRecursive 对映射按键名进行递归升序排序后遍历处理
func WalkByKSortRecursive(mapData map[string]interface{}, f func(k string, v interface{})) {
	sortedMap := KSort(mapData)
	for _, k := range sortedMap.Keys {
		if v, ok := sortedMap.MapData[k].(map[string]interface{}); ok {
			WalkByKSortRecursive(v, f)
		} else {
			f(k, sortedMap.MapData[k])
		}
	}
}

// MapTrim 删除映射字符串某些特定字符
func MapTrim(mapData map[string]interface{}, cutset string) map[string]interface{} {
	for k, v := range mapData {
		switch v.(type) {
		case string:
			mapData[k] = strings.Trim(cast.ToString(v), cutset)
		case map[string]interface{}:
			mapData[k] = MapTrim(cast.ToStringMap(v), cutset)
		}
	}
	return mapData
}

// SortIndictOrder 按字典序排序
func SortIndictOrder(mapData map[string]interface{}) (keys []string) {
	keys = make([]string, 0, len(mapData))
	for key, _ := range mapData {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}
