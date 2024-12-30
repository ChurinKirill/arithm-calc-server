package main

import (
	logger "arithm-calc-server/logger"
	server "arithm-calc-server/server"
)

func main() {
	err := logger.InitializeLogger("../log.txt")
	if err != nil {
		return
	}
	server.StartServer()
	logger.ShutdownLogger()
}
