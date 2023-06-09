// Package helpers 辅助函数包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 16:20
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package helpers

import (
	"crypto/rand"
	"fmt"
	"io"
	mathRand "math/rand"
	"reflect"
	"sort"
	"time"

	"github.com/fatih/structs"
	"github.com/yidejia/gofw/pkg/str"
)

// Empty 类似于 PHP 的 empty() 函数
func Empty(val interface{}) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(val, reflect.Zero(v.Type()).Interface())
}

// MicrosecondsStr 将 time.Duration 类型（nano seconds 为单位）
// 输出为小数点后 3 位的 ms （microsecond 毫秒，千分之一秒）
func MicrosecondsStr(elapsed time.Duration) string {
	return fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
}

// RandomNumber 生成长度为 length 的随机数字字符串
func RandomNumber(length int) string {
	table := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// FirstElement 安全地获取 args[0]，避免 panic: runtime error: index out of range
func FirstElement(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// RandomString 生成长度为 length 的随机字符串
func RandomString(length int) string {
	mathRand.Seed(time.Now().UnixNano())
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}
	return string(b)
}

// StructToMap 将结构体转换成映射
func StructToMap(obj interface{}, fields ...string) map[string]interface{} {

	_map := structs.Map(obj)
	newMap := make(map[string]interface{})

	oldFields := make([]string, len(fields))
	copy(oldFields, fields)
	// sort 包使用二分查找算法，为了提高查找效率，需要先进行排序
	if !sort.StringsAreSorted(fields) {
		sort.Strings(fields)
	}

	fieldsLen := len(fields)
	var k string
	var v interface{}
	var i int
	for k, v = range _map {
		k = str.Snake(k)
		i = sort.SearchStrings(fields, k)
		if i < fieldsLen && fields[i] == k {
			newMap[k] = v
		}
	}

	oldMap := make(map[string]interface{})
	var ok bool
	for _, k = range oldFields {
		if v, ok = newMap[k]; ok {
			oldMap[k] = v
		}
	}

	return oldMap
}

// MergeMaps 合并多个映射
func MergeMaps(_map map[string]interface{}, moreMaps ...map[string]interface{}) map[string]interface{} {
	if len(moreMaps) > 0 {
		var k string
		var v interface{}
		var moreMap map[string]interface{}
		for _, moreMap = range moreMaps {
			for k, v = range moreMap {
				_map[k] = v
			}
		}
	}
	return _map
}

// SearchStringInSlice 在切片中查找字符串
func SearchStringInSlice(_slice []string, str string) int {
	// sort 包使用二分查找算法，为了提高查找效率，需要先进行排序
	if !sort.StringsAreSorted(_slice) {
		sort.Strings(_slice)
	}
	if i := sort.SearchStrings(_slice, str); i < len(_slice) && _slice[i] == str {
		return i
	}
	return -1
}
