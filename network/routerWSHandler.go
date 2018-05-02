package network

import (
	"reflect"

	"github.com/golang/protobuf/proto"

	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/inconshreveable/log15"
	"github.com/ntfox0001/dbsvr/dberror"
	"github.com/ntfox0001/dbsvr/msgData"
)

type RouterWSHandler struct {
	upgrader   websocket.Upgrader
	server     *Server
	conn       *websocket.Conn
	msgMap     map[string]func(msg MsgData)
	jsonMsgMap map[string]func(msg interface{})
}

func NewRouterWSHandler() *RouterWSHandler {

	return &RouterWSHandler{
		upgrader:   websocket.Upgrader{},
		msgMap:     make(map[string]func(MsgData)),
		jsonMsgMap: make(map[string]func(interface{})),
	}
}
func (h *RouterWSHandler) setServer(svr *Server) {
	h.server = svr
}
func (h *RouterWSHandler) SendMsg(msg MsgData) error {
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
func (h *RouterWSHandler) SendJsonMsg(msg interface{}) error {
	h.conn.WriteJSON(msg)
	return nil
}
func (h *RouterWSHandler) DispatchJsonMsg(msg interface{}) error {

	if msgId, ok := msg.(map[string]interface{})["msgId"]; !ok {
		log.Error("network", "RoouterWSHandler not exist msgId.")
		return dberror.NewStringErr("RoouterWSHandler not exist msgId.")
	} else {
		if handler, ok := h.jsonMsgMap[msgId.(string)]; ok {
			// handler应该有缓存处理
			callJsonFunc(handler, msg)
		}
	}
	return nil
}
func callJsonFunc(handler func(interface{}), msg interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("network", "msg is panic!", msg)
		}
	}()
	handler(msg)
}
func (h *RouterWSHandler) RegisterJsonMsg(msgId string, handler func(interface{})) error {
	if _, ok := h.jsonMsgMap[msgId]; ok {
		return dberror.NewStringErr("msgId has exist:" + msgId)
	} else {
		h.jsonMsgMap[msgId] = handler
	}
	return nil
}
func (h *RouterWSHandler) DispatchMsg(msg MsgData) error {
	msgId := msg.GetMsgName()
	if msgId == "" {
		log.Error("network", "RoouterWSHandler msgId is nil.")
		return dberror.NewStringErr("RoouterWSHandler msgId is nil.")
	}
	if handler, ok := h.msgMap[msgId]; ok {
		// handler应该有缓存处理
		callFunc(handler, msg)
	}
	return nil
}
func callFunc(handler func(MsgData), msg MsgData) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("network", "msg is panic!", msg)
		}
	}()
	handler(msg)
}
func (h *RouterWSHandler) RegisterMsg(msgId string, handler func(MsgData)) error {
	if _, ok := h.msgMap[msgId]; ok {
		return dberror.NewStringErr("msgId has exist:" + msgId)
	} else {
		h.msgMap[msgId] = handler
	}
	return nil
}

func (h *RouterWSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("network", "upgradeError:", err.Error())
		return
	}
	h.conn = c
	defer c.Close()
	for {
		if err := h.processMsg(c); err != nil {
			break
		}
	}
}

func (h *RouterWSHandler) processMsg(c *websocket.Conn) error {
	var headMsg msgData.MsgHead
	mt, msg, err := c.ReadMessage()
	if err != nil {
		log.Error("network:", "readMessageErr:", err.Error())
		return dberror.NewCommErr(err.Error(), 100)
	}
	if mt == websocket.TextMessage {
		// is json msg

		var jsonMsg interface{}
		err := json.Unmarshal(msg, &jsonMsg)
		if err == nil {
			h.DispatchJsonMsg(jsonMsg)
		} else {
			log.Error("network", "invalid json format", err.Error())
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
	return nil
}
