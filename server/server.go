package server

import (
	"frame/context"
	"frame/middleware"
	"net/http"
)

type Routable interface {
	Router(method, pattern string, handleFunn middleware.HandleFun)
}

type Server interface {
	Routable
	Start(address string) error
}

type Handler interface {
	ServeHTTP(c *context.Conetxt)
	Routable
}

type HandlerBasedOnMap struct {
	handlers map[string]func(c *context.Conetxt)
}

func (h *HandlerBasedOnMap) Key(method, pattern string) string {
	return method + "#" + pattern
}

func (h *HandlerBasedOnMap) Router(method, pattern string, handleFun middleware.HandleFun) {
	key := h.Key(method, pattern)
	h.handlers[key] = handleFun
}

func (h *HandlerBasedOnMap) ServeHTTP(c *context.Conetxt) {
	key := h.Key(c.R.Method, c.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("404 - not found"))
	}
}

func NewHandlder() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(c *context.Conetxt)),
	}
}

type sdkHttpServer struct {
	handler Handler
	root    middleware.Filter
}

func (s *sdkHttpServer) Router(method, pattern string, handleFun middleware.HandleFun) {
	s.handler.Router(method, pattern, handleFun)
}

func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		c := context.NewContext(rw, r)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)

}

func NewServer(builders ...middleware.FilterBuilder) Server {
	handers := NewHandlder()
	var root middleware.Filter
	root = handers.ServeHTTP
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}
	return &sdkHttpServer{
		handler: handers,
		root:    root,
	}
}
