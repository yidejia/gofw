// Package mail email 包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 18:13
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package mail

import (
	"github.com/yidejia/gofw/pkg/config"
	"sync"
)

type From struct {
	Address string
	Name    string
}

// Email email 信息
type Email struct {
	From    From
	To      []string
	Bcc     []string
	Cc      []string
	Subject string
	Text    []byte // Plaintext message (optional)
	HTML    []byte // Html message (optional)
}

// Mailer email 操作对象
type Mailer struct {
	Driver Driver
}

var once sync.Once
var internalMailer *Mailer

// NewMailer 单例模式获取
func NewMailer() *Mailer {
	once.Do(func() {
		internalMailer = &Mailer{
			Driver: &SMTP{},
		}
	})

	return internalMailer
}

// Send 发布 email
func (mailer *Mailer) Send(email Email) bool {
	return mailer.Driver.Send(email, config.GetStringMapString("mail.smtp"))
}