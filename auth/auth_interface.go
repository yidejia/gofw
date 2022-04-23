package auth

// Authenticate 可认证接口
type Authenticate interface {
	// AuthId 授权的唯一 id
	AuthId() uint64
	// AuthName 授权的唯一名称
	AuthName() string
}

// UserResolver 获取用户模型的函数类型
type UserResolver func(userId string) (user Authenticate)
