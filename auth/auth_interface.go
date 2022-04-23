package auth

// Authenticate 可认证接口
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 15:30
// @copyright © 2010-2022 广州伊的家网络科技有限公司
type Authenticate interface {
	// AuthId 授权的唯一 id
	AuthId() uint64
	// AuthName 授权的唯一名称
	AuthName() string
}

// UserResolver 获取用户模型的函数类型
type UserResolver func(userId string) (user Authenticate)
