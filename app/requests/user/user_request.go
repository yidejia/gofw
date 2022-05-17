// Package user 封装用户模块请求和验证逻辑
// @author 余海坚 haijianyu10@qq.com
// @created 2022-05-16 19:27
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package user

import (
	"github.com/thedevsaddam/govalidator"
	gfReqs "github.com/yidejia/gofw/pkg/requests"
)

type CreateUserRequest struct {
	gfReqs.Request
	Name string `json:"name" form:"name" valid:"name"`
}

func (req *CreateUserRequest) Validate(data interface{}, extra ...interface{}) map[string][]string  {

	rules := govalidator.MapData{
		"name":          []string{"required", "alpha_num", "between:3,20"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
	}

	return gfReqs.ValidateStruct(data, rules, messages)
}

type UpdateUserRequest struct {
	Name string `json:"name" valid:"name"`
}

func (req *UpdateUserRequest) Validate(data interface{}, extra ...interface{}) map[string][]string  {

	rules := govalidator.MapData{
		"name":          []string{"required", "alpha_num", "between:3,20"},
	}

	messages := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
	}

	return gfReqs.ValidateStruct(data, rules, messages)
}
