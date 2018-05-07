package network

import (
	"fmt"
	"net/http"
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

func (f *RouterFileHandler) setServer(svr *Server) {
	f.server = svr
}
