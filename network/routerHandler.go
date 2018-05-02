package network

import (
	"net/http"
)

type RouterHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	DispatchMsg(msg MsgData) error
	RegisterMsg(msgId string, handler func(MsgData)) error
	SendJsonMsg(msg interface{}) error

	DispatchJsonMsg(msg interface{}) error
	RegisterJsonMsg(msgId string, handler func(interface{})) error
	SendMsg(msg MsgData) error

	setServer(svr *Server)
}
