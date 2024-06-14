// Package middleware 存放系统中间件
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-22 17:14
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/yidejia/gofw/pkg/helpers"
	"github.com/yidejia/gofw/pkg/logger"
	gfReqs "github.com/yidejia/gofw/pkg/requests"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// RequestLog 请求日志
type RequestLog struct {
	UserID    uint64 `json:"user_id,omitempty" gorm:"type:bigint unsigned;not null;index:user_id;comment:用户id"`
	UserName  string `json:"user_name,omitempty" gorm:"type:varchar(30);not null;index:user_name;comment:登录时的名称"`
	IP        string `json:"ip,omitempty" gorm:"type:varchar(60);not null;index:ip;comment:IP地址"`
	UserAgent string `json:"user_agent,omitempty" gorm:"type:varchar(1000);not null;comment:用户代理"`
	Method    string `json:"method,omitempty" gorm:"type:varchar(10);not null;index:method;comment:请求方法"`
	URL       string `json:"url,omitempty" gorm:"type:varchar(1000);not null;index:url;comment:请求URL"`
	Query     string `json:"query,omitempty" gorm:"type:varchar(1000);not null;comment:查询参数"`
	Status    int    `json:"status,omitempty" gorm:"type:int unsigned;not null;index:status;comment:响应状态"`
	Request   string `json:"request,omitempty" gorm:"type:mediumtext;not null;comment:请求体"`
	Response  string `json:"response,omitempty" gorm:"type:mediumtext;not null;comment:响应体"`
	Errors    string `json:"errors,omitempty" gorm:"type:varchar(1000);not null;comment:错误"`
	TakeTime  string `json:"take_time,omitempty" gorm:"type:varchar(30);not null;comment:消耗时间(毫秒)"`
}

// RequestLogHandler 请求日志处理器接口
type RequestLogHandler interface {
	// HandleRequestLog 处理请求日志
	HandleRequestLog(log RequestLog)
}

// requestLogHandler 请求日志处理器
var requestLogHandler RequestLogHandler

// RegisterRequestLogHandler 注册请求日志处理器
func RegisterRequestLogHandler(handler RequestLogHandler) {
	requestLogHandler = handler
}

// Logger 记录请求日志
func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 获取 response 内容
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 获取请求数据
		var requestBody []byte
		if c.Request.Body != nil {
			// c.Request.Body 是一个 buffer 对象，只能读取一次
			requestBody, _ = ioutil.ReadAll(c.Request.Body)
			// 读取后，重新赋值 c.Request.Body ，以供后续的其他操作
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 设置开始时间
		start := time.Now()

		c.Next()

		// 计算请求消耗时间
		takeTime := time.Since(start)

		// 记录请求关键信息
		requestLog := RequestLog{
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Method:    c.Request.Method,
			Status:    c.Writer.Status(),
			Errors:    c.Errors.ByType(gin.ErrorTypePrivate).String(),
			TakeTime:  helpers.MicrosecondsStr(takeTime),
		}

		// 接口对当前请求 URL 进行过自定义处理后，会缓存在请求上下文中，一般是对请求查询参数进行了脱敏处理，例如需要隐藏 URL 中的 token
		reqURL := gfReqs.GetRequestURLLog(c)
		if len(reqURL) > 0 {
			requestLog.URL = reqURL
		} else {
			// 默认记录原始请求 URL
			requestLog.URL = c.Request.URL.String()
		}

		// 接口对当前请求查询参数进行过自定义处理后，会缓存在请求上下文中，一般是对请求查询参数进行了脱敏处理，例如需要隐藏查询参数中的 token
		reqQuery := gfReqs.GetRequestQueryLog(c)
		if len(reqQuery) > 0 {
			requestLog.Query = reqQuery
		} else {
			// 默认记录原始请求查询参数
			requestLog.Query = c.Request.URL.RawQuery
		}

		// 接口对当前请求进行过自定义处理后，会缓存在请求上下文中，一般是对请求进行了脱敏处理，例如登录接口一般需要隐藏登录密码
		reqLog := gfReqs.GetRequestLog(c)
		if len(reqLog) > 0 {
			requestLog.Request = reqLog
		} else {
			// 默认记录原始请求数据
			requestLog.Request = string(requestBody)
		}

		// 接口对当前响应进行过自定义处理后，会缓存在请求上下文中，一般是对响应进行了脱敏处理，例如登录接口一般需要隐藏响应中的访问令牌
		respLog := gfReqs.GetResponseLog(c)
		if len(respLog) > 0 {
			requestLog.Response = respLog
		} else {
			// 默认记录原始响应数据
			requestLog.Response = w.body.String()
		}

		// 记录请求用户信息
		requestLog.UserID = gfReqs.GetUserID(c)
		requestLog.UserName = gfReqs.GetUserName(c)

		// 日志字段
		logFields := []zap.Field{
			zap.Int("status", requestLog.Status),
			zap.String("request", requestLog.Method+" "+requestLog.URL),
			zap.String("query", requestLog.Query),
			zap.String("ip", requestLog.IP),
			zap.String("user-agent", requestLog.UserAgent),
			zap.String("errors", requestLog.Errors),
			zap.String("time", requestLog.TakeTime),
		}

		// 异常请求才记录请求体和响应体
		if !(requestLog.Method == "GET" && requestLog.Status < 400) {
			// 请求体
			logFields = append(logFields, zap.String("Request Body", requestLog.Request))
			// 响应体
			logFields = append(logFields, zap.String("Response Body", requestLog.Response))
		}

		if requestLog.Status > 400 && requestLog.Status <= 499 {
			// 除了 StatusBadRequest 以外，warning 提示一下，常见的有 403 404，开发时都要注意
			logger.Warn("HTTP Warning "+cast.ToString(requestLog.Status), logFields...)
		} else if requestLog.Status >= 500 && requestLog.Status <= 599 {
			// 除了内部错误，记录 error
			logger.Error("HTTP Error "+cast.ToString(requestLog.Status), logFields...)
		} else {
			// 常规调试
			logger.Debug("HTTP Access Log", logFields...)
		}

		// 调用请求日志处理器对日志进行额外处理
		if !gfReqs.RequestLogIsCleared(c) && requestLogHandler != nil {
			requestLogHandler.HandleRequestLog(requestLog)
		}
	}
}
