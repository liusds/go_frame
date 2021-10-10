package main

import (
	"frame/server"
)

func main() {
	server := server.NewServer()
	if err := server.Start(":8080"); err != nil {
		panic("服务启动失败：" + err.Error())
	}
}
