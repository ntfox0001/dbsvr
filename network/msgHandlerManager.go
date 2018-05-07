package network

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/ntfox0001/dbsvr/dberror"
	"github.com/ntfox0001/dbsvr/msgData"
)

type MsgHandlerManager interface {
	// 首先调用，返回值控制是否继续
	CheckConn(w http.ResponseWriter, r *http.Request) bool
	Initial(mh MsgHandler) bool
	Release(mh MsgHandler)
}

type MsgHandler interface {
	ReadMessage() (messageType int, p []byte, err error)
	SendJsonMsg(msg interface{}) error
	DispatchJsonMsg(msg interface{}) error
	RegisterJsonMsg(msgId string, handler func(interface{})) error

	SendMsg(msg MsgData) error
	RegisterMsg(msgId string, handler func(MsgData)) error
	DispatchMsg(msg MsgData) error
}

func processMsg(h MsgHandler) (rtErr error) {
	defer func() {
		if err := recover(); err != nil {
			rtErr = dberror.NewCommErr(err.(error).Error(), NetErrorProcessMsg)
		}
	}()
	var headMsg msgData.MsgHead
	mt, msg, err := h.ReadMessage()
	if err != nil {
		// 读取错误，直接断开
		log.Error("network:", "readMessageErr:", err.Error())
		rtErr = dberror.NewCommErr(err.Error(), NetErrorReadMsg)
		return
	}
	if mt == websocket.TextMessage {
		// is json msg
		var jsonMsg interface{}
		err := json.Unmarshal(msg, &jsonMsg)
		if err == nil {
			if err := h.DispatchJsonMsg(jsonMsg); err != nil {
				//逻辑错误

			}
		} else {
			// 解析错误直接断开
			log.Error("network", "invalid json format", err.Error())
			return dberror.NewCommErr(err.Error(), NetErrorUnmarshal)
		}
	} else if mt == websocket.BinaryMessage {
		// is protobuf msg
		if headMsg.MsgName == "" {
			if err := headMsg.Unmarshal(msg); err != nil {
				log.Error("network", err.Error())
				headMsg.Reset()
			}
		} else {
			refType := proto.MessageType(headMsg.MsgName)
			v := reflect.New(refType)
			h.DispatchMsg(v.Elem().Interface().(MsgData))
			headMsg.Reset()
		}

	}
	rtErr = nil
	return
}
