// Package client http 客户端包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-17 17:29
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package client

import (
	"encoding/json"
	"net/url"

	"github.com/spf13/cast"
	"github.com/valyala/fasthttp"
)

// Post 请求
func Post(url string, param interface{}) (string, int, error) {
	return Request("POST", url, param)
}

// Get 请求
func Get(url string, param interface{}) (string, int, error) {
	return Request("GET", url, param)
}

// Patch 请求
func Patch(url string, param interface{}) (string, int, error) {
	return Request("PATCH", url, param)
}

// Put 请求
func Put(url string, param interface{}) (string, int, error) {
	return Request("PUT", url, param)
}

// Delete 请求
func Delete(url string, param interface{}) (string, int, error) {
	return Request("DELETE", url, param)
}

// Request 请求
func Request(method, url string, param interface{}) (string, int, error) {

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")

	buf, err := json.Marshal(param)
	if err != nil {
		return "", 500, err
	}
	req.SetBody(buf)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return "", 500, err
	}

	b := resp.Body()

	return string(b), resp.Header.StatusCode(), nil
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
