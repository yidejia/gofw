// Package verifycode 验证码包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 18:06
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package verifycode

type Store interface {
	// Set 保存验证码
	Set(id string, value string) bool

	// Get 获取验证码
	Get(id string, clear bool) string

	// Verify 检查验证码
	Verify(id, answer string, clear bool) bool
}