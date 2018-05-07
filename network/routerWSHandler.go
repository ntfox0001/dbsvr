package network

import (
	"context"

	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/ntfox0001/dbsvr/dberror"
)

type RouterWSHandler struct {
	upgrader websocket.Upgrader
	server   *Server
	msghMgr  MsgHandlerManager
}

func NewRouterWSHandler(msghMgr MsgHandlerManager) *RouterWSHandler {

	return &RouterWSHandler{
		upgrader: websocket.Upgrader{},
		msghMgr:  msghMgr,
	}
}

func (h *RouterWSHandler) setServer(svr *Server) {
	h.server = svr
}

func (h *RouterWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.msghMgr.CheckConn(w, r) == false {
		return
	}
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("network", "upgradeError:", err.Error())
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	ac := newAcceptConn(ctx, c)
	h.msghMgr.Initial(ac)
	defer func() {
		cancel()
		c.Close()
		ac.onError(dberror.NewCommErr("Conn closed.", NetErrorConnectClosed))
		h.msghMgr.Release(ac)
	}()
	for {
		if err := processMsg(ac); err != nil {
			break
		}
	}
}

// func (h *RouterWSHandler) processMsg(ac *AcceptConn) (rtErr error) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			rtErr = dberror.NewCommErr(err.(error).Error(), NetErrorProcessMsg)
// 		}
// 	}()
// 	var headMsg msgData.MsgHead
// 	mt, msg, err := ac.conn.ReadMessage()
// 	if err != nil {
// 		// 读取错误，直接断开
// 		log.Error("network:", "readMessageErr:", err.Error())
// 		rtErr = dberror.NewCommErr(err.Error(), NetErrorReadMsg)
// 		return
// 	}
// 	if mt == websocket.TextMessage {
// 		// is json msg
// 		var jsonMsg interface{}
// 		err := json.Unmarshal(msg, &jsonMsg)
// 		if err == nil {
// 			if err := ac.DispatchJsonMsg(jsonMsg); err != nil {
// 				//逻辑错误

// 			}
// 		} else {
// 			// 解析错误直接断开
// 			log.Error("network", "invalid json format", err.Error())
// 			return dberror.NewCommErr(err.Error(), NetErrorUnmarshal)
// 		}
// 	} else if mt == websocket.BinaryMessage {
// 		// is protobuf msg
// 		if headMsg.MsgName == "" {
// 			if err := headMsg.Unmarshal(msg); err != nil {
// 				log.Error("network", err.Error())
// 				headMsg.Reset()
// 			}
// 		} else {
// 			refType := proto.MessageType(headMsg.MsgName)
// 			v := reflect.New(refType)
// 			ac.DispatchMsg(v.Elem().Interface().(MsgData))
// 			headMsg.Reset()
// 		}

// 	}
// 	rtErr = nil
// 	return
// }
