package captcha

import (
	"errors"
	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/redis"
	"time"
)

// RedisStore 实现 base64Captcha.Store interface
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 17:49
// @copyright © 2010-2022 广州伊的家网络科技有限公司
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

// Set 实现 base64Captcha.Store interface 的 Set 方法
func (s *RedisStore) Set(key string, value string) error {

	ExpireTime := time.Minute * time.Duration(config.GetInt64("captcha.expire_time"))
	// 方便本地开发调试
	if app.IsLocal() {
		ExpireTime = time.Minute * time.Duration(config.GetInt64("captcha.debug_expire_time"))
	}

	if ok := s.RedisClient.Set(s.KeyPrefix+key, value, ExpireTime); !ok {
		return errors.New("无法存储图片验证码答案")
	}
	return nil
}

// Get 实现 base64Captcha.Store interface 的 Get 方法
func (s *RedisStore) Get(key string, clear bool) string {
	key = s.KeyPrefix + key
	val := s.RedisClient.Get(key)
	if clear {
		s.RedisClient.Del(key)
	}
	return val
}

// Verify 实现 base64Captcha.Store interface 的 Verify 方法
func (s *RedisStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}