package requests

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
	"github.com/yidejia/gofw/pkg/events"
	"github.com/yidejia/gofw/pkg/response"
)

// Dispatchable 可分发事件接口
type Dispatchable interface {
	Dispatch() events.Event
}

type IEventRequest interface {
	Validatable
	Signable
	Dispatchable
}

type NewEventRequestFunc func() IEventRequest

// eventRequests 事件请求集合
var eventRequests = make(map[string]NewEventRequestFunc)

// RegisterEventRequest 注册事件请求
func RegisterEventRequest(event events.Event, f NewEventRequestFunc) {
	eventRequests[event.EventCode()] = f
}

// NewEventRequest 创建事件请求
func NewEventRequest(c *gin.Context, eventCode string) IEventRequest {
	if f, ok := eventRequests[eventCode]; ok {
		return f()
	}
	response.ValidationError(c, map[string][]string{"code": {"无效事件"}})
	return nil
}

// EventRequest 事件请求
type EventRequest struct {
	SignRequest
	Code string `json:"code,omitempty" form:"code" valid:"code"` // 编码
}

// Validate 验证请求
func (req *EventRequest) Validate(extra ...interface{}) map[string][]string {

	rules := govalidator.MapData{
		"code": []string{"required", "min:1", "max:255", "exists:events,code"},
	}

	messages := govalidator.MapData{
		"code": []string{
			"required:事件编码为必填项",
			"min:事件编码长度至少为 1",
			"max:事件编码长度至少为 255",
			fmt.Sprintf("exists:事件编码 %s 不存在", req.Code),
		},
	}

	return req.MergeValidateErrors(
		req.ValidateStruct(req, rules, messages),
		req.SignRequest.Validate(extra...),
	)
}

// ParamsToSign 返回用于签名验证的请求参数
func (req *EventRequest) ParamsToSign() map[string]interface{} {
	params := req.SignRequest.ParamsToSign()
	params["code"] = req.Code
	return params
}
