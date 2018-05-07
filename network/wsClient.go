package network

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	msgHandler
	resp   *http.Response
	quitCh chan interface{}
}

func NewWsClient(url url.URL, header http.Header) (*WsClient, error) {
	conn, resp, err := websocket.DefaultDialer.Dial(url.String(), header)
	client := &WsClient{

		resp:   resp,
		quitCh: make(chan interface{}),
		msgHandler: msgHandler{
			conn:       conn,
			msgMap:     make(map[string]func(MsgData)),
			jsonMsgMap: make(map[string]func(interface{})),
		},
	}

	if err == nil {
		return client, nil
	} else {
		return nil, err
	}
}

func (w *WsClient) Start() {
	defer func() {
		close(w.quitCh)
	}()

	for {
		if err := processMsg(w); err != nil {
			break
		}
	}
}
