// Package json 工具包
// @author 余海坚 haijianyu10@qq.com
// @created 2023-06-15 20:20
// @copyright © 2010-2023 广州伊的家网络科技有限公司
package json

import (
	"errors"

	"github.com/tidwall/gjson"
)

// BindMap 解析 json 字符串并绑定到 map 结构上
func BindMap(jsonStr string, _map *map[string]interface{}) error {
	m, ok := gjson.Parse(jsonStr).Value().(map[string]interface{})
	if !ok {
		return errors.New("json can't bind to a map")
	}
	*_map = m
	return nil
}
