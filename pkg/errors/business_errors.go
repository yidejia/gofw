package errors

import "net/http"

// errorBadRequest 请求资源格式错误
type errorBadRequest struct {
	message string
	err error
}

func NewErrorBadRequest(err error, message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorBadRequest{message[0], err}
	} else {
		return &errorBadRequest{"请求资源格式不正确", err}
	}
}

func (err *errorBadRequest) HttpStatus() int {
	return http.StatusBadRequest
}

func (err *errorBadRequest) Message() string {
	return err.message
}

func (err *errorBadRequest) Error() error {
	return err.err
}

func (err *errorBadRequest) Errors() map[string][]string  {
	return nil
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

func (err *errorUnauthorized) Message() string {
	return err.message
}

func (err *errorUnauthorized) Error() error {
	return nil
}

func (err *errorUnauthorized) Errors() map[string][]string  {
	return nil
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

func (err *errorForbidden) Message() string {
	return err.message
}

func (err *errorForbidden) Error() error {
	return nil
}

func (err *errorForbidden) Errors() map[string][]string  {
	return nil
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

func (err *errorNotFound) Message() string {
	return err.message
}

func (err *errorNotFound) Error() error {
	return nil
}

func (err *errorNotFound) Errors() map[string][]string  {
	return nil
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

func (err *errorMethodNotAllowed) Message() string {
	return err.message
}

func (err *errorMethodNotAllowed) Error() error {
	return nil
}

func (err *errorMethodNotAllowed) Errors() map[string][]string  {
	return nil
}

// errorUnprocessableEntity 资源无法处理错误
type errorUnprocessableEntity struct {
	message string
	errors map[string][]string
}

func NewErrorUnprocessableEntity(errors map[string][]string, message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorUnprocessableEntity{message[0], errors}
	} else {
		return &errorUnprocessableEntity{"资源无法处理", errors}
	}
}

func (err *errorUnprocessableEntity) HttpStatus() int {
	return http.StatusUnprocessableEntity
}

func (err *errorUnprocessableEntity) Message() string {
	return err.message
}

func (err *errorUnprocessableEntity) Error() error {
	return nil
}

func (err *errorUnprocessableEntity) Errors() map[string][]string  {
	return err.errors
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

func (err *errorLocked) Message() string {
	return err.message
}

func (err *errorLocked) Error() error {
	return nil
}

func (err *errorLocked) Errors() map[string][]string  {
	return nil
}

// errorInternal 内部错误
type errorInternal struct {
	message string
	err error
}

func NewErrorInternal(err error, message ...string) ResponsiveError {
	if len(message) > 0 {
		return &errorInternal{message[0], err}
	} else {
		return &errorInternal{"请求资源时发生系统错误", err}
	}
}

func (err *errorInternal) HttpStatus() int {
	return http.StatusInternalServerError
}

func (err *errorInternal) Message() string {
	return err.message
}

func (err *errorInternal) Error() error {
	return err.err
}

func (err *errorInternal) Errors() map[string][]string  {
	return nil
}
