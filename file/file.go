// Package file 文件操作辅助函数
// @author 余海坚 haijianyu10@qq.com
// @created 2022-04-23 16:57
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package file

import (
	"io/ioutil"
	"os"
)

// Put 将数据存入文件
func Put(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Exists 判断文件是否存在
func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}