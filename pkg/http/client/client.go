// Package client http 客户端包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-17 17:29
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package client

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
)

// RequestOptions 请求选项
type RequestOptions struct {
	WithForm bool // 通过表单提交参数
}

// RequestOption 请求选项设置函数
type RequestOption func(*RequestOptions)

// NewRequestOptions 新建请求选项
func NewRequestOptions(options ...RequestOption) *RequestOptions {
	requestOption := &RequestOptions{}
	for _, option := range options {
		option(requestOption)
	}
	return requestOption
}

// RequestWithForm 设置通过表单提交参数
func RequestWithForm() RequestOption {
	return func(options *RequestOptions) {
		options.WithForm = true
	}
}

// Request 请求
func Request(method, url string, param interface{}, options ...*RequestOptions) (string, int, error) {

	var _options *RequestOptions
	if len(options) > 0 {
		_options = options[0]
	} else {
		_options = NewRequestOptions()
	}

	// 通过表单提交参数
	if _options.WithForm {
		// 表单数据需要以 key-value 映射对的形式提供
		form, ok := param.(map[string]interface{})
		if !ok {
			return "", 500, errors.New("提交表单数据格式不正确")
		}
		if method == fasthttp.MethodPost {
			return postForm(url, form)
		}
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")

	buf, err := json.Marshal(param)
	if err != nil {
		return "", 500, err
	}
	req.SetBody(buf)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err = fasthttp.Do(req, resp); err != nil {
		return "", 500, err
	}

	return string(resp.Body()), resp.Header.StatusCode(), nil
}

// Post 请求
func Post(url string, param interface{}, options ...*RequestOptions) (string, int, error) {
	return Request(fasthttp.MethodPost, url, param, options...)
}

// Get 请求
func Get(url string, param interface{}, options ...*RequestOptions) (string, int, error) {
	return Request(fasthttp.MethodGet, url, param, options...)
}

// Patch 请求
func Patch(url string, param interface{}, options ...*RequestOptions) (string, int, error) {
	return Request(fasthttp.MethodPatch, url, param, options...)
}

// Put 请求
func Put(url string, param interface{}, options ...*RequestOptions) (string, int, error) {
	return Request(fasthttp.MethodPut, url, param, options...)
}

// Delete 请求
func Delete(url string, param interface{}, options ...*RequestOptions) (string, int, error) {
	return Request(fasthttp.MethodDelete, url, param, options...)
}

// postForm 请求提交表单
func postForm(url string, form map[string]interface{}) (string, int, error) {

	args := &fasthttp.Args{}
	for k, v := range form {
		args.Add(k, cast.ToString(v))
	}

	status, resp, err := fasthttp.Post(nil, url, args)
	if err != nil {
		return "", 500, err
	}

	return string(resp), status, nil
}

// BuildQuery 构建查询字符串
func BuildQuery(params map[string]interface{}) (paramStr string) {
	q := (&url.URL{}).Query()
	for k, v := range params {
		q.Add(k, cast.ToString(v))
	}
	paramStr = q.Encode()
	return
}
