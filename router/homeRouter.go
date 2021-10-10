package router

import (
	"frame/controller"
	"frame/server"
	"net/http"
)

func HomeRouter() {
	hr := server.NewServer()
	hr.Router(http.MethodGet, "/", controller.HomePage)
}
