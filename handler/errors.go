package handler

import "errors"

type Error struct {
	Err error
}

func (e Error) Error() string {
	return e.Err.Error()
}

var WebContextNotFoundError = &Error{
	Err: errors.New("无法找到 WebContextMust，请检查 WebContextMiddleware 是否被正确"),
}

type ServerInternalError struct {
	code int
	Body interface{}
	err  error
}

func PanicServerInternalError(err ServerInternalError) {
	panic(err)
}
