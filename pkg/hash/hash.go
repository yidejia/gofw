// Package hash 哈希包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 18:27
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package hash

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/yidejia/gofw/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// Md5 返回一个 32 位 md5 加密后的字符串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Md5To16 返回一个16位md5加密后的字符串
func Md5To16(str string) string {
	return Md5(str)[8:24]
}

// BcryptHash 使用 bcrypt 对密码进行加密
// Bcrypt 是一个很流行的密码哈希算法，是 Niels Provos 和 David Mazières 基于 Blowfish 加密算法设计的密码哈希算法，于1999年在 USENIX 协会上提交。
// Bcrypt 在设计上包含了一个盐 Salt 来防御彩虹表攻击，还提供了一种自适应功能，可以随着时间的推移，通过增加迭代计数以使其执行更慢，使得即便在增加计算能力的情况下，
// Bcrypt 仍然能保持抵抗暴力攻击。
func BcryptHash(password string) string {
	// 先将密码转换成固定32位的的字符串，避免超出 Bcrypt 算法的长度限制
	// 通常为50～72字符，准确的长度限制取决于具体的 Bcrypt 实现。超过最大长度的密码将被截断
	password = Md5(password)
	// GenerateFromPassword 的第二个参数是 cost 值。建议大于 12，数值越大耗费时间越长
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	logger.LogIf(err)

	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	// 先将密码转换成固定32位的的字符串，避免超出 Bcrypt 算法的长度限制
	// 通常为50～72字符，准确的长度限制取决于具体的 Bcrypt 实现。超过最大长度的密码将被截断
	password = Md5(password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// BcryptIsHashed 判断字符串是否是哈希过的数据
func BcryptIsHashed(str string) bool {
	// bcrypt 加密后的长度等于 60
	return len(str) == 60
}
