package network

import (
	"net/http"
)

type RouterHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	setServer(svr *Server)
}
