// Package jwt JWT 认证包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 19:27
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/yidejia/gofw/pkg/app"
	"github.com/yidejia/gofw/pkg/config"
	"github.com/yidejia/gofw/pkg/logger"
	"github.com/yidejia/gofw/pkg/redis"

	"github.com/gin-gonic/gin"
	jwtPKG "github.com/golang-jwt/jwt"
)

var (
	ErrTokenMissing           = errors.New("请求缺少令牌")
	ErrTokenExpired           = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         = errors.New("请求令牌格式有误")
	ErrTokenInvalid           = errors.New("请求令牌无效")
	ErrHeaderEmpty            = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        = errors.New("请求头中令牌格式有误")
	ErrTokenUnsupported       = errors.New("令牌不支持解析")
)

// Driver JWT 驱动接口
type Driver interface {
	// PreParserToken 预处理 token
	PreParserToken(tokenString string, context interface{}) (string, error)
	// NewTokenID 生成新的 token id
	NewTokenID(claims *JWTCustomClaims) (string, error)
	// SaveToken 保存 token
	SaveToken(tokenStr string, claims *JWTCustomClaims) error
	// InvalidateToken 使 token 失效
	InvalidateToken(claims *JWTCustomClaims) error
	// TraverseInvalidTokensNotExpired 遍历未过期的失效 token
	TraverseInvalidTokensNotExpired(handler func(claims *JWTCustomClaims) error) error
}

// driver JWT 驱动实例
var driver Driver

// SetDriver 设置 JWT 驱动
func SetDriver(_driver Driver) {
	driver = _driver
}

// JWT 定义一个 jwt 对象
type JWT struct {
	// 秘钥，用以加密 JWT，读取配置信息 app.secret
	SignKey []byte
	// 刷新 token 的最大过期时间
	MaxRefresh time.Duration
}

// JWTCustomClaims 自定义载荷
type JWTCustomClaims struct {
	UserID       uint64 `json:"user_id"`
	UserName     string `json:"user_name"`
	ExpireAtTime int64  `json:"expire_time"`

	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwtPKG.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("jwt.sign_key")),
		MaxRefresh: time.Duration(config.GetInt64("jwt.max_refresh_time")) * time.Minute,
	}
}

// InitTokenBlacklist 初始化 token 黑名单
func InitTokenBlacklist() {
	jwt := NewJWT()
	err := driver.TraverseInvalidTokensNotExpired(func(claims *JWTCustomClaims) error {
		jwt.putTokenInBlacklist(claims)
		return nil
	})
	logger.LogIf(err)
}

// ParserToken 解析 token，中间件中调用
func (jwt *JWT) ParserToken(c *gin.Context) (*JWTCustomClaims, error) {

	// 1. 从请求中获取 token
	tokenString, parseErr := jwt.getTokenFromRequest(c)
	if parseErr != nil {
		return nil, parseErr
	}

	// 2. 预处理 token
	tokenString, parseErr = driver.PreParserToken(tokenString, c)
	if parseErr != nil {
		return nil, parseErr
	}

	// 3. 调用 jwt 库解析 token
	token, err := jwt.parseTokenString(tokenString)
	// 解析出错
	if err != nil {
		// 判断具体错误
		validationErr, ok := err.(*jwtPKG.ValidationError)
		if ok {
			if validationErr.Errors == jwtPKG.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtPKG.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		// 解析失败，token 失效
		return nil, ErrTokenInvalid
	}

	// 4. 将 token 中的 claims 信息解析出来和 JWTCustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		// 判断 token 是否已在黑名单，这种 token 也属于已失效 token
		if jwt.tokenInBlacklist(claims) {
			return nil, ErrTokenInvalid
		}
		// 解析成功，返回载荷数据
		return claims, nil
	}

	// 解析失败
	return nil, ErrTokenInvalid
}

