package main

import (
	"common/boot"
	common_server "common/server"
	"todo-backend/internal/server"
)

func main() {

	appServer := server.NewServer()

	httpServer := common_server.New(8088)
	httpServer.Handler = appServer.RegisterRoutes()

	boot.Run(appServer, httpServer)

}
