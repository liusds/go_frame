package server

import (
	"frame/context"
	"frame/middleware"
	"net/http"
	"strings"
)

type HandlerBased struct {
	root *node
}

type node struct {
	path     string
	children []*node

	handler middleware.HandleFun
}

func (h *HandlerBased) ServeHTTP(c *context.Conetxt) {
	handler, found := h.findRouter(c.R.URL.Path)
	if !found {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not - found"))
		return
	}
	handler(c)
}

func (h *HandlerBased) Router(method, pattern string, handleFun middleware.HandleFun) {
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")
	cur := h.root

	for index, path := range paths {
		mathChild, ok := h.findMatchChild(cur, path)
		if ok {
			cur = mathChild
		} else {
			h.createSub(cur, paths[index:], handleFun)
			return
		}
	}
}

func (h *HandlerBased) createSub(root *node, paths []string, handleFunc middleware.HandleFun) {
	cur := root
	for _, path := range paths {
		n := newNode(path)
		cur.children = append(cur.children, n)
		cur = n
	}
	cur.handler = handleFunc
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0, 2),
	}
}

func (h *HandlerBased) findMatchChild(root *node, path string) (*node, bool) {
	var wildcardNode *node
	for _, child := range root.children {
		if child.path == path && child.path != "*" {
			return child, true
		}
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
}

func (h *HandlerBased) findRouter(path string) (middleware.HandleFun, bool) {
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur := h.root
	for _, p := range paths {
		matchChild, found := h.findMatchChild(cur, p)
		if !found {
			return nil, false
		}
		cur = matchChild
	}
	if cur.handler == nil {
		return nil, false
	}
	return cur.handler, true
}
