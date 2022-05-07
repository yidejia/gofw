package errors

import "net/http"

// errorBadRequest 请求资源格式错误
type errorBadRequest struct {
	message string
}

func NewErrorBadRequest(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorBadRequest{message[0]}
	} else {
		return &errorBadRequest{"请求资源格式不正确"}
	}
}

func (err *errorBadRequest) HttpStatus() int {
	return http.StatusBadRequest
}

func (err *errorBadRequest) Error() string {
	return err.message
}

// errorUnauthorized 用户未授权错误
type errorUnauthorized struct {
	message string
}

func NewErrorUnauthorized(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorUnauthorized{message[0]}
	} else {
		return &errorUnauthorized{"用户未授权"}
	}
}

func (err *errorUnauthorized) HttpStatus() int {
	return http.StatusUnauthorized
}

func (err *errorUnauthorized) Error() string {
	return err.message
}

// errorForbidden 无权请求错误
type errorForbidden struct {
	message string
}

func NewErrorForbidden(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorForbidden{message[0]}
	} else {
		return &errorForbidden{"无权请求资源"}
	}
}

func (err *errorForbidden) HttpStatus() int {
	return http.StatusForbidden
}

func (err *errorForbidden) Error() string {
	return err.message
}

// errorNotFound 请求资源不存在错误
type errorNotFound struct {
	message string
}

func NewErrorNotFound(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorNotFound{message[0]}
	} else {
		return &errorNotFound{"请求资源不存在"}
	}
}

func (err *errorNotFound) HttpStatus() int {
	return http.StatusNotFound
}

func (err *errorNotFound) Error() string {
	return err.message
}

// errorMethodNotAllowed 请求方法不允许错误
type errorMethodNotAllowed struct {
	message string
}

func NewErrorMethodNotAllowed(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorMethodNotAllowed{message[0]}
	} else {
		return &errorMethodNotAllowed{"请求方法不允许"}
	}
}

func (err *errorMethodNotAllowed) HttpStatus() int {
	return http.StatusMethodNotAllowed
}

func (err *errorMethodNotAllowed) Error() string {
	return err.message
}

// errorUnprocessableEntity 资源无法处理错误
type errorUnprocessableEntity struct {
	message string
}

func NewErrorUnprocessableEntity(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorUnprocessableEntity{message[0]}
	} else {
		return &errorUnprocessableEntity{"资源无法处理"}
	}
}

func (err *errorUnprocessableEntity) HttpStatus() int {
	return http.StatusUnprocessableEntity
}

func (err *errorUnprocessableEntity) Error() string {
	return err.message
}

// errorLocked 资源已锁定错误
type errorLocked struct {
	message string
}

func NewErrorLocked(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorLocked{message[0]}
	} else {
		return &errorLocked{"资源已锁定"}
	}
}

func (err *errorLocked) HttpStatus() int {
	return http.StatusLocked
}

func (err *errorLocked) Error() string {
	return err.message
}

// errorInternal 内部错误
type errorInternal struct {
	message string
}

func NewErrorInternal(message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorInternal{message[0]}
	} else {
		return &errorInternal{"资源已锁定"}
	}
}

func (err *errorInternal) HttpStatus() int {
	return http.StatusInternalServerError
}

func (err *errorInternal) Error() string {
	return err.message
}
