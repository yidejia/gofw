// Package events 事件包
// @author 余海坚 haijianyu10@qq.com
// @created 2022-07-19 11:53
// @copyright © 2010-2022 广州伊的家网络科技有限公司
package events

// Event 事件接口
type Event interface {
	// EventCode 事件编码
	EventCode() string
}

// EventListener 事件监听器接口
type EventListener interface {
	// Handle 处理事件
	Handle(event Event)
}

var listen = make(map[string][]EventListener)

// AddListener 添加事件监听器
func AddListener(event Event, listener EventListener) {
	code := event.EventCode()
	listen[code] = append(listen[code], listener)
}

// Dispatch 分发事件
func Dispatch(event Event) {
	code := event.EventCode()
	for _, listener := range listen[code] {
		go listener.Handle(event)
	}
}

// DispatchSync 同步分发事件
func DispatchSync(event Event) {
	code := event.EventCode()
	for _, listener := range listen[code] {
		listener.Handle(event)
	}
}
