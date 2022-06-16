// Package validators 数据验证包
// 这个定义框架自定义验证规则
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 18:22
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package validators

import (
	"errors"
	"fmt"
	"github.com/yidejia/gofw/pkg/db"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/thedevsaddam/govalidator"
)

// Initialize 利用 Go 包加载机制，手动初始化验证包时调用，一般在main入口文件里调用
// 调用示例：
//  func init() {
//	    // 加载验证包并初始化自定义验证规则
//	    validators.Initialize()
//  }
func Initialize() {
}

// 此方法会在初始化时执行，注册自定义验证规则
func init() {

	// 自定义规则 not_exists，验证请求数据必须不存在于数据库中。
	// 常用于保证数据库某个字段的值唯一，如用户名、邮箱、手机号、或者分类的名称。
	// not_exists 参数可以有两种，一种是 2 个参数，一种是 3 个参数：
	// not_exists:users,email 检查数据库表里是否存在同一条信息
	// not_exists:users,email,32 排除用户掉 id 为 32 的用户
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {

		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		// 第一个参数，表名称，如 users
		tableName := rng[0]
		// 判断表名是否带了数据库连接名
		var dbConn string
		tableNames := strings.Split(tableName, ".")
		if len(tableNames) > 1 {
			dbConn = tableNames[0]
			tableName = tableNames[1]
		} else {
			dbConn = ""
		}
		// 第二个参数，字段名称，如 email 或者 phone
		dbFiled := rng[1]

		// 第三个参数，排除 ID
		var exceptID string
		if len(rng) > 2 {
			exceptID = rng[2]
		}

		// 用户请求过来的数据
		requestValue := value.(string)

		// 拼接 SQL
		query := db.Connection(dbConn).Table(tableName).Where(dbFiled+" = ?", requestValue)

		// 如果传参第三个参数，加上 SQL Where 过滤
		if len(exceptID) > 0 {
			query.Where("id != ?", exceptID)
		}

		// 查询数据库
		var count int64
		query.Count(&count)

		// 验证不通过，数据库能找到对应的数据
		if count != 0 {
			// 如果有自定义错误消息的话
			if message != "" {
				return errors.New(message)
			}
			// 默认的错误消息
			return fmt.Errorf("%v 已被占用", requestValue)
		}
		// 验证通过
		return nil
	})

	// max_cn:8 中文长度设定不超过 8
	govalidator.AddCustomRule("max_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "max_cn:"))
		if valLength > l {
			// 如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度不能超过 %d 个字", l)
		}
		return nil
	})

	// min_cn:2 中文长度设定不小于 2
	govalidator.AddCustomRule("min_cn", func(field string, rule string, message string, value interface{}) error {
		valLength := utf8.RuneCountInString(value.(string))
		l, _ := strconv.Atoi(strings.TrimPrefix(rule, "min_cn:"))
		if valLength < l {
			// 如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("长度需大于 %d 个字", l)
		}
		return nil
	})

	// 自定义规则 exists，确保数据库存在某条数据
	// 一个使用场景是创建话题时需要附带 category_id 分类 ID 为参数，此时需要保证
	// category_id 的值在数据库中存在，即可使用：
	// exists:categories,id
	govalidator.AddCustomRule("exists", func(field string, rule string, message string, value interface{}) error {

		rng := strings.Split(strings.TrimPrefix(rule, "exists:"), ",")

		// 第一个参数，表名称，如 categories
		tableName := rng[0]
		// 判断表名是否带了数据库连接名
		var dbConn string
		tableNames := strings.Split(tableName, ".")
		if len(tableNames) > 1 {
			dbConn = tableNames[0]
			tableName = tableNames[1]
		} else {
			dbConn = ""
		}
		// 第二个参数，字段名称，如 id
		dbFiled := rng[1]

		// 用户请求过来的数据
		requestValue := value.(string)

		// 查询数据库
		var count int64
		db.Connection(dbConn).Table(tableName).Where(dbFiled+" = ?", requestValue).Count(&count)
		// 验证不通过，数据不存在
		if count == 0 {
			// 如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("%v 不存在", requestValue)
		}
		return nil
	})

	// 验证是否是手机号
	govalidator.AddCustomRule("mobile", func(field string, rule string, message string, value interface{}) error {

		// 用户请求过来的数据
		requestValue := value.(string)

		if ok, _ := regexp.MatchString("^1[345789]\\d{9}$", requestValue); !ok {
			// 如果有自定义错误消息的话，使用自定义消息
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("%v 不是手机号", requestValue)
		}

		return nil
	})
}