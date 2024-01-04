package client

import (
	"encoding/json"
	"errors"

	"github.com/yidejia/gofw/pkg/logger"
)

// JSONClient 返回 json 数据格式的 http 客户端
// @author 余海坚 haijianyu10@qq.com
// @created 2023-06-07 11:59
// @copyright © 2010-2023 广州伊的家网络科技有限公司
type JSONClient struct {
}

// NewJSONClient 生成返回 json 数据格式的 http 客户端实例
func NewJSONClient() *JSONClient {
	return &JSONClient{}
}

// Request 发送请求
func (c *JSONClient) Request(method, url string, params map[string]interface{}, result interface{}, options ...*RequestOptions) (status int, err error) {

	// 发送请求并获取响应字符串
	var resp string
	if resp, status, err = Request(method, url, params, options...); err != nil {
		return
	}

	// 解码响应字符串
	if err = json.Unmarshal([]byte(resp), &result); err != nil {
		logger.ErrorString("http 包 JSON 响应客户端", "发送请求", "Request 发送请求错误："+err.Error())
		err = errors.New("发送请求失败")
		return
	}

	return
}

// Post 发送 Post 请求
func (c *JSONClient) Post(url string, params map[string]interface{}, result interface{}, options ...*RequestOptions) (status int, err error) {
	status, err = c.Request("POST", url, params, result, options...)
	return
}

// Get 发送 Get 请求
func (c *JSONClient) Get(url string, params map[string]interface{}, result interface{}, options ...*RequestOptions) (status int, err error) {
	status, err = c.Request("GET", url, params, result, options...)
	return
}

// Patch 发送 Patch 请求
func (c *JSONClient) Patch(url string, params map[string]interface{}, result interface{}, options ...*RequestOptions) (status int, err error) {
	status, err = c.Request("PATCH", url, params, result, options...)
	return
}

// Put 发送 Put 请求
func (c *JSONClient) Put(url string, params map[string]interface{}, result interface{}, options ...*RequestOptions) (status int, err error) {
	status, err = c.Request("PUT", url, params, result, options...)
	return
}

// Delete 发送 Delete 请求
func (c *JSONClient) Delete(url string, params map[string]interface{}, result interface{}, options ...*RequestOptions) (status int, err error) {
	status, err = c.Request("DELETE", url, params, result, options...)
	return
}
