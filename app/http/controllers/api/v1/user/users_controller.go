// Package user 存放用户模块相关控制器的包
package user

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/yidejia/gofw/app/http/controllers/api/v1"
	"github.com/yidejia/gofw/app/models/user"
	userSvcs "github.com/yidejia/gofw/app/services/user"
	"github.com/yidejia/gofw/pkg/response"
)

// UsersController 用户控制器
type UsersController struct {
	v1.APIController
}

// Store 创建用户
func (ctrl *UsersController) Store(c *gin.Context) {

	userSvc := userSvcs.NewUserService()

	if _user, err := userSvc.Create(); err != nil {
		response.AbortWithError(c, err)
	} else {
		response.Item(c, _user)
	}
}

// Index 用户列表
func (ctrl *UsersController) Index(c *gin.Context) {
	users := []user.User{
		{
			Name: "1",
		},
		{
			Name: "2",
		},
		{
			Name: "3",
		},
	}
	response.Data(c, users)
}
