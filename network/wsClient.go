package network

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type wsClientConnection struct {
	wsConn     websocket.Conn
	respWriter http.ResponseWriter
}
