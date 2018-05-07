package network

import (
	"runtime/debug"

	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/ntfox0001/dbsvr/dberror"
	"github.com/ntfox0001/dbsvr/msgData"
)

type msgHandler struct {
	conn       *websocket.Conn
	msgMap     map[string]func(msg MsgData)
	jsonMsgMap map[string]func(msg interface{})
}

func (h *msgHandler) ReadMessage() (messageType int, p []byte, err error) {
	return h.conn.ReadMessage()
}
func (h *msgHandler) SendMsg(msg MsgData) error {
	head := msgData.MsgHead{
		MsgName: msg.GetMsgName(),
	}

	if headBuf, err := head.Marshal(); err == nil {
		err := h.conn.WriteMessage(websocket.BinaryMessage, headBuf)
		if err != nil {
			log.Error("network", err.Error())
			return err
		} else {
			if msgBuf, err := msg.Marshal(); err == nil {
				err := h.conn.WriteMessage(websocket.BinaryMessage, msgBuf)
				if err != nil {
					log.Error("network", err.Error())
					return err
				}
			} else {
				log.Error("network", err.Error())
				return err
			}
		}
	} else {
		log.Error("network", err.Error())
		return err
	}

	return nil
}
func (h *msgHandler) SendJsonMsg(msg interface{}) error {
	h.conn.WriteJSON(msg)
	return nil
}

func (h *msgHandler) DispatchJsonMsg(msg interface{}) error {

	if msgId, ok := msg.(map[string]interface{})["msgId"]; !ok {
		log.Error("network", "RoouterWSHandler not exist msgId.")
		return dberror.NewCommErr("RoouterWSHandler not exist msgId.", NetErrorUnknowMsg)
	} else {
		if handler, ok := h.jsonMsgMap[msgId.(string)]; ok {
			// handler应该有缓存处理
			if err := callJsonFunc(handler, msg); err != nil {
				return err
			}
		}
	}
	return nil
}
func callJsonFunc(handler func(interface{}), msg interface{}) (rterr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("network", "logic error:", err.(error).Error(), msg)
			debug.PrintStack()
			rterr = dberror.NewCommErr("logic error:"+err.(error).Error(), NetErrorLogic)
		}
	}()
	handler(msg)
	return nil
}
func (h *msgHandler) RegisterJsonMsg(msgId string, handler func(interface{})) error {
	if _, ok := h.jsonMsgMap[msgId]; ok {
		return dberror.NewCommErr("msgId has exist:"+msgId, NetErrorExistMsg)
	} else {
		h.jsonMsgMap[msgId] = handler
	}
	return nil
}
func (h *msgHandler) DispatchMsg(msg MsgData) error {
	msgId := msg.GetMsgName()
	if msgId == "" {
		log.Error("network", "RoouterWSHandler msgId is nil.")
		return dberror.NewCommErr("RoouterWSHandler msgId is nil.", NetErrorUnknowMsg)
	}
	if handler, ok := h.msgMap[msgId]; ok {
		// handler应该有缓存处理
		if err := callFunc(handler, msg); err != nil {
			return err
		}
	}
	return nil
}
func callFunc(handler func(MsgData), msg MsgData) (rterr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("network", "logic error:", err.(error).Error(), msg)
			debug.PrintStack()
			rterr = dberror.NewCommErr("logic error:"+err.(error).Error(), NetErrorLogic)
		}
	}()
	handler(msg)
	return nil
}
func (h *msgHandler) RegisterMsg(msgId string, handler func(MsgData)) error {
	if _, ok := h.msgMap[msgId]; ok {
		return dberror.NewCommErr("msgId has exist:"+msgId, NetErrorExistMsg)
	} else {
		h.msgMap[msgId] = handler
	}
	return nil
}
