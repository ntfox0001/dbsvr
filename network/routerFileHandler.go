package network

import (
	"fmt"
	"net/http"

	"github.com/ntfox0001/dbsvr/dberror"
)

type RouterFileHandler struct {
	server      *Server
	path        string
	defaultFile string
}

func NewRouterFileHandler(path string, defFile string) *RouterFileHandler {
	return &RouterFileHandler{
		path:        path,
		defaultFile: defFile,
	}
}

func (f *RouterFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != f.path {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println(r)
	http.ServeFile(w, r, f.defaultFile)
}

func (f *RouterFileHandler) DispatchMsg(msg MsgData) error {
	return dberror.NewStringErr("not Implement")
}
func (f *RouterFileHandler) RegisterMsg(msgId string, handler func(MsgData)) error {
	return dberror.NewStringErr("not Implement")
}
func (f *RouterFileHandler) SendJsonMsg(msg interface{}) error {
	return dberror.NewStringErr("not Implement")
}

func (f *RouterFileHandler) DispatchJsonMsg(msg interface{}) error {
	return dberror.NewStringErr("not Implement")
}
func (f *RouterFileHandler) RegisterJsonMsg(msgId string, handler func(interface{})) error {
	return dberror.NewStringErr("not Implement")
}
func (f *RouterFileHandler) SendMsg(msg MsgData) error {
	return dberror.NewStringErr("not Implement")
}

func (f *RouterFileHandler) setServer(svr *Server) {
	f.server = svr
}
