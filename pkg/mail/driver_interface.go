package mail

// Driver email 接口
type Driver interface {
	// Send 检查验证码
	Send(email Email, config map[string]string) bool
}