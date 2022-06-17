// Package client http 客户端包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-06-17 17:29
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package client

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
)

// Post 请求
func Post(url string, param interface{}) (string, error) {
	return Request("POST", url, param)
}

// Get 请求
func Get(url string, param interface{}) (string, error) {
	return Request("GET", url, param)
}

// Request 请求
func Request(method, url string, param interface{}) (string, error) {

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")

	buf, err := json.Marshal(param)
	if err != nil {
		return "", err
	}
	req.SetBody(buf)

	err = fasthttp.Do(req, resp)

	if err != nil {
		return "", err
	}

	b := resp.Body()

	return string(b), nil
}