// RefreshToken 更新 token，用以提供 refresh token 接口
func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {

	// 1. 从请求中获取 token
	tokenString, parseErr := jwt.getTokenFromRequest(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 预处理 token
	tokenString, parseErr = driver.PreParserToken(tokenString, c)
	if parseErr != nil {
		return "", parseErr
	}

	// 3. 调用 jwt 库解析token
	token, err := jwt.parseTokenString(tokenString)
	// 解析出错
	if err != nil {
		validationErr, ok := err.(*jwtPKG.ValidationError)
		// token 已过期时仍可以尝试刷新 token，发生其它错误时不能刷新 token
		if !ok || validationErr.Errors != jwtPKG.ValidationErrorExpired {
			return "", err
		}
	}

	// 4. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*JWTCustomClaims)

	// 5. 检查是否过了『最大允许刷新的时间』
	if x := app.TimeNowInTimezone().Add(-jwt.MaxRefresh).Unix(); claims.IssuedAt <= x {
		// 刷新 token 失败，最大可刷新时间已过期
		return "", ErrTokenExpiredMaxRefresh
	}

	// 6. 判断 token 是否已在黑名单，这种 token 也属于已失效 token
	if jwt.tokenInBlacklist(claims) {
		return "", ErrTokenInvalid
	}

	// 7. 设置新 token 过期时间
	expireAtTime := jwt.expireAtTime()
	claims.ExpireAtTime = expireAtTime
	claims.StandardClaims.ExpiresAt = expireAtTime

	// 8. 设置新 token id
	var id string
	id, err = driver.NewTokenID(claims)
	if err != nil {
		return "", err
	}
	claims.Id = id

	// 9. 生成新 token
	var newToken string
	newToken, err = jwt.createToken(*claims)
	if err != nil {
		return "", err
	}

	// 10. 保存新 token
	if err = driver.SaveToken(newToken, claims); err != nil {
		return "", err
	}

	// 返回新 token
	return newToken, nil
}

// MakeToken 生成 token
func (jwt *JWT) MakeToken(userID uint64, userName string) string {

	// 1. 构造用户 claims 信息(负荷)
	expireAtTime := jwt.expireAtTime()
	claims := JWTCustomClaims{
		userID,
		userName,
		expireAtTime,
		jwtPKG.StandardClaims{
			NotBefore: app.TimeNowInTimezone().Unix(), // 签名生效时间
			IssuedAt:  app.TimeNowInTimezone().Unix(), // 首次签名时间（后续刷新 token 不会更新）
			ExpiresAt: expireAtTime,                   // 签名过期时间
			Issuer:    config.GetString("app.name"),   // 签名颁发者
		},
	}

	// 2. 设置新 token id
	id, err := driver.NewTokenID(&claims)
	if err != nil {
		logger.LogIf(err)
		return ""
	}
	claims.Id = id

	// 3. 根据 claims 生成 token 对象
	token, err := jwt.createToken(claims)
	if err != nil {
		logger.LogIf(err)
		return ""
	}

	// 4. 保存 token
	if err = driver.SaveToken(token, &claims); err != nil {
		logger.LogIf(err)
		return ""
	}

	return token
}

// createToken 创建 token，内部使用，外部请调用 MakeToken
func (jwt *JWT) createToken(claims JWTCustomClaims) (string, error) {
	// 使用 HS256 算法进行 token 生成
	token := jwtPKG.NewWithClaims(jwtPKG.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}

// expireAtTime 过期时间
func (jwt *JWT) expireAtTime() int64 {

	timeNow := app.TimeNowInTimezone()

	var expireTime int64
	if config.GetBool("app.debug") {
		expireTime = config.GetInt64("jwt.debug_expire_time")
	} else {
		expireTime = config.GetInt64("jwt.expire_time")
	}

	expire := time.Duration(expireTime) * time.Minute
	return timeNow.Add(expire).Unix()
}

// parseTokenString 使用 jwtPKG.ParseWithClaims 解析 token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtPKG.Token, error) {
	return jwtPKG.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtPKG.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

// getTokenFromHeader 从请求头中获取 token
// 优先提取 token
// 其次提取 Authorization:Bearer xxxxx
func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("token")
	if authHeader == "" {
		authHeader = c.Request.Header.Get("Authorization")
		if authHeader == "" {
			return "", ErrHeaderEmpty
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return "", ErrHeaderMalformed
		}
		return parts[1], nil
	}
	return authHeader, nil
}

// getTokenFromRequest 从请求中提取 token
// 优先从请求查询参数中提取，其次才人请求头部中提取
func (jwt *JWT) getTokenFromRequest(c *gin.Context) (string, error) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		authHeader, _ := jwt.getTokenFromHeader(c)
		if authHeader == "" {
			return "", ErrTokenMissing
		}
		return authHeader, nil
	}
	return tokenStr, nil
}

// Invalidate 使 token 失效
func (jwt *JWT) Invalidate(claims *JWTCustomClaims) bool {

	// 标记 token 失效
	if err := driver.InvalidateToken(claims); err != nil {
		logger.LogIf(err)
		return false
	}

	// 将 token 放入黑名单
	return jwt.putTokenInBlacklist(claims)
}

// PutTokenInBlacklist 将 token 放入黑名单
func (jwt *JWT) putTokenInBlacklist(claims *JWTCustomClaims) bool {
	return redis.
		Connection("jwt").
		Set(
			jwt.generateTokenKey(claims),
			0,
			time.Unix(claims.ExpireAtTime, 0).Sub(app.TimeNowInTimezone()),
		)
}

// generateTokenKey 生成缓存 token 的 key
func (jwt *JWT) generateTokenKey(claims *JWTCustomClaims) string {
	return config.GetString("app.name") + ":invalid-token:" + claims.Id
}

// tokenInBlacklist token 在黑名单里
func (jwt *JWT) tokenInBlacklist(claims *JWTCustomClaims) bool {
	return redis.Connection("jwt").Has(jwt.generateTokenKey(claims))
}
