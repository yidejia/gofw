package config

import "github.com/yidejia/gofw/pkg/config"

// 添加 email 配置
func init() {
	config.Add("mail", func() map[string]interface{} {
		return map[string]interface{}{

			// 默认是 Mailhog 的配置
			"smtp": map[string]interface{}{
				"host":     config.Env("MAIL_HOST", "localhost"),
				"port":     config.Env("MAIL_PORT", 1025),
				"username": config.Env("MAIL_USERNAME", ""),
				"password": config.Env("MAIL_PASSWORD", ""),
			},

			"from": map[string]interface{}{
				"address": config.Env("MAIL_FROM_ADDRESS", "gohub@example.com"),
				"name":    config.Env("MAIL_FROM_NAME", "Gohub"),
			},
		}
	})
}