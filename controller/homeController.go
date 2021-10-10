package controller

import (
	"frame/context"
	"net/http"
)

func HomePage(c *context.Conetxt) {
	c.W.WriteHeader(http.StatusOK)
	c.W.Write([]byte("hello home"))
}
