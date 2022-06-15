// Package str 字符串辅助方法包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 16:54
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package str

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// Plural 转为复数 user -> users
func Plural(word string) string {
	return pluralize.NewClient().Plural(word)
}

// Singular 转为单数 users -> user
func Singular(word string) string {
	return pluralize.NewClient().Singular(word)
}

// Snake 转为蛇形 snake_case，如 UserName -> user_name
func Snake(s string) string {
	return strcase.ToSnake(s)
}

// Camel 转为大驼峰 CamelCase，如 user_name -> UserName
func Camel(s string) string {
	return strcase.ToCamel(s)
}

// LowerCamel 转为小驼峰 lowerCamelCase，如 UserName -> userName
func LowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}