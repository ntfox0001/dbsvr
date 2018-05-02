package network

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/inconshreveable/log15"

	"github.com/gorilla/mux"
)

type Server struct {
	ip        string
	port      uint
	router    *mux.Router
	routerMap map[string]*RouterHandler
	quitCh    chan interface{}
	server    *http.Server
}

func NewServer(ip string, port uint) *Server {

	svr := &Server{
		ip:        ip,
		port:      port,
		router:    mux.NewRouter(),
		routerMap: make(map[string]*RouterHandler),
	}

	return svr
}

func (s *Server) Start() {
	addr := fmt.Sprintf("%s:%d", s.ip, s.port)
	s.server = &http.Server{Addr: addr, Handler: s.router}
	go func() {
		go s.server.ListenAndServe()
		<-s.quitCh

		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Error("network", "shutdown:", err.Error())
		}
	}()
}

func (s *Server) Close() {
	s.quitCh <- struct{}{}
}

func (s *Server) RegisterRouter(router string, handler RouterHandler) {
	handler.setServer(s)
	s.router.Handle(router, handler)
}

func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) Ip() string {
	return s.ip
}
func (s *Server) Port() uint {
	return s.port
}
