package network

import (
	"context"

	"github.com/gorilla/websocket"
	"github.com/ntfox0001/dbsvr/dberror"
)

type AcceptConn struct {
	msgHandler
	keyValMap map[string]interface{}
	ctx       context.Context
	onError   func(err dberror.CommError) int
}

func defNetError(err dberror.CommError) int {
	return NetErrRtBreak
}
func newAcceptConn(ctx context.Context, conn *websocket.Conn) *AcceptConn {
	return &AcceptConn{
		ctx:       ctx,
		keyValMap: make(map[string]interface{}),
		onError:   defNetError,
		msgHandler: msgHandler{
			conn:       conn,
			msgMap:     make(map[string]func(MsgData)),
			jsonMsgMap: make(map[string]func(interface{})),
		},
	}
}

func (a *AcceptConn) SetErrorFunc(errFunc func(err dberror.CommError) int) {
	a.onError = errFunc
}
func (a *AcceptConn) GetContext() context.Context {
	return a.ctx
}
